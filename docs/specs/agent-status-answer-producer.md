# Produsent-Spec: nyttige svar fra eksisterende agenter

T661E488B · 2026-10-28 · Task `agent-status-answer-producer-spec-20261028`.
Status: spesifisert konsept; ingen produktimplementering eller aktivering i denne
Tasken. Administrator eier Ply/Herdr-avklaringer og senere byggkvalifisering;
app-root eier samlet verdiaksept og varsling. Leveransen er dokumentasjon lokalt
til Epic `ply-agentic-workflow-support`, ikke avslutning av Epicen.

## 1. Mål, autoritet og kontrakt

Et eksplisitt spørsmål i PWA skal nå den valgte eksisterende agenten og få et eget,
nyttig svar om oppdrag, fremdrift og neste steg eller hindring. Åpning, refresh av
lesedata, fokus, reconnect og timere skal ikke sende spørsmål. Arbeid, svart
statusspørsmål og transport er tre ulike fakta. Metadata og en disponibel prøveagent
oppfyller ikke brukerens verdiaksept.

Denne forfatter-Taskens native oppdrag og user-mandate styrer omfanget. Levert
server-Spec og client-reconciliation er avstemt felles kontrakt; øvrige dokumenter
er kildebelegg og konsumentkrav. Eksakte frosne dokumenter, revisjoner og
kildeankre er registrert i [kildegrunnlaget](agent-status-answer-producer/sources.md).
Henvisninger S/H/P nedenfor peker dit. Nåværende oppførsel er lest i kilde;
ingen live agent, HTTP-reise eller providerkø er prøvd i denne Tasken.

`agent-status-v1` beholdes: samme tre POST-ruter, closed DTO-er, enumverdier,
UTF-16-/bytegrenser, auth/principal/Origin/scope, idempotens, latch og recovery.
Ingen nye felt, typed blocked-kategorier, ruter, reason-parsing eller adapter.
PWA eier spørrehandling og svarvisning. Server eier HTTP-resultat innen 3000 ms
per kall; agentens svar kommer asynkront. Produsent eier brief, transportfakta,
korrelasjon og recheck-validering. Serverens fail-closed decode beholdes.

| Punkt | Valgt for senere implementering | Nåværende leveringsstatus |
| --- | --- | --- |
| P1 | Spørsmål først; faktisk aktiv kvalifisert absolutt callback-CLI | Ikke implementert |
| P2 | `fulfilled` betyr besvart statusspørsmål, også ved pågående/blokkert jobb | Native semantikk finnes; ny brief/veiledning gjenstår |
| P3 | Skrivingsvern for nye `summary/done/next/blocker` | Ikke implementert |
| P7(a) | Utelat privat valgfri `tab_label` i begge offentlige projeksjoner | Ikke implementert; ulik eksisterende filtrering |
| P4 | Behold faktisk full identitet; avgrens manglende identitet per provider | Åpent for berørte eksisterende sesjoner, se §5 |
| P6 | Behold `working_deferred` inntil trygg providerbane er kvalifisert | Åpent, se §6 |
| P5 / P7(b) | Utvidet re-ask / server-side label-dropp | Ikke valgt |

Ingen produktkode, tester, CLI-bygg/installasjon, global PATH-endring, runtime-
konfigurasjon, restart, Serve, push, PR eller main/master-integrasjon inngår her.
Fremtidige produktpakker nedenfor er en implementerbar plan under eget mandat.

## 2. P1 og P2: brief og mottakerens svar

Nå bruker `Service.Dispatch` en protokollorientert `fmt.Sprintf` med bare `ply`
til show/recheck (P1). Bytt den teksten og tilhørende veiledning, uten å endre
reservasjon, sending eller kontrolltilstander. Foreslått brief; plassholdere skal
fylles fra faktisk dispatch, aldri fra vennlig navn eller eksempelidentiteter:

