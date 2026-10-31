# Sparkwing staging grid fixtures

This throwaway repository supplies GitHub App tests for staging. It contains no
production credentials or deployment target.

The default branch uses the Go project in `.sparkwing`. It declares hello,
secret/OIDC readers, an intentional failure, child/pinned-source fixtures and
controller schedules. `fixtures/go/.sparkwing` preserves the original Go seed.
The `fx-pinned` branch supplies the mutable branch pin used by `fx-pin`.

`ts-public/.sparkwing` uses public npm `@sparkwing/sdk` version 0.68.2.
`ts/.sparkwing` uses a private local package stamped 0.68.0-qa.20, packed from
engine commit f8f72274c6d5c17bcf72285bb30098c55486f538. The package stamp selects
the matching candidate engine. Both projects define build/publish, marketing
build/link-check/preview and PR fixtures. `fx-typescript` preserves the commit
with this TypeScript project at the repository root, because the root project
selects one runtime at a time.

A lane creates its own shared `QA_SHARED` and repo/pipeline-scoped `QA_PIPELINE`
values through stdin, then removes them afterward. Go readers log lengths and
decoded OIDC claims, never secret values or raw tokens. On qa.20, optional secret
declarations still hit the required-secret preflight when absent.

The root declares hello every five minutes, hello hourly in America/Boise and an
every-minute 180-second job with global concurrency capacity one. A staging admin
must enable repository automation and push the default branch to arm them.
After testing, the lane pauses App schedules and removes automation consent.
Historical paused App rows can remain because the App owns those declarations.

Temporary PRs, tags and test branches are removed after verification. Runs remain
as terminal history. Trigger subscriptions, created secrets and temporary runner
registrations are removed. The `fx-pinned` and `fx-typescript` source branches
remain reusable.

## Dependencies

- Go `github.com/sparkwing-dev/sparkwing` v0.68.0-qa.20 supplies the tested SDK
  protocol. This candidate pin stays until the staging candidate changes.
- npm `@sparkwing/sdk` supplies the TypeScript implementation of that protocol.
  The public project uses 0.68.2. The staging fixture uses the matching qa.20
  package so the test executes the candidate engine.
