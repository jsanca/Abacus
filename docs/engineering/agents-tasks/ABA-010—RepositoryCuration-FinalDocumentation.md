ABA-010 — Repository Curation & Final Documentation
Field	Value
Status	Ready
Owner	Elito
Role	Knowledge Curator · Technical Writer · Frontend Developer
Target	45–60 minutes
Hard stop	75 minutes
Commit	No
Repository	/Users/jonathan/code/Abacus
Objective

Perform the final documentation and repository curation pass before submission.

This task does not introduce new functionality.

Its objective is to ensure the repository tells a coherent engineering story and satisfies every deliverable requested by Sezzle.

Think like an external reviewer opening the repository for the first time.

Required Skills

Apply:

.claude/skills/process/osk-execution-timebox/SKILL.md
.claude/skills/process/osk-engineering-reporting/SKILL.md
.claude/skills/documentation/osk-knowledge-kurator/SKILL.md

Follow:

CLAUDE.md
AGENTS.md
Primary Goal

The repository should answer four questions without requiring code inspection:

What is this project?
How do I run it?
Why was it designed this way?
How do I verify it works?
Scope
1. README Final Pass

Review the README from top to bottom.

Ensure it contains, at minimum, the items explicitly requested by Sezzle.

Required

✅ Project overview

✅ Features

✅ Setup instructions

Prerequisites
Docker
Go
Node

or whichever are actually required.

Running the project

Explain clearly:

make build

make run

make verify

make e2e

or the canonical commands actually implemented.

Document:

frontend
backend
Docker
Playwright
REST examples

Document:

GET /api/v1/operations

POST /api/v1/calculations

Include request and response examples.

Do not duplicate the full API Contract; link to it where appropriate.

Design decisions

Summarize the architectural ideas:

immutable operation registry;
backend-generated manifest;
declarative validation;
manifest-driven frontend;
operation-owned execution and expression formatting;
graceful shutdown;
Docker-first execution;
layered testing strategy.

Keep this concise (approximately one page or less).

Testing

Explain the testing pyramid:

Go unit tests

↓

React component tests

↓

Playwright E2E

Reference the acceptance evidence.

Screenshots

Embed:

desktop screenshot;
mobile screenshot.

Use the curated evidence already produced.

Repository structure

Briefly explain:

client/
server/
docs/

No need to describe every folder.

AI Usage

The assignment explicitly requests the prompts used.

Document:

which AI tools participated;
their roles;
where prompts can be found;
the engineering process followed.

Do not overemphasize AI; the repository should remain centered on engineering.

2. Documentation Integrity Review

Review all documentation under:

docs/

Verify:

broken links;
outdated filenames;
renamed reports;
stale references;
missing screenshots;
invalid Markdown links;
duplicate sections.

Repair every broken internal reference.

3. Engineering Log

Review:

docs/engineering/ENGINEERING_LOG.md

Verify:

chronological consistency;
correct task identifiers;
working links;
final entries;
no duplicated checkpoints.

The engineering log should tell the story of the project from bootstrap to completion.

4. Roadmap

Update:

docs/engineering/ROADMAP.md

Mark completed initiatives appropriately.

Do not delete completed work.

The roadmap should reflect the actual delivery journey.

5. Technical Debt

Review:

docs/engineering/TECHNICAL_DEBT.md

Keep only genuine accepted trade-offs.

Do not list hypothetical future enhancements.

6. API Contract

Review:

docs/knowledge/api/API_CONTRACT.md

Ensure it matches the implemented API exactly.

Verify:

request examples;
response examples;
percentage expression;
error contract;
validation responses.
7. Testing Documentation

Review:

docs/knowledge/testing/

Ensure consistency between:

TEST_STRATEGY
ACCEPTANCE_TEST_RESULTS
screenshots
README

No contradictory commands.

8. Architecture Diagrams

Add two lightweight diagrams.

Component View
Browser
    │
React
    │
Nginx
    │
Go API
    │
Operation Registry
Sequence
User

↓

React

↓

Manifest

↓

Validation

↓

Calculation

↓

Result

Markdown Mermaid is acceptable.

Keep diagrams intentionally simple.

9. Repository Polish

Verify:

consistent naming;
no temporary files;
no TODOs forgotten;
no debug comments;
no generated artifacts committed;
.gitignore correctness.
10. Prompt Inventory

Create:

docs/knowledge/AI_PROMPTS_USED.md

Summarize the prompts used during the project.

Do not dump entire conversations.

Instead describe categories such as:

UX exploration
Go architecture
Registry design
Playwright generation
Documentation

This satisfies Sezzle's requirement while keeping the repository professional.

Verification

Execute:

make verify

Run:

make e2e-local

Verify:

README commands work;
documentation links resolve;
screenshots render correctly on GitHub Markdown.
Deliverables

Final curated repository including:

polished README;
repaired documentation;
completed roadmap;
verified engineering log;
API contract alignment;
architecture diagrams;
testing documentation;
AI prompt inventory;
repository cleanup.
Report

Create:

docs/engineering/agents-reports/ABA-010-repository-curation.md

Include:

documentation reviewed;
links repaired;
README improvements;
diagrams added;
Sezzle requirements checklist;
final repository observations;
files changed;
no-commit confirmation.
Definition of Done
README satisfies every Sezzle requirement.
Setup instructions are accurate.
Frontend, backend and Docker execution are documented.
REST examples are included.
Design decisions are summarized.
Architecture diagrams are present.
Engineering Log is consistent.
Roadmap reflects the completed project.
API Contract matches implementation.
Testing documentation is coherent.
Internal links are valid.
Prompt inventory is documented.
Repository contains no obvious documentation inconsistencies.
No commits were created.
