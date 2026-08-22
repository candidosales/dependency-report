package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
)

func (app *AppConfig) setupRoutes() {
	app.fiber.Get("/generate-report", app.generateReport)
	app.fiber.Get("/ping", app.ping)
	app.fiber.Get("/", app.root)
}

func (app *AppConfig) generateReport(c fiber.Ctx) error {
	app.getPackageJSONs()
	projects, components, projectsClientData, componentsClientData := app.splitProjectsComponents()

	app.log.Info("Generate data to graphs ... \n")
	countDependenciesByVersions := app.statsCountDependenciesByVersions(*projects)
	countComponentsByFilters := app.statsCountComponentsByFilters(*components, *componentsClientData)
	countProjectsByFilters := app.statsCountProjectsByFilters(*projects, *projectsClientData)

	app.statusProjectsByComponents(*projects, *projectsClientData, *componentsClientData)

	summary := app.generateSummary(*projectsClientData)

	clientData := &ClientData{
		GeneratedAt: time.Now(),
		Summary:     summary,
		Projects:    projectsClientData,
		Components:  componentsClientData,
		GraphData: map[string][]FilterCount{
			"projectsByFilters":   countProjectsByFilters,
			"componentsByFilters": countComponentsByFilters,
		},
		DependenciesByVersions: countDependenciesByVersions,
	}

	if app.config.OutputFile != "" {
		clientDataJSON, err := json.MarshalIndent(clientData, "", " ")
		err = os.WriteFile(app.config.OutputFile+fileOutput, clientDataJSON, 0644)

		if err == nil {
			app.log.Info("Output file generated and sent to " + app.config.OutputFile + fileOutput)
		}
	}

	return c.JSON(clientData)
}

func (app *AppConfig) ping(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"pong": "ok"})
}

func (app *AppConfig) root(c fiber.Ctx) error {
	return c.SendFile("./index.html", fiber.SendFile{Compress: true})
}
