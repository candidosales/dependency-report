package main

import (
	"regexp"
	"time"
)

const (
	TypeProject     = "project"
	TypeComponent   = "component"
	pathFileInput   = "./config.json"
	pathFileOutput  = "../client/src/assets/config/data-test.json"
	concurrentLimit = 25
)

type Config struct {
	Filters      []string     `json:"filters"`
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	URL         string       `json:"url"`
	Type        string       `json:"type"`
	Topics      []string     `json:"topics"`
	PackageJSON *PackageJSON `json:"packageJSON"`
	Alias       string       `json:"alias"`
}

func (r *Repository) getRepositoryClientData() *RepositoryClientData {
	// TODO add filters
	return &RepositoryClientData{
		Name:    r.PackageJSON.Name,
		Version: r.PackageJSON.Version,
		URL:     r.URL,
	}
}

type PackageJSON struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	License          string            `json:"license"`
	Private          bool              `json:"private"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	Keywords         []string          `json:"keywords"`
}

func (p *PackageJSON) Prepare() {
	p.clearDependenciesVersions()
}

// clearDependenciesVersions - remove special characters. Ex: ^6.0.2 => 6.0.2
func (p *PackageJSON) clearDependenciesVersions() error {
	r, err := regexp.Compile(`([\d\.\-\w]+)`)
	if err != nil {
		return err
	}

	for key, value := range p.Dependencies {
		p.Dependencies[key] = r.FindString(value)
	}

	for key, value := range p.DevDependencies {
		p.DevDependencies[key] = r.FindString(value)
	}

	for key, value := range p.PeerDependencies {
		p.PeerDependencies[key] = r.FindString(value)
	}
	return nil
}
func (p *PackageJSON) getAlias() string {
	return GetAlias(p.Name, p.Version)
}

func GetAlias(name string, version string) string {
	return name + "_" + version
}

type ClientData struct {
	GeneratedAt time.Time `json:"generatedAt"`
	Projects             *[]RepositoryClientData                        `json:"projects"`
	Components           *[]RepositoryClientData                        `json:"components"`
	GraphData            map[string][]interface{}                       `json:"graphData"`
	DependenciesByVersions map[string]map[string]StatsDependencyVersion `json:"dependenciesByVersions"`
}

// RepositoryClientData - simplest object to display in the UI
type RepositoryClientData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Filter  string `json:"filter"`
	URL     string `json:"url"`
}

type StatsDependencyVersion struct {
	Quantity int      `json:"quantity"`
	Projects []string `json:"projects"`
}

type Results struct {
	Rows []interface{} `json:"results"`
}
