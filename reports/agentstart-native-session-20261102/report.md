# Rett observert Codex native-session-sperre

Task `ply-agentstart-native-session-20261102`, Goal `agentstart-native-session-v1` rev. 1.

## Kilde

| | commit | tree |
|---|---|---|
| Sluttkilde for review | `5d11d9bd7dcf708e5f89b2aa2588f75b94739920` | `f8133ac68529683711600f91e50722bb2881f78e` |
| Bevart delta (b23fed1..58033ae, blobber kontrollert like for alle 8 filer) | `b27df89` | |

Rapporten ligger i en egen commit oppå sluttkilden. Review `b27df89..5d11d9b`.

## Funn

Native ID er ikke nødvendig for dette startløpet:

- Prompten adresseres til pane, ikke til native ID. Goal-start har ingen native report-callbacks.
- `Session` leses ikke av noen utenfor `internal/agentstart`; `cmd` bruker bare receipten som output.
- Ifølge oppdraget tømmes Herdr-navnet når agenten erstattes. Navn, provider, pane, tab, terminal,
  workspace og cwd kontrolleres ved hver observasjon. Navnetømmingen er ikke testet mot ekte Herdr her.
- Goal-start på b23fed1 krevde ikke native ID.

Bootstrap-prompten er også overflødig. Eneste formål var å få provideren til å opprette en sesjon.
Smoke viser at den ikke tilfører noe vern:

- Codex tok imot `agent prompt ... --wait --until working` i samme tilstand (settled, interactive_ready,
  ingen sesjon). Oppdragsprompten bruker samme argv, bare med annen tekst.
- Claude-receipten har ikke `bootstrap`-felt, så bootstrap ble ikke sendt der.
- Å beholde bootstrap ville gitt en ekstra provider-tur for hver Codex-start. Ved tapt ack kunne
  oppdragsprompten også truffet en agent som fortsatt behandlet bootstrap.

## Valg

Uten observert native sesjon sendes oppdragsprompten når agenten er settled, `interactive_ready=true`
og rapporterer `foreground_cwd`. Med sesjon er kravet uendret.

| Vern | Status |
|---|---|
| Fullstendig Herdr-tuple, navn, provider | Uendret, nå testet både med og uten sesjon |
| Kjent feil cwd | Uendret, nå testet også for Codex uten sesjon |
| Manglende Codex-cwd | **Strammet inn uten sesjon:** gir `not_ready` og ingen prompt. Med sesjon fortsatt akseptert som før |
| Native konflikt (feil kind/provider, endret ID) | Uendret |
| Observert native sesjon som forsvinner | **Nytt:** `agent_response` og ingen prompt. Ville ellers åpnet en ny sesjonsløs vei |
| Én prompt per nøkkel, pending før effekt, ingen retry | Uendret, nå testet uten sesjon (tapt svar og krasjpunkter) |
| Ukjent forblir ukjent | `Session` og `session_bound` skrives bare når sesjon faktisk er observert |

Fjernet: bootstrap-prompten, `Receipt.Bootstrap` og tilhørende konstanter.
Omdøpt fordi de gamle navnene ble usanne: `await_session` til `await_ready`, `session_unbound` til `not_ready`.
Ingen bypass-flagg, falske IDer eller ny motor. `inputDigest` er uendret.

## Rød og grønn

- Rød før produktretting: `red.log`. 5 tester feilet (`TestAgentStartPromptsSettledCodexWithoutNativeSession`
  med agentformen fra `codex-final-agent.json`, cwd ikke rapportert, forsvunnet sesjon, tapt svar uten sesjon,
  krasjpunkter uten sesjon). Hovedfeilen var den observerte: `session_unbound` etter deadline.
- Grønn på `5d11d9b`-innhold: `focused.log` og `race.log` for `internal/agentstart`, 125 PASS og 0 FAIL hver.
- `cmd-goal.log` og `cmd-goal-race.log` (`go test ./cmd/ -run Goal`, også med `-race`) er grønne.
- `go build .`, `go vet ./internal/agentstart`, `git diff --check` og `gofmt -l` er rene. Bygget lå i scratchpad, ikke installert.
- Testen for receipt fra bootstrap-kandidaten ble lagt til etter rettingen og beskriver eksisterende lesing. Den har ingen rød fase.
- Fullsuite er ikke kjørt, som avtalt.

## Netto endring (b27df89..5d11d9b)

| Fil | + | - |
|---|---|---|
| `internal/agentstart/agentstart.go` | 29 | 54 |
| `internal/agentstart/herdr.go` | 0 | 1 |
| `internal/agentstart/README.md` | 18 | 19 |
| `docs/goal-agent-start.md` | 1 | 1 |
| `internal/agentstart/agentstart_test.go` | 125 | 117 |

Utenom tester: +48/-75.

## Begrensninger

- Ingen ekte Herdr- eller providerkjøring. Dette er testdobler, og koordinator eier ny Claude/Codex-kvalifisering og fresh review.
- Lang oppdragsprompt til Codex 0.164.0 uten bootstrap er ikke observert med denne kilden. Den må dekkes av kvalifiseringen.
- Claude uten rapportert sesjon får nå prompten direkte i stedet for bootstrap. Claude-smoke hadde sesjon, så denne grenen er bare testet med double.
- Receipt-lesing av den bevarte smoke-kvitteringen: `bootstrap`-feltet ignoreres ved lesing og vises ikke lenger i output.
  Filen endres ikke, `failure.effect_unknown=true` består og samme nøkkel avvises. Klassen `session_unbound` er fortsatt lesbar som data.
- Forhåndseksisterende og uendret: med sesjon kreves bare settled, ikke `interactive_ready`.
- Etter at prompten er sendt observeres ikke identitet videre.
- Modell, profil og sandbox for denne sesjonen er ikke attestert. Ingen TaskResult, QA, integrasjon, push, install eller varsel er utført.
