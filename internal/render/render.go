package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/AthirsonSilva/golang-net-http-restapi/internal/config"
	"github.com/AthirsonSilva/golang-net-http-restapi/internal/models"
	"github.com/justinas/nosurf"
)

// Creates empty variable for the System-wide configuration typ
var appConfig *config.AppConfig

// Creates a new instance of the Templates function
func NewTemplates(ac *config.AppConfig) {
	appConfig = ac
}

// AddDefaultData adds data for all templates
func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {
	templateData.Flash = appConfig.Session.PopString(r.Context(), "flash")
	templateData.Error = appConfig.Session.PopString(r.Context(), "error")
	templateData.Warning = appConfig.Session.PopString(r.Context(), "warning")
	templateData.CSRFToken = nosurf.Token(r)
	return templateData
}

// RenderTemplate renders templates using html/template with caching
func RenderTemplate(w http.ResponseWriter, r *http.Request, templateFile string, templateData *models.TemplateData) error {
	var templateCache map[string]*template.Template

	// Check if the cache is being used
	if appConfig.UseCache {
		templateCache = appConfig.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// Get the template set
	templateCache = appConfig.TemplateCache
	template, ok := templateCache[templateFile]

	if !ok {
		log.Fatal("Template not found => ", templateFile)
	}

	buffer := new(bytes.Buffer)
	templateData = AddDefaultData(templateData, r)
	_ = template.Execute(buffer, templateData)
	_, err := buffer.WriteTo(w)

	if err != nil {
		log.Println("Error writing template to browser => ", err)
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return templateCache, err
	}

	// Loop through all files ending with page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)

		if err != nil {
			return templateCache, err
		}

		// Look for any layout files
		layoutMatches, err := filepath.Glob("./templates/*.layout.tmpl")

		if err != nil {
			return templateCache, err
		}

		// If there are any layout files
		if len(layoutMatches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")

			if err != nil {
				return templateCache, err
			}
		}

		// Add template set to the cache
		templateCache[name] = templateSet
	}

	return templateCache, nil
}