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
	githubClient *github.Client
}

func main() {

	ctx := context.Background()

	app := &App{
		githubClient: githubClient(ctx),
	}

	config, err := readConfig(pathFile)
	if err != nil {
		fmt.Errorf("error", err)
	}

	for i := 0; i < len(config.Repositories); i++ {
		info, err := splitRepositoryURL(config.Repositories[i])
		if err != nil {
			fmt.Errorf("error[%#v]", err)
			continue
		}
		config.Repositories[i].PackageJSON = app.fetchPackageJson(ctx, info)
		config.Repositories[i].Topics = app.fetchTopics(ctx, info)
	}

	projects, _ := app.splitProjectsComponents(config.Repositories)

	// statsCountComponentsByProject(projects, components)
	statsCountComponentsByVersionAllProjects(projects)
}

func githubClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_AUTH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func splitRepositoryURL(repository Repository) (map[string]string, error) {

	// r, err := regexp.Compile(`https:\/\/github\.com\/([\D\d]+)\/([\D\d]+)`)
	// if err != nil {
	// 	return nil, err
	// }

	// values := r.FindStringSubmatch(url)

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

func (app *App) splitProjectsComponents(repositories []Repository) ([]Repository, []Repository) {
	projects := []Repository{}
	components := []Repository{}

	for _, r := range repositories {
		if r.Type == TypeProject {
			projects = append(projects, r)
		}

		if r.Type == TypeComponent {
			components = append(components, r)
		}
	}

	return projects, components
}

func statsCountComponentsByProject(projects []Repository, components []Repository) {
	statsData := &StatsData{}

	statsData.Datasets = append(statsData.Datasets, StatsDataset{
		Values: make([]int, len(components)),
	})

	for i, c := range components {
		statsData.Labels = append(statsData.Labels, c.PackageJSON.getAlias())
		for _, p := range projects {
			if p.PackageJSON.Dependencies[c.PackageJSON.Name] == c.PackageJSON.Version {
				statsData.Datasets[0].Values[i] = statsData.Datasets[0].Values[i] + 1
			}
		}
	}
	statsDataJSON, _ := json.Marshal(statsData)
	fmt.Println(string(statsDataJSON))
}

func statsCountComponentsByVersionAllProjects(projects []Repository) {
	dependencies := map[string]string{}
	statsData := StatsData{}

	statsData.Datasets = append(statsData.Datasets, StatsDataset{
		Values: []int{},
	})

	for _, p := range projects {
		index := 0
		for name, version := range p.PackageJSON.Dependencies {
			label := name + "_" + version
			// fmt.Printf("label[%#v] \n", label)

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

	statsDataJSON, _ := json.Marshal(statsData)
	fmt.Println(string(statsDataJSON))
}
