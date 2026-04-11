# Canton Tooling API

Use this command for new or refactored Canton deployment authoring work under `deployment/`.

Treat the user text after `/canton-tooling-api` as the task to execute.

Examples:

- `/canton-tooling-api write a new changeset following this pattern`
- `/canton-tooling-api add a token pool registration changeset`
- `/canton-tooling-api refactor this deployment flow to follow the tooling api pattern`

## Router

Always use the Canton Tooling API hierarchy:

- changeset -> sequence -> operation

For the task:

- keep datastore and environment access in `deployment/changesets`
- keep single-chain orchestration in `deployment/sequences`
- keep one-side-effect contract work in `deployment/operations` or `deployment/utils/operations`

## Required References

Use these as the canonical references:

- `.cursor/rules/references/tooling-api-authoring.md`
- `.cursor/skills/canton-tooling-api-changesets/SKILL.md`
- `.cursor/skills/canton-tooling-api-operations/SKILL.md`

Use the upstream docs referenced in `tooling-api-authoring.md` when adapter boundaries or hierarchy are unclear.

## Response Shape

- Prefer the minimal layer changes required by the task.
- Do not bury datastore writes or orchestration inside leaf operations.
- Keep sequences serializable and single-chain.
- Match existing local naming, versioning, and validation patterns.
