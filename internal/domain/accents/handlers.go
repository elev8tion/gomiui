package accents

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/kcdacre8tor/gomiui/internal/web/render"
)

type Handlers struct {
	service  *Service
	renderer *render.Renderer
}

func NewHandlers(service *Service, renderer *render.Renderer) *Handlers {
	return &Handlers{service: service, renderer: renderer}
}

func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	h.page(w, r, "Accent — gomiui", "pages/accents/index-content", "")
}

func (h *Handlers) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	preset, err := h.service.Find(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "load accent preset", http.StatusInternalServerError)
		return
	}
	palette := paint(preset.Accent, preset.Strong, preset.Ink, true)
	data := h.view(r, palette, "")
	data["SelectedPreset"] = preset
	data["SelectedCode"] = palette.Code
	h.renderer.HTML(w, http.StatusOK, "layouts/base", render.ViewData{
		Title:   preset.Name + " — gomiui",
		Content: "pages/accents/show-content",
		Data:    data,
	})
}

func (h *Handlers) Studio(w http.ResponseWriter, r *http.Request) {
	h.renderStudio(w, r, http.StatusOK, Current(r, themeOf(r)), "")
}

func (h *Handlers) Apply(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse form", http.StatusBadRequest)
		return
	}
	palette, err := resolve(r)
	if err != nil {
		h.fail(w, r, "Enter a hex color such as #6ee7b7.")
		return
	}
	WriteCookie(w, palette)
	h.finish(w, r, palette)
}

func (h *Handlers) Reset(w http.ResponseWriter, r *http.Request) {
	ClearCookie(w)
	h.finish(w, r, Default(themeOf(r)))
}

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse form", http.StatusBadRequest)
		return
	}
	palette, err := resolve(r)
	if err != nil {
		h.fail(w, r, "The current accent is not a valid color.")
		return
	}
	if err := h.service.Save(r.Context(), r.FormValue("name"), palette); err != nil {
		h.fail(w, r, saveMessage(err))
		return
	}
	WriteCookie(w, palette)
	h.finish(w, r, palette)
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	if err := h.service.Remove(r.Context(), id); err != nil && !errors.Is(err, ErrNotFound) {
		http.Error(w, "delete accent preset", http.StatusInternalServerError)
		return
	}
	h.finish(w, r, Current(r, themeOf(r)))
}

func (h *Handlers) fail(w http.ResponseWriter, r *http.Request, message string) {
	if isHTMX(r) {
		h.renderStudio(w, r, http.StatusUnprocessableEntity, Current(r, themeOf(r)), message)
		return
	}
	h.page(w, r, "Accent — gomiui", "pages/accents/index-content", message)
}

func (h *Handlers) finish(w http.ResponseWriter, r *http.Request, palette Palette) {
	if isHTMX(r) {
		h.renderStudio(w, r, http.StatusOK, palette, "")
		return
	}
	http.Redirect(w, r, back(r), http.StatusSeeOther)
}

func (h *Handlers) page(w http.ResponseWriter, r *http.Request, title, content, flash string) {
	palette := Current(r, themeOf(r))
	h.renderer.HTML(w, http.StatusOK, "layouts/base", render.ViewData{
		Title:   title,
		Content: content,
		Flash:   flash,
		Data:    h.view(r, palette, flash),
	})
}

func (h *Handlers) renderStudio(w http.ResponseWriter, r *http.Request, status int, palette Palette, flash string) {
	h.renderer.HTML(w, status, "partials/accents/studio", render.ViewData{
		Flash: flash,
		Data:  h.view(r, palette, flash),
	})
}

func (h *Handlers) view(r *http.Request, palette Palette, flash string) map[string]any {
	data, err := h.service.Studio(r.Context(), palette)
	if err != nil {
		data = map[string]any{"Palette": palette, "Presets": markActive(builtinsOnly(), palette)}
		if flash == "" {
			flash = "Saved presets could not be loaded."
		}
	}
	data["Theme"] = themeOf(r)
	data["ThemeURL"] = themeURL(r)
	data["AccentFlash"] = flash
	return data
}

func resolve(r *http.Request) (Palette, error) {
	strong := r.FormValue("strong")
	ink := r.FormValue("ink")
	if strong == "" && ink == "" {
		return Derive(r.FormValue("accent"))
	}
	return Exact(r.FormValue("accent"), strong, ink)
}

func themeOf(r *http.Request) string {
	theme := r.FormValue("theme")
	if theme == "" {
		theme = r.URL.Query().Get("theme")
	}
	if strings.EqualFold(theme, "light") {
		return "light"
	}
	return "dark"
}

func themeURL(r *http.Request) string {
	path := r.URL.Path
	if path == "" {
		path = "/components"
	}
	query := r.URL.Query()
	if themeOf(r) == "dark" {
		query.Set("theme", "light")
	} else {
		query.Set("theme", "dark")
	}
	return path + "?" + query.Encode()
}

func back(r *http.Request) string {
	ref := r.Referer()
	if ref == "" {
		return "/components"
	}
	parsed, err := url.Parse(ref)
	if err != nil {
		return "/components"
	}
	if parsed.Host != "" && !strings.EqualFold(parsed.Host, r.Host) {
		return "/components"
	}
	if parsed.Path == "" {
		return "/components"
	}
	if parsed.RawQuery == "" {
		return parsed.Path
	}
	return parsed.Path + "?" + parsed.RawQuery
}

func saveMessage(err error) string {
	switch {
	case errors.Is(err, ErrBadName):
		return "Name the preset with 1 to 40 characters."
	case errors.Is(err, ErrNameTaken):
		return "A preset already uses that name."
	case errors.Is(err, ErrBadColor):
		return "That accent color is not valid."
	default:
		return "The preset could not be saved."
	}
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