> Statusspørsmål fra brukeren: Hva er ditt faktiske oppdrag nå, hva har du gjort,
> og hva er neste steg eller hindringen? Svar om oppgaven din. Fortsett innen ditt
> eksisterende mandat; spørsmålet gir ingen ny myndighet og starter ikke nytt arbeid.
>
> Les `{QUALIFIED_PLY} workflow followup show {ID} --format json` i ditt eksisterende
> workspace. Kontroller target_ref={TARGET}, generation={G}, cycle={C} og waiting.
> Ignorer meldingen hvis kontrollen er kansellert, erstattet eller gjelder en annen
> mottaker. Bruk fersk forventet revisjon, eier og mål fra show og en stabil request_id.
> Lever med `{QUALIFIED_PLY} workflow followup recheck {ID} --file <mutation.json> --format json`.
> `summary`: oppdrag og arbeid nå; `done`: siste konkrete resultat; `next`: neste
> steg; `blocker`: faktisk hindring eller «ingen». Hvert felt har høyst 256 UTF-16-
> kodeenheter; emoji kan telle to. Bruk oppgavekode og relative filstier, aldri
> private/absolutte stier, UUID-er eller hemmeligheter i svarfeltene.
>
> Et komplett svar bruker outcome=fulfilled selv om jobben fortsatt pågår eller er
> blokkert. Beskriv jobbtilstanden i feltene. Ikke velg waiting fordi jobben venter:
> waiting krever next_check_at_utc og åpner en ny planlagt syklus. Cancelled er ikke
> et svar. Dette ber aldri om å gjenta en avvist handling.

`QUALIFIED_PLY` er den faktisk aktive, administrator-kvalifiserte absolutte
binærstien. Samme fil brukes til show og recheck. Produsentens lokale oppsett skal
levere denne stien og dens kvalifiseringsbevis til briefbyggingen i `Dispatch`/
`System`; dette krever ingen ny offentlig DTO. Det konkrete lokale lagrings-
formatet velges i implementerings-Tasken. Manglende eller avvikende kvalifisert
verktøy må avklares før transportreservasjon/sending; ingen stille PATH-fallback.
Ingen ny rigid likhetsregel mot historisk `Runtime.PlyExecutable` velges.

K05 er bare kvalifisert verktøy for denne dokument-Taskens native prosess. Det har
ikke P3 og skal ikke hardkodes i fremtidige brief. Det senere bygget som faktisk
utfører recheck må inneholde den kvalifiserte P3-guarden. Korrekt ny producer som
instruerer mottakeren til å skrive gjennom gammel K05 oppfyller ikke P3.

Fremtidig bevis må bevare fysisk callbackbinær, absolutt sti og SHA-256, samsvar
med administratorens kvalifiseringskvittering, faktisk show/recheck-kommando,
exit/utfall og korrelert nyttig svar fra samme mottaker. Byggets P3-tester og
faktisk bruk må begge vises. Kjør også fra et eksisterende miljø der PATH `ply`
peker et annet sted. Manglende/endrede bytes gir et konkret avvik; en prompt alene
eller transport=`sent` beviser ikke callbackbruk. Bevis lagres privat, ikke i
svarfeltene. Ukjent callbackeffekt observeres før avhengig handling; ikke rearm,
ny request_id eller replay for å få et grønt resultat.

Native `Mutate` setter `fulfilled` som kontrolltilstand uten `advanceCycle`.
`waiting` lagrer recheck og kaller `advanceCycle`, som arkiverer syklusen og lager
ny `not_sent`-syklus. Denne semantikken endres ikke. Senere eksplisitt spørsmål
bruker dagens `canRequestAfterFulfilled` med eksakt current-cycle reporter; ingen
utvidelse etter unknown, cancelled eller ubesvart konsumert kontroll (P5).

## 3. P3: skrivingsvern uten omskriving

I `Service.Mutate`, recheck-grenen i [followup.go](../../internal/agentwork/followup.go),
valider `summary`, `done`, `next`, `blocker` med samme privacy-predikat som
serverens `privateText` (S1), i tillegg til dagens felt- og korrelasjonsgrenser.
Samle predikatet i en liten pakkeintern hjelper slik at P7(a) bruker identisk regel.
Nåværende `reportText` og separat summary-kontroll dekker bare UTF-16 og kontrolltegn.

