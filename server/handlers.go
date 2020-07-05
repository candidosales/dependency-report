package main
//
//import (
//	"encoding/json"
//	"io/ioutil"
//	"net/http"
//	"time"
//)
//
//// RootHandler - Route to root
//func (app *AppConfig) RootHandler(w http.ResponseWriter, r *http.Request) {
//	http.ServeFile(w, r, "./index.html")
//
//}
//
//// GenerateReportHandler - Route to generate report
//func (app *AppConfig) GenerateReportHandler(w http.ResponseWriter, r *http.Request) {
//	app.getPackageJSONs()
//	projects, components, projectsClientData, componentsClientData := app.splitProjectsComponents()
//
//	app.log.Info("Generate data to graphs ... \n")
//	countDependenciesByVersions := app.statsCountDependenciesByVersions(*projects)
//	countComponentsByFilters := app.statsCountComponentsByFilters(*components, *componentsClientData)
//	countProjectsByFilters := app.statsCountProjectsByFilters(*projects, *projectsClientData)
//
//	app.statusProjectsByComponents(*projects, *projectsClientData, *componentsClientData)
//
//	summary := app.generateSummary(*projectsClientData)
//
//	clientData := &ClientData{
//		GeneratedAt: time.Now(),
//		Summary:     summary,
//		Projects:    projectsClientData,
//		Components:  componentsClientData,
//		GraphData: map[string][]interface{}{
//			"projectsByFilters":   countProjectsByFilters,
//			"componentsByFilters": countComponentsByFilters,
//		},
//		DependenciesByVersions: countDependenciesByVersions,
//	}
//
//	if app.config.OutputFile != "" {
//		clientDataJSON, err := json.MarshalIndent(clientData, "", " ")
//		err = ioutil.WriteFile(app.config.OutputFile+fileOutput, clientDataJSON, 0644)
//
//		if err == nil {
//			app.log.Info("Output file generated and sent to " + app.config.OutputFile + fileOutput)
//		}
//	}
//
//	json.NewEncoder(w).Encode(clientData)
//}
