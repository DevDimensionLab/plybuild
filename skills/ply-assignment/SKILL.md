---
name: ply-assignment
description: Work on a Ply Assignment named at session start, e.g. `/ply-assignment ASSIGNMENT-7` or `/ply-assignment ASSIGNMENT-7 --project ply`. Read it with `ply assignment show`, take it before working, mark it done when finished, and follow a Task it mentions read-only.
argument-hint: ASSIGNMENT-<n> [--project <project-id>]
---

# Ply Assignment

Arguments: $ARGUMENTS

An Assignment is the human's work description. Its text, the user and the
repository's AGENTS.md/CLAUDE.md decide the work and its authority; this skill
only says how to read, take and finish it. It grants nothing beyond them.

Use `--project <project-id>` in every `ply assignment` command below only when
the arguments give one.

## 1. Read the Assignment

```sh
ply assignment show ASSIGNMENT-<n> [--project <project-id>] --format json
```

Done when you have the JSON. On an error, report the class and detail as
printed and stop for the human: on `workspace_project_invalid_arguments` ask
which Project; otherwise do not guess the content, read the JSON files
directly, or create another Assignment.

Then look at `lifecycle`:

- `state` is `done`: report that it is done, by whom and when, and stop.
- `state` is `taken` and `lifecycle.holder.actor` is not this session's own
  label: report the holder and `taken_at` and stop. Do not work in parallel.
- `state` is `taken` by this session's own label (a take whose answer was
  lost, or a resumed session): adopt `lifecycle.holder.claim` and go to step 3.
- `state` is `available`: go to step 2 with `r` = `lifecycle.revision`.

## 2. Take it before any work

Your label is `claude:${CLAUDE_SESSION_ID}`. Claude Code fills in the session
id when it loads this skill and keeps it when the session is resumed. If the
label still shows the placeholder or an empty id, use `claude:` plus random hex
instead. Use the same label for the whole session.

```sh
ply assignment take ASSIGNMENT-<n> --revision <r> --as <label> [--project <project-id>] --format json
```

Keep `lifecycle.holder.claim` from the answer for the session. On
`assignment_conflict`, re-run step 1 once and follow it; never retry take in a
loop and never start work without a claim. If the answer is uncertain
(`result unknown, run show`), run step 1 before anything else.

## 3. Follow references when the work needs them

When the text names a Ply Task, read it only through:

```sh
ply workspace task show <task-id> --format json
ply workspace task spec list <task-id> --format json
ply workspace task spec show <task-id> --spec <spec-id> --revision <n> --format json
```

and the documents that readback lists. This is background, not a new mandate:
register, start or change nothing in the Task. If a needed reference cannot be
read, say so and ask before relying on it.

## 4. Do the work

Work in the current session. `repos` says which repositories the Assignment
concerns. Write where the Assignment text, the human or the repository rules
already authorize, also in another repository; ask only when the needed write
scope is not covered by them. Do not create worktrees or branches,
start agents, add or change Assignments or Tasks, or write extra report files
unless the Assignment text or the human asks. You may re-run step 1 to see the
current text.

## 5. Finish

When the work the Assignment text describes is complete, mark it done. This
moves it to the archive view; there is no separate archive step.

```sh
ply assignment done ASSIGNMENT-<n> --claim <claim> [--project <project-id>]
```

If you stop unfinished (blocked, out of scope, or the human says stop), give
it back instead and say why:

```sh
ply assignment release ASSIGNMENT-<n> --claim <claim> [--project <project-id>]
```

On an uncertain result, run step 1 before anything else.

Never revoke (`release --revoke` is a human recovery step), take a second
Assignment, or use a claim you did not get from your own take or adopt by your
own label.

## 6. Report

In your final answer, name the Assignment ID, its Project, its final state and
your claim, and each Task and Spec revision you read.
