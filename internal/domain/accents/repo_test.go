package accents

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kcdacre8tor/gomiui/internal/db"
	"github.com/kcdacre8tor/gomiui/internal/web/render"
)

func TestPresetRoundTripAndStudioRender(t *testing.T) {
	database, err := db.Open(filepath.Join(t.TempDir(), "app.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })
	if err := db.RunMigrations(database, filepath.Join("..", "..", "..", "migrations")); err != nil {
		t.Fatal(err)
	}

	service := NewService(NewRepo(database))
	palette, err := Derive("#60a5fa")
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Save(context.Background(), "Signal blue", palette); err != nil {
		t.Fatal(err)
	}
	if err := service.Save(context.Background(), "Signal blue", palette); err != ErrNameTaken {
		t.Fatalf("duplicate = %v", err)
	}

	renderer, err := render.New(filepath.Join("..", "..", "..", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/components?theme=dark", nil)
	data := service.Decorate(request.Context(), request, map[string]any{"Theme": "dark"})
	var body bytes.Buffer
	renderer.HTML(bodyWriter{&body}, http.StatusOK, "partials/accents/studio", render.ViewData{Data: data})
	html := body.String()
	if strings.Contains(html, "template render error") || strings.Contains(html, "ZgotmplZ") {
		t.Fatalf("studio template failed: %s", html)
	}
	for _, want := range []string{"Mint", "Signal blue", "--accent: #6ee7b7;", "accent-studio", "Custom color", "accent-hex"} {
		if !strings.Contains(html, want) {
			t.Fatalf("studio missing %q\n%s", want, html)
		}
	}

	request.AddCookie(&http.Cookie{Name: cookieName, Value: "8b5cf6,6d28d9,ffffff"})
	data = service.Decorate(request.Context(), request, map[string]any{"Theme": "dark"})
	body.Reset()
	renderer.HTML(bodyWriter{&body}, http.StatusOK, "partials/accents/studio", render.ViewData{Data: data})
	html = body.String()
	for _, want := range []string{"--accent: #8b5cf6;", "html { --accent: #8b5cf6;"} {
		if !strings.Contains(html, want) {
			t.Fatalf("custom studio missing %q\n%s", want, html)
		}
	}
	if strings.Contains(html, "ZgotmplZ") {
		t.Fatal("style was rejected by the template escaper")
	}
}

type bodyWriter struct{ buf *bytes.Buffer }

func (w bodyWriter) Header() http.Header         { return http.Header{} }
func (w bodyWriter) Write(p []byte) (int, error) { return w.buf.Write(p) }
func (w bodyWriter) WriteHeader(int)             {}
