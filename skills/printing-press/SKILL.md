---
name: printing-press
description: Set up a new integration, connector, or CLI binding for any API. Wrap or generate a ship-ready Go CLI from an OpenAPI, HAR, or Postman spec via the lean research -> generate -> build -> shipcheck loop. Use when the user says build a CLI, wrap this API, set up a new integration, add a connector, integrate with a service, or names an API by domain.
version: 2.0.0
min-binary-version: "4.0.0"
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Glob
  - Grep
  - WebFetch
  - WebSearch
  - AskUserQuestion
  - Agent
created_by: user
---

# /printing-press

Generate the best useful CLI for an API without burning an hour on phase theater.

```bash
/printing-press Notion
/printing-press Discord codex
/printing-press --spec ./openapi.yaml
/printing-press --har ./capture.har --name MyAPI
/printing-press https://postman.com/explore
/printing-press https://postman.com
```

## What Changed In v2

The old skill inflated the path to ship:
- too many mandatory research documents before code existed
- too many separate late-stage validation phases after code existed
- too many chances to discover obvious failures late

This version uses one lean loop:
1. Resolve the spec and write one research brief
2. Generate
3. Build the highest-value gaps
4. Run one shipcheck block
5. Optionally run live API smoke tests

Artifacts are still written, but only the ones that materially help the next step.

## Modes

### Default

Normal mode. Claude does research, generation orchestration, implementation, and verification.

### Codex Mode

If the arguments include `codex` or `--codex`, offload pure code-writing tasks to Codex CLI.

Use Codex for:
- writing store/data-layer code
- writing workflow commands
- fixing dead flags / dead code / path issues
- README cookbook edits

Keep on Claude:
- research and product positioning
- choosing which gaps matter
- verification results and ship decisions

If Codex fails 3 times in a row, stop delegating and finish locally.

### Polish Mode (Standalone Skill)

For second-pass improvements to an existing CLI, use the standalone polish skill:

```bash
/printing-press-polish redfin
```

See the `printing-press-polish` skill for details. It runs diagnostics, fixes verify failures, removes dead code, cleans up descriptions and README, and offers to publish.

## Rules

- **Do not ship a CLI that hasn't been behaviorally tested against real targets.** `go build` and `verify` pass-rate are structural signals, not correctness signals. Phase 5's mechanical test matrix runs every subcommand + `--json` + error paths; if that matrix was not executed, the CLI is not shippable. Quick Check is the floor; Full Dogfood is required when the user asks for thoroughness.
- **Bugs found during dogfood are fix-before-ship, not "file for v0.2".** If a 1-3 file edit resolves it, do it now. `ship-with-gaps` is deprecated as a default verdict (see Phase 4). Context is freshest in-session; a v0.2 backlog that may never be revisited ships known-broken CLIs.
- **Features approved in Phase 1.5 are shipping scope.** Do not downgrade a shipping-scope feature to a stub mid-build. If implementation becomes infeasible, return to Phase 1.5 with a revised manifest and get explicit re-approval.
- **Do not quote human-time estimates for sub-tasks** ("~15-30 min", "~1 hour", "quick fix") in `AskUserQuestion` options, phase descriptions, or reference docs. The agent does the work, not the user; agent-fabricated estimates are notoriously bad and train users to distrust the prompt. Describe scope instead (lines of code, files touched, relative size). The carve-outs are wall-clock estimates for genuinely time-bound things: the whole-CLI run (set the user's expectation up front — most CLIs take 30+ minutes), tool installs (`go install` takes ~10 seconds), and printing-press subcommands that do network-bound work (crowd-sniff scans npm + GitHub, ~5-10 minutes). Anything bounded by agent reasoning time is not time-bound — describe scope.
- **Use raw captures for contract research.** When reading official docs, auth/error/rate-limit pages, endpoint references, OpenAPI/Postman links, or source pages whose exact identifiers affect the generated CLI, read [references/fetch-docs.md](references/fetch-docs.md) and use its `fetch-docs.sh` helper. Reserve `WebFetch` for quick TL;DR reads where losing field-level details is acceptable.
- Optimize for time-to-ship, not time-to-document.
- Reuse prior research whenever it is already good enough.
- Do not split one idea across multiple mandatory artifacts.
- Durable files produced by this skill go under `$PRESS_RUNSTATE/` (working state) or `$PRESS_MANUSCRIPTS/` (archived). Short-lived command captures may use `/tmp/printing-press/` and must be removed after use.
- Do not create a separate narrative phase for dogfood, dead-code audit, runtime verification, and final score. Treat them as one shipcheck block.
- Run cheap, high-signal checks early.
- Fix blockers and high-leverage failures first.
- Reuse the same spec path across `generate`, `dogfood`, `verify`, and `scorecard`.
- YAML, JSON, local paths, and URLs are all valid spec inputs for the verification tools.
- Maximum 2 verification fix loops unless the user explicitly asks for more.

## Secret & PII Protection (Cardinal Rules)

**These rules are non-negotiable. They apply at ALL times during a run.**

API key **values**, token **values**, passwords, and session cookies must NEVER
appear in any artifact: source code, manuscripts, proofs, READMEs, HARs, or
anything committed to git. Env var **names** (e.g., `STEAM_API_KEY`) and
placeholders (e.g., `"your-key-here"`) are safe.

During Phase 5.6 (archiving) and before publishing, read and apply
[references/secret-protection.md](references/secret-protection.md) for:
- Exact-value scanning and auto-redaction of artifacts
- HAR auth stripping (headers, query strings, cookies)
- API key handling rules during the run
- Session state cleanup ordering


## Preflight

**This section MUST run before any user-facing prompt — including the Orientation and Briefing flow below.** A missing binary or available upgrade is information the user needs *before* they commit to an API. Do not invoke `AskUserQuestion`, print the orientation prose, or otherwise engage the user until preflight has completed and any signals from `references/setup-checks.md` have been handled.


Execute the setup contract by running the bundled script in a single Bash tool call from the session's working directory (substitute this skill's directory for `<skill-dir>`):

