# SL01 wire contract v2 — forslag til felles lås

Disse filene er designartefakter, ikke schemas levert av ordinær Ply. Ingen
eksempelidentitet er registrert, og ingen eksempelpath er opprettet. Kontrakten
gjelder bare den nye `week_codex_same_owner_turn_v1`, med `schema_version: 2`.
Den gamle planlagte MR04/R4 v1 og moden `week_codex_same_thread_v1` er ikke
kompatible aliaser. Capability-ID-ene kan gjenbrukes, men må annonsere akkurat
støttet schema/mode; appen krever begge før en effekt.

[Schema](workflow-bound-answer-v2.schema.json) har fire lukkede røtter:
`question_read`, `preview`, `request` og `status` under `$defs`.
[Eksempler](examples.json) inneholder gyldige, schema-ugyldige, rå ugyldige og
semantisk ugyldige verdier. Alle objektfelter er påkrevde. Ukjent felt på
hvilket som helst objektnivå avvises; fravær og `null` er ulike verdier.

## Faste operasjoner og grenser

```text
ply workflow execute question RUN_ID --workspace ABS --format json
ply workflow execute answer preview RUN_ID --workspace ABS --question-id Q
  --question-revision 1 --option-id OPTION --request-id ID --format json
ply workflow execute answer apply --file ABS_APPROVED_REQUEST --format json
ply workflow execute answer status ID --workspace ABS --format json
```

Read/preview/status har ingen effekter. Apply leverer `status` med samme
requestidentitet; det er ikke en ny resultattype. Semantisk avvisning bruker
denne JSON-responsen med `bound_answer_*`-årsak og nonzero exit. Argument-/fil-
eller parserfeil som hindrer en troverdig requestidentitet gir stabil
`bound_answer_invalid_arguments`/`bound_answer_invalid_request` på stderr og
ingen oppdiktet status på stdout. En exitkode beviser aldri fravær av eldre
effekt. Dette er en avgrenset answer-feilkontrakt, ikke et globalt error-API.

Maksimal request er 64 KiB UTF-8. Maksimal enkeltrespons er 256 KiB.
Avvis BOM, ugyldig UTF-8, enslige surrogater, duplikatnøkler på alle nivåer,
flere JSON-verdier, flyttall/eksponenter/negativ null og ukjente felter før
schema/semantisk validering. Integerfelter bruker JSON-heltall i intervallet
angitt i schema; bool kan ikke tolkes som integer. Strenger er data, aldri
shell, og normaliseres ikke etter godkjenning.

`workspace` er et absolutt, fysisk og normalisert workspace, ikke Task-worktree
eller cwd-fallback. Native validerer ingen symlinks eller `.`/`..`-segmenter.
`scope.task_worktree` er separat. Lokal appregistrering binder workspace,
Run/request-digest, prosjekt/repo/Epic/Task/Goal og parent-week/workset samt
kvalifisert binær/schema. Den leser status fra dette nye spørsmålstargetet;
gammel control-journal gir ingen svarmyndighet.

Question-objektets `continuation.dispatch_policy` er den uforanderlige
`subject_to_parent_control`. Parent-pause endrer ikke spørsmålets innhold eller
digest. Den aktuelle previewen velger og binder i requestet enten
`dispatch_policy=when_parent_allows` eller
`dispatch_policy=deferred_until_explicit_parent_resume`. Det siste kreves når
preview uttrykkelig godtar et svar under parent-pause. Parentendring før første
apply gjør requestet stale; parentendring etter lagret svar bevarer svaret og
vurderes foran dispatch.

## Tid og digester

Alle tidspunkt er UTC RFC3339 med `Z` og høyst ni desimaler. Question/readback
viser `original_deadline_utc`; `parent_week.original_deadline_utc` er den
opprinnelige parent-fristen. Den effektive originalfristen er tidligste
allerede bundne Run/worker/parent-frist. Den kan ikke flyttes ved preview,
restart, svar eller nytt prosessforsøk.

Requestet har både originalfristen og
`validity.{issued_at_utc,expires_at_utc}`. Preview velger
`expires_at = min(issued_at + 300s, original_deadline)`. Ingen `ready` request
produseres hvis vinduet er tomt. Første apply krever
`issued_at <= now < expires_at <= original_deadline` og høyst 300 sekunders
gyldighet. Tidsstemplene inngår i bindingen. `observed_at_utc` er bare
observasjonstid og gir ingen gyldighet. Gyldigheten gjelder ny apply, ikke
senere lesing av samme allerede registrerte request. Samme-ID replay er
lesende selv om preview har utløpt.

`request_sha256` er `sha256:` pluss lowercase SHA-256 av hele requestet i
kanonisk JSON, inkludert binding, validity og mode. Bruk d5s
`internal/canonicaljson`-domene: RFC8785-strenger/sortering uten whitespace,
ingen BOM eller slutt-newline, og bare de ovennevnte heltallene. Hash aldri
den pretty-printede filen i stedet. Appen bevarer originale requestbytes og
beregnet digest; endret kropp under samme ID gir konflikt. JSON-keyrekkefølge
og ubetydelig whitespace påvirker ikke kanonisk requestidentitet.

