# Go Template tokens

[styles.css](../web/src/styles.css) is the source of truth. Its Tailwind `@theme`
exposes these light-theme decisions directly; no second token-value file is generated.

| Token | Type | Role and usage |
| --- | --- | --- |
| `--font-sans` | Font family | Local Martian Mono with monospace fallbacks |
| `--color-surface` | Color | Page background and secondary-action hover |
| `--color-panel` | Color | Fields, cards and secondary actions |
| `--color-ink` | Color | Primary content and headings |
| `--color-muted` | Color | Supporting copy and background-work labels |
| `--color-border` | Color | Neutral field/card boundaries |
| `--color-accent` | Color | Primary actions and visible focus |
| `--color-accent-hover` | Color | Hovered primary action |
| `--color-danger` | Color | Error text and notice border |

Tailwind generates utilities such as `bg-panel` and `text-muted` from these names.
Keep roles distinct even when values match. Primary buttons retain white foreground
as defined by their component source; other consumers cannot infer that every surface
may use white text. Status is also expressed with words, not only color.

Components own their internal spacing, rounding and focus treatment. Pages use the
existing Tailwind scale for outer gaps and widths. This starter has one light theme;
there is no saved theme preference or dark-mode behavior. When introducing a new
combination, verify its rendered contrast and focus rather than relying on token names.
