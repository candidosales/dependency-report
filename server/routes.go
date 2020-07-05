package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gofiber/fiber"
)

func (app *AppConfig) setupRoutes() {
	app.fiber.Get("/generate-report", app.generateReport)
	app.fiber.Get("/ping", app.ping)
	app.fiber.Get("/", app.root)
}

func (app *AppConfig) generateReport(c *fiber.Ctx) {
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
		GraphData: map[string][]interface{}{
			"projectsByFilters":   countProjectsByFilters,
			"componentsByFilters": countComponentsByFilters,
		},
		DependenciesByVersions: countDependenciesByVersions,
	}

	if app.config.OutputFile != "" {
		clientDataJSON, err := json.MarshalIndent(clientData, "", " ")
		err = ioutil.WriteFile(app.config.OutputFile+fileOutput, clientDataJSON, 0644)

		if err == nil {
			app.log.Info("Output file generated and sent to " + app.config.OutputFile + fileOutput)
		}
	}

	c.JSON(clientData)
}

func (app *AppConfig) ping(c *fiber.Ctx) {
	if err := c.JSON(fiber.Map{"pong": "ok"}); err != nil {
		c.Next(err)
	}
}

func (app *AppConfig) root(c *fiber.Ctx) {
	if err := c.SendFile("./index.html", true); err != nil {
		c.Next(err)
	}
}
