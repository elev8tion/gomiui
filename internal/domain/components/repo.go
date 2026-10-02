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
		{
			Slug: "textarea", Name: "Textarea", Category: "controls", Summary: "A multiline field with the same label, hint, and error contract as Input.", Description: "miui's completed textarea is a native control. This port keeps the visible label, the error state, and vertical resizing. Auto-grow stays a progressive enhancement, not a requirement.", Status: "stable", Tags: []string{"form", "text", "validation"},
			Markup: `<label class="ui-field">
  <span class="ui-field__label">Release notes</span>
  <textarea class="ui-textarea" name="notes" rows="4">Shipped the completed form primitives.</textarea>
  <span class="ui-field__hint">Plain text. The server stores exactly what is submitted.</span>
</label>`,
			Accessibility: []string{"Use a visible label, not a placeholder.", "Associate errors with aria-describedby and aria-invalid.", "Let the user resize the field."},
		},
		{
			Slug: "checkbox", Name: "Checkbox", Category: "controls", Summary: "A native checkbox with a drawn box, a label, and an optional description.", Description: "The completed miui checkbox keeps a real input in the accessibility tree. Checked, unchecked, and disabled are CSS states. Indeterminate still needs a script, so this port does not fake it.", Status: "stable", Tags: []string{"form", "choice", "native"},
			Markup: `<label class="ui-check">
  <input type="checkbox" name="index" checked>
  <span class="ui-check__box" aria-hidden="true"></span>
  <span class="ui-choice"><strong>Index on save</strong><small>Rebuild the graph after each write.</small></span>
</label>`,
			Accessibility: []string{"The native checkbox remains the control.", "The text, not only the box, is the hit target.", "Do not replace the check mark with color alone."},
		},
		{
			Slug: "radio-group", Name: "Radio group", Category: "controls", Summary: "A single-choice set rendered as a fieldset of native radios.", Description: "miui's radio group is one named set of native inputs. The fieldset and legend carry the question, so the group works before any script.", Status: "stable", Tags: []string{"form", "choice", "fieldset"},
			Markup: `<fieldset class="ui-radios">
  <legend>Default theme</legend>
  <label class="ui-radio"><input type="radio" name="theme" value="dark" checked><span class="ui-radio__mark" aria-hidden="true"></span><span>Dark</span></label>
  <label class="ui-radio"><input type="radio" name="theme" value="light"><span class="ui-radio__mark" aria-hidden="true"></span><span>Light</span></label>
</fieldset>`,
			Accessibility: []string{"Group radios in a fieldset with a legend.", "Share one name so the browser enforces a single choice.", "Do not use radios for on/off. That is a switch."},
		},
		{
			Slug: "select", Name: "Select", Category: "controls", Summary: "A native menu with the kit's field chrome and a visible chevron.", Description: "The completed miui select is a styled native select, not a custom listbox. Keyboard, typeahead, and form posts stay with the browser.", Status: "stable", Tags: []string{"form", "menu", "native"},
			Markup: `<label class="ui-field">
  <span class="ui-field__label">Model</span>
  <span class="ui-select">
    <select name="model">
      <option>Claude</option>
      <option>GPT</option>
      <option>Local</option>
    </select>
  </span>
</label>`,
			Accessibility: []string{"Keep the native select. Do not rebuild the popup.", "The label names the control.", "Do not rely on the chevron as the only affordance."},
		},
		{
			Slug: "slider", Name: "Slider", Category: "controls", Summary: "A native range with the current value written next to the label.", Description: "miui's slider is a range input. The displayed value is server text, so a form post or an HTMX swap can update it without a client store.", Status: "stable", Tags: []string{"form", "range", "value"},
			Markup: `<label class="ui-slider">
  <span>Density <strong>64</strong></span>
  <input type="range" name="density" min="0" max="100" value="64">
</label>`,
			Accessibility: []string{"Use a label that includes the current value.", "Set min, max, and value on the native range.", "Do not make the thumb the only way to read the number."},
		},
		{
			Slug: "accordion", Name: "Accordion", Category: "navigation", Summary: "Stacked disclosures that open and close without a script.", Description: "The completed miui accordion manages exclusive panels in React. This port uses native details elements, so each panel is keyboard-operable and works with no JavaScript. Exclusive mode can be added later with a small server or HTMX rule.", Status: "stable", Tags: []string{"disclosure", "native", "content"},
			Markup: `<div class="ui-accordion">
  <details open>
    <summary>Tokens</summary>
    <p>Color, type, and radius live in one stylesheet.</p>
  </details>
  <details>
    <summary>Motion</summary>
    <p>Effects are optional and collapse under reduced motion.</p>
  </details>
</div>`,
			Accessibility: []string{"Use summary as the control and put the panel text after it.", "Do not hide the disclosure triangle without a replacement.", "Each panel's content stays in the page, not in a script."},
		},
		{
			Slug: "dialog", Name: "Dialog", Category: "feedback", Summary: "A native dialog panel with a title, explanation, and explicit actions.", Description: "miui's dialog traps focus when opened with script. The Go kit uses the dialog element. The open attribute shows it in place. showModal is the progressive step when a real modal is required.", Status: "stable", Tags: []string{"overlay", "confirm", "native"},
			Markup: `<dialog class="ui-dialog" open>
  <form method="dialog">
    <header><h3>Publish this kit?</h3></header>
    <p>The catalogue will be visible to everyone on this machine.</p>
    <div class="ui-dialog__actions">
      <button class="ui-button ui-button--quiet" value="cancel">Cancel</button>
      <button class="ui-button ui-button--primary" value="publish">Publish</button>
    </div>
  </form>
</dialog>`,
			Accessibility: []string{"Give the dialog a heading.", "Provide a cancel action as well as the primary one.", "Use showModal when the rest of the page must be inert."},
		},
		{
			Slug: "tooltip", Name: "Tooltip", Category: "feedback", Summary: "A short hint that appears on hover and keyboard focus.", Description: "The completed miui tooltip is a positioned note, not a replacement for a label. This port shows it with CSS on hover and focus-within, so it does not depend on a portal.", Status: "stable", Tags: []string{"hint", "focus", "hover"},
			Markup: `<span class="ui-tip">
  <button class="ui-button ui-button--secondary" type="button">Density</button>
  <span class="ui-tip__bubble" role="tooltip">How tightly cards pack on the page.</span>
</span>`,
			Accessibility: []string{"Do not put essential instructions only in a tooltip.", "The hint must appear on keyboard focus, not only hover.", "Keep the text short."},
		},
		{
			Slug: "avatar", Name: "Avatar", Category: "data", Summary: "Initials, stacked people, and an overflow count.", Description: "miui's avatar falls back from an image to initials. This port uses initials so the catalogue has no remote image dependency, and stacks them the same way.", Status: "stable", Tags: []string{"people", "identity", "stack"},
			Markup: `<div class="ui-avatars" aria-label="Ada, Ken, and 3 more">
  <span class="ui-avatar">A</span>
  <span class="ui-avatar ui-avatar--b">K</span>
  <span class="ui-avatar ui-avatar--c">M</span>
  <span class="ui-avatar ui-avatar--more">+3</span>
</div>`,
			Accessibility: []string{"Name the group in text or aria-label.", "Initials are a fallback, not the person's only name in a profile.", "Do not use color as the only difference between people."},
		},
		{
			Slug: "kbd", Name: "Keyboard hint", Category: "data", Summary: "A row of keys for a shortcut that is also written in words nearby.", Description: "miui's kbd renders modifier and key caps. The hint is decorative documentation. The action still needs a real control.", Status: "stable", Tags: []string{"shortcut", "keys", "docs"},
			Markup:        `<p class="ui-kbd-line"><span>Command palette</span> <span class="ui-kbd"><kbd>⌘</kbd><kbd>K</kbd></span></p>`,
			Accessibility: []string{"Write the action in words as well as key caps.", "Do not use a shortcut hint as the only name of a button.", "Platform-specific glyphs need a text equivalent."},
		},
		{
			Slug: "separator", Name: "Separator", Category: "data", Summary: "A quiet rule, with or without a label, between related groups.", Description: "The completed miui separator is a non-interactive divider. A labeled rule uses text, not only a line.", Status: "stable", Tags: []string{"divider", "layout", "group"},
			Markup: `<hr class="ui-separator">
<div class="ui-separator ui-separator--label" role="separator">or</div>`,
			Accessibility: []string{"Use hr or role=separator for a real thematic break.", "Do not announce decorative rules that sit inside a labeled group.", "A labeled separator must include the label text."},
		},
		{
			Slug: "scroll-area", Name: "Scroll area", Category: "data", Summary: "A bounded region with its own scroll and a fade that signals more content.", Description: "miui's scroll area clips a region without taking over the page. The fade is a hint. The scrollbar and keyboard scrolling remain native.", Status: "stable", Tags: []string{"overflow", "region", "clip"},
			Markup: `<div class="ui-scroll" tabindex="0">
  <p>First finding. The graph has two edges.</p>
  <p>Second finding. Hubs sit in the renderer and the handlers.</p>
  <p>Third finding. Experimental canvas effects stay out of this kit.</p>
</div>`,
			Accessibility: []string{"Make the region focusable when the mouse is not the only way to scroll.", "Do not clip content without a visible overflow cue.", "The fade is decorative. The text is the content."},
		},
		{
			Slug: "sparkline", Name: "Sparkline", Category: "data", Summary: "A one-series trend line with a text summary beside it.", Description: "miui ships sparkline as its own chart, separate from bars and donuts. This port is an inline SVG plus a caption, with no chart runtime.", Status: "stable", Tags: []string{"chart", "trend", "svg"},
			Markup: `<figure class="ui-sparkline">
  <svg class="ui-spark" viewBox="0 0 120 36" role="img" aria-label="Latency fell from 42 to 18">
    <polygon points="0,28 16,24 32,26 48,18 64,20 80,12 96,14 120,6 120,36 0,36"></polygon>
    <polyline points="0,28 16,24 32,26 48,18 64,20 80,12 96,14 120,6"></polyline>
  </svg>
  <figcaption>Latency <strong>18 ms</strong></figcaption>
</figure>`,
			Accessibility: []string{"Give the SVG a useful aria-label or a caption.", "The number in the caption is the accessible value.", "Do not use the line as the only copy of the metric."},
		},
		{
			Slug: "donut", Name: "Donut chart", Category: "data", Summary: "A three-part share chart with a center value and a text legend.", Description: "miui's donut is a completed data primitive next to bar and sparkline. The Go port uses a conic gradient and a legend so the shares exist as text.", Status: "stable", Tags: []string{"chart", "share", "legend"},
			Markup: `<div class="ui-donut-block">
  <div class="ui-donut" style="--a: 62%; --b: 86%" role="img" aria-label="Passed 62%, running 24%, queued 14%"><span>62%</span></div>
  <ul class="ui-legend">
    <li><i style="--swatch: var(--accent)"></i>Passed 62%</li>
    <li><i style="--swatch: var(--warning)"></i>Running 24%</li>
    <li><i style="--swatch: var(--line-strong)"></i>Queued 14%</li>
  </ul>
</div>`,
			Accessibility: []string{"Repeat every share in the legend.", "The center number matches the leading share.", "Color supports the label. It does not replace it."},
		},
		{
			Slug: "cyber-button", Name: "Cyber button", Category: "controls", Summary: "A completed miui action with a beveled face and a hairline accent.", Description: "Cyber button is a finished primitive in miui, separate from the plain button. The port keeps the mono label and the bottom hairline, and maps the glow onto the kit accent instead of a second cyan palette.", Status: "stable", Tags: []string{"action", "accent", "button"},
			Markup:        `<button class="ui-cyber" type="button"><span>Deploy</span></button>`,
			Accessibility: []string{"The visible word is the action.", "The hairline and glow are decorative.", "Keep a visible focus ring."},
		},
		{
			Slug: "text-shimmer", Name: "Text shimmer", Category: "effects", Summary: "A highlight that sweeps across a heading and disappears when motion is reduced.", Description: "miui replaced a framer-motion shimmer with one CSS gradient. The same treatment is here. Reduced motion keeps solid text.", Status: "experimental", Tags: []string{"motion", "type", "css"},
			Markup:        `<p class="ui-shimmer">Indexing the graph</p>`,
			Accessibility: []string{"The words stay readable with motion off.", "Do not shimmer body copy.", "The effect is decoration, not a status."},
		},
		{
			Slug: "marquee", Name: "Marquee", Category: "effects", Summary: "A looping strip of labels, with a static wrap when motion is reduced.", Description: "The completed miui marquee is a CSS translation of a duplicated track. Nothing important should live only inside the moving strip.", Status: "experimental", Tags: []string{"motion", "loop", "css"},
			Markup: `<div class="ui-marquee">
  <div class="ui-marquee__track">
    <span>Go</span><span>HTMX</span><span>SQLite</span><span>CSS</span>
    <span>Go</span><span>HTMX</span><span>SQLite</span><span>CSS</span>
  </div>
</div>`,
			Accessibility: []string{"Duplicate the track only for the loop, not as extra information.", "Reduced motion shows the labels without scrolling.", "Do not put a required action inside a marquee."},
		},
		{
			Slug: "border-beam", Name: "Border beam", Category: "effects", Summary: "A thin accent that travels the top edge of one featured surface.", Description: "miui's border beam is a CSS pass over a card edge. Use it once. The card's text does not depend on the beam.", Status: "experimental", Tags: []string{"motion", "border", "accent"},
			Markup:        `<article class="ui-card ui-beam"><p class="ui-kicker">Live</p><h3 class="ui-card__title">Edge light</h3><p class="ui-card__copy">One moving line. The content stays still.</p></article>`,
			Accessibility: []string{"The beam is decorative.", "Keep text contrast on the card itself.", "Use one beam in a group."},
		},
		{
			Slug: "counter", Name: "Counter", Category: "effects", Summary: "A large final number with its unit written beside it.", Description: "miui animates a counter up to a value. The server-rendered port prints the final number so the metric exists without waiting for an animation.", Status: "stable", Tags: []string{"metric", "number", "text"},
			Markup:        `<p class="ui-counter"><strong>1,284</strong><span>files indexed</span></p>`,
			Accessibility: []string{"The number is real text, not a canvas.", "Include the unit.", "Do not rely on the count-up to reveal the value."},
		},
		{
			Slug: "typewriter", Name: "Typewriter", Category: "effects", Summary: "A single line that types itself, and shows the full line when motion is reduced.", Description: "miui's typewriter is a finished CSS effect. The full sentence is in the HTML. The steps animation only reveals it.", Status: "experimental", Tags: []string{"motion", "type", "css"},
			Markup:        `<p class="ui-type">Server rendered, then enhanced.</p>`,
			Accessibility: []string{"The full string is in the element.", "Reduced motion shows the whole line at once.", "Do not type out instructions the user must read to proceed."},
		},
		{
			Slug: "system-pulse", Name: "System pulse", Category: "effects", Summary: "A status dot with a word for healthy, waiting, and failed.", Description: "miui's system pulse pairs a dot with a label. The port keeps three tones and never uses the dot alone.", Status: "stable", Tags: []string{"status", "dot", "health"},
			Markup: `<p class="ui-pulse"><i aria-hidden="true"></i>API healthy</p>
<p class="ui-pulse ui-pulse--warn"><i aria-hidden="true"></i>Queue waiting</p>
<p class="ui-pulse ui-pulse--bad"><i aria-hidden="true"></i>Worker failed</p>`,
			Accessibility: []string{"The word is the status.", "The dot is aria-hidden.", "Do not use green, amber, or red alone."},
		},
		{
			Slug: "liquid-tick", Name: "Liquid tick", Category: "effects", Summary: "A checkbox whose box fills when it is checked.", Description: "miui's liquid tick is a completed effect on top of a native checkbox. The check state is the input. The fill is decoration.", Status: "experimental", Tags: []string{"motion", "checkbox", "css"},
			Markup: `<label class="ui-tick">
  <input type="checkbox" checked>
  <span aria-hidden="true"></span>
  <b>Mark reviewed</b>
</label>`,
			Accessibility: []string{"Use a native checkbox.", "The label names the action.", "The fill must not be required to see the checked state."},
		},
		{
			Slug: "stat-card", Name: "Stat card", Category: "blocks", Summary: "One metric, a direction, and a comparison written in words.", Description: "miui's stat card is a completed block. The delta includes a sign and an arrow so a drop can be good or bad without depending on color.", Status: "stable", Tags: []string{"metric", "delta", "block"},
			Markup: `<article class="ui-stat">
  <p class="ui-kicker">Latency</p>
  <strong>18<small>ms</small></strong>
  <em>↓ 12% vs last week</em>
</article>`,
			Accessibility: []string{"Include the sign and the comparison, not only a color.", "The unit sits with the number.", "One card, one metric."},
		},
		{
			Slug: "feature-card", Name: "Feature card", Category: "blocks", Summary: "A short capability card with a kicker, a title, and one sentence.", Description: "The completed miui feature card is a composed surface, not a new primitive. It reuses the card and can take a single border beam.", Status: "stable", Tags: []string{"block", "marketing", "card"},
			Markup: `<article class="ui-card ui-feature">
  <p class="ui-kicker">01 / Index</p>
  <h3 class="ui-card__title">Know the graph</h3>
  <p class="ui-card__copy">Files, symbols, and the few edges between them.</p>
</article>`,
			Accessibility: []string{"Use a real heading.", "Keep the sentence short enough to scan.", "Do not make the whole card a link unless the entire surface is one action."},
		},
		{
			Slug: "pricing-card", Name: "Pricing card", Category: "blocks", Summary: "A plan, a price, a short list, and one action.", Description: "miui's pricing card is a finished block. The featured plan is marked with text, not only a brighter border.", Status: "stable", Tags: []string{"block", "price", "plan"},
			Markup: `<article class="ui-price is-featured">
  <p class="ui-kicker">Pro · featured</p>
  <b>12<small>/mo</small></b>
  <ul><li>Unlimited catalogues</li><li>Shared themes</li><li>Email support</li></ul>
  <a class="ui-button ui-button--primary" href="/components">Choose Pro</a>
</article>`,
			Accessibility: []string{"Say which plan is featured in text.", "The price includes its period.", "The action is a real link or button."},
		},
		{
			Slug: "model-card", Name: "Model card", Category: "blocks", Summary: "A model name, its provider, and a written status.", Description: "The completed miui model card summarizes one model. Status is a badge with a word, matching the kit's rule that color never stands alone.", Status: "stable", Tags: []string{"block", "model", "status"},
			Markup: `<article class="ui-model">
  <header><h3>Claude</h3><span class="ui-badge ui-badge--accent">ready</span></header>
  <p>Anthropic · 200k context</p>
</article>`,
			Accessibility: []string{"The status word is visible.", "Provider and context are text.", "Use a heading for the model name."},
		},
		{
			Slug: "code-window", Name: "Code window", Category: "blocks", Summary: "A tabbed code surface with a real pre and code element.", Description: "miui's code window is a finished block beside Terminal. Tabs here are static labels. The code itself is selectable text.", Status: "stable", Tags: []string{"block", "code", "tabs"},
			Markup: `<section class="ui-code" aria-label="Sample handler">
  <header><span class="is-active">handlers.go</span><span>repo.go</span></header>
  <pre><code>func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
  h.renderer.HTML(w, http.StatusOK, "pages/components/index", data)
}</code></pre>
</section>`,
			Accessibility: []string{"Use pre and code for the sample.", "Name the window.", "Tabs that do not switch must not pretend to be selected controls."},
		},
		{
			Slug: "label", Name: "Label", Category: "controls", Summary: "The visible name of a field, with a required mark that is not the only signal.", Description: "miui's label is a completed form primitive. The asterisk is decorative. The required state still belongs on the control.", Status: "stable", Tags: []string{"form", "label", "required"},
			Markup: `<label class="ui-label" for="project">Project name <em aria-hidden="true">*</em></label>
<input id="project" class="ui-input" name="project" required>`,
			Accessibility: []string{"Tie the label to the control with for and id.", "Put required on the control, not only on the asterisk.", "Do not use placeholder text as the label."},
		},
		{
			Slug: "low-power", Name: "Low-power toggle", Category: "controls", Summary: "An explicit switch for taking the cheap visual path.", Description: "miui's low-power toggle is a finished escape hatch, separate from the system switch. The label says which mode is on.", Status: "stable", Tags: []string{"toggle", "motion", "preference"},
			Markup: `<label class="ui-switch">
  <input type="checkbox" name="low_power">
  <span class="ui-switch__track" aria-hidden="true"><span></span></span>
  <span>Full motion</span>
</label>`,
			Accessibility: []string{"The label names the current or resulting mode.", "Use a native checkbox.", "Do not hide the control in a corner without a name."},
		},
		{
			Slug: "theme-selector", Name: "Theme selector", Category: "navigation", Summary: "A row of named theme swatches with one selected choice.", Description: "miui's theme selector is a completed radiogroup of color choices. Each swatch has a text name. Color is not the only label.", Status: "stable", Tags: []string{"theme", "radio", "swatch"},
			Markup: `<div class="ui-themes" role="radiogroup" aria-label="Theme">
  <label><input type="radio" name="kit-theme" value="forest" checked><i style="--swatch:#34d399"></i><span>Forest</span></label>
  <label><input type="radio" name="kit-theme" value="observatory"><i style="--swatch:#a78bfa"></i><span>Observatory</span></label>
  <label><input type="radio" name="kit-theme" value="sunset"><i style="--swatch:#fb923c"></i><span>Sunset</span></label>
</div>`,
			Accessibility: []string{"Use a radiogroup and one name.", "Every swatch has a text name.", "The selected ring is extra, not the name."},
		},
		{
			Slug: "edge-handle", Name: "Edge handle", Category: "navigation", Summary: "A button pinned to a panel edge that stays available to the keyboard.", Description: "miui's edge handle reveals a side panel. Auto-hide is off. The handle remains a button with an accessible name in every state.", Status: "stable", Tags: []string{"panel", "edge", "button"},
			Markup: `<div class="ui-edge">
  <input id="ui-edge-open" type="checkbox">
  <label for="ui-edge-open"><span class="sr-only">Open inspector</span><span aria-hidden="true">›</span></label>
  <div class="ui-edge__panel"><p class="ui-kicker">Inspector</p><p>Density, motion, and the active accent.</p></div>
</div>`,
			Accessibility: []string{"The handle is a checkbox with a text name.", "Checking it opens the panel. Unchecking closes it.", "Do not remove the handle on a timer."},
		},
		{
			Slug: "limelight-nav", Name: "Limelight nav", Category: "navigation", Summary: "A compact nav whose active item carries a soft accent glow.", Description: "miui measures the glow onto the active item. This port places the glow with CSS on the current link, so it works with no layout script.", Status: "stable", Tags: []string{"nav", "active", "glow"},
			Markup: `<nav class="ui-lime" aria-label="Primary">
  <a href="#home">Home</a>
  <a class="is-active" href="#kit" aria-current="page">Kit</a>
  <a href="#lab">Lab</a>
</nav>`,
			Accessibility: []string{"Mark the current page with aria-current.", "The glow is decoration. The text is the link name.", "Keep the links real."},
		},
		{
			Slug: "smoke-nav", Name: "Smoke nav", Category: "navigation", Summary: "A frosted top bar with a brand, links, and one action.", Description: "miui's smoke nav is a completed sticky bar. The Go port keeps the blur, the current link, and a wrapping row so it does not overflow a phone.", Status: "stable", Tags: []string{"nav", "bar", "blur"},
			Markup: `<header class="ui-smoke">
  <strong>gomiui</strong>
  <nav aria-label="Product">
    <a class="is-active" href="#components" aria-current="page">Components</a>
    <a href="#tasks">Lab</a>
    <a href="#docs">Docs</a>
  </nav>
  <a class="ui-button ui-button--primary ui-button--small" href="#components">Open</a>
</header>`,
			Accessibility: []string{"Name the navigation.", "Mark the current link.", "Let the row wrap instead of forcing a horizontal scroll."},
		},
		{
			Slug: "flow-button", Name: "Flow button", Category: "effects", Summary: "A pill that fills from the center and slides an arrow through on hover.", Description: "miui's flow button is a finished CSS effect: a growing circle, a label, and an arrow. Reduced motion drops the travel and leaves the action.", Status: "experimental", Tags: []string{"motion", "button", "hover"},
			Markup:        `<button class="ui-flow" type="button"><i aria-hidden="true"></i><b>Continue</b><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M5 12h14M12 5l7 7-7 7"></path></svg></button>`,
			Accessibility: []string{"The word is the action.", "The circle and arrow are decorative.", "Focus shows the same fill as hover."},
		},
		{
			Slug: "save-switch", Name: "Save switch", Category: "effects", Summary: "A pill that reads Saved when the item is kept.", Description: "miui's save switch floods the pill when the bookmark is on. This port shows the saved state in the label. The motion is optional.", Status: "experimental", Tags: []string{"motion", "save", "toggle"},
			Markup:        `<button class="ui-save is-saved" type="button" aria-pressed="true"><span>Saved</span></button>`,
			Accessibility: []string{"Use aria-pressed for the saved state.", "The visible word changes with the state.", "Do not rely on the bookmark animation alone."},
		},
		{
			Slug: "film-grain", Name: "Film grain", Category: "effects", Summary: "A quiet noise tile over a surface. It never takes pointer or reading focus.", Description: "miui's overlay set includes a film-grain tile made from an inline noise image. It is aria-hidden and does not block clicks.", Status: "experimental", Tags: []string{"overlay", "texture", "decoration"},
			Markup:        `<div class="ui-grain" aria-hidden="true"><p>Texture only. The words stay in the page.</p></div>`,
			Accessibility: []string{"Hide decorative grain from assistive tech.", "Do not put the only copy of a message inside the overlay.", "Keep the noise light enough to preserve contrast."},
		},
		{
			Slug: "not-found", Name: "Not found", Category: "blocks", Summary: "A missing-page block with a reason and one way back.", Description: "miui's ghost 404 is a completed block. The canvas particles are skipped. The number, the sentence, and the return link remain.", Status: "stable", Tags: []string{"empty", "error", "block"},
			Markup: `<section class="ui-void" aria-labelledby="void-title">
  <b>404</b>
  <h3 id="void-title">That page is not in the kit.</h3>
  <p>The address does not match a component.</p>
  <a class="ui-button ui-button--secondary" href="/components">Back to the catalogue</a>
</section>`,
			Accessibility: []string{"Explain what is missing in text.", "Give one way back.", "Do not use a particle field as the message."},
		},
		{
			Slug: "weather-widget", Name: "Weather widget", Category: "blocks", Summary: "A city, a temperature, and the condition written out.", Description: "miui's weather widget is a completed block fed by data. This sample is server text: city, degrees, humidity, and wind. The icon is not the only condition.", Status: "stable", Tags: []string{"widget", "weather", "metric"},
			Markup: `<article class="ui-weather">
  <p class="ui-kicker">San Francisco</p>
  <strong>18°</strong>
  <span>Partly cloudy · humidity 62% · wind 14 km/h</span>
</article>`,
			Accessibility: []string{"Write the condition, not only an icon.", "Include the unit.", "The city names the place."},
		},
		{
			Slug: "research-progress", Name: "Research progress", Category: "blocks", Summary: "A written list of research steps with one current line.", Description: "miui's research progress cycles status phrases in script. The port lists the steps and marks the current one, so the state is visible without a timer.", Status: "stable", Tags: []string{"status", "steps", "research"},
			Markup: `<ol class="ui-research" aria-label="Research">
  <li class="is-done">Gathering sources</li>
  <li class="is-now" aria-current="step">Reading documentation</li>
  <li>Writing the report</li>
</ol>`,
			Accessibility: []string{"Use an ordered list for the sequence.", "Mark only one current step.", "Do not replace the step name with a spinner."},
		},
		{
			Slug: "ai-input", Name: "AI input", Category: "blocks", Summary: "A prompt field with a submit action and a visible label.", Description: "miui's AI input morphs from a pill into a field. The useful completed surface is the field itself: a label, a textarea, and a submit button that works as a normal form.", Status: "stable", Tags: []string{"prompt", "form", "block"},
			Markup: `<form class="ui-ask" action="/components" method="get">
  <label class="ui-field">
    <span class="ui-field__label">Ask the kit</span>
    <textarea class="ui-textarea" name="q" rows="3" placeholder="Which control opens a native menu?"></textarea>
  </label>
  <button class="ui-button ui-button--primary" type="submit">Send</button>
</form>`,
			Accessibility: []string{"The field has a visible label.", "Submit is a real button.", "Placeholder text is a hint, not the name."},
		},
		{
			Slug: "agent-grid", Name: "Agent grid", Category: "blocks", Summary: "A responsive grid of agents with role, status, and one action.", Description: "miui's agent grid is a completed block of profile cards. Status is a word. The run action is a button, not the whole card.", Status: "stable", Tags: []string{"grid", "agent", "card"},
			Markup: `<div class="ui-agents">
  <article class="ui-agent">
    <header><h3>Research</h3><span class="ui-badge ui-badge--accent">active</span></header>
    <p>Reads sources and writes a cited note.</p>
    <button class="ui-button ui-button--secondary ui-button--small" type="button">Run</button>
  </article>
</div>`,
			Accessibility: []string{"Each agent has a heading and a status word.", "The action is a button inside the card.", "The grid wraps instead of squeezing the cards."},
		},
		{
			Slug: "permission-dialog", Name: "Permission dialog", Category: "feedback", Summary: "A consent panel that lists each scope and asks for an explicit yes or no.", Description: "miui's permission dialog is a completed consent block. Required scopes are marked in text. Authorize and deny are separate buttons.", Status: "stable", Tags: []string{"consent", "dialog", "scope"},
			Markup: `<dialog class="ui-dialog" open>
  <form method="dialog">
    <header><h3>Allow workspace access?</h3></header>
    <label class="ui-check"><input type="checkbox" checked disabled><span class="ui-check__box" aria-hidden="true"></span><span class="ui-choice"><strong>Read files</strong><small>Required to answer questions about this repo.</small></span></label>
    <label class="ui-check"><input type="checkbox"><span class="ui-check__box" aria-hidden="true"></span><span class="ui-choice"><strong>Run tools</strong><small>Optional. Lets the agent write and test.</small></span></label>
    <div class="ui-dialog__actions">
      <button class="ui-button ui-button--quiet" value="deny">Deny</button>
      <button class="ui-button ui-button--primary" value="allow">Allow</button>
    </div>
  </form>
</dialog>`,
			Accessibility: []string{"Name every scope in text.", "Say which scopes are required.", "Offer deny as well as allow."},
		},
		{
			Slug: "platform-faq", Name: "Platform FAQ", Category: "blocks", Summary: "A titled list of questions that open without a script.", Description: "miui's platform FAQ is a completed accordion block. Native details keep every answer in the page.", Status: "stable", Tags: []string{"faq", "accordion", "content"},
			Markup: `<section class="ui-faq" aria-labelledby="faq-title">
  <p class="ui-kicker">Platform</p>
  <h3 id="faq-title">Common questions</h3>
  <div class="ui-accordion">
    <details open><summary>Where do keys live?</summary><p>On this machine. They are not sent to a second store.</p></details>
    <details><summary>Does it work offline?</summary><p>The catalogue is server-rendered and has no animation runtime.</p></details>
  </div>
</section>`,
			Accessibility: []string{"The section has a heading.", "Each question is a summary.", "Answers are text in the page, not loaded later."},
		},
		{
			Slug: "model-comparison", Name: "Model comparison", Category: "blocks", Summary: "Two model answers side by side, with latency written on each.", Description: "miui's comparison block is a finished layout for a prompt and two results. The sample is static server text, so both columns exist before any request runs.", Status: "stable", Tags: []string{"compare", "model", "layout"},
			Markup: `<section class="ui-compare" aria-label="Model comparison">
  <p class="ui-kicker">Prompt · explain the catalogue</p>
  <div class="ui-compare__grid">
    <article><header><h3>Claude</h3><span>1.2s</span></header><p>A server-rendered kit. HTML first.</p></article>
    <article><header><h3>Local</h3><span>2.8s</span></header><p>Same answer, run on this machine.</p></article>
  </div>
</section>`,
			Accessibility: []string{"Name both models.", "Put the latency in text.", "Stack the columns on a narrow screen."},
		},
		{
			Slug: "folder", Name: "Folder", Category: "effects", Summary: "A folder that lifts its papers on hover and opens from a checkbox.", Description: "The papers and flap are CSS. Checking the box opens the folder without a script. Reduced motion keeps every state and drops the travel.", Status: "experimental", Tags: []string{"motion", "folder", "css"},
			Markup: `<label class="ui-folder">
  <input type="checkbox" aria-label="Open folder">
  <span class="ui-folder__stage">
    <span class="ui-folder__back"></span>
    <span class="ui-folder__cards" aria-hidden="true"><i></i><i></i><i></i></span>
    <span class="ui-folder__flap" aria-hidden="true"><svg viewBox="0 0 321 241"><path d="M0 25C0 11.1929 11.1929 0 25 0H136.084C143.044 0 149.689 2.90139 154.42 8.00608L178.08 33.5343C182.811 38.639 189.456 41.5404 196.416 41.5404H296C309.807 41.5404 321 52.7333 321 66.5404V216C321 229.807 309.807 241 296 241H25C11.1929 241 0 229.807 0 216V25Z"></path></svg></span>
  </span>
</label>`,
			Accessibility: []string{"The control is a checkbox with an accessible name.", "Open and closed are both real states, not hover-only.", "Reduced motion still shows the open folder."},
		},
		{
			Slug: "button-stack", Name: "Button stack", Category: "controls", Summary: "Three glowing actions with visible names.", Description: "Create, export, and confirm stay real buttons. The glow and icon motion are decoration. The names stay readable without a hover.", Status: "stable", Tags: []string{"button", "glow", "actions"},
			Markup: `<div class="ui-stack">
  <div class="ui-stack__item"><button class="ui-stack__btn ui-stack__btn--teal" type="button" aria-label="Create"><svg viewBox="0 0 24 24" aria-hidden="true"><line x1="12" y1="4" x2="12" y2="20"></line><line x1="4" y1="12" x2="20" y2="12"></line></svg></button><span>Create</span></div>
  <div class="ui-stack__item"><button class="ui-stack__btn ui-stack__btn--silver" type="button" aria-label="Export"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M7 17 L17 7 M17 7 L9 7 M17 7 L17 15"></path></svg></button><span>Export</span></div>
  <div class="ui-stack__item"><button class="ui-stack__btn ui-stack__btn--mauve" type="button" aria-label="Confirm"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M5 13 L10 17.5 L19 7"></path></svg></button><span>Confirm</span></div>
</div>`,
			Accessibility: []string{"Each action is a button.", "The name is visible text, not only a tooltip.", "Focus uses the same glow as hover."},
		},
		{
			Slug: "branched-menu", Name: "Branched menu", Category: "navigation", Summary: "A section list whose selected branch draws itself.", Description: "Sections are checkboxes and items are radios, so the menu opens and selects with no script. The accent line shows which item is current.", Status: "stable", Tags: []string{"menu", "tree", "css"},
			Markup: `<nav class="ui-branch" aria-label="Guide">
  <div class="ui-branch__section">
    <input class="ui-branch__toggle" type="checkbox" id="ui-branch-sec-1" checked>
    <span class="ui-branch__marker" aria-hidden="true"></span>
    <label class="ui-branch__head" for="ui-branch-sec-1">Getting started</label>
    <div class="ui-branch__body"><div class="ui-branch__fold"><div class="ui-branch__tree">
      <svg class="ui-branch__lines" viewBox="0 0 40 156" aria-hidden="true"><path class="ui-branch__base" d="M14 0 V122"></path><path class="ui-branch__base" d="M14 14 A10 10 0 0 0 24 24 H32"></path><path class="ui-branch__base" d="M14 50 A10 10 0 0 0 24 60 H32"></path><path class="ui-branch__base" d="M14 86 A10 10 0 0 0 24 96 H32"></path><path class="ui-branch__base" d="M14 122 A10 10 0 0 0 24 132 H32"></path><path class="ui-branch__reach reach-1" d="M14 0 V14 A10 10 0 0 0 24 24 H32"></path><path class="ui-branch__reach reach-2" d="M14 0 V50 A10 10 0 0 0 24 60 H32"></path><path class="ui-branch__reach reach-3" d="M14 0 V86 A10 10 0 0 0 24 96 H32"></path><path class="ui-branch__reach reach-4" d="M14 0 V122 A10 10 0 0 0 24 132 H32"></path></svg>
      <input class="ui-branch__radio" type="radio" name="ui-branch-item" id="ui-branch-1" checked>
      <label class="ui-branch__item" for="ui-branch-1">Installation</label>
      <input class="ui-branch__radio" type="radio" name="ui-branch-item" id="ui-branch-2">
      <label class="ui-branch__item" for="ui-branch-2">Quick start</label>
      <input class="ui-branch__radio" type="radio" name="ui-branch-item" id="ui-branch-3">
      <label class="ui-branch__item" for="ui-branch-3">Configuration</label>
      <input class="ui-branch__radio" type="radio" name="ui-branch-item" id="ui-branch-4">
      <label class="ui-branch__item" for="ui-branch-4">Theming</label>
    </div></div></div>
  </div>
  <div class="ui-branch__section">
    <input class="ui-branch__toggle" type="checkbox" id="ui-branch-sec-2">
    <label class="ui-branch__head" for="ui-branch-sec-2">Components</label>
    <div class="ui-branch__body"><div class="ui-branch__fold"><div class="ui-branch__tree">
      <input class="ui-branch__radio" type="radio" name="ui-branch-item" id="ui-branch-5">
      <label class="ui-branch__item" for="ui-branch-5">Buttons</label>
      <input class="ui-branch__radio" type="radio" name="ui-branch-item" id="ui-branch-6">
      <label class="ui-branch__item" for="ui-branch-6">Typography</label>
      <input class="ui-branch__radio" type="radio" name="ui-branch-item" id="ui-branch-7">
      <label class="ui-branch__item" for="ui-branch-7">Overlays</label>
      <input class="ui-branch__radio" type="radio" name="ui-branch-item" id="ui-branch-8">
      <label class="ui-branch__item" for="ui-branch-8">Toasts</label>
    </div></div></div>
  </div>
</nav>`,
			Accessibility: []string{"Sections are checkboxes and items are radios.", "Every control has a visible label.", "The current item is selected, not only colored."},
		},
	}
}