```bash
bash <skill-dir>/scripts/preflight.sh
```

The script is the contract bash block, extracted verbatim; it emits the same stdout markers documented below and in [references/setup-checks.md](references/setup-checks.md). Do not paraphrase or re-implement it inline.


**MANDATORY: Read and apply [references/setup-checks.md](references/setup-checks.md) immediately after the setup contract bash block runs, before any other action.** It handles the contract output signals: `[setup-error]` (refuse to run, surface the install instructions), optional `[local-binary-stale]` / `[local-binary-rebuilt]` repo-mode rebuild markers, `[go-toolchain-old]` / `[low-disk]` advisories, `[repo-upgrade-available]` (interactive `AskUserQuestion` prompt + optional repo pull), `PRESS_REPO_MODE=<true|false>` plus the targeted global open-agent-skills freshness check, the min-binary-version compatibility check (hard stop if binary is too old), `[upgrade-required]` (hard gate below the published currency floor — interactive upgrade-or-abort, no skip), `[upgrade-available]` (interactive `AskUserQuestion` prompt + optional standalone binary upgrade), `[browser-tools-missing]` (interactive `AskUserQuestion` prompt + optional install of browser-use and/or agent-browser), and the `PRINTING_PRESS_BIN=<abs-path>` marker plus optional `[binary-shadow]` warning (capture the path; use it for every subsequent generator invocation). Skipping the reference will cause the skill to proceed with a missing or out-of-date binary, run with stale global skill text when the session is managed by open-agent-skills, hit a mid-flight install prompt if browser-sniff is later needed, or invoke the wrong binary because a stale global or the public-library installer on `PATH` shadowed the local build. Do not skip.

**Absolute-path rule.** The preflight contract always emits `PRINTING_PRESS_BIN=<absolute path>` to stdout. Capture this value and substitute it (the resolved absolute path, not the literal `$PRINTING_PRESS_BIN` token) for every subsequent `cli-printing-press ...` invocation in this skill, references, and any sub-skill you delegate to. The `export PATH=...` line inside the contract only affects the single Bash tool call it runs in; later Bash tool calls open fresh shells and resolve bare `cli-printing-press` against the user's default `PATH`, where a stale globally-installed binary (`$HOME/go/bin/cli-printing-press`, Homebrew copy, etc.) will silently shadow the local build the preflight just chose. Bash code examples below are written `cli-printing-press generate ...` for readability — replace `cli-printing-press` with the captured absolute path each time you actually run one.

Only after preflight completes successfully (no `[setup-error]`; no `[upgrade-required]` left unresolved — the user either upgraded or the run was aborted; no global skill update that requires restart; any `[repo-upgrade-available]`, `[upgrade-available]`, or `[browser-tools-missing]` was offered to the user; `PRINTING_PRESS_BIN` is captured) should you proceed to the Orientation & Briefing section below.


## Outputs

Every run writes up to 5 concise artifacts under the current managed run and archives them to `$PRESS_MANUSCRIPTS/<api-slug>/<run-id>/`:

1. `research/<stamp>-feat-<api>-pp-cli-brief.md`
2. `research/<stamp>-feat-<api>-pp-cli-absorb-manifest.md`
3. `proofs/<stamp>-fix-<api>-pp-cli-build-log.md`
4. `proofs/<stamp>-fix-<api>-pp-cli-shipcheck.md`
5. `proofs/<stamp>-fix-<api>-pp-cli-live-smoke.md` (only if live testing runs)

These do not need to be 200+ lines. Keep them dense, evidence-backed, and directly useful.


## Workflow

The phase playbooks live in `references/`, one file per phase, listed here in execution order. **Before starting a phase, read its reference file in full** — the file is the phase; the line below is only a map entry. Skipping the read means running the phase without its gates and halt conditions.