Eksakt serverregel ved S1 (case-insensitive regulært uttrykk):

```text
(?i)(/Users/|/private/|/home/|/tmp/|~[/\\]|xox[baprs]-|xapp-|token=|password=|secret=|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})
```

Recheck beholder maks 256 UTF-16-kodeenheter per felt, ingen kontrolltegn og minst
ett ikke-tomt svarfelt samlet. Behold dagens outcome-, tid-, request-id-, eier-,
revisjons- og identitetssjekker. Et brudd skal gi vanlig korrigerbar
`followup_invalid_recheck` før ny `Recheck`, revisjon, historikk eller request-
resultat lagres. Ikke trunker, rediger eller sanitiser agentens tekst. Agenten
kan rette et avvist payload etter frisk show; en lagret request_id er fortsatt
bundet til identiske bytes og kan ikke brukes til å bytte et allerede lagret svar.

Guarden gjelder nye verdier i disse fire svarfeltene. Eksisterende identisk
idempotent resultat kan fortsatt være et gammelt svar; det må ikke gis en falsk
påstand om ny sikker validering. Ingen migrering eller sletting av historikk er valgt.

| Bane | Dekning etter P3/P7(a) | Restansvar |
| --- | --- | --- |
| Nye recheck-svarfelt | P3 avviser privacy-brudd før lagring | Mottaker retter teksten og leverer faktisk svar |
| Live `tab_label` | P7(a) i begge projeksjoner | Herdr-tabben beholdes |
| Gammel recheck/historikk/idempotent retur | Ingen reparasjon | Administrator + app-root avklarer konkret håndtering før aktivering |
| Work, capability, source, transport og annen fritekst | Ikke generelt dekket | Serverens gjentatte lesekontroll og fail-closed decode; vurder konkrete treff |

Et validert svar beviser ikke at hele overview er trygt. S1 leser alle strenger,
og ett privat historikkfelt kan fortsatt blokkere alle agenters HTTP-lesing.
S2 beskriver reproduksjon med fake producer, ikke livebevis utført i denne Tasken.

## 4. P7(a): begge labelprojeksjoner

Rik `Service.Overview` i [status.go](../../internal/agentwork/status.go) setter nå
`TabLabel` fra `tabLabels` uten metadata-only-filteret. `tabLabels` bruker
`validTabLabel` (Unicode, ikke blank, maks 256 runer og kontrolltegn).
`StatusMetadataContext` i [metadata.go](../../internal/agentwork/metadata.go) bruker
`validPublicTabLabel` med strengere plain-/UTF-16-/markørregler, men mangler hele
S1-regelen, blant annet UUID-er. Sessionløs metadata er lesbar, ikke dispatch-identitet.

Behold dagens innhentingsgrenser, rene labels og agentrader. Legg S1-predikatet
til den offentlige labelvalideringen og bruk den ved BÅDE rik og metadata-only
projeksjon før `tab_label` eksporteres. Hver bane beholder egne strengere
kontraktgrenser; felles minste offentlig labelregel skal oppfylle S1s 256 UTF-16
og kontrolltegnbegrensning. Ingen rett til å utvide tidligere avviste metadata-labels.
`validTabLabel` alene er utilstrekkelig offentlig filter. Ugyldig valgfri label
utelates (`omitempty`); behold agentidentitet, kapabiliteter og resten av raden.
Ikke erstatt med et påfunnet navn, endre Herdr-tabben eller serverdecode (P7(b)).

Bevis må dekke `~/github/repo`, `/tmp/build`, blandet case `/uSeRs/`, credential-
markører, UUID i ellers vanlig tekst, blank/control og >256 UTF-16 med emoji i
begge faktiske kommandoresultater. «Vanlig agentnavn» beholdes byteidentisk.
Test også metadata-only for sessionløs rad. At én bane er grønn beviser ikke den andre.

## 5. P4: faktisk identitet etter daemonoppdatering

