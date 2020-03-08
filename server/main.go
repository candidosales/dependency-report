package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

type App struct {
	config       Config
	githubClient *github.Client
}

func main() {

	ctx := context.Background()

	config, err := readConfig(pathFile)
	if err != nil {
		fmt.Errorf("error", err)
	}

	app := &App{
		config:       config,
		githubClient: githubClient(ctx),
	}

	for i := 0; i < len(config.Repositories); i++ {
		info, err := splitRepositoryURL(config.Repositories[i])
		if err != nil {
			fmt.Errorf("error[%#v]", err)
			continue
		}
		packageJSON := app.fetchPackageJson(ctx, info)

		if packageJSON != nil {
			config.Repositories[i].PackageJSON = packageJSON
			config.Repositories[i].Topics = app.fetchTopics(ctx, info)
		}
	}

	projects, components, projectsClientData, componentsClientData := app.splitProjectsComponents(config.Repositories)

	countComponentsByProject := statsCountComponentsByProject(*projects, *components)
	countComponentsByVersionAllProjects := statsCountComponentsByVersionAllProjects(*projects)
	countProjectsByFilters := app.statsCountProjectsByFilters(*projects)
	countComponentsByFilters := app.statsCountComponentsByFilters(*components)

	clientData := &ClientData{
		Projects:   projectsClientData,
		Components: componentsClientData,
		GraphData: map[string]*StatsDataFrappe{
			"componentsByProject":            countComponentsByProject,
			"componentsByVersionAllProjects": countComponentsByVersionAllProjects,
			"projectsByFilters":              countProjectsByFilters,
			"componentsByFilters":            countComponentsByFilters,
		},
	}

	clientDataJSON, _ := json.MarshalIndent(clientData, "", " ")
	fmt.Println(string(clientDataJSON))
	_ = ioutil.WriteFile("../client/src/assets/config/data-test.json", clientDataJSON, 0644)

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
		fmt.Errorf("error[%#v]", err)
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
		fmt.Errorf("error[%#v]", err)
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

func statsCountComponentsByProject(projects []Repository, components []Repository) *StatsDataFrappe {
	statsData := &StatsDataFrappe{}

	statsData.Datasets = append(statsData.Datasets, StatsDataset{
		Values: make([]int, len(components)),
	})

	for i, c := range components {
		statsData.Labels = append(statsData.Labels, c.PackageJSON.getAlias())
		for _, p := range projects {
			if p.PackageJSON.Dependencies[c.PackageJSON.Name] == c.PackageJSON.Version {
				statsData.Datasets[0].Values[i] = statsData.Datasets[0].Values[i] + 1
				continue
			}
		}
	}
	return statsData
}

func statsCountComponentsByVersionAllProjects(projects []Repository) *StatsDataFrappe {
	dependencies := map[string]string{}
	statsData := &StatsDataFrappe{}

	statsData.Datasets = append(statsData.Datasets, StatsDataset{
		Values: []int{},
	})

	for _, p := range projects {
		index := 0
		for name, version := range p.PackageJSON.Dependencies {
			label := name + "_" + version
			if dependencies[label] == "" {
				dependencies[label] = version
				statsData.Labels = append(statsData.Labels, label)
				statsData.Datasets[0].Values = append(statsData.Datasets[0].Values, 1)
				index = index + 1
			} else {
				statsData.Datasets[0].Values[index] = statsData.Datasets[0].Values[index] + 1
			}
		}
	}
	return statsData
}

func (app *App) statsCountProjectsByFilters(repositories []Repository) *StatsDataFrappe {
	statsData := &StatsDataFrappe{}

	statsData.Datasets = append(statsData.Datasets, StatsDataset{
		Values: make([]int, len(app.config.Filters)),
	})

	for i, f := range app.config.Filters {
		statsData.Labels = append(statsData.Labels, f)
		for _, r := range repositories {
			for key, value := range r.PackageJSON.Dependencies {
				if strings.Contains(GetAlias(key, value), f) {
					statsData.Datasets[0].Values[i] = statsData.Datasets[0].Values[i] + 1
					continue
				}
			}
		}
	}
	return statsData
}

func (app *App) statsCountComponentsByFilters(repositories []Repository) *StatsDataFrappe {
	statsData := &StatsDataFrappe{}

	statsData.Datasets = append(statsData.Datasets, StatsDataset{
		Values: make([]int, len(app.config.Filters)),
	})

	for i, f := range app.config.Filters {
		statsData.Labels = append(statsData.Labels, f)
		for _, r := range repositories {
			for key, value := range r.PackageJSON.PeerDependencies {
				if strings.Contains(GetAlias(key, value), f) {
					statsData.Datasets[0].Values[i] = statsData.Datasets[0].Values[i] + 1
					continue
				}
			}
		}
	}
	return statsData
}
