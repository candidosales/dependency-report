package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// RootHandler - Route to root
func (app *App) RootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")

}

// PingHandler - Route to ping
func (app *App) PingHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
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

	if app.config.OutputFile != "" {
		clientDataJSON, err := json.MarshalIndent(clientData, "", " ")
		err = ioutil.WriteFile(app.config.OutputFile + fileOutput, clientDataJSON, 0644)

		if err == nil {
			app.log.Info("Output file generated and sent to " + app.config.OutputFile + fileOutput)
		}
	}

	json.NewEncoder(w).Encode(clientData)
}