Ply krever hele provider/session/workspace/tab/pane/terminal i `Identity.valid`.
`wireAgent.observation` fyller session bare fra Herdr `agent_session` med samme
provider og kind=`id`; `Ref` er hash av hele identiteten. `find`, `Get` og native
binding/veto bruker hele tuple. Verken label, pane alene, metadata-only eller et
Task/Run-ID er en erstatning (P2–P4 i kilderegisteret).

Herdr-kildene H1–H3 viser eksisterende reelle mekanismer:

| Provider | Faktisk kilde i Herdr | Hva den ikke beviser for eksisterende brukeragent |
| --- | --- | --- |
| Codex | Hook bruker faktisk `session_id`; SessionStart med transcript_path rapporterer `pane.report_agent_session`; working/idle-hendelser rapporterer agent med session. Arvet CODEX_THREAD_ID avvik avvises. | At en allerede aktiv sesjon uten metadata har fått en gyldig ny hendelse etter daemonoppdatering, eller at dagens installerte daemon tilsvarer lest kilde |
| Claude | SessionStart-hook rapporterer faktisk session_id, eventuell transcript_path og kilde; subagent/Cursor forkastes | At en eksisterende sesjon kan gjenpublisere identitet uten ny SessionStart, eller at daemonen har bevart bindingen |

`handle_pane_report_agent_session` kontrollerer nyere rapport og aktuell sesjon
før resume-informasjon beholdes. Det finnes en offentlig rapportmekanisme, men
den er ikke fullmakt til manuell rebind, fabrikkert rapport eller kunstig SessionStart.
Persisted launch-/resume-data finnes i Herdr-kilden; dette er ingen bevist
backfill-mekanisme for de berørte levende sesjonene.

**Beslutning:** Behold Ply-vernet. Støttet full `agent list/get`-identitet kan
brukes som i dag. For eksisterende agenter som fortsatt mangler kind=id er P4 et
presist åpent Herdr-gap, ikke løst av Specen. Eier administrator/Herdr må først
lese faktisk offentlig list/get og installert daemonproveniens for berørt daglig
Codex-orkestrator og aktuell Claude-sesjon, og følge ekte providerhendelse dersom
den finnes innen eksisterende mandat. Registrer før/etter hele tuple privat og
hvor bindingen tapes. Deretter avklar minst mulig støttet bevaring/gjenpublisering
fra faktisk providerlivssyklus med Herdr-eier under separat mandat.

Ingen agentmeldinger, restart, oppstart, rebind eller identitetsreparasjon er utført
her. Historisk S2-snapshot med 7 av 10 agenter uten komplett identitet er historikk,
ikke dagens liveinventar. Hvis legitim mekanisme ikke finnes, rapporter berørte
sesjoner som uferdig støtte med neste tiltak; disabled UI eller ny prøveagent
lukker verken D1 eller samlet verdiaksept.

## 6. P6: trygg mottakelse mens agenten arbeider

`Service.eligibility` returnerer nå `working_deferred` før sending for working;
ellers kreves idle/done, readiness, ingen launch pending og klar native veto.
`HerdrRuntime.Prompt` bruker offentlig `agent prompt` og godtar bare
`agent_prompted`. Dispatch lagrer unknown før én mulig sending (P1/P4).

Herdr H4 viser `queue_agent_prompt`: avviser blocked, ukjent/ikke-foreground
provider og launch pending, og køer tekst + Enter med 300 ms submit-delay til
terminalens `queue_user_input_submission`. Fullføring kvitterer terminalskriving.
Dette er ikke en bevist providersemantisk kø ved trygg turn-grense. Ingen callback
for agentens tankeprosess eller ny arbeidskø-API utledes av funksjonsnavnet.
Windows Codex paste-boundary-koden er ikke bevis for macOS-agentens working-kø.

