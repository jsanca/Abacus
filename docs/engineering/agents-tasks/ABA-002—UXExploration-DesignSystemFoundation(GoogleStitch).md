# ABA-002 — UX Exploration & Design System Foundation (Google Stitch)

**Status:** READY
**Owner:** Google Stitch
**Role:** Product Designer · UX/UI Exploration
**Target:** 20–30 minutes
**Hard Stop:** 45 minutes

---

# Objective

Design the initial user experience and visual design system for **Abacus**, a modern full-stack calculator application built as a take-home assignment for Sezzle.

The objective is **not** to replicate Sezzle's application, but to produce a polished interface that could naturally fit within the visual quality and product philosophy of a modern fintech.

The output should provide a production-quality visual foundation for React implementation.

---

# Product Context

Abacus is a calculator centered around **operations**, not a physical keypad.

Users provide one or two operands depending on the selected operation.

The backend performs every calculation.

The frontend consumes backend capabilities dynamically.

The UI should communicate:

* clarity
* confidence
* simplicity
* trust
* precision

---

# Primary Users

Developers, reviewers and engineering teams evaluating the assignment.

The interface should immediately communicate:

> "This application was designed with care."

---

# Functional Requirements

The interface must support the following operations:

* Addition
* Subtraction
* Multiplication
* Division
* Exponentiation
* Square Root
* Percentage

Operations are selected explicitly.

Do **not** use a traditional calculator keypad.

---

# Interaction Model

The calculator is based on:

```text
Operand 1

↓

Operation

↓

Operand 2 (when required)

↓

Calculate

↓

Result
```

The interface must support continuous calculations.

Example:

```text
144 ÷ 12 = 12
```

When the user selects another operator:

```text
12 × [ ]
```

The previous result automatically becomes the first operand.

The second operand is cleared.

Focus moves naturally to the second operand.

This interaction should feel similar to the macOS calculator while preserving the cleaner operation-based layout.

---

# Layout Requirements

Desktop:

* centered card
* generous whitespace
* strong visual hierarchy
* operation selector always visible
* result prominently displayed

Mobile:

* mobile-first adaptation
* stacked layout
* comfortable touch targets
* no horizontal scrolling

---

# Operator Selection

Operators must always remain visible.

Preferred presentation:

```text
[ + ] [ − ] [ × ] [ ÷ ] [ xʸ ] [ √ ] [ % ]
```

Avoid dropdowns.

Selection should require a single click.

The selected operator must have a clear active state.

---

# Inputs

Two numeric inputs.

Square Root hides the second operand.

Percentage displays:

```text
Value

Percentage [%]
```

The percent symbol should be visually fixed as an input suffix.

---

# Result

The result should have the strongest visual emphasis.

It should remain visible until the next operation begins.

It should support continuation naturally.

---

# Required UI States

Design the following states:

* Initial
* Editing
* Loading
* Success
* Validation Error
* Backend Error

Loading should be subtle.

Validation should clearly identify the affected operand.

---

# Accessibility

Prioritize accessibility.

Include:

* keyboard navigation
* visible focus states
* sufficient contrast
* ARIA-friendly interaction patterns
* large touch targets
* responsive typography

---

# Visual Language

The application should be inspired by:

* Sezzle
* Stripe
* Linear
* Vercel
* Apple Human Interface Guidelines

Do **not** copy any proprietary interface.

Instead, capture similar principles:

* premium
* minimal
* calm
* modern
* approachable

---

# Color Palette

Neutral background.

One primary accent color inspired by Sezzle's indigo/purple identity.

Restrained color usage.

Avoid:

* excessive gradients
* glassmorphism
* neon
* gaming aesthetics
* skeuomorphic calculator designs

---

# Typography

Use clean modern typography.

Differentiate numeric results from interface labels.

Favor readability.

---

# Motion

Use subtle motion only.

Examples:

* operator selection
* successful calculation
* state transitions

Avoid decorative animations.

---

# Components

Produce a small design system including:

* Buttons
* Operator Buttons
* Inputs
* Cards
* Result Display
* Error Message
* Loading State
* Color Palette
* Typography Scale
* Spacing Scale

---

# Deliverables

Produce:

1. High-fidelity desktop design.
2. High-fidelity mobile design.
3. Basic design system.
4. Interaction notes for important behaviors.
5. Component inventory.

---

# Success Criteria

The design should feel like:

* a modern fintech product;
* a premium engineering portfolio project;
* simple rather than simplistic;
* approachable without being playful;
* elegant without unnecessary decoration.

A reviewer should immediately perceive that usability, accessibility and long-term maintainability guided the design decisions.

---

# Constraints

Do not imitate a physical calculator.

Do not include a numeric keypad.

Do not introduce features outside the assignment scope.

Do not prioritize visual novelty over usability.

Favor clarity, consistency and engineering practicality.
