(() => {
  document.addEventListener('click', async (event) => {
    const button = event.target.closest('[data-copy-target]');
    if (!button) return;
    const target = document.querySelector(button.dataset.copyTarget);
    if (!target) return;
    try {
      await navigator.clipboard.writeText(target.innerText);
      const previous = button.textContent;
      button.textContent = 'Copied';
      window.setTimeout(() => { button.textContent = previous; }, 1200);
    } catch (_) {
      button.textContent = 'Select manually';
    }
  });

  let accentSaveTimer = 0;
  let accentSaveRequest = null;

  function accentChannel(value) {
    const c = value / 255;
    if (c <= 0.04045) return c / 12.92;
    return ((c + 0.055) / 1.055) ** 2.4;
  }

  function accentFromHex(raw) {
    const match = raw.trim().match(/^#?([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$/);
    if (!match) return null;
    const digits = match[1].length === 3 ? match[1].replace(/./g, (char) => char + char) : match[1];
    const accent = '#' + digits.toLowerCase();
    const r = parseInt(digits.slice(0, 2), 16);
    const g = parseInt(digits.slice(2, 4), 16);
    const b = parseInt(digits.slice(4, 6), 16);
    const strong = [r, g, b].map((channel) => Math.max(0, Math.min(255, Math.round(channel * 0.78))).toString(16).padStart(2, '0')).join('');
    const luminance = (0.2126 * accentChannel(r)) + (0.7152 * accentChannel(g)) + (0.0722 * accentChannel(b));
    return { accent, strong: '#' + strong, ink: luminance > 0.62 ? '#06251a' : '#ffffff' };
  }

  function paintAccent(palette) {
    const style = document.getElementById('accent-override');
    if (style) style.textContent = 'html { --accent: ' + palette.accent + '; --accent-strong: ' + palette.strong + '; --accent-ink: ' + palette.ink + '; }';
    const code = document.querySelector('#accent-code code');
    if (code) code.textContent = 'html {\n  --accent: ' + palette.accent + ';\n  --accent-strong: ' + palette.strong + ';\n  --accent-ink: ' + palette.ink + ';\n}\n';
    const save = document.querySelector('.accent-save');
    if (!save) return;
    save.querySelector('[name="accent"]').value = palette.accent;
    save.querySelector('[name="strong"]').value = palette.strong;
    save.querySelector('[name="ink"]').value = palette.ink;
  }

  function rememberAccent(value) {
    window.clearTimeout(accentSaveTimer);
    accentSaveTimer = window.setTimeout(async () => {
      if (accentSaveRequest) accentSaveRequest.abort();
      accentSaveRequest = new AbortController();
      const theme = document.querySelector('#accent-apply [name="theme"]');
      const body = new URLSearchParams({ accent: value, theme: theme ? theme.value : 'dark' });
      try {
        await fetch('/accents/apply', {
          method: 'POST',
          headers: { 'Content-Type': 'application/x-www-form-urlencoded', 'HX-Request': 'true' },
          body,
          signal: accentSaveRequest.signal,
        });
      } catch (_) {}
    }, 350);
  }

  function hsvToHex(h, s, v) {
    const c = v * s;
    const hh = ((h % 360) + 360) % 360 / 60;
    const x = c * (1 - Math.abs(hh % 2 - 1));
    let r = 0;
    let g = 0;
    let b = 0;
    if (hh < 1) [r, g, b] = [c, x, 0];
    else if (hh < 2) [r, g, b] = [x, c, 0];
    else if (hh < 3) [r, g, b] = [0, c, x];
    else if (hh < 4) [r, g, b] = [0, x, c];
    else if (hh < 5) [r, g, b] = [x, 0, c];
    else [r, g, b] = [c, 0, x];
    const m = v - c;
    const hex = [r, g, b].map((channel) => Math.round((channel + m) * 255).toString(16).padStart(2, '0')).join('');
    return '#' + hex;
  }

  function hexToHsv(hex) {
    const raw = parseInt(hex.slice(1), 16);
    const r = ((raw >> 16) & 255) / 255;
    const g = ((raw >> 8) & 255) / 255;
    const b = (raw & 255) / 255;
    const max = Math.max(r, g, b);
    const min = Math.min(r, g, b);
    const d = max - min;
    let h = 0;
    if (d !== 0) {
      if (max === r) h = ((g - b) / d) % 6;
      else if (max === g) h = (b - r) / d + 2;
      else h = (r - g) / d + 4;
      h *= 60;
      if (h < 0) h += 360;
    }
    return { h, s: max === 0 ? 0 : d / max, v: max };
  }

  function pickerParts() {
    const picker = document.querySelector('[data-accent-picker]');
    if (!picker) return null;
    return {
      picker,
      panel: picker.querySelector('[data-accent-panel]'),
      sv: picker.querySelector('[data-accent-sv]'),
      hue: picker.querySelector('[data-accent-hue]'),
      opener: picker.querySelector('[data-accent-open]'),
      hex: picker.querySelector('[data-accent-hex]'),
    };
  }

  function commitAccent(palette) {
    const parts = pickerParts();
    if (parts?.hex) parts.hex.value = palette.accent;
    if (parts?.opener) parts.opener.style.setProperty('--swatch', palette.accent);
    paintAccent(palette);
    rememberAccent(palette.accent);
  }

  function paintFromHsv(parts, close) {
    const h = Number(parts.panel.dataset.h || parts.hue.value || 0);
    const s = Number(parts.panel.dataset.s || 1);
    const v = Number(parts.panel.dataset.v || 1);
    parts.sv.style.setProperty('--hue', String(h));
    parts.sv.style.setProperty('--s', String(s));
    parts.sv.style.setProperty('--v', String(v));
    parts.hue.value = String(Math.round(h));
    const palette = accentFromHex(hsvToHex(h, s, v));
    if (!palette) return;
    if (close) {
      parts.panel.hidden = true;
      parts.opener.setAttribute('aria-expanded', 'false');
      commitAccent(palette);
      return;
    }
    if (parts.hex) parts.hex.value = palette.accent;
    if (parts.opener) parts.opener.style.setProperty('--swatch', palette.accent);
    paintAccent(palette);
  }

  function pointToSv(parts, event) {
    const rect = parts.sv.getBoundingClientRect();
    const s = Math.min(1, Math.max(0, (event.clientX - rect.left) / rect.width));
    const v = Math.min(1, Math.max(0, 1 - ((event.clientY - rect.top) / rect.height)));
    parts.panel.dataset.s = String(s);
    parts.panel.dataset.v = String(v);
    parts.panel.dataset.h = parts.hue.value;
  }

  document.addEventListener('click', (event) => {
    const parts = pickerParts();
    if (!parts?.panel) return;
    const open = event.target.closest('[data-accent-open]');
    if (open) {
      const palette = accentFromHex(parts.hex.value) || accentFromHex(getComputedStyle(document.documentElement).getPropertyValue('--accent'));
      const hsv = hexToHsv(palette ? palette.accent : '#6ee7b7');
      parts.panel.dataset.h = String(hsv.h);
      parts.panel.dataset.s = String(hsv.s);
      parts.panel.dataset.v = String(hsv.v);
      const next = parts.panel.hidden;
      parts.panel.hidden = !next;
      parts.opener.setAttribute('aria-expanded', next ? 'true' : 'false');
      if (next) paintFromHsv(parts, false);
      else paintFromHsv(parts, true);
      return;
    }
    if (event.target.closest('[data-accent-done]')) {
      paintFromHsv(parts, true);
      return;
    }
    if (!parts.panel.hidden && !event.target.closest('[data-accent-picker]')) {
      paintFromHsv(parts, true);
    }
  });

  document.addEventListener('pointerdown', (event) => {
    const sv = event.target.closest('[data-accent-sv]');
    if (!sv) return;
    const parts = pickerParts();
    sv.setPointerCapture(event.pointerId);
    pointToSv(parts, event);
    paintFromHsv(parts, false);
  });

  document.addEventListener('pointermove', (event) => {
    const sv = event.target.closest('[data-accent-sv]');
    if (!sv || !sv.hasPointerCapture?.(event.pointerId)) return;
    const parts = pickerParts();
    pointToSv(parts, event);
    paintFromHsv(parts, false);
  });

  document.addEventListener('pointerup', (event) => {
    const sv = event.target.closest('[data-accent-sv]');
    if (!sv) return;
    const parts = pickerParts();
    pointToSv(parts, event);
    paintFromHsv(parts, true);
  });

  document.addEventListener('input', (event) => {
    const parts = pickerParts();
    if (!parts) return;
    if (event.target === parts.hue) {
      parts.panel.dataset.h = parts.hue.value;
      paintFromHsv(parts, false);
      return;
    }
    if (event.target !== parts.hex) return;
    const palette = accentFromHex(parts.hex.value);
    if (!palette) return;
    const hsv = hexToHsv(palette.accent);
    parts.panel.dataset.h = String(hsv.h);
    parts.panel.dataset.s = String(hsv.s);
    parts.panel.dataset.v = String(hsv.v);
    commitAccent(palette);
  });

  document.addEventListener('htmx:afterSwap', (event) => {
    if (event.target?.id !== 'component-grid') return;
    document.querySelectorAll('.filter-chip').forEach((chip) => chip.classList.remove('is-active'));
  });
})();
