package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/google/go-github/v29/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type App struct {
	ctx          context.Context
	environment  string
	config       Config
	githubClient *github.Client
	log          *logrus.Logger
}

func main() {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Start server ... ")
	ctx := context.Background()

	config, err := readConfig(pathFileInput)
	if err != nil {
		log.Error("error", err)
	}

	app := &App{
		ctx:         ctx,
		environment: strings.TrimSpace(os.Getenv("APP_ENV")),
		config:      config,
		log:         log,
	}

	app.setUpGithubClient()

	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)
	router.Use(mux.CORSMethodMiddleware(router))

	router.HandleFunc("/generate-report", app.GenerateReportHandler).Methods("GET")
	router.HandleFunc("/", app.RootHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))

}

// GenerateReportHandler - Route to generate report
func (app *App) GenerateReportHandler(w http.ResponseWriter, r *http.Request) {
	app.getPackageJSONs(app.ctx)
	projects, components, projectsClientData, componentsClientData := app.splitProjectsComponents()

	app.log.Info("Generate data to graphs ... \n")
	countDependenciesByVersions := app.statsCountDependenciesByVersions(*projects)
	countComponentsByFilters := app.statsCountComponentsByFilters(*components, *componentsClientData)
	countProjectsByFilters := app.statsCountProjectsByFilters(*projects, *projectsClientData)

	clientData := &ClientData{
		GeneratedAt: time.Now(),
		Projects:    projectsClientData,
		Components:  componentsClientData,
		GraphData: map[string][]interface{}{
			"projectsByFilters":   countProjectsByFilters,
			"componentsByFilters": countComponentsByFilters,
		},
		DependenciesByVersions: countDependenciesByVersions,
	}

	//if app.environment != "production" {
	//	clientDataJSON, err := json.MarshalIndent(clientData, "", " ")
	//	err = ioutil.WriteFile(pathFileOutput, clientDataJSON, 0644)
	//
	//	if err == nil {
	//		app.log.Info("Output file generated and sent to " + pathFileOutput)
	//	}
	//}
	json.NewEncoder(w).Encode(clientData)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// RootHandler - Route to root
func (app *App) RootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

// getPackageJSONs - Get Package Json for each repository from config file
func (app *App) getPackageJSONs(ctx context.Context) {
	app.log.Printf("Get the package.json for the %d repositories ... \n", len(app.config.Repositories))

	// this buffered channel will block at the concurrency limit
	semaphoreChan := make(chan struct{}, concurrentLimit)

	// this channel will not block and collect the packageJSON results
	resultsChan := make(chan *PackageJSON)

	// make sure we close these channels when we're done with them
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	for i := 0; i < len(app.config.Repositories); i++ {

		// start a go routine with the index and url in a closure
		go func(i int, repository Repository) {

			// this sends an empty struct into the semaphoreChan which
			// is basically saying add one to the limit, but when the
			// limit has been reached block until there is room
			semaphoreChan <- struct{}{}

			// send the request and put the response in a result struct
			// along with the index so we can sort them later along with
			// any error that might have occoured
			info, err := splitRepositoryURL(repository)
			if err != nil {
				app.log.Errorf("error[%#v]", err)
			}

			packageJSON := app.fetchPackageJson(ctx, info)

			if packageJSON != nil {
				app.config.Repositories[i].PackageJSON = packageJSON
				app.config.Repositories[i].Topics = app.fetchTopics(ctx, info)
			}

			// now we can send the result struct through the resultsChan
			resultsChan <- packageJSON

			// once we're done it's we read from the semaphoreChan which
			// has the effect of removing one from the limit and allowing
			// another goroutine to start
			<-semaphoreChan

		}(i, app.config.Repositories[i])
	}

	var results []PackageJSON
	// start listening for any results over the resultsChan
	// once we get a result append it to the result slice
	for {
		result := <-resultsChan
		results = append(results, *result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(app.config.Repositories) {
			break
		}
	}
}

func (app *App) setUpGithubClient() {
	githubAuthToken := strings.TrimSpace(os.Getenv("GITHUB_AUTH_TOKEN"))
	if githubAuthToken == "" {
		app.log.Fatal("GITHUB_AUTH_TOKEN env is empty")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_AUTH_TOKEN")},
	)
	tc := oauth2.NewClient(app.ctx, ts)

	app.githubClient = github.NewClient(tc)
}

func splitRepositoryURL(repository Repository) (map[string]string, error) {

	values := strings.Split(repository.URL, "/")
	packageJSON := "package.json"

	if repository.Type == TypeComponent {
		packageJSON = "/" + strings.Join(values[7:], "/") + "/package.json"
	}

	return map[string]string{
		"owner":       values[3],
		"repo":        values[4],
		"packageJSON": packageJSON,
	}, nil

}

func readConfig(filePath string) (Config, error) {
	var config Config
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return config, err
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return config, err
	}

	json.Unmarshal(bytes, &config)

	return config, err
}

func (app *App) fetchPackageJson(ctx context.Context, info map[string]string) *PackageJSON {
	var packageJSON PackageJSON

	repContent, _, _, err := app.githubClient.Repositories.GetContents(
		ctx,
		info["owner"],
		info["repo"],
		info["packageJSON"],
		&github.RepositoryContentGetOptions{Ref: "master"},
	)
	if err != nil {
		app.log.Error("error: ", err)
	}

	if repContent == nil {
		return &packageJSON
	}

	content, _ := repContent.GetContent()

	json.Unmarshal([]byte(content), &packageJSON)
	packageJSON.Prepare()

	return &packageJSON
}

func (app *App) fetchTopics(ctx context.Context, info map[string]string) []string {
	topics, _, err := app.githubClient.Repositories.ListAllTopics(
		ctx,
		info["owner"],
		info["repo"],
	)
	if err != nil {
		app.log.Error("error: ", err)
	}

	return topics
}

// splitProjectsComponents - split in different arrays the projects and components
func (app *App) splitProjectsComponents() (*[]Repository, *[]Repository, *[]RepositoryClientData, *[]RepositoryClientData) {
	projects := &[]Repository{}
	components := &[]Repository{}

	projectsClientData := &[]RepositoryClientData{}
	componentsClientData := &[]RepositoryClientData{}

	for _, r := range app.config.Repositories {
		if r.Type == TypeProject {
			*projects = append(*projects, r)
			projectClientData := r.getRepositoryClientData()
			*projectsClientData = append(*projectsClientData, *projectClientData)
		}

		if r.Type == TypeComponent {
			*components = append(*components, r)
			componentClientData := r.getRepositoryClientData()
			*componentsClientData = append(*componentsClientData, *componentClientData)
		}
	}

	return projects, components, projectsClientData, componentsClientData
}

// statsCountProjectsByFilters - Count how many projects there are per filter
func (app *App) statsCountProjectsByFilters(projects []Repository, projectsClientData []RepositoryClientData) []interface{} {
	array := []interface{}{}
	for i, f := range app.config.Filters {
		array = append(array, []interface{}{f, 0})
		for j, p := range projects {
			for key, value := range p.PackageJSON.Dependencies {
				if strings.Contains(GetAlias(key, value), f) {
					array[i].([]interface{})[1] = array[i].([]interface{})[1].(int) + 1
					projectsClientData[j].Filter = f
					continue
				}
			}
		}
	}
	return array
}

// statsCountComponentsByFilters - Count how many dependencies are used per version in different projects
func (app *App) statsCountComponentsByFilters(components []Repository, componentsClientData []RepositoryClientData) []interface{} {
	array := []interface{}{}
	for i, filter := range app.config.Filters {
		array = append(array, []interface{}{filter, 0})
		for j, c := range components {
			for key, value := range c.PackageJSON.PeerDependencies {
				if strings.Contains(GetAlias(key, value), filter) {
					array[i].([]interface{})[1] = array[i].([]interface{})[1].(int) + 1
					componentsClientData[j].Filter = filter
					continue
				}
			}
		}
	}
	return array
}

// statsCountDependenciesByVersions - Count how many components there are per filter
func (app *App) statsCountDependenciesByVersions(projects []Repository) map[string]map[string]StatsDependencyVersion {
	var statsComponentsByVersion = map[string]map[string]StatsDependencyVersion{}

	for _, project := range projects {
		for key, value := range project.PackageJSON.Dependencies {
			if statsComponentsByVersion[key] == nil && len(statsComponentsByVersion[key][value].Projects) == 0 && statsComponentsByVersion[key][value].Quantity == 0 {
				statsComponentsByVersion[key] = map[string]StatsDependencyVersion{}
				statsComponentsByVersion[key][value] = StatsDependencyVersion{
					Quantity: 1,
					Projects: []string{project.PackageJSON.Name},
				}
			} else {
				projects := []string{}
				if !contains(statsComponentsByVersion[key][value].Projects, project.PackageJSON.Name) {
					for _, p := range statsComponentsByVersion[key][value].Projects {
						projects = append(projects, p)
					}
					projects = append(projects, project.PackageJSON.Name)
				}

				quantity := statsComponentsByVersion[key][value].Quantity + 1

				statsComponentsByVersion[key][value] = StatsDependencyVersion{
					Quantity: quantity,
					Projects: projects,
				}
			}

		}
	}
	return statsComponentsByVersion
}

// contains -  Checks whether a string exists in an array of string
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
