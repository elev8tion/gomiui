package app

import (
	"net/http"

	"github.com/kcdacre8tor/gomiui/internal/domain/accents"
	"github.com/kcdacre8tor/gomiui/internal/domain/components"
	"github.com/kcdacre8tor/gomiui/internal/domain/tasks"
)

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	accentRepo := accents.NewRepo(a.database)
	accentService := accents.NewService(accentRepo)
	accentHandlers := accents.NewHandlers(accentService, a.renderer)

	componentRepo := components.NewRepo()
	componentService := components.NewService(componentRepo)
	componentHandlers := components.NewHandlers(componentService, accentService, a.renderer)

	taskRepo := tasks.NewRepo(a.database)
	taskService := tasks.NewService(taskRepo)
	taskHandlers := tasks.NewHandlers(taskService, accentService, a.renderer)

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/components", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /components", componentHandlers.Index)
	mux.HandleFunc("GET /components/{slug}", componentHandlers.Show)
	mux.HandleFunc("GET /fragments/components/list", componentHandlers.ListFragment)

	mux.HandleFunc("GET /accents", accentHandlers.Index)
	mux.HandleFunc("GET /accents/{id}", accentHandlers.Show)
	mux.HandleFunc("POST /accents", accentHandlers.Create)
	mux.HandleFunc("POST /accents/apply", accentHandlers.Apply)
	mux.HandleFunc("POST /accents/reset", accentHandlers.Reset)
	mux.HandleFunc("POST /accents/{id}/delete", accentHandlers.Delete)
	mux.HandleFunc("GET /fragments/accents/studio", accentHandlers.Studio)

	// The task screen stays as a small persistence example: it demonstrates
	// the same HTML-first + HTMX contract against SQLite.
	mux.HandleFunc("GET /tasks", taskHandlers.Index)
	mux.HandleFunc("POST /tasks", taskHandlers.Create)
	mux.HandleFunc("POST /tasks/{id}/complete", taskHandlers.Complete)
	mux.HandleFunc("POST /tasks/{id}/delete", taskHandlers.Delete)
	mux.HandleFunc("GET /fragments/tasks/list", taskHandlers.ListFragment)

	return mux
}
