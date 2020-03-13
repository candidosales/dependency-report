package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type App struct {
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

	log.Info("Start ... \n")
	ctx := context.Background()

	config, err := readConfig(pathFileInput)
	if err != nil {
		log.Error("error", err)
	}

	app := &App{
		config:       config,
		githubClient: githubClient(ctx),
		log:          log,
	}

	// this buffered channel will block at the concurrency limit
	semaphoreChan := make(chan struct{}, concurrentLimit)

	// this channel will not block and collect the packageJSON results
	resultsChan := make(chan *PackageJSON)

	// make sure we close these channels when we're done with them
	defer func() {
		close(semaphoreChan)
		close(resultsChan)
	}()

	app.log.Printf("Get the package.json for the %d repositories ... \n", len(config.Repositories))
	for i := 0; i < len(config.Repositories); i++ {

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
				app.log.Error("error[%#v]", err)
			}

			// now we can send the result struct through the resultsChan
			packageJSON := app.fetchPackageJson(ctx, info)

			if packageJSON != nil {
				config.Repositories[i].PackageJSON = packageJSON
				config.Repositories[i].Topics = app.fetchTopics(ctx, info)
			}

			resultsChan <- packageJSON

			// once we're done it's we read from the semaphoreChan which
			// has the effect of removing one from the limit and allowing
			// another goroutine to start
			<-semaphoreChan

		}(i, config.Repositories[i])
	}

	var results []PackageJSON
	// start listening for any results over the resultsChan
	// once we get a result append it to the result slice
	for {
		result := <-resultsChan
		results = append(results, *result)

		// if we've reached the expected amount of urls then stop
		if len(results) == len(config.Repositories) {
			break
		}
	}

	projects, components, projectsClientData, componentsClientData := app.splitProjectsComponents(config.Repositories)

	app.log.Info("Generate data to graphs ... \n")
	// countComponentsByProject := statsCountComponentsByProject(*projects, *components)
	// countComponentsByVersionAllProjects := statsCountComponentsByVersionAllProjects(*projects)
	countComponentsByVersions := app.statsCountComponentsByVersions(*projects, *components)
	countComponentsByFilters := app.statsCountComponentsByFilters(*components, *componentsClientData)
	countProjectsByFilters := app.statsCountProjectsByFilters(*projects, *projectsClientData)

	clientData := &ClientData{
		Projects:   projectsClientData,
		Components: componentsClientData,
		GraphData: map[string][]interface{}{
			"projectsByFilters":   countProjectsByFilters,
			"componentsByFilters": countComponentsByFilters,
		},
		ComponentsByVersions: countComponentsByVersions,
	}

	clientDataJSON, _ := json.MarshalIndent(clientData, "", " ")
	err = ioutil.WriteFile(pathFileOutput, clientDataJSON, 0644)

	if err != nil {
		app.log.Error("error[%#v]", err)
	} else {
		app.log.Info("Output file generated and sent to " + pathFileOutput)
	}
}

func githubClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_AUTH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
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
		app.log.Error("error[%#v]", err)
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
		app.log.Error("error[%#v]", err)
	}

	return topics
}

// Stats

func (app *App) splitProjectsComponents(repositories []Repository) (*[]Repository, *[]Repository, *[]RepositoryClientData, *[]RepositoryClientData) {
	projects := &[]Repository{}
	components := &[]Repository{}

	projectsClienteData := &[]RepositoryClientData{}
	componentsClientData := &[]RepositoryClientData{}

	for _, r := range repositories {
		if r.Type == TypeProject {
			*projects = append(*projects, r)
			projectClientData := r.getRepositoryClientData()
			*projectsClienteData = append(*projectsClienteData, *projectClientData)

		}

		if r.Type == TypeComponent {
			*components = append(*components, r)
			componentClientData := r.getRepositoryClientData()
			*componentsClientData = append(*componentsClientData, *componentClientData)
		}
	}

	return projects, components, projectsClienteData, componentsClientData
}

// func statsCountComponentsByProject(projects []Repository, components []Repository) *StatsDataFrappe {
// 	statsData := &StatsDataFrappe{}

// 	statsData.Datasets = append(statsData.Datasets, StatsDataset{
// 		Values: make([]int, len(components)),
// 	})

// 	for i, c := range components {
// 		statsData.Labels = append(statsData.Labels, c.PackageJSON.getAlias())
// 		for _, p := range projects {
// 			if p.PackageJSON.Dependencies[c.PackageJSON.Name] == c.PackageJSON.Version {
// 				statsData.Datasets[0].Values[i] = statsData.Datasets[0].Values[i] + 1
// 				continue
// 			}
// 		}
// 	}
// 	return statsData
// }

// func statsCountComponentsByVersionAllProjects(projects []Repository) *StatsDataFrappe {
// 	dependencies := map[string]string{}
// 	statsData := &StatsDataFrappe{}

// 	statsData.Datasets = append(statsData.Datasets, StatsDataset{
// 		Values: []int{},
// 	})

// 	for _, p := range projects {
// 		index := 0
// 		for name, version := range p.PackageJSON.Dependencies {
// 			label := name + "_" + version
// 			if dependencies[label] == "" {
// 				dependencies[label] = version
// 				statsData.Labels = append(statsData.Labels, label)
// 				statsData.Datasets[0].Values = append(statsData.Datasets[0].Values, 1)
// 				index = index + 1
// 			} else {
// 				statsData.Datasets[0].Values[index] = statsData.Datasets[0].Values[index] + 1
// 			}
// 		}
// 	}
// 	return statsData
// }

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

func (app *App) statsCountComponentsByFilters(components []Repository, componentsClientData []RepositoryClientData) []interface{} {
	array := []interface{}{}
	for i, f := range app.config.Filters {
		array = append(array, []interface{}{f, 0})
		for j, c := range components {
			for key, value := range c.PackageJSON.PeerDependencies {
				if strings.Contains(GetAlias(key, value), f) {
					array[i].([]interface{})[1] = array[i].([]interface{})[1].(int) + 1
					componentsClientData[j].Filter = f
					continue
				}
			}
		}
	}
	return array
}

func (app *App) statsCountComponentsByVersions(projects []Repository, components []Repository) map[string]map[string]StatsComponentVersion {
	var statsComponentsByVersion = map[string]map[string]StatsComponentVersion{}

	for _, project := range projects {
		for key, value := range project.PackageJSON.Dependencies {
			if statsComponentsByVersion[key] == nil && len(statsComponentsByVersion[key][value].Projects) == 0 && statsComponentsByVersion[key][value].Quantity == 0 {
				statsComponentsByVersion[key] = map[string]StatsComponentVersion{}
				statsComponentsByVersion[key][value] = StatsComponentVersion{
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

				statsComponentsByVersion[key][value] = StatsComponentVersion{
					Quantity: quantity,
					Projects: projects,
				}
			}

		}
	}
	return statsComponentsByVersion
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