| Provider | Støttet beslutning nå | Gjenstående bevis og neste tiltak |
| --- | --- | --- |
| Codex | Behold working-defer; lever én gang ved senere kvalifisert idle/done | Herdr-/Codex-eier må dokumentere faktisk støttet turn-grense og vise samme eksisterende arbeidsagent, langvarig verktøy, ingen interrupt eller Enter som godkjenner dialog, én mottakelse og korrelert nyttig svar |
| Claude | Samme vern; Herdr transport-ACK er ikke agentsvar | Herdr-/Claude-eier må kvalifisere providerens egen kø/mottak ved trygg grense, uten at statusspørsmål blir tillatelse eller ny jobb |

P6 er åpent per provider. Ikke fjern working-veto fra denne kildeobservasjonen.
Den nåværende reisen kan spørre mens agenten arbeider, vente uten ETA, og sende
når samme agent senere er kvalifisert idle/done; det beviser ikke direkte levering
under working. Kvalifisering må skille «arbeid pågår» fra terminalstatus og vise
brukernytte mens oppdraget fortsatt pågår. Påstand om working-køstøtte krever egen
faktisk kvalifisering, ikke bare denne fallbacken. Kvotestopp hos Claude gir
uferdig providerbevis, ikke lov til å restarte, probe kapasitet eller hevde Codex
beviser Claude. Administrator samordner gapet med Herdr-eier.

## 7. Korrelasjon og ukjent effekt bevares

Behold `owner_ref`, followup-id, expected_revision, generation, cycle, target_ref,
reporter og stabil request_id. Briefen sjekker cycle fra show; dagens Mutation
har ingen ny cycle-parameter. Revisjon/generation og fersk identitetsoppløsning
hindrer at sen callback blir svar fra en ny mottaker. Kansellert/erstattet mål
ignoreres før svararbeid. Test cancel/retarget mellom brief og callback.

Tre bevisnivåer skal rapporteres hver for seg:

1. Request er durabelt registrert; server gir meningsfullt utfall innen 3000 ms.
2. Herdr `sent` betyr kvittert eller køet input. Det er ikke mottakerens egen ACK.
3. Agentens recheck er egenrapport. «Besvart» krever fulfilled + eksakt korrelasjon
   med beholdt request/latch og aktuell rapport etter felles/server-Spec §4.3.
   `matching_report` alene definerer ikke alle svar eller fulfilled-recovery.

Dispatch reserverer unknown på disk før sideeffekt. Avbrudd etter claim eller tap
av ACK kan bety sendt; refresh, polling, daemonoppdatering eller ny brief skal
ikke sende på nytt. Ingen automatisk rearm/replay; bevar historikk og observer
faktisk effekt. Bare dagens eksplisitte request etter korrekt fulfilled kan åpne
ny syklus. Ikke bruk waiting som kvittering for «mottatt, jobben går videre»:
waiting lager faktisk ny syklus og kan sende et nytt spørsmål.

## 8. Minste implementeringsrekkefølge og senere akseptanse

Hver produktpakke krever eget faktisk mandat. Følg TDD for atferdsendring og
bevar relevante eksisterende tester. Ingen av testene i tabellen er kjørt her.

| Pakke | Minste berørte filer/symboler | Observerbar akseptanse før pakken kan kvalifiseres |
| --- | --- | --- |
| I1 P3 | `internal/agentwork/followup.go`: recheck/reportText og liten delt privacy-hjelper; `followup_test.go` | Avvis hvert S1-mønster i hvert av fire felt uten endret lagring/revisjon/request-cache; godta relativ sti/oppgavekode, 256 UTF-16, avvis 257/control. Rett avvist tekst og få ekte fulfilled. Behold konflikt- og idempotensvern. |
| I2 P7(a) | `status.go:Overview`, `metadata.go:validPublicTabLabel/StatusMetadataContext`, eventuelt delt hjelper; `tab_labels_test.go`, `metadata_test.go` | Begge offentlige baner utelater privat label, beholder ren label og rad. Sessionløs metadata beholdes. Ingen Herdr-tabskriving. Kjør grense-/tids-/prosessregresjoner for eksisterende innhenting. |
| I3 P1/P2 | `followup.go:Dispatch`, `native.go:System` for lokal kvalifisert callbackkilde; `docs/agent-work/v1.md`; `followup_test.go` | Spørsmål først; samme kvalifiserte absolutte verktøy i show/recheck; manglende verktøy oppdages før mulig send. Working/blokkert jobb kan besvares fulfilled uten advanceCycle/rearm; ingen produktmyndighet fra prompt. |
| I4 Produktkvalifisering | Administrator-kvalifisert bygg med I1–I3; server SW-3; samordning med PWA | Faktisk binærsti/hash i mottakerens kommando og P3-testbevis, faktisk callbackutfall og nyttig svar; kontroller historikk og øvrig tekst før pinning/aktivering under eget mandat. |
| I5 P4/P6 avklaringer | Herdr-eier + administrator; ingen forhåndsvalgt produktpatch | Eksisterende daglige sesjoner og providerbaner kvalifiseres særskilt etter §5–§6; eller presis uferdig støtte registreres. |

