package page

import (
	"log"
	"net/http"
	"text/template"

	"github.com/hacdias/caddy-hugo/assets"
	"github.com/hacdias/caddy-hugo/utils"
)

const (
	templateExtension = ".tmpl"
)

var funcMap = template.FuncMap{
	"splitCapitalize": utils.SplitCapitalize,
	"isMarkdown":      utils.IsMarkdownFile,
}

// Page type
type Page struct {
	Name  string
	Class string
	Body  interface{}
}

// Render the page
func (p *Page) Render(w http.ResponseWriter, r *http.Request, templates ...string) (int, error) {
	tpl, err := GetTemplate(r, templates...)

	if err != nil {
		log.Print(err)
		return 500, err
	}

	tpl.Execute(w, p)
	return 200, nil
}

// GetTemplate is used to get a ready to use template based on the url and on
// other sent templates
func GetTemplate(r *http.Request, templates ...string) (*template.Template, error) {
	// If this is a pjax request, use the minimal template to send only
	// the main content
	if r.Header.Get("X-PJAX") == "true" {
		templates = append(templates, "base_minimal")
	} else {
		templates = append(templates, "base_full")
	}

	var tpl *template.Template

	// For each template, add it to the the tpl variable
	for i, t := range templates {
		// Get the template from the assets
		page, err := assets.Asset("templates/" + t + templateExtension)

		// Check if there is some error. If so, the template doesn't exist
		if err != nil {
			log.Print(err)
			return new(template.Template), err
		}

		// If it's the first iteration, creates a new template and add the
		// functions map
		if i == 0 {
			tpl, err = template.New(t).Funcs(funcMap).Parse(string(page))
		} else {
			tpl, err = tpl.Parse(string(page))
		}

		if err != nil {
			log.Print(err)
			return new(template.Template), err
		}
	}

	return tpl, nil
}
