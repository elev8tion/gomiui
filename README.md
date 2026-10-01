# gomiui

A Go + HTMX port of the strongest ideas in `/Users/kcdacre8tor/Developer/miui`.

gomiui is a server-rendered component catalogue, not a React clone. It keeps miui's dark-first token system, semantic component families, CSS-native effects, accessibility notes, and copy-ready markup while replacing client state with HTML, links, forms, and HTMX fragments.

## Run

```bash
go run ./cmd/app
```

Open <http://localhost:8080/components>.

## What is included

- searchable, filterable component catalogue
- foundations: Button, Badge, Card, Input
- controls: Switch, Textarea, Checkbox, Radio group, Select, Slider, Cyber button
- feedback: Alert, Progress, Dialog, Tooltip
- navigation: Tabs, Stepper, Accordion
- data: Table, Chart, Avatar, Keyboard hint, Separator, Scroll area, Sparkline, Donut
- blocks: Terminal, Empty state, Stat card, Feature card, Pricing card, Model card, Code window
- CSS-first effects: Liquid button, Color orb, Spotlight, Text shimmer, Marquee, Border beam, Counter, Typewriter, System pulse, Liquid tick
- component detail pages with live preview, copy-ready markup, and accessibility contracts
- `/tasks` as a small SQLite + HTMX persistence example
- light/dark themes, reduced-motion support, responsive layout, keyboard focus states

## Architecture

- `internal/domain/components` owns the static UI-kit catalogue and filtering.
- `templates/components/kit.html` is the reusable template component layer.
- `templates/pages/components` contains full pages.
- `templates/partials/components` contains HTMX fragment responses.
- `static/css/app.css` contains the token system and CSS-native effects.
- `internal/domain/tasks` demonstrates the factory's SQLite CRUD contract.

## Checks

```bash
go test ./...
go build ./...
```

The project deliberately avoids React, Tailwind, an ORM, and JSON endpoints for normal UI updates. HTMX requests return HTML fragments, and every important interaction still works as ordinary HTML.
