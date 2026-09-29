package components

// Component is a documented, server-rendered building block in the Go UI kit.
type Component struct {
	Slug          string
	Name          string
	Category      string
	Summary       string
	Description   string
	Status        string
	Tags          []string
	Markup        string
	Accessibility []string
}

// Category is used by the catalog navigation and filter.
type Category struct {
	Slug  string
	Label string
	Count int
}
