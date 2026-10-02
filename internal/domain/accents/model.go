package accents

// Palette is the accent token set currently on screen.
type Palette struct {
	Accent string
	Strong string
	Ink    string
	Custom bool
	Code   string
	Style  string
}

// Preset is a named accent, either built in or saved by the user.
type Preset struct {
	ID      int64
	Name    string
	Accent  string
	Strong  string
	Ink     string
	Builtin bool
	Active  bool
}

func paint(accent, strong, ink string, custom bool) Palette {
	palette := Palette{Accent: accent, Strong: strong, Ink: ink, Custom: custom}
	palette.Code = "html {\n  --accent: " + accent + ";\n  --accent-strong: " + strong + ";\n  --accent-ink: " + ink + ";\n}\n"
	if custom {
		palette.Style = "--accent:" + accent + ";--accent-strong:" + strong + ";--accent-ink:" + ink + ";"
	}
	return palette
}
