package main

import "regexp"

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
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	License      string            `json:"license"`
	Private      bool              `json:"private"`
	Dependencies map[string]string `json:"dependencies"`
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
	return nil
}
