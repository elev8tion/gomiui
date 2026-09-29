package components

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service { return &Service{repo: repo} }

func (s *Service) Catalog(ctx context.Context, query, category string) ([]Component, []Category, error) {
	items, err := s.repo.List(ctx, query, category)
	if err != nil {
		return nil, nil, fmt.Errorf("list components: %w", err)
	}
	categories, err := s.repo.Categories(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list component categories: %w", err)
	}
	return items, categories, nil
}

func (s *Service) Find(ctx context.Context, slug string) (Component, error) {
	component, ok, err := s.repo.Get(ctx, strings.TrimSpace(slug))
	if err != nil {
		return Component{}, fmt.Errorf("find component: %w", err)
	}
	if !ok {
		return Component{}, fmt.Errorf("component %q not found", slug)
	}
	return component, nil
}
