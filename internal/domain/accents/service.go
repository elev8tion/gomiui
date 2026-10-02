package accents

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const cookieName = "gomiui_accent"

var (
	ErrBadColor  = errors.New("bad accent color")
	ErrBadName   = errors.New("bad preset name")
	ErrNameTaken = errors.New("preset name taken")
	ErrNotFound  = errors.New("preset not found")
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Decorate(ctx context.Context, r *http.Request, data map[string]any) map[string]any {
	if data == nil {
		data = map[string]any{}
	}
	theme, _ := data["Theme"].(string)
	palette := Current(r, theme)
	studio, err := s.studio(ctx, palette)
	if err != nil {
		data["AccentFlash"] = "Saved presets could not be loaded."
		studio = markActive(builtinsOnly(), palette)
	}
	data["Palette"] = palette
	data["Presets"] = studio
	return data
}

func (s *Service) Studio(ctx context.Context, palette Palette) (map[string]any, error) {
	presets, err := s.studio(ctx, palette)
	if err != nil {
		return nil, err
	}
	return map[string]any{"Palette": palette, "Presets": presets}, nil
}

func (s *Service) Save(ctx context.Context, name string, palette Palette) error {
	name = strings.TrimSpace(name)
	if name == "" || len([]rune(name)) > 40 || strings.ContainsAny(name, "\r\n\t") {
		return ErrBadName
	}
	_, err := s.repo.Create(ctx, Preset{
		Name:   name,
		Accent: palette.Accent,
		Strong: palette.Strong,
		Ink:    palette.Ink,
	})
	return err
}

func (s *Service) Remove(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Find(ctx context.Context, id int64) (Preset, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) studio(ctx context.Context, palette Palette) ([]Preset, error) {
	saved, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return markActive(append(Builtins(), saved...), palette), nil
}

func builtinsOnly() []Preset {
	return append([]Preset(nil), Builtins()...)
}

func markActive(presets []Preset, palette Palette) []Preset {
	for i := range presets {
		presets[i].Active = palette.Custom && same(presets[i].Accent, palette.Accent) && same(presets[i].Strong, palette.Strong) && same(presets[i].Ink, palette.Ink)
	}
	return presets
}

func Current(r *http.Request, theme string) Palette {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value == "" {
		return Default(theme)
	}
	parts := strings.Split(cookie.Value, ",")
	if len(parts) != 3 {
		palette, deriveErr := Derive(cookie.Value)
		if deriveErr != nil {
			return Default(theme)
		}
		return palette
	}
	palette, exactErr := Exact(parts[0], parts[1], parts[2])
	if exactErr != nil {
		return Default(theme)
	}
	return palette
}

func WriteCookie(w http.ResponseWriter, palette Palette) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    strings.TrimPrefix(palette.Accent, "#") + "," + strings.TrimPrefix(palette.Strong, "#") + "," + strings.TrimPrefix(palette.Ink, "#"),
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 365,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
