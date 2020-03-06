package main

import "regexp"

const (
	TypeProject   = "project"
	TypeComponent = "component"
	pathFile      = "./config.json"
)

type Config struct {
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	URL         string       `json:"url"`
	Type        string       `json:"type"`
	Topics      []string     `json:"topics"`
	PackageJSON *PackageJSON `json:"packageJSON"`
	Alias       string       `json:"alias"`
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
	return p.Name + "_" + p.Version
}

type StatsData struct {
	Labels   []string       `json:"labels"`
	Datasets []StatsDataset `json:"datasets"`
}

type StatsDataset struct {
	Values []int `json:"values"`
}
