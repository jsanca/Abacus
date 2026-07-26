# Engineering Roadmap

This roadmap describes intended delivery order. It is not a historical changelog; completed work remains here only to show the foundation on which later work depends.

| Order | ID | Initiative | Intent | Status |
| --- | --- | --- | --- | --- |
| 1 | ABA-001 | Repository Bootstrap | Establish repository structure and the engineering foundation. | Complete |
| 2 | ABA-SKILL-001 | Go Engineering Skill | Provide reusable Go implementation and review guidance. | Complete |
| 3 | ABA-002 | UX Exploration (Google Stitch) | Establish the visual direction and design-system foundation. | Complete |
| 4 | ABA-004 | Responsive React Frontend | Deliver the manifest-driven calculator UI against a mocked API boundary. | Complete |
| 5 | ABA-004-FIX-001 | Frontend Review Remediation | Resolve lifecycle, recovery, focus, and accessibility findings from ABA-004R. | Complete |
| 6 | ABA-003 | Backend Foundation | Establish the Go service composition, configuration, lifecycle, and observability boundary. | Complete |
| 7 | ABA-004.1 | Client Containerization | Build and serve the production React client alongside the backend through Docker Compose. | Complete |
| 8 | ABA-005 | Operation Registry | Implement the immutable operation registry and backend-generated manifest. | Complete |
| 9 | ABA-006 | Validation Model | Define and apply declarative operation validation. | Complete |
| 9.1 | ABA-006.1 | Cohesive Validation Refinement | Polymorphic expression evaluation and operation-owned validation. | Complete |
| 10 | ABA-007 | REST API & Frontend Integration | Expose operations and calculation endpoints; replace mock adapter with live HTTP client. | Complete |
| 10.1 | ABA-007.1 | Smoke-Test Remediation | Fix percentage expression notation; verify and test calculation continuation focus; add all-operation format tests. | Complete |
| 11 | A-008 | Reproducible Playwright E2E Acceptance Suite | Exercise the complete browser-to-backend calculation flow through a Compose-owned acceptance suite. | Complete |
| 12 | ABA-009R | Final Full-Stack Architecture Delivery Review | Independently assess delivery readiness across architecture, testing, containerization, and documentation. | Complete |
| 13 | ABA-010 | Repository Curation & Final Documentation | Curate submission-facing documentation, evidence, prompt inventory, and repository consistency. | Complete |
