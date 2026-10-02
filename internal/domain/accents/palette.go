package accents

import (
	"math"
	"strconv"
	"strings"
)

func Default(theme string) Palette {
	if strings.EqualFold(theme, "light") {
		return paint("#087f5b", "#066b4c", "#ffffff", false)
	}
	return paint("#6ee7b7", "#34d399", "#06251a", false)
}

func Builtins() []Preset {
	return []Preset{
		{Name: "Mint", Accent: "#6ee7b7", Strong: "#34d399", Ink: "#06251a", Builtin: true},
		{Name: "Grove", Accent: "#087f5b", Strong: "#066b4c", Ink: "#ffffff", Builtin: true},
		{Name: "Observatory", Accent: "#8b5cf6", Strong: "#6d28d9", Ink: "#ffffff", Builtin: true},
		{Name: "Sunset", Accent: "#f97316", Strong: "#c2410c", Ink: "#3b1604", Builtin: true},
		{Name: "Signal", Accent: "#60a5fa", Strong: "#2563eb", Ink: "#041018", Builtin: true},
		{Name: "Rose", Accent: "#fb7185", Strong: "#e11d48", Ink: "#2a0a12", Builtin: true},
	}
}

// Derive turns one picked color into the three tokens the kit actually uses.
func Derive(raw string) (Palette, error) {
	accent, ok := NormalizeHex(raw)
	if !ok {
		return Palette{}, ErrBadColor
	}
	r, g, b := rgb(accent)
	strong := format(scale(r, 0.78), scale(g, 0.78), scale(b, 0.78))
	ink := "#ffffff"
	if luminance(r, g, b) > 0.62 {
		ink = "#06251a"
	}
	return paint(accent, strong, ink, true), nil
}

// Exact keeps a saved or built-in triple instead of deriving it again.
func Exact(accentRaw, strongRaw, inkRaw string) (Palette, error) {
	accent, ok := NormalizeHex(accentRaw)
	if !ok {
		return Palette{}, ErrBadColor
	}
	strong, ok := NormalizeHex(strongRaw)
	if !ok {
		return Derive(accent)
	}
	ink, ok := NormalizeHex(inkRaw)
	if !ok {
		return Derive(accent)
	}
	return paint(accent, strong, ink, true), nil
}

func NormalizeHex(raw string) (string, bool) {
	value := strings.TrimSpace(strings.TrimPrefix(raw, "#"))
	if len(value) == 3 {
		value = strings.Repeat(value[0:1], 2) + strings.Repeat(value[1:2], 2) + strings.Repeat(value[2:3], 2)
	}
	if len(value) != 6 {
		return "", false
	}
	for _, char := range value {
		if !isHex(char) {
			return "", false
		}
	}
	return "#" + strings.ToLower(value), true
}

func rgb(hex string) (float64, float64, float64) {
	n, _ := strconv.ParseUint(hex[1:], 16, 32)
	return float64(n >> 16), float64((n >> 8) & 0xff), float64(n & 0xff)
}

func scale(channel, amount float64) int {
	return int(math.Round(channel * amount))
}

func format(r, g, b int) string {
	return "#" + byteHex(r) + byteHex(g) + byteHex(b)
}

func byteHex(value int) string {
	if value < 0 {
		value = 0
	}
	if value > 255 {
		value = 255
	}
	const digits = "0123456789abcdef"
	return string([]byte{digits[value>>4], digits[value&0x0f]})
}

func luminance(r, g, b float64) float64 {
	return 0.2126*channel(r) + 0.7152*channel(g) + 0.0722*channel(b)
}

func channel(value float64) float64 {
	c := value / 255
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func isHex(char rune) bool {
	return (char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')
}

func same(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}
