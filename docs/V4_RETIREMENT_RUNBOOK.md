# Runbook: Retiring a cmdguard Major (v4 Case)

**Status:** v4 retirement EXECUTED — frozen at `v4.3.1` on the `release/v4`
branch; master develops v5 (`github.com/larsartmann/cmdguard/v5`). This
runbook documents how it was done so the v5 retirement (much later) follows
the same path.

## The four steps, in order

### 1. Freeze the branch

Cut `release/vN` from the last commit that still develops that major:

```bash
git switch -c release/vN master        # from the pre-next-major state
git push -u origin release/vN
```

After the freeze, the branch takes ONLY:

- security fixes
- regressions that break vN consumers
- documentation corrections

Everything else waits for the next major. If a fix lands on `release/vN`,
it must be mirrored (or already superseded) on master — the branch is a
maintenance line, not a fork.

### 2. Cut the final tag

Tag the last release ON the frozen branch, not from master:

```bash
scripts/pre-tag-check.sh vN.X.Y        # full battery first
git tag -a vN.X.Y origin/release/vN -m "Final vN release; branch frozen"
git push origin vN.X.Y
```

The release workflow runs tests + lint on the tag and publishes. For v4 the
final tag was **v4.3.1**. Verify the release run is green before declaring
the branch frozen — a red run on the final tag is forever (the v4.3.0
release check is red permanently; v4.3.1 exists so consumers have a green
final version).

### 3. Mark the CHANGELOG and docs

- CHANGELOG: the final version section notes the retirement
  ("last v4 release; branch frozen").
- README: the freeze claim must state the TRUE final version — the
  2026-11-11 audit found the README claiming v4.2.1 while truth was v4.3.1.
- Website: a retirement notice on the site linking the frozen branch and the
  migration guide (`/guides/migrating-v4-to-v5/`).

### 4. Point discovery at the new major

- `lint` rule **CG002** fires when consumer code imports a retired major
  (walk.go `CurrentMajor` must be bumped when the new major ships).
- The migration guide + `docs/migration_v4_v5/codemod.sh` automate the path.
- `git tag -l` keeps ALL historical tags — never delete them; consumers pin
  them and the module proxy serves them permanently.

## What NOT to do

- **Never delete a retired branch or its tags.** The proxy caches; deletion
  breaks checksum-pinned consumers.
- **Never re-tag a version.** v4.3.1 exists; a force-moved v4.3.1 poisons
  every `go.sum` that pinned it.
- **Never tag the retired branch from master.** Tags must contain the
  retired major's own module path and code.
- **Never let the frozen branch drift its sub-module pins.** The release/v4
  hygiene question (D10 in the 2026-11-11 plan) is exactly this class.

## Verification checklist

- [ ] `release/vN` exists and is pushed
- [ ] Final tag's release workflow run is GREEN (check before announcing)
- [ ] `go get github.com/larsartmann/cmdguard/vN@vN.X.Y` resolves from the proxy
- [ ] CHANGELOG + README state the true final version
- [ ] Website retirement notice live (smoke: `python3 scripts/check-website-smoke.py`)
- [ ] CG002 `CurrentMajor` bumped in the lint module when the NEXT major ships
