# ABA-SKILL-001R — Refine `osk-go-guide` & Align Project Documentation

**Status:** READY
**Owner:** Elito
**Role:** Knowledge Curator · Documentation Engineer
**Target:** 30–45 minutes
**Hard Stop:** 60 minutes
**Repository:** `/Users/jonathan/code/Abacus`

---

# Objective

Refine the recently created `osk-go-guide` skill and consolidate the project's engineering documentation based on the architectural review performed after the initial bootstrap.

This task focuses exclusively on improving documentation quality, engineering consistency, and the OSK knowledge base.

No application code should be modified.

---

# Required OSK Skills

Apply the project's existing documentation and knowledge-curation skills.

This task is documentation-first.

---

# Part 1 — Project Documentation Refinement

Update the project documentation to reflect the engineering decisions agreed after repository bootstrap.

## 1. Documentation Structure

Ensure the repository documentation is organized as:

```text
docs/

engineering/

    ROADMAP.md

    ENGINEERING_PROCESS.md

    ENGINEERING_LOG.md

    TECHNICAL_DEBT.md

    adr/

knowledge/

    use-cases/

    acceptance-tests/

    architecture/

    glossary/
```

Remove any obsolete documentation that no longer matches the engineering process.

---

## 2. Remove AI_PROMPTS

The standalone AI prompts document is no longer part of the engineering model.

Remove it.

Prompt history should instead be recorded inside **ENGINEERING_LOG.md** whenever relevant.

---

## 3. ROADMAP

Create (or update) the engineering roadmap.

The roadmap should represent implementation intent rather than history.

Initial entries should include at least:

* ABA-001 Repository Bootstrap
* ABA-SKILL-001 Go Engineering Skill
* ABA-002 UX Exploration (Google Stitch)
* Backend Foundation
* Operation Registry
* Validation Model
* REST API
* React Integration
* Docker & Runtime
* Documentation & Final Review

Additional items may be added if they improve planning clarity.

---

## 4. ENGINEERING_LOG

Update the Engineering Log using the same philosophy established in Arbitrier.

The log should be a durable engineering record rather than a conversational diary.

Each entry should reference engineering artifacts rather than discussions.

Preserve the existing table-oriented style.

Record:

* bootstrap
* repository structure
* documentation decisions
* Go skill creation
* accepted architectural decisions

Prompt usage should be recorded here whenever it materially influences the project.

---

## 5. Knowledge Base

Populate the knowledge structure with placeholder documents.

Create:

```text
knowledge/

use-cases/

acceptance-tests/

architecture/

glossary/
```

Include concise README/index files describing the responsibility of each section.

Do not author the complete functional documentation yet.

---

# Part 2 — Refine `osk-go-guide`

Review the generated Go Engineering Skill and apply the architectural review feedback.

The goal is to evolve the skill from a curated guide into a reusable OSK engineering skill.

---

## 1. Rename

Rename the primary title from:

> Go Code Quality Skill

to:

> Go Engineering Skill

The skill now represents engineering guidance rather than code style alone.

---

## 2. Add Engineering Philosophy

Introduce a short section describing the OSK engineering philosophy.

Suggested principles:

* Keep the happy path obvious.
* Prefer explicit ownership.
* Favor small composable abstractions.
* Prefer immutable configuration.
* Keep architecture proportional to problem size.
* Optimize for maintainability before cleverness.

This section should complement—not replace—the Uber guidance.

---

## 3. Add Decision Heuristics

Introduce a new section:

> Decision Heuristics

This section should help implementation agents decide **when** to apply a rule rather than merely stating the rule.

Examples include:

* When to introduce an interface.
* When a concrete type is preferable.
* When composition is sufficient.
* When function types are preferable to interface hierarchies.
* When immutability is the correct choice.
* When additional abstraction is justified.

These heuristics should teach engineering judgment rather than enforce mechanical rules.

---

## 4. Add "When to Deviate"

Create a short section explaining that engineering principles occasionally require deviating from generic style guidance.

Examples:

* API consistency may outweigh an internal style recommendation.
* Simplicity may outweigh theoretical extensibility.
* Product requirements may justify exceptions.

The goal is to reinforce engineering reasoning rather than dogmatism.

---

## 5. Add Ownership Rules

Expand lifecycle guidance.

Introduce a concise ownership model covering:

* Who creates a resource.
* Who owns the resource.
* Who shuts it down.
* Why ownership boundaries matter.

Relate this to:

* runtimes
* HTTP servers
* goroutines
* registries
* telemetry

---

## 6. Add BAD / GOOD Examples

Introduce small illustrative examples.

Examples should remain concise.

Potential topics:

* speculative interfaces
* registry ownership
* concrete vs interface injection
* descriptive naming
* function-based strategies
* unnecessary hierarchy

The examples should educate rather than overwhelm.

---

## 7. Preserve Existing Quality

Retain the existing:

* Severity Model
* Review Checklist
* Anti-patterns
* Review Output
* Source Attribution

Only improve them where necessary.

---

# Validation

Before completing the task verify that:

* documentation structure reflects the agreed engineering model;
* Engineering Log remains navigable;
* roadmap is coherent;
* AI prompt history is no longer maintained separately;
* the Go Engineering Skill remains concise;
* no substantial portions of the Uber Guide have been copied;
* the skill feels like an OSK artifact rather than a summarized external guide.

---

# Deliverables

Expected modified artifacts include:

* `docs/engineering/ROADMAP.md`
* `docs/engineering/ENGINEERING_LOG.md`
* `docs/engineering/ENGINEERING_PROCESS.md` (if needed)
* `docs/knowledge/*`
* `.claude/skills/osk-go-guide/SKILL.md`
* supporting reference notes only if required

---

# Report

Produce a standard Engineering Report summarizing:

* documentation updates;
* structural changes;
* skill improvements;
* rationale behind each refinement;
* validation performed;
* outstanding recommendations.

---

# Constraints

* Do not modify application code.
* Do not introduce new skills.
* Do not expand project scope.
* Keep documentation proportional to the size of the assignment.
* Preserve the project's engineering philosophy.

---

# Definition of Done

* Documentation structure matches the agreed engineering model.
* Engineering Log follows the Arbitrier-style artifact index.
* Roadmap exists and reflects planned work.
* Knowledge folders are initialized.
* `osk-go-guide` incorporates the agreed refinements.
* The project documentation and Go skill present a consistent OSK engineering philosophy.
