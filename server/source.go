package main

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/go-github/v29/github"
	"regexp"
	"time"
)

const (
	TypeProject     = "project"
	TypeComponent   = "component"
	pathFileInput   = "./config.json"
	pathOutput  = "../client/src/assets/config/"
	fileOutput = "data.json"
	concurrentLimit = 25
)

type Config struct {
	Filters      []string     `json:"filters"`
	Repositories []Repository `json:"repositories"`
	OutputFile string `json:"outputFile"`
}

type Repository struct {
	URL           string        `json:"url"`
	Type          string        `json:"type"`
	Topics        []string      `json:"topics"`
	PackageJSON   *PackageJSON  `json:"packageJSON"`
	Alias         string        `json:"alias"`
	Documentation Documentation `json:"documentation"`
	Notifications []*github.Notification `json:"notifications"`
}

func (r *Repository) getRepositoryClientData() *RepositoryClientData {
	// TODO add filters
	return &RepositoryClientData{
		Name:          r.PackageJSON.Name,
		Version:       r.PackageJSON.Version,
		URL:           r.URL,
		Documentation: r.Documentation,
		Notifications: r.Notifications,
		Updates: []*UpdateComponent{},
		Tags: []string{},
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
	err := p.clearDependenciesVersions()
	if err != nil {

	}
}

func (p PackageJSON) Validate() error {
	return validation.ValidateStruct(&p,
		// Name cannot be empty, and the length must between 5 and 50
		validation.Field(&p.Name, validation.Required),
	)
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
	GeneratedAt            time.Time                                    `json:"generatedAt"`
	Summary SummaryData `json:"summary"`
	Projects               *[]RepositoryClientData                      `json:"projects"`
	Components             *[]RepositoryClientData                      `json:"components"`
	GraphData              map[string][]interface{}                     `json:"graphData"`
	DependenciesByVersions map[string]*DependencyVersion `json:"dependenciesByVersions"`
}

// RepositoryClientData - simplest object to display in the UI
type RepositoryClientData struct {
	Name          string        `json:"name"`
	Version       string        `json:"version"`
	Filter        string        `json:"filter"`
	URL           string        `json:"url"`
	Documentation Documentation `json:"documentation"`
	Notifications []*github.Notification `json:"notifications"`
	Updates []*UpdateComponent `json:"updates"`
	Tags []string `json:"tags"`
}

type UpdateComponent struct {
	Name string `json:"name"`
	Current string `json:"current"`
	Update string `json:"update"`
}

type SummaryData struct {
	Updated []string `json:"updated"`
	Inconsistent  []string  `json:"inconsistent"`
	Vulnerable  []string  `json:"vulnerable"`
}

type Documentation struct {
	Frontend string `json:"frontend"`
	Design   string `json:"design"`
}

type DependencyVersion struct {
	Type string      `json:"type"`
	Versions map[string]StatsDependencyVersion `json:"versions"`
}

type StatsDependencyVersion struct {
	Quantity int      `json:"quantity"`
	Projects []string `json:"projects"`
}

type FilterNotificationsGetOptions struct {
	Reason string
	Unread bool
}
