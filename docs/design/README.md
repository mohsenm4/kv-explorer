# KV-Explorer — Design

The visual source of truth for KV-Explorer lives in **Figma**:
<https://www.figma.com/design/QAe3rjM2Mrb5dCUZgDr4Dm/Untitled?node-id=1-1997&p=f>

This folder holds the engineering-side reference that complements it.

## Where things live

| What you want                          | Where to look                                          |
| -------------------------------------- | ------------------------------------------------------ |
| Screens, components, layouts (visual)  | Figma file (link maintained by the design owner)       |
| Color/spacing/typography token values  | [`spec.md`](spec.md) §2–4                              |
| Mapping each component to a Fyne widget| [`spec.md`](spec.md) §6                                |
| Which states each component must have  | [`spec.md`](spec.md) §8                                |
| Keyboard shortcuts                     | [`spec.md`](spec.md) §9                                |
| Accessibility rules                    | [`spec.md`](spec.md) §10                               |
| What Fyne can / can't render           | [`spec.md`](spec.md) §12 — read before designing       |

## Workflow

- **Designer**: works in Figma. Every screen and component must match the
  tokens and Fyne constraints defined in `spec.md`. If a design choice
  conflicts with §12 (e.g. button gradients, per-corner radius), the spec
  wins — bring it up before shipping.
- **Engineer**: implements the theme layer directly from §2–4 of
  `spec.md`. Every other UI package reads colors and sizes from the theme
  only — no hardcoded values. Use Figma for the visual reference of a
  screen, but pull all numbers from the spec.

## When the spec and Figma disagree

`spec.md` is authoritative for **tokens, behavior, and constraints**.
Figma is authoritative for **visual composition**. Anything else
(spacing of a specific element, exact layout of a dialog) — Figma wins.
