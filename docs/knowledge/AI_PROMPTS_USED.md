# AI prompt inventory

This inventory satisfies the take-home requirement to document AI assistance without treating conversation transcripts as product documentation. The durable task briefs and review records are retained under `docs/engineering/agents-tasks/` and `docs/engineering/agents-reports/`.

| Category | Tools and role | Durable evidence |
| --- | --- | --- |
| UX exploration | Google Stitch explored the responsive calculator visual direction and design-system ideas. | `docs/knowledge/design/` and ABA-002 task record |
| Go architecture | AI assistance helped shape the proportional Go/Chi service, graceful shutdown, and observability boundary; decisions were reviewed against project conventions. | `docs/engineering/adr/`, ABA-003 report |
| Registry and validation | AI-assisted implementation/review work developed the immutable registry, manifest projection, declarative validation, and operation formatting. | ABA-005 through ABA-007.1 reports and reviews |
| Frontend and accessibility | AI assistance supported the manifest-driven React UI, recovery behavior, focus semantics, and accessibility remediation. | ABA-004 and ABA-004-FIX-001 reports |
| Playwright acceptance | Codex assisted with a reproducible desktop/mobile Chromium suite, live-stack routes, evidence capture, and test documentation. | ABA-008 report and `docs/knowledge/testing/` |
| Documentation curation | Codex assisted this final coherence audit, link review, diagrams, and submission-facing README. | ABA-010 report |

Human review remained responsible for scope, repository conventions, architecture decisions, and acceptance of changes. Prompt details are summarized by category because task records and repository evidence are more durable and useful than full conversational transcripts.
