package accents

import "testing"

func TestDeriveMintUsesDarkInk(t *testing.T) {
	palette, err := Derive("#6ee7b7")
	if err != nil {
		t.Fatal(err)
	}
	if palette.Accent != "#6ee7b7" || palette.Ink != "#06251a" || !palette.Custom {
		t.Fatalf("palette = %#v", palette)
	}
	if palette.Strong == palette.Accent {
		t.Fatal("strong should differ from the picked accent")
	}
	if !containsAll(palette.Code, "--accent: #6ee7b7;", "--accent-strong: "+palette.Strong+";", "--accent-ink: #06251a;") {
		t.Fatalf("code = %s", palette.Code)
	}
}

func TestDeriveDarkColorUsesWhiteInk(t *testing.T) {
	palette, err := Derive("#087f5b")
	if err != nil {
		t.Fatal(err)
	}
	if palette.Ink != "#ffffff" {
		t.Fatalf("ink = %s", palette.Ink)
	}
}

func TestNormalizeShortHex(t *testing.T) {
	got, ok := NormalizeHex("0F8")
	if !ok || got != "#00ff88" {
		t.Fatalf("got %s %v", got, ok)
	}
}

func TestRejectBadHex(t *testing.T) {
	if _, err := Derive("green"); err == nil {
		t.Fatal("expected an error")
	}
}

func TestDefaultCodeMatchesStylesheet(t *testing.T) {
	dark := Default("dark")
	if dark.Custom || dark.Accent != "#6ee7b7" || dark.Strong != "#34d399" || dark.Ink != "#06251a" {
		t.Fatalf("dark = %#v", dark)
	}
	light := Default("light")
	if light.Accent != "#087f5b" || light.Ink != "#ffffff" {
		t.Fatalf("light = %#v", light)
	}
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !contains(value, part) {
			return false
		}
	}
	return true
}

func contains(value, part string) bool {
	return len(part) == 0 || (len(value) >= len(part) && (value == part || len(value) > 0 && stringIndex(value, part) >= 0))
}

func stringIndex(value, part string) int {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}
