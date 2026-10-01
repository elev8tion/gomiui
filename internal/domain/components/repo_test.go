package components

import (
	"bytes"
	"context"
	"html/template"
	"strings"
	"testing"
)

func TestRepoListFiltersByCategoryAndQuery(t *testing.T) {
	repo := NewRepo()
	items, err := repo.List(context.Background(), "liquid button", "effects")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Slug != "liquid-button" {
		t.Fatalf("got %#v, want liquid-button", items)
	}
}

func TestCompletedMiuiPortsAreCatalogued(t *testing.T) {
	repo := NewRepo()
	want := []string{
		"textarea", "checkbox", "radio-group", "select", "slider", "accordion",
		"dialog", "tooltip", "avatar", "kbd", "separator", "scroll-area",
		"sparkline", "donut", "cyber-button", "counter", "system-pulse",
		"stat-card", "feature-card", "pricing-card", "model-card", "code-window",
	}
	for _, slug := range want {
		item, ok, err := repo.Get(context.Background(), slug)
		if err != nil || !ok {
			t.Fatalf("missing completed port %s", slug)
		}
		if item.Status != "stable" {
			t.Fatalf("%s status = %s, want stable", slug, item.Status)
		}
	}
}

func TestEveryCatalogPreviewRenders(t *testing.T) {
	tmpl := template.Must(template.ParseFiles(
		"../../../templates/components/kit.html",
		"../../../templates/partials/components/row.html",
	))
	items, err := NewRepo().List(context.Background(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) < 30 {
		t.Fatalf("catalog has %d components, want the completed miui ports included", len(items))
	}
	for _, item := range items {
		var preview bytes.Buffer
		if err := tmpl.ExecuteTemplate(&preview, "kit/preview", item); err != nil {
			t.Fatalf("preview %s: %v", item.Slug, err)
		}
		body := strings.TrimSpace(preview.String())
		if body == `<div class="kit-preview"></div>` || (!strings.Contains(body, "ui-") && !strings.Contains(body, "mini-")) {
			t.Fatalf("preview %s rendered no component markup: %s", item.Slug, body)
		}
		var row bytes.Buffer
		if err := tmpl.ExecuteTemplate(&row, "partials/components/row", item); err != nil {
			t.Fatalf("row %s: %v", item.Slug, err)
		}
		if !strings.Contains(row.String(), item.Name) {
			t.Fatalf("row %s missing name", item.Slug)
		}
	}
}

func TestRepoCategoriesCountCatalog(t *testing.T) {
	repo := NewRepo()
	categories, err := repo.Categories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	got := 0
	for _, category := range categories {
		got += category.Count
	}
	if got != len(catalog()) {
		t.Fatalf("category count = %d, catalog count = %d", got, len(catalog()))
	}
}
