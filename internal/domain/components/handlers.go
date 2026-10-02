package components

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/kcdacre8tor/gomiui/internal/domain/accents"
	"github.com/kcdacre8tor/gomiui/internal/web/render"
)

type Handlers struct {
	service  *Service
	accents  *accents.Service
	renderer *render.Renderer
}

func NewHandlers(service *Service, accents *accents.Service, renderer *render.Renderer) *Handlers {
	return &Handlers{service: service, accents: accents, renderer: renderer}
}

func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	items, categories, err := h.service.Catalog(r.Context(), query, category)
	if err != nil {
		http.Error(w, "load component catalog", http.StatusInternalServerError)
		return
	}

	h.renderer.HTML(w, http.StatusOK, "layouts/base", render.ViewData{
		Title:   "gomiui — server-rendered interface kit",
		Content: "pages/components/index-content",
		Data:    h.withAccent(r, indexData(items, categories, query, category, r.URL.Query().Get("theme"))),
	})
}

func (h *Handlers) Show(w http.ResponseWriter, r *http.Request) {
	component, err := h.service.Find(r.Context(), r.PathValue("slug"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	items, categories, catalogErr := h.service.Catalog(r.Context(), "", "")
	if catalogErr != nil {
		http.Error(w, "load component catalog", http.StatusInternalServerError)
		return
	}

	h.renderer.HTML(w, http.StatusOK, "layouts/base", render.ViewData{
		Title:   component.Name + " — gomiui",
		Content: "pages/components/show-content",
		Data: h.withAccent(r, map[string]any{
			"Selected":   component,
			"Categories": categories,
			"Components": items,
			"Theme":      theme(r.URL.Query().Get("theme")),
			"ThemeURL":   themeURL(r.URL, r.URL.Query().Get("theme")),
		}),
	})
}

func (h *Handlers) ListFragment(w http.ResponseWriter, r *http.Request) {
	items, _, err := h.service.Catalog(r.Context(), r.URL.Query().Get("q"), r.URL.Query().Get("category"))
	if err != nil {
		http.Error(w, "filter components", http.StatusInternalServerError)
		return
	}
	h.renderer.HTML(w, http.StatusOK, "partials/components/list", render.ViewData{
		Data: map[string]any{"Components": items},
	})
}

func (h *Handlers) withAccent(r *http.Request, data map[string]any) map[string]any {
	if h.accents == nil {
		return data
	}
	return h.accents.Decorate(r.Context(), r, data)
}

func indexData(items []Component, categories []Category, query, active, requestedTheme string) map[string]any {
	return map[string]any{
		"Components": items,
		"Categories": categories,
		"Query":      query,
		"Active":     active,
		"Theme":      theme(requestedTheme),
		"ThemeURL":   themeURL(nil, requestedTheme),
		"Stats":      map[string]any{"Components": len(catalog()), "Categories": len(categoryLabels), "Motion": "CSS-first"},
	}
}

func theme(value string) string {
	if strings.EqualFold(value, "light") {
		return "light"
	}
	return "dark"
}

func themeURL(current *url.URL, currentTheme string) string {
	path := "/components"
	if current != nil && current.Path != "" {
		path = current.Path
	}
	query := url.Values{}
	if current != nil {
		query = current.Query()
	}
	if theme(currentTheme) == "dark" {
		query.Set("theme", "light")
	} else {
		query.Set("theme", "dark")
	}
	return path + "?" + query.Encode()
}
