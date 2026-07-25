ABA-001 — Bootstrap Project & Engineering Foundation

Owner: Elito (Knowledge Curator)

Role: Project Bootstrap · Documentation · Repository Structure

Target: 30–45 min

Hard Stop: 60 min

Objective

Bootstrap the Abacus repository with the initial engineering structure, documentation skeleton, and project layout so implementation can begin immediately.

This task does not implement functionality. It establishes the engineering foundation for the project.

Deliverables
Repository structure

Initialize the following high-level structure:

abacus/
├── .claude/
│   └── skills/
├── .codex/
│   └── skills/
├── .opencode/
│   └── skills/
├── client/
├── server/
├── docs/
│   ├── engineering/
│   │   ├── ENGINEERING_PROCESS.md
│   │   ├── AI_PROMPTS.md
│   │   ├── TECHNICAL_DEBT.md
│   │   └── adr/
│   └── knowledge/
├── .gitignore
├── docker-compose.yml
├── AGENTS.md
├── CLAUDE.md
├── README.md
└── LICENSE (optional)

The frontend and backend folders should also contain an initial placeholder README describing their intended responsibility.

Documentation

Create documentation skeletons (headings only where appropriate).

README.md

Include sections for:

Project Overview
Objectives
Architecture Overview
Technology Stack
Repository Structure
Running the Project
Testing
Docker
Engineering Process
AI-assisted Development
Design Decisions
Future Improvements
ENGINEERING_PROCESS.md

Document the agreed engineering workflow.

Include at least:

Engineering philosophy
Product-first approach
Architecture before implementation
AI-assisted engineering workflow
Human responsibilities
Agent responsibilities
Verification process
Documentation strategy
Trade-off management
Definition of Done
AI_PROMPTS.md

Explain that prompts are versioned as part of the engineering process.

Create placeholders for:

Stitch prompts
Clio prompts
Deep review prompts
Documentation prompts
TECHNICAL_DEBT.md

Create an initially empty debt register.

Columns:

ID
Description
Decision
Status
ADR folder

Create the folder and a README explaining the ADR numbering convention.

No ADRs should be authored yet.

CLAUDE.md

Bootstrap project instructions.

Include engineering conventions already agreed:

Naming
Prefer descriptive names.
Avoid abbreviations.
Avoid one-letter variables except trivial loop indices.
Documentation
Exported Go symbols must include Go documentation comments.
Public React components should include concise documentation when appropriate.
Logging
Use structured logging (log/slog).
Avoid fmt.Println for application logging.
Telemetry
Respect the observability abstraction.
Use the no-op implementation unless otherwise specified.
Errors
Prefer explicit domain errors.
Avoid anonymous string comparisons.
Architecture
Keep abstractions proportional.
Favor composition.
Avoid speculative interfaces.
AGENTS.md

Document the engineering team.

Clio

Main Developer

Responsibilities:

Backend implementation
Frontend implementation
Unit tests
Docker
Deep Pro

Architecture Reviewer

Responsibilities:

Go idiomatic review
React review
Design review
Trade-off review
Final implementation review
Elito

Knowledge Curator

Responsibilities:

Documentation
ADRs
README
Engineering Process
Prompt curation
Technical debt tracking
Knowledge consistency
Knowledge

Create an initial docs/knowledge README explaining that this folder contains long-lived engineering knowledge rather than implementation artifacts.

Explicit decisions already made

Document (without implementing) the following architectural decisions:

React + TypeScript frontend.
Go backend.
Chi HTTP router.
Docker Compose for local execution.
Runtime with graceful shutdown.
Immutable operation registry.
Backend-generated operation manifest.
Shared declarative validation model.
Function-based operation execution.
Operation registry as the single source of truth.
Responsive UI.
Mobile support.
Operation-based calculator (not keypad-based).
Architecture proportional to assignment scope.
Out of Scope

Do not implement:

Backend logic
React components
Dockerfiles
Runtime
Registry
API
Tests

Only prepare the engineering foundation.

Acceptance Criteria
Repository structure is complete.
Documentation placeholders exist.
Engineering conventions are documented.
Agent responsibilities are documented.
Initial architecture decisions are recorded.
The repository is ready for implementation by Clio.