| Order | Phase | Read first | One-line map |
|---|---|---|---|
| 1 | Orientation & Briefing | [references/orientation-briefing.md](references/orientation-briefing.md) | No-args orientation vs. with-args briefing; multi-source priority gate. |
| 2 | Run Initialization | [references/run-initialization.md](references/run-initialization.md) | Create the managed run dirs ($API_RUN_DIR, $SESSION_DIR) with the exact umask/containment semantics in the file. |
| 3 | Phase 0: Resolve And Reuse | [references/phase-0-resolve-and-reuse.md](references/phase-0-resolve-and-reuse.md) | Resolve the API to a spec source; reuse existing library work before any research. |
| 4 | Phase 1: Research Brief | [references/phase-1-research-brief.md](references/phase-1-research-brief.md) | Produce the brief using the template sections (API Identity … Build Priorities) in the file. |
| 5 | Phase 1.6: Pre-Browser-Sniff Auth Intelligence | [references/phase-1-6-auth-intelligence.md](references/phase-1-6-auth-intelligence.md) | Auth reconnaissance before any browser capture. |
| 6 | Phase 1.7: Browser-Sniff Gate | [references/phase-1-7-browser-sniff-gate.md](references/phase-1-7-browser-sniff-gate.md) | Decide and run authenticated browser capture; session-state containment rules. |
| 7 | Phase 1.8: Crowd-Sniff Gate | [references/phase-1-8-crowd-sniff-gate.md](references/phase-1-8-crowd-sniff-gate.md) | Crowd-sourced capture fallback gate. |
| 8 | Phase 1.5: Ecosystem Absorb Gate | [references/phase-1-5-absorb-gate.md](references/phase-1-5-absorb-gate.md) | Absorb manifest of incumbent features; **generation is blocked until this gate is approved**. |
| 9 | Phase 1.9: API Reachability Gate | [references/phase-1-9-reachability-gate.md](references/phase-1-9-reachability-gate.md) | Reachability probe + risk classification before generate. |
| 10 | Phase 2: Generate | [references/phase-2-generate.md](references/phase-2-generate.md) | Drive `cli-printing-press generate` (via the captured absolute path) to a compiling tree. |
| 11 | Phase 3: Build The GOAT | [references/phase-3-build-the-goat.md](references/phase-3-build-the-goat.md) | Build every absorb-manifest feature, then transcend; hand-edit durability rules. |
| 12 | Phase 4: Shipcheck | [references/phase-4-shipcheck.md](references/phase-4-shipcheck.md) | Run shipcheck; interpret the scorecard. |
| 13 | Phase 4.7–4.95: Review Gates | [references/phase-4-review-gates.md](references/phase-4-review-gates.md) | Sync param-drop gate, agentic SKILL/output reviews, README/SKILL/AGENTS audit, local code review. |
| 14 | Phase 5: Dogfood Testing | [references/phase-5-dogfood.md](references/phase-5-dogfood.md) | Live use of the CLI as proof; smoke-log artifact. |
| 15 | Phase 5.5: Polish | [references/phase-5-5-polish.md](references/phase-5-5-polish.md) | Bounded cosmetic pass — no late-phase gauntlet rerun. |
| 16 | Phase 5.6: Promote and Archive | [references/phase-5-6-promote-archive.md](references/phase-5-6-promote-archive.md) | Promote into $PRESS_LIBRARY (swap vs. regen-merge preservation) and archive manuscripts. |
| 17 | Phase 6: Next Steps | [references/phase-6-next-steps.md](references/phase-6-next-steps.md) | Post-ship follow-ups and handoff. |

Notes that survive across all phases: the execution order intentionally runs 1.6 → 1.7 → 1.8 **before** 1.5; the absolute-path rule above applies to every `cli-printing-press` invocation in every reference file; Secret & PII cardinal rules apply at all times.

## Fast Guidance

### When to use `cli-printing-press print`

Use `cli-printing-press print <api>` only when the user explicitly wants a resumable on-disk pipeline with phase seeds. It is optional.

The fast path for `/printing-press <API>` is:
- brief
- generate
- build
- shipcheck

### When to stop researching

Stop when you can answer:
- what to build first
- what data to persist
- what incumbent features cannot be missing

If the next research step does not change those answers, stop and generate.

### What not to do

Do not:
- write 5 separate mandatory research documents
- defer all workflows to "future work"
- skip verification because the CLI compiles
- treat scorecard alone as ship proof
- discover YAML/URL spec incompatibility late and manually convert specs if the tools can already consume them
- rerun the whole late-phase gauntlet for cosmetic README polish
- skip features because "the MCP already handles that" (absorb everything, beat it with offline + agent-native)
- build only "top 3-5 workflows" when the absorb manifest has 15+ (build them ALL, then transcend)
- generate before the Phase 1.5 Ecosystem Absorb Gate is approved
- call a CLI "GOAT" without matching every feature the best competitor has

### What counts as success

Success is:
- a generated CLI that gets to shipcheck without generator blockers
- verification tools working against the same spec the user generated from
- one or two fix loops, not a maze of re-entry phases
- a CLI that is plausibly shippable today, not a perfect design memo