Produktakseptansereisen skal bevare privat bevis for:

- **E1, reell nytte:** velg eksisterende daglig working-agent (særlig berørt Codex-
  orkestrator), eksplisitt spør én gang, observer faktisk mottak ved støttet punkt,
  samme session/tuple og request/generation/cycle. Svaret navngir faktisk oppdrag,
  konkret fremdrift, neste steg/hindring. En probe, metadata eller working alene
  er utilstrekkelig. Hver provider som hevdes støttet må ha eget bevis.
- **E2, svart versus jobb:** komplett svar om fortsatt working/blokkert jobb bruker
  fulfilled. Ingen planlagt ny syklus eller utsending oppstår av svaret. Senere
  eksplisitt spørrehandling bruker dagens støttede request, ikke P5-utvidelse.
- **E3, faktisk P3-utfører:** mottakerens aktive kvalifiserte binær har guarden;
  privat svar avvises, rettes og leveres med faktisk callbackkommando/utfall.
  Feil PATH-verktøy blir ikke brukt. Unknown gir ingen replay.
- **E4, label og full lesing:** begge projeksjoner viser bevarte rader/rene labels
  og utelatte private labels. Skill svar- og labeltest. Les alle overview-strenger
  inkludert gammel historikk med serverens eksisterende §5.3-kontroll før og etter
  byggpinning/aktivering og ved senere source_unavailable. Null treff er et
  tidspunktbevis, ikke garanti for fremtidige fritekstendringer.
- **E5, livsløp:** cancel/retarget/identitetsbytte, duplikatklikk, transporttap og
  fulfilled-recovery beholder sannferdig korrelasjon, latch og no-replay. Gamle
  rapporter vises som gamle; sent svar fra gammelt mål blir ikke svar fra nytt.
- **E6, samordnet HTTP:** bevar maks responstid per kall, høyst 3000 ms, på faktisk
  server/PWA-reise. Pending/queued betyr ventende, ikke mottatt/svart. Ingen ny
  aktivering eller ingressendring autoriseres av denne Specen.

D1 (eksisterende identitet), D3 (nyttig prompt/callback), E1 (daglig agentsvar) og
D8 (tekst-/historikktrygghet) er fortsatt samlet uferdig verdi. P4/P6, gammel
historikk og annen fritekst er ikke waived eller fjernet fra pilotens verdiaksept.
En levert dokument-Spec eller lokal UI er fremdrift, ikke ferdig brukerreise.

## 9. Akseptanse for denne dokument-Tasken

T661E488B leverer bare denne Specen og kilderegisteret. Privat acceptance.sh
kontrollerer docs-only diff mot base, gyldige dokumentlenker, faktiske kildeankre
ved eksakte revisjoner og at alle avtalte beslutnings-/bevisavsnitt finnes.
Semantisk fullstendighet krever separat native review av ren committet kandidat;
ordtreff alene er ikke et produktbevis. Ingen dummy candidate-binary brukes.
Etter retting/review kvalifiseres eksakt kandidat source-only og integreres med
native callback til avtalt lokal Epic. Teknisk akseptanse, menneskelig vurdering
og faktisk levering rapporteres som separate fakta.
