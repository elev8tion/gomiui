package components

import (
	"context"
	"strings"
)

type Repo struct{}

func NewRepo() *Repo { return &Repo{} }

func (r *Repo) List(ctx context.Context, query, category string) ([]Component, error) {
	_ = ctx
	query = strings.ToLower(strings.TrimSpace(query))
	category = strings.ToLower(strings.TrimSpace(category))
	result := make([]Component, 0)
	for _, component := range catalog() {
		if category != "" && component.Category != category {
			continue
		}
		if query != "" && !matches(component, query) {
			continue
		}
		result = append(result, component)
	}
	return result, nil
}

func (r *Repo) Get(ctx context.Context, slug string) (Component, bool, error) {
	_ = ctx
	for _, component := range catalog() {
		if component.Slug == slug {
			return component, true, nil
		}
	}
	return Component{}, false, nil
}

func (r *Repo) Categories(ctx context.Context) ([]Category, error) {
	_ = ctx
	counts := map[string]int{}
	for _, component := range catalog() {
		counts[component.Category]++
	}
	categories := make([]Category, 0, len(categoryLabels))
	for _, item := range categoryLabels {
		categories = append(categories, Category{Slug: item.slug, Label: item.label, Count: counts[item.slug]})
	}
	return categories, nil
}

func matches(component Component, query string) bool {
	if strings.Contains(strings.ToLower(component.Name), query) || strings.Contains(strings.ToLower(component.Summary), query) || strings.Contains(strings.ToLower(component.Category), query) {
		return true
	}
	for _, tag := range component.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

type categoryLabel struct{ slug, label string }

var categoryLabels = []categoryLabel{
	{slug: "foundations", label: "Foundations"},
	{slug: "controls", label: "Controls"},
	{slug: "feedback", label: "Feedback"},
	{slug: "navigation", label: "Navigation"},
	{slug: "data", label: "Data display"},
	{slug: "effects", label: "Effects"},
	{slug: "blocks", label: "Blocks"},
}

func catalog() []Component {
	return []Component{
		{
			Slug: "button", Name: "Button", Category: "controls", Summary: "A small, honest action surface with clear variants and states.", Description: "The base control keeps every state visible: focus, active, loading, disabled, destructive, and quiet. It works as a normal form button before HTMX arrives.", Status: "stable", Tags: []string{"action", "form", "focus"},
			Markup: `<button class="ui-button ui-button--primary" type="button">Save changes</button>
<button class="ui-button ui-button--secondary" type="button">Preview</button>
<button class="ui-button ui-button--quiet" type="button">Cancel</button>`,
			Accessibility: []string{"Use a native button and an explicit type.", "Loading states retain aria-busy and a readable label.", "Focus-visible rings are never removed."},
		},
		{
			Slug: "badge", Name: "Badge", Category: "foundations", Summary: "Compact status language for labels, metadata, and system state.", Description: "Badges are deliberately restrained. Color supports the label; it never replaces it.", Status: "stable", Tags: []string{"label", "status", "token"},
			Markup: `<span class="ui-badge ui-badge--accent">stable</span>
<span class="ui-badge ui-badge--neutral">server rendered</span>
<span class="ui-badge ui-badge--warning">needs review</span>`,
			Accessibility: []string{"Never communicate status by color alone.", "Keep labels short enough to scan in a dense list."},
		},
		{
			Slug: "card", Name: "Card", Category: "foundations", Summary: "A quiet surface for grouping content without turning everything into a tile.", Description: "Use a card only when it creates a real boundary. The kit has solid, glass, and accent-edge treatments that all share the same radius token.", Status: "stable", Tags: []string{"surface", "layout", "glass"},
			Markup: `<article class="ui-card">
  <p class="ui-kicker">Pattern / 01</p>
  <h3 class="ui-card__title">A useful boundary</h3>
  <p class="ui-card__copy">The surface gives related content room to breathe.</p>
</article>`,
			Accessibility: []string{"Use a heading inside a card when it is a meaningful section.", "Do not make a non-interactive card pretend to be a link."},
		},
		{
			Slug: "input", Name: "Input", Category: "controls", Summary: "Labeled form fields with helper, error, and success affordances.", Description: "Inputs are plain HTML first. Add HTMX attributes to the form or field without changing the visual contract.", Status: "stable", Tags: []string{"form", "label", "validation"},
			Markup: `<label class="ui-field">
  <span class="ui-field__label">Project name</span>
  <input class="ui-input" name="project_name" placeholder="e.g. atlas" autocomplete="off">
  <span class="ui-field__hint">Use a short, memorable name.</span>
</label>`,
			Accessibility: []string{"Labels are visible and associated with their control.", "Errors sit next to the field and are announced with aria-describedby.", "Never use placeholder text as the label."},
		},
		{
			Slug: "switch", Name: "Switch", Category: "controls", Summary: "A native checkbox styled as a calm on/off control.", Description: "The kit prefers native inputs. This switch is a checkbox with a visible label, so it remains usable without JavaScript.", Status: "stable", Tags: []string{"toggle", "form", "native"},
			Markup: `<label class="ui-switch">
  <input type="checkbox" checked>
  <span class="ui-switch__track" aria-hidden="true"><span></span></span>
  <span>Reduce motion</span>
</label>`,
			Accessibility: []string{"The native checkbox remains in the accessibility tree.", "The text label is the click target, not only the small track."},
		},
		{
			Slug: "alert", Name: "Alert", Category: "feedback", Summary: "Inline feedback that explains what happened and what to do next.", Description: "Alerts are persistent context, not a decorative toast. The tone is carried by a border and a label as well as color.", Status: "stable", Tags: []string{"feedback", "error", "notice"},
			Markup: `<div class="ui-alert ui-alert--success" role="status">
  <strong>Saved.</strong>
  <span>Your preferences are now active.</span>
</div>
<div class="ui-alert ui-alert--warning" role="alert">
  <strong>Review required.</strong>
  <span>This action needs your attention.</span>
</div>`,
			Accessibility: []string{"Use role=status for non-urgent updates and role=alert for urgent errors.", "Include text that explains the next action."},
		},
		{
			Slug: "progress", Name: "Progress", Category: "feedback", Summary: "A measured progress rail with a readable value and active state.", Description: "Pair the visual rail with a value or description. The component can update through an HTMX fragment with no client state.", Status: "stable", Tags: []string{"loading", "status", "meter"},
			Markup: `<div class="ui-progress" role="progressbar" aria-valuenow="68" aria-valuemin="0" aria-valuemax="100" aria-label="Indexing progress">
  <span style="--value: 68%"></span>
</div>
<p class="ui-progress__meta"><span>Indexing files</span><strong>68%</strong></p>`,
			Accessibility: []string{"Expose min, max, and current values on the progressbar.", "Do not rely on animation to communicate completion."},
		},
		{
			Slug: "tabs", Name: "Tabs", Category: "navigation", Summary: "A compact way to switch between related views while keeping the page legible.", Description: "For server-rendered pages, tabs are links when they change the URL and buttons when they only change local content. This preview uses links so it degrades cleanly.", Status: "stable", Tags: []string{"navigation", "links", "views"},
			Markup: `<nav class="ui-tabs" aria-label="Component views">
  <a class="ui-tab is-active" href="#preview" aria-current="page">Preview</a>
  <a class="ui-tab" href="#markup">Markup</a>
  <a class="ui-tab" href="#tokens">Tokens</a>
</nav>`,
			Accessibility: []string{"Use aria-current for the active URL view.", "Prefer real links when the selected tab is addressable."},
		},
		{
			Slug: "stepper", Name: "Stepper", Category: "navigation", Summary: "A progress rail that makes the current step and completed work unmistakable.", Description: "The visual language comes from miui's animated stepper, translated into server-rendered state. Completed, active, loading, and future steps each have a distinct treatment.", Status: "stable", Tags: []string{"workflow", "progress", "forms"},
			Markup: `<ol class="ui-stepper" aria-label="Setup progress">
  <li class="ui-step is-complete"><span>✓</span><strong>Choose</strong></li>
  <li class="ui-step is-active" aria-current="step"><span>2</span><strong>Configure</strong></li>
  <li class="ui-step"><span>3</span><strong>Publish</strong></li>
</ol>`,
			Accessibility: []string{"Use an ordered list for sequential work.", "Mark only one item aria-current=step."},
		},
		{
			Slug: "table", Name: "Table", Category: "data", Summary: "A dense, readable data surface with semantic headers and quiet row states.", Description: "Tables keep alignment and hierarchy doing the work. Use them for comparison, not for every list.", Status: "stable", Tags: []string{"data", "comparison", "responsive"},
			Markup: `<div class="ui-table-wrap">
<table class="ui-table">
  <caption>Recent builds</caption>
  <thead><tr><th scope="col">Build</th><th scope="col">Status</th><th scope="col">Duration</th></tr></thead>
  <tbody><tr><th scope="row">gomiui-042</th><td><span class="ui-badge ui-badge--accent">passed</span></td><td>02:18</td></tr></tbody>
</table>
</div>`,
			Accessibility: []string{"Use caption, th, scope, and row headers for real data.", "Wrap tables horizontally instead of shrinking type below legibility."},
		},
		{
			Slug: "chart", Name: "Chart", Category: "data", Summary: "Small inline visualizations for trend and distribution, with text fallbacks.", Description: "The React version shipped bar, donut, and sparkline primitives. The Go kit keeps the same visual grammar with CSS bars and an adjacent data summary.", Status: "stable", Tags: []string{"chart", "metrics", "css"},
			Markup: `<div class="ui-chart" role="img" aria-label="Weekly requests: 32, 48, 41, 62, 74, 68, 86">
  <span style="--bar: 36%"></span><span style="--bar: 53%"></span><span style="--bar: 45%"></span><span style="--bar: 69%"></span><span style="--bar: 82%"></span><span style="--bar: 75%"></span><span style="--bar: 94%"></span>
</div>
<p class="ui-chart__caption">Requests / last 7 days <strong>+18.4%</strong></p>`,
			Accessibility: []string{"Provide a useful text summary or table alongside a visual chart.", "The chart never carries the only copy of a metric."},
		},
		{
			Slug: "terminal", Name: "Terminal", Category: "blocks", Summary: "An app-level block for commands, logs, and live system output.", Description: "Terminal is a block rather than a primitive. It keeps the chrome, status, and monospace rhythm from miui without pretending to be a shell.", Status: "stable", Tags: []string{"block", "code", "logs"},
			Markup: `<section class="ui-terminal" aria-label="Build output">
  <header><span class="ui-status-dot"></span><strong>build / latest</strong><span class="ui-terminal__time">02:18</span></header>
  <pre><code><span class="code-muted">$</span> go test ./...
<span class="code-good">ok</span>   github.com/acme/gomiui   0.42s</code></pre>
</section>`,
			Accessibility: []string{"Use a real pre/code pair for log output.", "Do not use color alone for success or failure."},
		},
		{
			Slug: "liquid-button", Name: "Liquid button", Category: "effects", Summary: "A CSS-only fill effect that keeps the motion optional and the action obvious.", Description: "The original miui effect used rising liquid, wave layers, and bubbles. The Go port keeps that idea in CSS pseudo-elements and native hover/focus states.", Status: "experimental", Tags: []string{"motion", "css", "hover"},
			Markup:        `<button class="ui-liquid" type="button"><span aria-hidden="true"></span><b>Run build</b></button>`,
			Accessibility: []string{"The effect is decorative; the label carries the action.", "Reduced motion collapses the transition and leaves the control usable."},
		},
		{
			Slug: "orb", Name: "Color orb", Category: "effects", Summary: "A restrained ambient accent for empty states, heroes, and command surfaces.", Description: "miui's orb was a conic-gradient and mask stack. This version uses the same CSS-native approach with no canvas, WebGL, or animation runtime.", Status: "experimental", Tags: []string{"ambient", "gradient", "motion"},
			Markup:        `<div class="ui-orb" aria-hidden="true"></div>`,
			Accessibility: []string{"Keep decorative orbs aria-hidden.", "Never put essential text over a moving effect without a solid contrast surface."},
		},
		{
			Slug: "empty-state", Name: "Empty state", Category: "blocks", Summary: "A useful pause that tells people why the surface is empty and what to do next.", Description: "A blank screen is not a state. This block gives the page a quiet visual anchor and one clear next step.", Status: "stable", Tags: []string{"block", "onboarding", "copy"},
			Markup: `<section class="ui-empty" aria-labelledby="empty-title">
  <div class="ui-empty__mark">∅</div>
  <h3 id="empty-title">Nothing indexed yet</h3>
  <p>Connect a repository to see its first component graph.</p>
  <a class="ui-button ui-button--primary" href="/components">Browse the kit</a>
</section>`,
			Accessibility: []string{"Explain the empty state in text, not only an icon.", "Give the user one obvious next action."},
		},
		{
			Slug: "spotlight", Name: "Spotlight", Category: "effects", Summary: "A soft border and glow that gives focus to one important surface.", Description: "Use the spotlight as a single emphasis layer. It is deliberately not a permanent glow on every card.", Status: "experimental", Tags: []string{"focus", "border", "accent"},
			Markup:        `<article class="ui-card ui-spotlight"><p class="ui-kicker">Featured primitive</p><h3 class="ui-card__title">Make one thing matter</h3><p class="ui-card__copy">The accent edge creates hierarchy without adding another color.</p></article>`,
			Accessibility: []string{"Treat the glow as decoration and preserve normal contrast.", "Use one spotlight in a group, not a glow wall."},
		},
	}
}
