# Engineering Log

This is a durable index of engineering artifacts and accepted decisions. It records outcomes, not conversational history. Prompt usage is recorded here only when it materially influences an engineering artifact or decision.

| Date | Work item | Artifact or accepted decision | Evidence | Status |
| --- | --- | --- | --- | --- |
| 2026-07-24 | ABA-001 — Repository Bootstrap | Repository structure, client/server placeholders, engineering documentation, ADR convention, team responsibilities, and project conventions established. | `README.md`, `AGENTS.md`, `CLAUDE.md`, `docs/engineering/`, `client/README.md`, `server/README.md` | Accepted |
| 2026-07-24 | ABA-001 — Architecture foundation | React + TypeScript client; Go + Chi backend; Docker Compose; graceful shutdown; immutable operation registry; backend-generated manifest; declarative validation; function-based execution; responsive operation-based UI; proportional architecture. | `README.md`, `CLAUDE.md`, `docs/engineering/ENGINEERING_PROCESS.md` | Accepted |
| 2026-07-24 | ABA-SKILL-001 — Go Engineering Skill | Reusable Go implementation and review skill created, curated from the Uber Go Style Guide and adapted for Abacus. | `.claude/skills/osk-go-guide/SKILL.md`, `.claude/skills/osk-go-guide/references/uber-go-guide-notes.md`, `docs/engineering/agents-reports/report-ABA-SKILL-001-clio.txt` | Accepted |
| 2026-07-24 | ABA-SKILL-001R — Documentation alignment | Engineering documentation reorganized around roadmap, artifact log, debt, ADRs, and long-lived knowledge. Standalone prompt history removed; material prompt influence belongs in this log. | `docs/engineering/ROADMAP.md`, `docs/engineering/ENGINEERING_LOG.md`, `docs/knowledge/`, `.claude/skills/osk-go-guide/SKILL.md` | Accepted |
