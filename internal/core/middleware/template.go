package middleware

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type TemplateMiddleware struct {
	Tmpl *template.Template
}

func NewTemplateMiddleware(templateDir string) *TemplateMiddleware {
	tmpl := template.New("")

	// Walk through the directory and parse all .html files
	filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(path) == ".html" {
			tmpl.ParseFiles(path)
		}
		return nil
	})

	return &TemplateMiddleware{Tmpl: tmpl}
}

func (tm *TemplateMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pass the template object into the context
		ctx := context.WithValue(r.Context(), templateKey{}, tm.Tmpl)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type templateKey struct{}

// GetTemplate returns the template from the context
func GetTemplate(ctx context.Context) *template.Template {
	if tmpl, ok := ctx.Value(templateKey{}).(*template.Template); ok {
		return tmpl
	}
	return nil
}
