package main

import "time"

var pathToTemplates = "./cmd/web/templates"

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]float64
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
}
