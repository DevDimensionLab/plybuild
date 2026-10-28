# Kildegrunnlag for T661E488B

Lest 2026-10-28. Kildeutsagn gjelder de eksakte revisjonene nedenfor, ikke en
påstand om installert daemon eller live produktstøtte. Ingen aktiv agentmelding,
providerprobe, restart eller HTTP-reproduksjon inngår i denne forfatter-Tasken.

## Frosne dokumenter

Det normative oppdragsgrunnlaget er den native assignmenten (`spec`) og
`user-mandate`; server-Spec og client-reconciliation er avstemt felles kontrakt.
Klient-Spec/source-findings er konsumentkrav/belegg; prior-producer-assignment
bevarer informativ historikk. Alle ble lest fra frosne objekter i
`/Users/perottochristensen/github/ply_projects/.ply/task-content/v1/objects/sha256/`.
Filnavnet er SHA-256 uten prefiks. Frossen Goal-manifest
`a20dfa2d4e3316a43d72dd3d9672f21b33f2cdc7bd609952e6292f77025c2987`
binder dokumentene; execution-Spec er revisjon 2,
`ee6bd9cb9874e58fc23bea690668d79d4f60d223a10ad0e0cd71faeeec40a308`.

| Dokument | SHA-256 | Brukt til |
| --- | --- | --- |
| spec (native-spec-assignment.md) | `3e9cac6241b6a4f6cac95517612823ca0bf72107a6a0376d036f304675aba142` | Docs-only, norsk, P1/P2/P3/P7(a), P4/P6, aktiv callbackutfører |
| user-mandate | `d883d82ac40c5931f18040583d530d3432b908d3ce4872e0523a55f9d2f07191` | Daglig eksisterende agent, eksplisitt spørsmål, nyttig svar, 3000 ms, app-root-varsling |
| server-spec | `79bf2b3490c25b8d8d662efbd91e3b819ef69dc98371118f022176caeb089c55` | §4 fellessemantikk; §4.9 privacy; §5.1–§5.4 produsent; §8 A-1–A-11 |
| client-reconciliation | `6ff8dd85e9c58048b3d06ecfa0860cd93c56a74adff0ec8d210be2bd593bd150` | Uendret agent-status-v1, D1–D8 og lokal UI versus samlet verdi |
| client-spec | `7222077e6713d4ca00e973e811a0277d7ff78bde9480dfe4c5e61cc3033fff07` | Klientens latch/recovery, fulfilled-korrelasjon og must-pass daglig-agent-reise |
| source-findings (S2) | `671c6d9da77964b2fa99a922d6589f1b308a8f58b433290800009fa48557c1af` | Historisk snapshot og fake-producer H1–H4; ingen nye lokale livebevis |
| prior-producer-assignment | `6dddf4330c61e90380fb4051993017ac3dee6cc3199c5b20ad44e9fa678a0d4e` | Tidligere oppdrag, rettet Git-sperre og eksplisitt P3-callbackpresisering |

Server-Spec er levert ved server `8ee03a9797b4e17cb8d50f88a961abd9d2404714`.
Klientens lokale UI-levering er oppgitt i assignmenten som
`dd999e712bcf972c27aec0aee908b905cddbed75`; dette er oppdragskilde, ikke en ny
verifikasjon av klientleveransen her. Tidligere serverprobe brukte `5a85c24` og
disponibel agent; den beviser ikke vanlig agents nytte eller dagens liveinventar.

## Lest kilde og konkrete ankere

Ply-repoet er denne Taskens kilde ved
`8d606716b7e3c4f3a38099aec83826a8f0529331`. Server-repoet er
`/Users/perottochristensen/github/ply_projects/plyapp-server` ved
`8ee03a9797b4e17cb8d50f88a961abd9d2404714`. Herdr er lest i
`/Users/perottochristensen/github/herdr` ved
`5da0a01e1eedda054db0c81dd3a780000c40d9f0` (rent ved inspeksjon).
Dette angir kildeproveniens; Herdr-kilden er ikke likhetsbevis for kjørende daemon.
`git show REV:PATH` gir de presise funksjonene. Tabellen brukes også av privat
dokumentkontroll: hvert eksakte anker må finnes i den oppgitte filen/revisjonen.

