---
name: Abacus Design System
colors:
  surface: '#f8f9fb'
  surface-dim: '#d9dadc'
  surface-bright: '#f8f9fb'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#f2f4f6'
  surface-container: '#edeef0'
  surface-container-high: '#e7e8ea'
  surface-container-highest: '#e1e2e4'
  on-surface: '#191c1e'
  on-surface-variant: '#474556'
  inverse-surface: '#2e3132'
  inverse-on-surface: '#f0f1f3'
  outline: '#787588'
  outline-variant: '#c8c4d9'
  surface-tint: '#543bef'
  primary: '#411cde'
  on-primary: '#ffffff'
  primary-container: '#5a42f5'
  on-primary-container: '#e2ddff'
  inverse-primary: '#c6c0ff'
  secondary: '#5f5e5e'
  on-secondary: '#ffffff'
  secondary-container: '#e5e2e1'
  on-secondary-container: '#656464'
  tertiary: '#853100'
  on-tertiary: '#ffffff'
  tertiary-container: '#ac4200'
  on-tertiary-container: '#ffd9ca'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#e4dfff'
  primary-fixed-dim: '#c6c0ff'
  on-primary-fixed: '#150066'
  on-primary-fixed-variant: '#3b0ed9'
  secondary-fixed: '#e5e2e1'
  secondary-fixed-dim: '#c8c6c5'
  on-secondary-fixed: '#1c1b1b'
  on-secondary-fixed-variant: '#474646'
  tertiary-fixed: '#ffdbcd'
  tertiary-fixed-dim: '#ffb596'
  on-tertiary-fixed: '#360f00'
  on-tertiary-fixed-variant: '#7c2e00'
  background: '#f8f9fb'
  on-background: '#191c1e'
  surface-variant: '#e1e2e4'
typography:
  display-results:
    fontFamily: Geist
    fontSize: 48px
    fontWeight: '600'
    lineHeight: 56px
    letterSpacing: -0.02em
  display-results-mobile:
    fontFamily: Geist
    fontSize: 36px
    fontWeight: '600'
    lineHeight: 44px
    letterSpacing: -0.02em
  headline-md:
    fontFamily: Geist
    fontSize: 18px
    fontWeight: '600'
    lineHeight: 24px
  body-lg:
    fontFamily: Geist
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
  body-sm:
    fontFamily: Geist
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
  label-caps:
    fontFamily: Geist
    fontSize: 12px
    fontWeight: '600'
    lineHeight: 16px
    letterSpacing: 0.05em
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  container-max: 480px
  gutter: 1.5rem
  stack-sm: 0.5rem
  stack-md: 1rem
  stack-lg: 2rem
  section-padding: 2.5rem
---

## Brand & Style
The design system is rooted in a **Modern Corporate** aesthetic with strong influences from **Minimalism**. It prioritizes clarity and precision to evoke a sense of calm and institutional trust, essential for a fintech-focused calculator. 

The visual language avoids the clutter of traditional utility apps, opting instead for a "software-as-a-service" (SaaS) feel. By utilizing generous whitespace, a restricted color palette, and high-quality typography, the interface transforms a utility into a premium experience. The emotional response should be one of confidence and effortless control.

## Colors
The palette is dominated by "Financial White" and a scale of "Slate Grays" to ensure the interface feels grounded and neutral. 

- **Primary:** The Indigo (#5A42F5) is used sparingly for high-intent actions, active states, and focus indicators.
- **Surface Scale:** Use #F8F9FB for the background and #FFFFFF for the primary calculator card to create a subtle layered depth.
- **Text Scale:** Primary text uses a deep near-black (#121212), while secondary labels use a medium gray (#64748B) to maintain a clear hierarchy.
- **Accents:** Use a highly desaturated Indigo tint (#F0EEFF) for hover states and secondary buttons to keep the interface soft.

## Typography
This design system utilizes **Geist** for its technical precision and developer-centric aesthetic, which aligns perfectly with fintech's requirement for legible numerals.

The hierarchy is strictly bifurcated: **Display** styles are reserved for numeric outputs and current calculations, using tight letter-spacing to feel "engineered." **Label** styles use uppercase and increased tracking to provide clear navigation for input fields and settings without distracting from the primary data.

## Layout & Spacing
The layout follows a **Fixed Grid** approach for the central calculator module, ensuring the tool remains compact and utility-focused on large screens. 

- **Desktop:** The calculator card is centered horizontally and vertically, with a maximum width of 480px.
- **Mobile:** The card expands to fill the width of the screen with 16px side margins. 
- **Rhythm:** An 8px linear scale is used for all internal spacing. Elements like input groups and operator rows use `stack-md` (16px) to maintain distinct separation without breaking the visual flow.

## Elevation & Depth
The design system employs **Tonal Layers** combined with **Ambient Shadows** to create a sophisticated, tactile sense of "stacking."

1. **Base Layer:** The application background (#F8F9FB).
2. **Surface Layer:** The main calculator card (#FFFFFF). It features a subtle 1px border (#E2E8F0) and a multi-stop shadow: `0 1px 3px rgba(0,0,0,0.05), 0 20px 25px -5px rgba(0,0,0,0.03)`.
3. **Interactive Layer:** Elements like input fields and operator buttons use a slight inset shadow or a solid 2px focus ring of the primary Indigo to indicate activity. 
4. **Overlay Layer:** Modals or tooltips use a higher elevation with a more pronounced blur to sit clearly above the calculation area.

## Shapes
The shape language is consistently **Rounded**, avoiding both the severity of sharp corners and the playfulness of full pills.

- **Standard Elements:** Buttons and input fields use a 0.5rem (8px) radius.
- **Large Elements:** The main calculator card uses a 1rem (16px) radius to feel modern and "object-like."
- **Small Elements:** Tooltips and status badges use 0.25rem (4px) to maintain clarity at small scales.

## Components

### Buttons
- **Primary:** Solid #5A42F5 background with white text. High-contrast, bold, used for the final "Calculate" or "Apply" action.
- **Operator:** Secondary style with a very light gray background (#F1F5F9) and dark text. On `active` or `selected`, these transition to a Primary style or an Indigo outline.
- **Ghost:** Used for auxiliary actions like "Clear" or "Settings," featuring no background and a secondary text color.

### Inputs
- **Numeric Fields:** Large, clear type with a placeholder color of #94A3B8. Focus state is a 2px solid Indigo ring.
- **Suffixes:** Icons or text (like %) are right-aligned within the input container, rendered in a muted gray to not compete with the user's input.

### Result Display
- A dedicated section at the top or bottom of the card. It features a subtle background tint (#F8F9FB) to set it apart from the input area. The font size is the largest in the system (`display-results`).

### Cards
- The primary container for the tool. It should have consistent internal padding (`section-padding`) to ensure the content doesn't feel cramped.

### Status Indicators
- **Errors:** Use a soft red background (#FEE2E2) with dark red text (#991B1B) for inline validation.
- **Loading:** A thin, animated indeterminate progress bar at the very top of the card surface, using the primary Indigo color.