`question_sha256` dekker hele question-objektet uten selve
`question_sha256`-feltet, i samme kanoniske form. Alternativenes rekkefølge
bevares. `binding` er opaque for appen; produsenten binder requestet uten
binding sammen med komplett autoritativ basis i domenet
`ply.workflow.bound-answer/v2`. Basis omfatter også private session-/policy-
og Git-identiteter uten å publisere innholdet. Observasjonstid og en avsluttet
PID er ikke alene autoritetsrevisjoner. Syntetisk binding i eksempelpakken
attesterer ingen faktisk basis.

## Nuller og effektbevis

En `ready` question/preview har komplett scope, parent, worker og aktivt
spørsmål. `owner_state=confirmed` er en datert låse-/eierobservasjon, ikke
provider-liveness. Ikke-støttet/ukjent grunnlag kan ha avtalte nullable
objekter. V1-revisjon av spørsmål er alltid 1; endret innhold krever ny ID.
Schema tillater 2–5 alternativer, mens SL01-targetet binder akkurat to og
én continuation; konsumenten avviser en uventet ekstra valgmulighet.
Stoppet parent og brukt continuation-grense gir aldri ready. Ready-preview
krever både current basis og current top-level freshness.

`not_recorded` kjenner bare workspace, forespurt request-ID og fast action fra
kallet. Digest/binding/scope/Run/worker/spørsmål/svar og begge effektbooler er
`null`, `execution_effect=unknown`. Det betyr ingen lesbar request, ikke
null effekt. `outcome_unknown` kan beholde kjent svaridentitet og positivt
lagret `answer_recorded=true`, men uobserverbare deler er `null`.

`answer_recorded` betyr varig vinner. `continuation_requested=true` begynner
først ved varig reservasjon av ett continuation-forsøk. `rejected` betyr at
akkurat denne invokasjonen ikke publiserte en ny intent; feltverdiene er ikke
en attest om en tidligere invokasjon. Native må først forsøke eksakt oppslag
av samme identitet, slik at kjente/ukjente tidligere effekter ikke blir rejected.

Pause etter reservasjon, før launch, beholder `continuation_reserved`,
`continuation_requested=true` og samme job-/intentidentitet; aktuell pause
forklares i `reasons` og evidens. Ikke nullstill dette til `continuation_deferred`.
Sistnevnte brukes bare før reservasjon. Etter en lovlig separat parent-resume
kan samme supervisor starte det samme reserverte forsøket høyst én gang.

I status er `scope`, `parent_week` og `worker` den bevarte, godkjente basisen,
ikke løpende eierskaps-/kontrollavlesning. Derfor er eksempelverdien
`worker.continuations_dispatched=0` fortsatt riktig etter faktisk start.
Nåværende kjente sperrer står i `reasons`/evidens; faktisk continuation ligger
i receipt. Question/preview inneholder derimot sin ferskt observerte basis.
Ingen historisk `owner_state=confirmed` blir en påstand om dagens levende eier.

`receipt` har flat, lukket form. Ved `continuation_observed` er faktisk nytt
prosessforsøk og samme session observert; alle consumed-/sluttresultatfelter
er fortsatt `null`. Ved `continuation_completed` er disse obligatoriske:

- `receipt.continuation_session_id` lik `receipt.source_session_id` og bundet
  `worker.provider_session_id`;
- `receipt.consumed_answer_request_id` lik øverste `request_id` og appens ID;
- `receipt.consumed_answer_request_sha256` lik øverste `request_sha256` og
  appens kanoniske digest;
- konkret `result_sha256`, sluttid og positivt validert terminalresultat fra
  dette eksakte continuation-forsøket.

Alle øvrige gjentatte receipt-/request-/question-/workspace-/Run-identiteter
må samsvare. JSON Schema uttrykker form/nuller; native og app må også teste
disse feltlikhetene, ordnede tidspunkter (inkludert at readback-observasjonen
ikke er før det påståtte sluttresultatet), digest, unike option-ID-er og at valgt
option finnes i det bundne spørsmålet. Et mismatch er unknown, aldri grunn til
ny dispatch. `evidence` inneholder ID/type/digest/tid, ingen privat locator som
konsumenten skal åpne.

Consumed-feltene viser hva samme avsluttede modellrunde rapporterte å ha brukt.
De beviser ikke dokumentets riktighet, global writer-ro eller integrasjon.
Påstand om levert valgt innhold krever kontroll av nøyaktig testet/reviewet/
integrert C/T. App-main attesteres separat av app-eierens faktiske kvittering;
ingen ny main-API inngår. `freshness=current` betyr fersk lesing av bevarte
bevis, ikke levende provider eller rett til ny effekt.

## Venting og beholdt myndighet

Den nye moden bruker positivt avsluttet providerturn og komplett validert
stream med samme Task/eier/reservasjon/session. `completed_model_turn_same_owner`
er uttrykkelig ikke global `paused` eller bevis på at alle tool-writere er borte.
Uklar source-turn kan ikke produsere et `ready` spørsmål. Kontroll-/svarlager
er beskyttet mot provider-skriving; identitets-/fildrift revalideres før apply
og launch. Parent pause/stopp/frist styrer fortsatt ny dispatch.

Restart alene gir ingen ny skrivemyndighet. Eventuell videreføring etter
restart må bindes til positivt bekreftet native eierskap og samme bevarte
reservasjon; ukjent tidligere forsøk sperrer. Ingen Task frigjøres fordi den
venter på svar. Sen tool-effekt kan ikke omtolkes som nytt svar eller brukes
til å attestere det endelige leverte innholdet. Produktakseptanse, review,
integrasjon og kildeavvikling beholder egne kontroller.