| ID | Repo | Fil | Eksakt anker | Belegg |
| --- | --- | --- | --- | --- |
| P1 | ply | internal/agentwork/followup.go | `func (s *Service) Dispatch(id string)` | Linje 304: claim unknown før input; linje 354: fast prompt med bare ply; sent er Herdr-ACK |
| P2 | ply | internal/agentwork/types.go | `func (i Identity) valid() bool` | Linje 36: full identitet; Ref hasher hele tuple |
| P3 | ply | internal/agentwork/runtime.go | `func (w wireAgent) observation() Observation` | Linje 94: session bare fra samme provider og kind=id |
| P4 | ply | internal/agentwork/runtime.go | `func (h HerdrRuntime) Prompt(i Identity, message string) error` | Linje 129: agent prompt mot pane, bare agent_prompted godtas |
| P5 | ply | internal/agentwork/followup.go | `func (s *Service) Mutate(id, op string, m Mutation)` | Linje 187: request-digest, owner/revision/generation; recheck validerer tekst og full identitet |
| P6 | ply | internal/agentwork/followup.go | `func (s *Service) eligibility(r record)` | Linje 271: working_deferred, så readiness/launch/native veto |
| P7 | ply | internal/agentwork/followup.go | `func canRequestAfterFulfilled(f Followup) bool` | Linje 387: eksplisitt ny statusrequest etter korrelert fulfilled; advanceCycle arkiverer |
| P8 | ply | internal/agentwork/status.go | `func (s *Service) Overview() (Overview, error)` | Rik offentlig rad beholdes uten session; register sperres; label settes uten validPublicTabLabel |
| P9 | ply | internal/agentwork/tab_labels.go | `func validTabLabel(label string) bool` | Linje 195: Unicode/rune/control, ingen full privacy-regel |
| P10 | ply | internal/agentwork/metadata.go | `func validPublicTabLabel(label string) bool` | Linje 90: strengere metadata-filter, men ingen UUID-regel; StatusMetadataContext bruker det |
| P11 | ply | internal/agentwork/native.go | `func System(root string) (*Service, error)` | Linje 20: adapterkonstruksjon; readNative bygger full binding fra faktisk inventory |
| S1 | server | internal/agentwork/service.go | `var privateText = regexp.MustCompile` | Linje 370: eksakt regexp; Decode/textBounds linje 60/75 går rekursivt over alle strenger |
| S3 | server | internal/ply/agentwork.go | `case agentwork.Agents, agentwork.AgentsMetadata:` | Rik agents-bane, se faktisk kommandovalg nedenfor |
| H1 | herdr | src/integration/assets/codex/herdr-agent-state.sh | `inherited_session_id = os.environ.get("CODEX_THREAD_ID")` | Linje 64: reell session fra hook; avvis arvet avvik; SessionStart versus working/idle |
| H2 | herdr | src/integration/assets/claude/herdr-agent-state.sh | `if hook_event_name != "SessionStart":` | Linje 54: faktisk SessionStart/session_id; subagent/Cursor ekskluderes |
| H3 | herdr | src/app/api/panes.rs | `pub(super) fn handle_pane_report_agent_session(` | Linje 1599: nyere rapport, sesjonsvalidering og resume bundet til anvendt sesjon |
| H4 | herdr | src/app/api/agents.rs | `fn queue_agent_prompt(` | Linje 111: blocked/launch/foreground-vern; PTY tekst+Enter-kø, ingen bevist providerturn-kø |
| H5 | herdr | src/terminal/runtime.rs | `pub fn queue_user_input_submission(` | Linje 481: transportdelegering til PTY-actor |
| H6 | herdr | src/app/agents.rs | `agent_session: pane.agent_session,` | Linje 390: offentlig agent-info eksporterer observert pane-session |

S3s faktiske `Agents`-valg leses uten `--metadata-only`; denne banen er grunnen til
at kun metadata-filteret ikke beskytter serveren. P8 og P10 er to separate
projeksjonssteder. S1/S2 skal ikke tolkes som en ny global tilgjengelighetsgaranti.

## Verktøyproveniens for dokument-Tasken

Valgt native verktøy er absolutt K05
`/Users/perottochristensen/github/ply_projects/ply/ply-data/builds/k05/ply`, observert
SHA-256 `a8ad90d00710e3a90a1e753b2e9b766718ab05b22f227bf3e212ccd658859a26`.
Første faktiske check-in ble kjørt som PATH `ply` før Goal-kravet var lest; den
kvitterte og bevares uten replay. Det er ikke bevis for at check-in brukte K05.
Etterpå brukes eksplisitt kvalifisert K05, også til rapport/review/verify/integrate.
K05-proveniens er ikke bevis for et fremtidig P3-bygg; §2 i Specen krever det nye
aktive byggets faktiske path/hash og callbackutfall.
