package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

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

	config, err := readConfig("./config-test.json")
	if err != nil {
		fmt.Errorf("error", err)
	}

	for i := 0; i < len(config.Repositories); i++ {
		info, err := splitRepositoryURL(config.Repositories[i].URL)
		if err != nil {
			fmt.Errorf("error[%#v]", err)
			continue
		}

		fmt.Printf("info[%#v] \n", info)

		config.Repositories[i].PackageJSON = app.fetchPackageJson(ctx, info)
		config.Repositories[i].Topics = app.fetchTopics(ctx, info)

		config.Repositories[i].PackageJSON.Prepare()
	}

	app.generateStats(config.Repositories)
}

func githubClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_AUTH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func splitRepositoryURL(url string) (map[string]string, error) {

	r, err := regexp.Compile(`https:\/\/github\.com\/([\D\d]+)\/([\D\d]+)`)
	if err != nil {
		return nil, err
	}

	values := r.FindStringSubmatch(url)

	return map[string]string{
		"owner": values[1],
		"repo":  values[2],
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
	repContent, _, _, err := app.githubClient.Repositories.GetContents(
		ctx,
		info["owner"],
		info["repo"],
		"package.json",
		&github.RepositoryContentGetOptions{Ref: "master"},
	)
	if err != nil {
		fmt.Errorf("error[%#v]", err)
	}

	content, _ := repContent.GetContent()

	var packageJSON PackageJSON
	json.Unmarshal([]byte(content), &packageJSON)

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

func (app *App) generateStats(repositories []Repository) {
	fmt.Println("Generate stats")
}
