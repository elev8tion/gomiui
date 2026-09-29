package components

import (
	"context"
	"testing"
)

func TestRepoListFiltersByCategoryAndQuery(t *testing.T) {
	repo := NewRepo()
	items, err := repo.List(context.Background(), "liquid", "effects")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Slug != "liquid-button" {
		t.Fatalf("got %#v, want liquid-button", items)
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
