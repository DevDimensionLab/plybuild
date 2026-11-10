// Package quota reads the remaining subscription quota that Codex and Claude
// report for the whole account, without reading any credentials.
//
// ReadCodex runs one short-lived `codex app-server` (experimental upstream)
// exchange: initialize, initialized, account/rateLimits/read. It starts no
// thread, turn or inference, uses one total deadline (DefaultCodexTimeout),
// skips notifications and invalid lines, limits each line to 1 MiB and kills
// and waits for its own process group on every return. Raw results, stderr and
// upstream messages are never returned; the message is only used to detect
// ErrNotSignedIn.
//
// Claude Code reports rate_limits only to its status line command. The
// claude-statusline-relay executable calls RecordClaudeStatusLine, which keeps
// only five_hour/seven_day in the v1 snapshot shared with AI Usage, and ReadClaude
// reads that snapshot when ~/.claude/settings.json (or ClaudeOptions) names the
// relay as the first shell word of statusLine.command. Project settings and
// CLAUDE_CONFIG_DIR are not discovered. Nothing here changes settings.
//
// Windows returns ErrUnsupportedPlatform. Default staleness: Codex after
// StaleAfterCodex, Claude after StaleAfterClaude. Remaining is never derived
// from a missing value, and Select never adds or averages windows.
//
// Discovery on 2026-11-07 (read-only, schema generated to a temporary dir):
// codex-cli 0.164.0 GetAccountRateLimitsResponse requires usedPercent and has
// nullable windowDurationMins/resetsAt; Claude Code 2.1.309 documents
// rate_limits.five_hour/seven_day. Fixtures are sanitized; live parity with the
// providers' own usage pages is not verified by the tests.
package quota
