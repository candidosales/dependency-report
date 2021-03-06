package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/go-github/v29/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type AppConfig struct {
	ctx          context.Context
	config       Config
	githubClient *github.Client
	fiber        *fiber.App
	log          *zap.SugaredLogger
}

func main() {

	outputFile := flag.String("output-file", pathOutput, "File path where the JSON generated by the report will be saved")

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	log := zapLogger.Sugar()

	ctx := context.Background()

	appFiber := fiber.New()
	appFiber.Use(logger.New())
	appFiber.Use(cors.New())

	config, err := readConfig(pathFileInput)
	config.OutputFile = *outputFile
	if err != nil {
		log.Error("error", err)
	}

	appConfig := &AppConfig{
		ctx:    ctx,
		config: config,
		fiber:  appFiber,
		log: log,
	}

	appConfig.setUpGithubClient()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	appConfig.setupRoutes()
	err = appConfig.fiber.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}

// getPackageJSONs - Get Package Json for each repository from config file
func (app *AppConfig) getPackageJSONs() {
	app.log.Infof("Get the package.json for the %d repositories ... \n", len(app.config.Repositories))

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
				app.log.Errorf("error[%#v] to info[%#v]", err, info)
			}

			packageJSON := app.fetchPackageJson(info)
			err = packageJSON.Validate()
			if err != nil {
				app.log.Errorf("validate - error[%s] to info[%#v]", err, info)
			}

			notifications := app.fetchNotifications(info, &FilterNotificationsGetOptions{
				Reason: "security_alert",
				Unread: true,
			})

			if packageJSON != nil {
				app.config.Repositories[i].PackageJSON = packageJSON
				app.config.Repositories[i].Topics = app.fetchTopics(info)
				app.config.Repositories[i].Notifications = notifications
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

func (app *AppConfig) setUpGithubClient() {
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

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	return config, err
}

func (app *AppConfig) fetchPackageJson(info map[string]string) *PackageJSON {
	var packageJSON PackageJSON

	repContent, _, _, err := app.githubClient.Repositories.GetContents(
		app.ctx,
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

	content, err := repContent.GetContent()
	if err != nil {
		app.log.Error("error: ", err)
	}

	err = json.Unmarshal([]byte(content), &packageJSON)
	if err != nil {
		app.log.Error("error: ", err)
	}
	packageJSON.Prepare()

	return &packageJSON
}

func (app *AppConfig) fetchNotifications(info map[string]string, opts *FilterNotificationsGetOptions) []*github.Notification {

	notifications, _, err := app.githubClient.Activity.ListRepositoryNotifications(
		app.ctx,
		info["owner"],
		info["repo"],
		&github.NotificationListOptions{
			All:           false,
			Participating: false,
		},
	)
	if err != nil {
		app.log.Error("error: ", err)
	}

	if opts.Reason != "" {
		notificationsFiltered := []*github.Notification{}
		for _, notification := range notifications {
			if opts.Reason == *notification.Reason && opts.Unread == *notification.Unread {
				notificationsFiltered = append(notificationsFiltered, notification)
			}
		}
		return notificationsFiltered
	}

	return notifications
}

func (app *AppConfig) fetchTopics(info map[string]string) []string {
	topics, _, err := app.githubClient.Repositories.ListAllTopics(
		app.ctx,
		info["owner"],
		info["repo"],
	)
	if err != nil {
		app.log.Error("error: ", err)
	}

	return topics
}

// splitProjectsComponents - split in different arrays the projects and components
func (app *AppConfig) splitProjectsComponents() (*[]Repository, *[]Repository, *[]RepositoryClientData, *[]RepositoryClientData) {
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
func (app *AppConfig) statsCountProjectsByFilters(projects []Repository, projectsClientData []RepositoryClientData) []interface{} {
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
func (app *AppConfig) statsCountComponentsByFilters(components []Repository, componentsClientData []RepositoryClientData) []interface{} {
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
func (app *AppConfig) statsCountDependenciesByVersions(projects []Repository) map[string]*DependencyVersion {
	var statsDependenciesByVersion = map[string]*DependencyVersion{}

	for _, project := range projects {
		for key, value := range project.PackageJSON.Dependencies {
			if statsDependenciesByVersion[key] == nil {
				//if (len(statsDependenciesByVersion[key].Versions[value].Projects) == 0 && statsDependenciesByVersion[key].Versions[value].Quantity == 0) {
				statsDependenciesByVersion[key] = &DependencyVersion{
					Type:     "good",
					Versions: map[string]StatsDependencyVersion{},
				}
				statsDependenciesByVersion[key].Versions[value] = StatsDependencyVersion{
					Quantity: 1,
					Projects: []string{project.PackageJSON.Name},
				}
				//}

			} else {
				projects := append(statsDependenciesByVersion[key].Versions[value].Projects, project.PackageJSON.Name)

				statsDependenciesByVersion[key].Type = getTypeDependency(len(statsDependenciesByVersion[key].Versions))
				quantity := statsDependenciesByVersion[key].Versions[value].Quantity + 1

				statsDependenciesByVersion[key].Versions[value] = StatsDependencyVersion{
					Quantity: quantity,
					Projects: projects,
				}
			}

		}
	}
	return statsDependenciesByVersion
}

func getTypeDependency(quantity int) string {
	switch {
	case quantity <= 2:
		return "good"
	case quantity > 2 && quantity <= 5:
		return "warning"
	case quantity > 5 && quantity <= 10:
		return "bad"
	case quantity > 10:
		return "terrible"
	default:
		return ""
	}
}

// statusProjectsByComponents - Count how many projects there are per filter
func (app *AppConfig) statusProjectsByComponents(projects []Repository, projectsClientData []RepositoryClientData, components []RepositoryClientData) {
	for i, p := range projects {
		for key, value := range p.PackageJSON.Dependencies {
			for _, c := range components {
				if key == c.Name && value != c.Version {
					projectsClientData[i].Updates = append(projectsClientData[i].Updates, &UpdateComponent{
						Name:    c.Name,
						Current: value,
						Update:  c.Version,
					})
				}
			}
		}
	}
}

func (app *AppConfig) generateSummary(projectsClientData []RepositoryClientData) SummaryData {
	summaryData := SummaryData{
		Updated:      []string{},
		Inconsistent: []string{},
		Vulnerable:   []string{},
	}

	for _, p := range projectsClientData {
		if !isVulnerable(&p) && !isInconsistent(&p) {
			summaryData.Updated = append(summaryData.Updated, p.Name)
			continue
		}

		if isVulnerable(&p) {
			summaryData.Vulnerable = append(summaryData.Vulnerable, p.Name)
		}
		if isInconsistent(&p) {
			summaryData.Inconsistent = append(summaryData.Inconsistent, p.Name)
		}

	}
	return summaryData
}

func isInconsistent(repository *RepositoryClientData) bool {
	return repository.Updates != nil && len(repository.Updates) > 0
}

func isVulnerable(repository *RepositoryClientData) bool {
	return repository.Notifications != nil && len(repository.Notifications) > 0
}
