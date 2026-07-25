# Engineering Process

## Engineering Philosophy

Build the smallest clear solution that meets the product need, preserves correctness, and remains easy to evolve.

## Product-first Approach

Start from the user outcome and acceptance criteria. Technical choices must support the calculator experience rather than exist for their own sake.

## Architecture Before Implementation

Agree on boundaries, ownership, data flow, and meaningful trade-offs before adding implementation. Record durable decisions in ADRs when needed.

## AI-assisted Engineering Workflow

Treat AI as a collaborative engineering tool. Prompts, outputs, review findings, and decisions remain subject to human judgment and repository conventions.

## Human Responsibilities

- Set product intent, priorities, and acceptance criteria.
- Review design trade-offs and implementation outcomes.
- Approve material scope, security, or operational decisions.

## Agent Responsibilities

- Implement assigned work within the agreed architecture.
- Keep documentation and technical debt records current.
- Surface assumptions, risks, and decisions requiring human input.
- Verify work before handoff.

## Verification Process

Every change should be checked at an appropriate level: formatting, static analysis, unit tests, integration checks, and manual behaviour verification where relevant. Review the resulting diff and document any limitations.

## Documentation Strategy

Keep implementation guidance close to the code, durable architectural decisions in ADRs, process guidance in this directory, and reusable long-lived knowledge in `docs/knowledge`.

## Trade-off Management

Make trade-offs explicit. Prefer simple, reversible choices when uncertainty is high, and record meaningful or long-lived choices in an ADR or the technical debt register.

## Definition of Done

Work is done when it meets its acceptance criteria, is reviewed at the appropriate level, is verified proportionately, follows project conventions, and has the required documentation and debt records updated.
