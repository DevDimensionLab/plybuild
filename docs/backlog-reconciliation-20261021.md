# Backloggavstemming 21.10.2026

Avstemt mot frosne 86 `koe/*.md` (README utelatt), tatt 19:04:46Z,
og Epic `d5c4d6aa5482d577c37c3ec71e7f604a22cc7410`, tree
`e8c7774e58bcf97360d95891b96db8061fdbc57d`. Dette er et datert bevisbilde,
ikke en ny kø eller vedlikeholdt Task-database. Tidligere whole-koe-mandat
består; historiske «ubestilt»/«ikke valgt» er ikke nye godkjenningsporter.

**Filregnskap: 35 C (dekket/avsluttet), 48 R (rest/delvis), 3 U (spesifikt uverifisert) = 86.**
R og U teller også overlapp med uferdig arbeid; ingen uferdig gruppe er
trukket fra som levert. C inkluderer historisk fysisk avslutning (rad12) og
ferdig design (rad82) og historisk runtimebinding (rad9), ikke bare produktkode.
30 grupper G01–G30 beskriver unike rester/avklaringer. Dette er
avgrensningsenheter, ikke 30 ferdigdefinerte Tasks eller et timeestimat.
G24 er én batchbar restmatrise fordelt på seks eksisterende delsystemkilder;
G28 omfatter registrering av **to** problemer. Produktarbeid, evidensavklaring,
appansvar og ekstern/human aksept er derfor ikke én ensartet mengde arbeid.

Full SHA256 i hver rad gjelder originalfilens UTF-8-bytes i det frosne
[snapshotet](/Users/perottochristensen/github/.ply/task-content/v1/objects/sha256/1d26b07d4d278f36d7f06de0b1a5eddb8d6aa9ca0b27ff2feb62f0c5e0a014ce). Snapshotets egen SHA256 er
`1d26b07d4d278f36d7f06de0b1a5eddb8d6aa9ca0b27ff2feb62f0c5e0a014ce`.
Filnavnene nedenfor er eksakte identiteter; relativ `koe/` finnes i planrepoet,
ikke dette kodecheckoutet.

## Bevis og hva som faktisk er kontrollert

Native `task list`, 41 `task show` og tilhørende 33 `execute show`-readbacks
ble hentet skrivebeskyttet med installert AL21-CLI. Antall og rå streams/exit
står i [privat bevisrot](/Users/perottochristensen/github/.ply/task-executions/ccfecb1e990f1016d369d7e0bdbf1da6fcaec652e06190d584eea3f5e8ed5ef7/al22-evidence) og
[native indeks](/Users/perottochristensen/github/.ply/task-executions/ccfecb1e990f1016d369d7e0bdbf1da6fcaec652e06190d584eea3f5e8ed5ef7/al22-evidence/native-index.json). Native Delivery-listen var tom for
IL01, AL18 og AC02: execute-leveransene bevises derfor av Runens kandidat,
integrasjonskvittering og Task-closeout, ikke en oppdiktet `dlv_`-registrering.
Fullførte historiske readbacks kan være stale etter kildesletting; det
opphever ikke en varig leveringskvittering.

AL01–AL12: les [tidligere audit](approval-requirements-audit.md) som datert
historikk, ikke dagens restliste. AL03 trust, AL04 kontinuitet, AL05 input,
AL06B rapportmodus, AL07 aktør, AL08 install, AL09 follow, AL10 uferdig rapport
og AL11 launch-kontrakt har native resultat/lukking. AL12 leverte dokumentasjon
på `c91e188d`; installert AL11-binær ble beholdt. AL13–AL21 lukker flere av
AL12s daværende restpåstander; AL22 endrer ikke den historiske auditteksten.

| Bevis | Eksakt levert kandidat / Run og nåværende kildeegenskap |
| --- | --- |
| AL01 | `d335e70e0c16c8acab34d30315207fe35078bc76`; `wfr_397376b22e6fe550e5efec3c46b728881c9c8049ae5e5cd7ba7c9029370b05cf`; `trs_4a470545d376a381d2b37d75102424c0`; native closeout `complete`. [report_startup.go](../internal/taskrun/report_startup.go) |
| AL02 | `922961b702156908786f41cd9b8e682aaf2d5242`; `wfr_3686f4b27a27f281ab3fc42d0b6490374455e97a15244e1be1e6d11c6b52391e`; `trs_1708e7b3785342b72d579725febf50fc`; native closeout `complete`. [approval-requirements-audit.md](../internal/taskrun/../../docs/approval-requirements-audit.md) |
| AL03 | `f5535ef88b8bf2e2b75720285234c2fa71070a34`; `wfr_cc51d1cff84e982ec3b1c2071dce37d370d0b9153b368c60f366fec52978d2c3`; `trs_60e89de2880a043cd8a51988a4527cdb`; native closeout `complete`. [herdr_codex_trust.go](../internal/taskrun/herdr_codex_trust.go) |
| AL04 | `7e049edfc69057b6f15e1187b2740cb47861f4bc`; `wfr_d98c20b634d254982a6abfc2c91f8385a96beeb2d550c4d1db3525516c29a1f7`; `trs_9b4617ae6ef14e79b90301e25b3b1d60`; native closeout `complete`. [herdr_store.go](../internal/taskrun/herdr_store.go) |
| AL05 | `c02e3ee474fb94da61a4d375f6b0d437c1c5c4df`; `wfr_09b3c4d6e24885dde5fec7f450f39a3ab8fc53dabdf3ca945fe8cf5518447fe5`; `trs_caa1316c676ebcff53a4b027035d4496`; native closeout `complete`. [schema.go](../internal/taskrun/schema.go) |
| AL06B | `104098764fb2f96cbdc24fb7d142938c111faae3`; `wfr_989d78a8f7b45c3aab474955dd42662a6a8622adb8c29071ed178267d54d0a83`; `trs_8fcbf89789999f7122789284732e43d5`; native closeout `complete`. [report_claude_mode_test.go](../internal/taskrun/report_claude_mode_test.go) |
| AL07 | `7fde35bc9855d2064da34e31dccef6b174f2b3e0`; `wfr_2e4c21adc24a27498b736bc99127516d3bd3eb2b8b8375e5aecdb5089ff88e67`; `trs_0d7d8a9c0ef264e4a81a273112fab137`; native closeout `complete`. [read_closeout.go](../internal/taskrun/read_closeout.go) |
| AL08 | `2941fcc9b0ca67fe218ea59e4d924d8d517460dc`; `wfr_74d202fb56a2f259e9c3c73e02d2859ea554bd3276566e5ae3afb46dfdec0417`; `trs_9809fb7edbc14f60145ef36a5047d0b2`; native closeout `complete`. [delivery_cli_install.go](../internal/taskrun/delivery_cli_install.go) |
| AL09 | `3e09cab01e430ee93b843eb5bc7fc543aac3fa2e`; `wfr_e5395f2fae5045a47f07389b24c41a7249cc9ca25dd959071d7cde0d3835176f`; `trs_0abd46f6aa161e91c2352ea87d38d85f`; native closeout `complete`. [follow_observation.go](../internal/taskrun/follow_observation.go) |
| AL10 | `4f82278923a3fff4895187a9a367ab2de8e584c3`; `wfr_73bf0e4fe9af269d260bc3cad5716f594b56eeca270f272e5da0c99d6aea4e7e`; `trs_cc0f18454e706dbe829f95d2a6b94148`; native closeout `complete`. [report_incomplete_test.go](../internal/taskrun/report_incomplete_test.go) |
| AL11 | `c6477d6094628d74eb6a1d04b41207174d9f93ea`; `wfr_66f0f2cb1b804352897824dc3364452a8e598846d810f4af5cdb6488bb81d75f`; `trs_ca562895f079c35def899b107d7b0bb6`; native closeout `complete`. [report_launch_contract_test.go](../internal/taskrun/report_launch_contract_test.go) |
| AL12 | `c91e188da203eda4adb8307d49c68bb7bb95a1fa`; `wfr_0ec60ba7192443e92320660ee925d1a187c834f78b17d6a0e8be5cfaaf13dc93`; `trs_343d255179a554d0e88235def66986bf`; native closeout `complete`. [approval-requirements-audit.md](../internal/taskrun/../../docs/approval-requirements-audit.md) |
| AL13 | `2a2fa2c0070d4084ffc000caaaa318dfe6a9ab4d`; `wfr_2f041770d51eec8677d3fa159d953e7fbebc6fb19f9d6be7a8cb7bd335a29653`; `trs_2fd5264b224c653a8bc50b064eb5088c`; native closeout `complete`. [delivery_review_prose_test.go](../internal/taskrun/delivery_review_prose_test.go) |
| AL14 | `5266d8db6f7c7c3a43c551304c2ad25fa908ca19`; `wfr_56a81c6d8663d8d368614beae4c7e9768ee35f51fc3c77e12f9d433566fdbecb`; `trs_704af9d34a7d33f956ae04d78eabec7e`; native closeout `complete`. [report_retirement_poststart_test.go](../internal/taskrun/report_retirement_poststart_test.go) |
| AL15 | `09708a92947effd39521d64ec96ec1cf9cbe4392`; `wfr_b51d08a1a146710e2d84ef17aa0e6628a3e460a647d3597c5839d520e44956c4`; `trs_989baf30dce064222138e622ee6c3bc6`; native closeout `complete`. [delivery_review_mandate_test.go](../internal/taskrun/delivery_review_mandate_test.go) |
| AL16 | `203d2e5c9dddc6671806e87f09a4c84ccd6d04da`; `wfr_9faee9c5b04920be2753439d07b3aa9046929379abb21e99635fef02287b85d4`; `trs_ad9b319bc8f42afb135c7863fcd00939`; native closeout `complete`. [Makefile](../internal/taskrun/../../Makefile) |
| AL17 | `63a2d0031aa1cfe08088cb60a9eec759d6cb7f96`; `wfr_45153d44658c57014ac328686bd03c053c3415c9ed3ba7cc8c34ba2329c993ce`; `trs_47984f1aebf79f15974f8948364db185`; native closeout `complete`. [delivery_review_mandate_test.go](../internal/taskrun/delivery_review_mandate_test.go) |
| AL18 | `bc7893b524581c5811ff8d44631a8bc71345a3d4`; `wfr_76f29f9ca42c5cd0742f0a960a2ed1cb259da42b957fa774c3093064ed27304d`; `trs_e45794a3a1c198d3a309f70a14cb8415`; native closeout `complete`. [delivery_review_round_test.go](../internal/taskrun/delivery_review_round_test.go) |
| AL19 | `d9a43fa2d6d0876e39d56e3c51b2fca281218ad9`; `wfr_f68413f19dce4525d96a893e3513874a7bd463c2470ff428a5715967faa8a3e3`; `trs_c54843796e564b6c13bd2aeabf49e18e`; native closeout `complete`. [delivery_startup_cli_test.go](../internal/taskrun/delivery_startup_cli_test.go) |
| AL20 | `f32c2c1628b466612e8d7e5b6f2f24ffe358930a`; `wfr_4fb37f20e6fd6eac01bdbc36d1344786b149b8f244585ab878eca7a10cb2ec65`; `trs_0c0d0f457bb6b8b2a27789a4dfdf1a94`; native closeout `complete`. [absence_freshness_test.go](../internal/taskrun/absence_freshness_test.go) |
| AL21 | `d5c4d6aa5482d577c37c3ec71e7f604a22cc7410`; `wfr_daf3f23f22a30a8251fbfc0fc42b2110ee9961b80f7c28da38f08cfd2236089f`; `trs_957f25f0eaf5fd5d6fc5522d76113588`; native closeout `complete`. [al21_closeout_observation_test.go](../internal/taskrun/al21_closeout_observation_test.go) |

Andre beviskoder i filtabellen er avgrensede kilder som faktisk er lest:

- E-native: [native-index.json](/Users/perottochristensen/github/.ply/task-executions/ccfecb1e990f1016d369d7e0bdbf1da6fcaec652e06190d584eea3f5e8ed5ef7/al22-evidence/native-index.json).
- E-Claude: [menneskeprøvd trust og egne classifier-begrensninger](/Users/perottochristensen/github/ply/planning/design/evidence/213-claude-trust-integration/report.json).
- E-closure: [faktisk fysisk avslutning](/Users/perottochristensen/.local/share/ply/planning/rest-closure-20261016-01a1431a/native-result.json).
- E-first: [første integrasjon og faktisk human QA](/Users/perottochristensen/github/ply/planning/design/evidence/201-first-task-integration/control.json).
- E-WF: [akseptert WF01-resultat](/Users/perottochristensen/github/ply/planning/design/39-ply-workflow-handoff-integrasjonsresultat-og-neste-produktvalg.md).
- E-WS04: [kontrollert WS04-resultat og modefeil](/Users/perottochristensen/github/ply/planning/design/57-ws04-lokal-integrasjon-fullfort-og-modekontroll.md).
- E-Slack: [kontrollert Slack-integrasjon](/Users/perottochristensen/github/ply/planning/design/evidence/77-workflow-notification-local-integration-result-control-v1.json).
- E-app: [W1–W9-resultat med egne QA-begrensninger](/Users/perottochristensen/github/ply/planning/design/evidence/211-plyapp-read-support/report.json).
- E-trace: [korrigert visualisering og faktisk human pass](/Users/perottochristensen/github/ply/planning/design/evidence/workflow-visualization-execution-2026-10-07/report-v2.json).
- E-design: [ferdig én-node-design](/Users/perottochristensen/github/ply/planning/design/06-feature-utforing-tilstandsmaskin.md).

E-ark er kildeinspeksjon av `cmd/workspace{,_epic,_project,_task}.go`,
`cmd/workflow.go`, `internal/workspace/{work_item,epic_base,task_spec_guard}.go`,
`internal/taskrun` og `internal/workflowhandoff`. Ingen
`internal/workflow/model/catalog/fsstore` eller WorkflowDefinition/Profile-kjerne
ble funnet. AC01/AC02 native Tasks er archived/unbound uten resultater.
Arkivref/ancestry alene brukes ikke som leveringsbevis. Task-kommandoene
registrerer problem, løsning, valg, forberedelse og resultat; Epic `adopt`
er ikke design79s `create`, Project `add` er ikke `init`/medlemsmutasjon.
Workspace/workflow-grensen støttes av dagens separate kommandoer og pakker.
E-WS03 kombinerer dette med [akseptert design40](/Users/perottochristensen/github/ply/planning/design/40-workspace-eid-epic-task-og-task-worktree.md)
(`design_ready_for_spec`, skiller workspace fra utføring) og bootstrap-Taskens
archived native readback; rad70 erklærer bare denne strukturelle grensen dekket,
ikke en ny native produktleveranse.

E-test er lest eksisterende kode: `test/cli_surface_contract_test.sh`
(unconditional cleanup på EXIT), `internal/taskrun/schema.go` (`readFileAs`,
SameFile), `delivery_review_round_test.go`, `worker*.go` og eksisterende
workspace-/handoff-faulttester. E-schema er offentlig
`schemas/read-contract/ply-capabilities-v1.schema.json` og bevart SL01-pass (se kontroll under).
E-PR er lest `internal/delivery` og dagens delivery/closeout-kontrakter:
PR-registrering/eksekvering finnes, men det gamle PR103-gapet er ikke derved
bevist lukket. Ingen GitHub-operasjon ble gjort.

E-provider/E-pilot/E-human/E-mobile/E-notify/E-Trip bygger på de konkrete
frosne filenes hendelsesbevis og readbacks, supplert med inspeksjon av
`internal/taskrun/{report_startup,delivery_retire,delivery_continuation}.go`,
`internal/workflownotification`, `internal/humaninput`, `internal/weekrun`
og `cmd/workflow_execute*.go`. De gir **ikke** nye provider-/telefonprøver.
AL06C2s aksepterte rapport er rapport-only; AL06 var kvotestoppet, og AL06C1
ble ressursavsluttet gjennom AL14 uten at tidligere avslag/arbeid ble grønt.
Task-listens «active» på bevarte report Tasks er ikke fullført produktlevering
eller bevis for levende provider. AL21s installasjonsreceipt rapporterer
SHA `b4b1ffba09ea9e5f68c226b5f3a5e9b73a1ae0134d2fe7276c35989ee9080472`.
Den eksisterende kvalifiserte binæren gjenbrukes i AL22; den er bygget fra
AL21, **ikke** fra dette dokumentets nye commit.


E-policy: [menneskelevert /status](/Users/perottochristensen/github/ply/planning/design/evidence/60-codex-active-profile-status.txt)
angir session `01a0eea0-3fbc-7dc3-818b-13aadbe1fbeb`, terminal og
`local-development`; [preparationkontrollen](/Users/perottochristensen/github/ply/planning/design/evidence/60-ws05-delivery-preparation.json)
binder dette med samme versions-/permissionbytes og bevart positiv commit.
Dette lukker rad9 på historisk nivå, uten ny probe. E-designs modellspor simulerer
T-17 og et fiktivt Trip-blad; snapshotets siste utfall lukker selve designoppgaven.
Ingen agent, Git eller Trip-workflow ble kjørt av modelltesten.

E-human: [faktisk telefonreise](/Users/perottochristensen/.local/share/ply/planning/human-input-delivery-20261018-01a15003/final-result.json)
bevarer `hiq_03a8925ecf67e792b1b372f4ffff860b`, agentens bruk av samme svar og
separat teknisk levering. E-mobile: [telefonobservasjonen](/Users/perottochristensen/.local/share/ply/planning/slack-followup-delivery1-20261016-01a1431a/phone-observation.json)
bekrefter prosjekt/tittel, men sier uttrykkelig at Check status-klikk ikke ble
særskilt gjengitt og tidligere kontrollutfall fortsatt er unknown.
E-schema: [bevart faktisk SL01-testpass](/Users/perottochristensen/.local/share/ply/planning/backlog-autonomous-20261020-01a15776/cs01-existing-regression.log)
(`cmd 0.558s`) sammenholdt med schemaets eksisterende effect-enum; ingen ny kjøring.
E-journal: [akseptert lokal integrasjonskontroll](/Users/perottochristensen/github/ply/planning/design/evidence/73-task-process-journal-local-integration-result-control-v1.json)
og nåværende `internal/taskjournal` viser første snitt, uten å lukke halen.

Andre native leveranser som inngår i C/deldekket-dommene (alle har eksakt
passed TaskResult og complete closeout i fersk native readback):

| Task | Kandidat / Run | Nåværende kilde |
| --- | --- | --- |
| `ply-automatic-acceptance` | `4d9759842c114c0d66c36297c151d4e49f91143b` / `wfr_c6b4fe82368000b00c1745c8533c447df241b38ba81c194248df220833d08d36` | [delivery_automatic_verification.go](../internal/taskrun/delivery_automatic_verification.go) |
| `ply-cl02-default-worktree-cleanup-20261019` | `c82b132ee8ef8557b35e3fa9ad6d02193ab05b6b` / `wfr_52aa777878af9200ae0852c50e56b3cf7fb6d1ca77ce5ec38ed17c430c94f880` | [delivery_default_cleanup_test.go](../internal/taskrun/delivery_default_cleanup_test.go) |
| `ply-il01-repository-integration-lock` | `f001fb437d7769ec42b65e75d19416546130cae5` / `wfr_287c28982e5494f2e73ee3e09b33d26b4d7ba2ce5bb59325eea54938cfc74ea7` | [task_integration.go](../internal/workspace/task_integration.go) |
| `ply-informative-evidence-retention-20261019` | `e176188c0a4913e5b4c32f1899201f8d53f3fdff` / `wfr_230a4dcd603e90fcabb0227fb9193287b299a64666cfd4536e115d25ba38adeb` | [task_closeout_retention.go](../internal/workspace/task_closeout_retention.go) |
| `ply-na01-current-action-after-answer-20261019` | `de8596e8abe3c232e6a3fe96507d5d539bb76e80` / `wfr_655bf506cdffd9884f67826deae4e4102f1eeddd180d5ee31f04dfd7276bdbea` | [read_inventory.go](../internal/taskrun/read_inventory.go) |
| `ply-resume01-interrupted-task-20261019` | `9d4eee0712e228472da6dcb515ce84f1ab6face5` / `wfr_511f4094a4168845208634567038255c15584655e34c643da7ea8031d3f00625` | [delivery_continuation.go](../internal/taskrun/delivery_continuation.go) |
| `ply-ro04b-confirmed-generation-completion` | `007bac90da1f99aa65d0aa320c83cc407a7b32d6` / `wfr_abecfdb6d7deae187326455499079eaddcc3ad8b34f941190459e4aff8c49c18` | [task_spec_guard.go](../internal/workspace/task_spec_guard.go) |
| `ply-ro04r-codex-capacity-reconciliation` | `32b904d9f116b800095af039f04f7d396445cce8` / `wfr_7a0b4838f529eb980aa926de810a4f4ce2dc9fa79ed97ee9f48fd9a2303579d2` | [delivery_capacity_reconcile.go](../internal/taskrun/delivery_capacity_reconcile.go) |
| `ply-rv01b-codex-completion` | `7553519b533008b2eae2f980bc7aab04a075ac0a` / `wfr_0fa6c3f5003c8fa6828050815856ad04bcc85f8ed931ddac7759dc33bb177c31` | [delivery_review_round.go](../internal/taskrun/delivery_review_round.go) |
| `ply-rv02-blocked-review-reassessment` | `a12ecbecc1b2004cbd340a653241eda6c5156734` / `wfr_af5b99e4f7f521a95f0e2d6aba86e1aeb3c67fa958cc4e3e7b66cb9f91ce8fff` | [delivery_review_round.go](../internal/taskrun/delivery_review_round.go) |
| `ply-sp01-quota-stop-reconciliation` | `8a2db93d070d1b2365a6584f8b3083544944ee58` / `wfr_7a0b29b6f3106c0847473f1979fe5de342c29cd1855b3cd06f917264a5464ca0` | [delivery_quota_reconcile.go](../internal/taskrun/delivery_quota_reconcile.go) |
| `ply-tm01-native-verifier-temp` | `40247f55365ca70c3e5d93f847b928aad6bc53b3` / `wfr_ed2a85b8957847e8636eb984841e3f6d5fdcbb88ee42fb4da049fa643a7c8c26` | [delivery_automatic_verification.go](../internal/taskrun/delivery_automatic_verification.go) |

AL21s produksjonsgren er [delivery_service.go](../internal/taskrun/delivery_service.go);
regresjonen er [al21_closeout_observation_test.go](../internal/taskrun/al21_closeout_observation_test.go).

## Full filavstemming

| # | Eksakt fil / SHA256 | Dom | Restgruppe | Bevis, dekket del og konkret gap |
| --- | --- | --- | --- | --- |
| 1 | `koe/ac-02-spec-scenario-qa-og-avhengighetsport.md`<br>`7c12d887868d96bdf3f73886f0bb6dd1596800a180e704dd1fbcb0fe6073aba5` | U | G30 | AC01/AC02 har ingen native resultater; katalogpathene finnes ikke. Mangler landet, relevant avhengighet og regenerert scenario-QA. E-ark. |
| 2 | `koe/agentic-workflow-support-som-selvhostet-arbeidshierarki.md`<br>`2f246680e7f49ff78f37cb983ed17f9afd6bbe676bd99fe22f5fdfd78857a192` | R | G16 | Epic/Task og faktisk selvhosting finnes; Subtask-hierarki er ikke dokumentert levert. Bruk eksisterende Task/Spec, ikke gammel katalog som premiss. E-ark. |
| 3 | `koe/claude-auto-avviser-delegert-oppgavestart.md`<br>`3a7303e924c4495c33e9aaf03393de8d99593de0c6ab3fe8cc28b56860f0b9e6` | R | G05 | Trust er levert, men første auto-evidenswrite ble avvist; ingen påvist classifier-rotårsak eller full legitim Claude-leveringsreise. E-Claude. |
| 4 | `koe/claude-innstillinger-for-agentstyrte-herdr-faner.md`<br>`8d00abf387b790818caa843ec44feda9262730274e8db5d3f2bc16c85302faae` | R | G05 | AL06B/C2 lukker rapportmodusdelen; settings/closeout og konkrete avslag fra øvrige reiser er fortsatt egne grenser. E-Claude, AL06B. |
| 5 | `koe/claude-opus-medium-ved-arbeidsstart.md`<br>`ae17eb6c642394b7e24e95d6ab0bc718941a82f9401c62c2b8aada1ee7237034` | R | G05 | Execute-modell/effort og menneskeprøvd Opus/medium finnes; eldre startflater er ikke samlet verifisert. Ikke hardkod datert modell. E-Claude. |
| 6 | `koe/claude-review-samme-oppstartsmodus.md`<br>`167cbd5393d90e6a02aed7a13d2517544a0f0ee8f1bd3b99681dfc3e516b3381` | R | G05 | P3 er teknisk levert og faktisk read-only rapportprøvd; dette lukker ikke hele review/verify/integrate-reisen med Claude. AL06B, AL12, E-pilot. |
| 7 | `koe/claude-trust-ved-ply-start.md`<br>`00fa2e8742e665754884aeaa5cfaa409a19795d1d639d12ce95651f99003caa8` | C | — | Mappetillit er menneskeprøvd og integrert; senere auto-avslag er rad 3. E-Claude + herdr_claude_trust.go. |
| 8 | `koe/claude-valgfri-sesjonstilbakemelding.md`<br>`530df240606a8422038ed8d3f99719a5a6979640fd5b7c3ac99f45480baf402c` | U | G05 | Innstilling ble lest som satt 18.10; gjennomgått historikk inneholder ingen bundet etterobservasjon av survey-fravær. Observer neste ordinære Claude-runde, uten testagent. E-Claude. |
| 9 | `koe/codex-runtime-policybinding-for-agentwrite.md`<br>`f1f2ebe3c23349a7e3dfa349f63d69c8a9522721851bfeeb6858eefa90291195` | C | — | Historisk økt87-binding er kontrollert mot menneskelevert /status for session 01a0eea0 og kildebundet preparationrapport. Dette er gammel profilbinding, ingen attestasjon av dagens runtime. E-policy. |
| 10 | `koe/fs01-reviewtest-stopper-i-opprydding-pa-levert-base.md`<br>`f085f7fa48c711a816c586bafc940d2fcac0466d07dc108f82da80c95d3dc4fb` | C | — | AL18 reparerer pane-get-fixturen; full FS01-reise og ukjent-skriververn er bevart. Ingen svekket produksjonscleanup. AL18. |
| 11 | `koe/integrasjonslas-skal-skille-repositories.md`<br>`7c25239fe64b6fd07f388f72beead539e943e75495c47c90f13d4c5432a9b972` | C | — | IL01 binder integrasjonslås til repository/common-dir og ref; native resultat og lukking finnes. E-native. |
| 12 | `koe/karanteneworktree-erfaringsuttak-og-rollback.md`<br>`33ad2e8c0ab89f136b767a25a7e2520f852c822b36339f4177288676b0e628be` | C | — | Læring og eksakt fysisk forkasting er bevart; ingen gammel dirty-kjerne adopteres. E-closure. Dette er historisk avslutning, ikke produktimplementasjon. |
| 13 | `koe/menneskeinput-fra-agent-til-mobil-uten-lokal-handling.md`<br>`f279de61a58bcc8b9c5531c1bdb7b092a173da6bdd8ecb47d579687044a3635f` | R | G19 | HI01/HI02 og faktisk Slack/passkeysvar er levert; kort inngang finnes. Varsel→app-lenke og publisert nødvendig spørsmål ved stopp er ikke bevist fullført av denne leveransen. E-human. |
| 14 | `koe/native-stopp-med-uprøvde-krav.md`<br>`11bd12c823705e6df9cb166f4b7c8bb9a374fbde6dcff7b2dd192194b4f0cd95` | C | — | Samme konkrete uferdig-rapport-gap som rad45: AL10 bevarer failed/not_run/unknown uten grønn kvalifisering. AL10. |
| 15 | `koe/native-verifier-temp-skaper-nestede-workspaces.md`<br>`d99cb53a260fa98fc834a4c9abe6a2f0b5525e9541e8aa786e1405da8b041252` | R | G22 | TM01 retter native verifier-TMPDIR; headless-pretest/indirekte temp og cleanup-begrensninger må fortsatt avgrenses. E-native. |
| 16 | `koe/ply-agentloop-kontrollskatt-og-pflyt-friksjonsreduksjon.md`<br>`40b6af9bcec237772a593e1271550612ec48c18774a7590de5171787670362ae` | R | G03,G17 | Mye kontrollskatt er levert bort AL01–AL21; rest er plannerens fremtidige porter, valgfri progressrapport/ID-diagnostikk og samlet mekanisk flyt. Ingen ny generell kontrollplattform. AL12. |
| 17 | `koe/ply-automatisk-akseptansetest-med-orkestratoransvar.md`<br>`a5f2dc017332a04bbc260df2c631e608457c58f27048a00b00596603179a923f` | C | — | Automatisk eksakt bytebundet akseptanse er native levert; human QA forblir separat. E-native: ply-automatic-acceptance. |
| 18 | `koe/ply-avslutning-av-bevart-worktree-etter-levering.md`<br>`7f48f4c794a5d799de3a7e4382fa59c32d16c3e81d530361abf3939c76abbc9b` | R | G11 | CL02 leverer standardcleanup og nye AL-Tasks er faktisk fjernet; etterfølgende complete/keep-avslutning, fjernede historiske kilder og blocked-week må ikke erklæres generelt løst. E-native. |
| 19 | `koe/ply-cli-svartboks-acceptance-mot-binaer.md`<br>`522dc44f6005a977b4aed64faf2cf6f47ec016b4d17277106540c15b1989d7a7` | R | G02 | Binært lag finnes allerede i test/cli_surface_contract_test.sh; unconditional EXIT sletter feilroot. Undersøk minste faktisk feilbevaringsgap; ingen ny harness. E-test. |
| 20 | `koe/ply-closeout-informative-json-som-avhengigheter.md`<br>`6f7ad11fcb7492e138266ca99f95e3f7b931668c44080f11e81c6d95632d1a6a` | C | — | LD03 native levert/lukket e176188c; informative bytes behandles separat fra ekte kontrollavhengigheter. Historisk LD02B-effekt gjentas ikke. E-native. |
| 21 | `koe/ply-codex-kapasitetsstopp-og-etterfolgende-kommandoretur.md`<br>`18146b23a4dfc8c7ca8576b190c56342bb8493869db48af74e64a90d3e6b711c` | C | — | 32b904d9 leverer kapasitet/kommandoretur-avstemming, RO04 retired og RO04B levert; dette er ikke automatisk backup. E-native. |
| 22 | `koe/ply-codex-tillitsgrant-skal-ikke-stanse-ordinart-arbeid.md`<br>`96b471250d1c3d6a1a709e3c0c310684404694d0f480c4945f45c6a56713d7a4` | C | — | AL03 tillater redundant/withheld trustgrant med advarsel, samme request og én start. AL03. |
| 23 | `koe/ply-codex-valgfrie-produktsporsmal.md`<br>`f6676e38cc5ce22eeccf61a6b07fc93f3d3607a2e9d3ef5d3cb9d5f077ebd8bc` | R | G06 | Valgfrie Codex-modell/update-dialoger og proaktiv trust-UX er ikke samlet løst av AL03. Installerte journeys er modellerte, ikke ny live dialogprøve. E-provider. |
| 24 | `koe/ply-eksisterende-task-worktree-binding-og-epic-refresh.md`<br>`4a6871cc3bc12fa9d015b01cc57d2870a279dae7a48b954e5d7a0e5d15a6937a` | R | G11 | Epic refresh/basehistorikk finnes; å adoptere eksisterende Task-worktree er fortsatt eget gap. Ingen supported adoption-kommando påvist. E-ark. |
| 25 | `koe/ply-enkel-workflow-define-plan-build-verify-review-ship.md`<br>`42deed4aa884f3466a391e6ef532aec7c88244277bfbfa309128985c8beffc86` | R | G17 | Execute/verify/review/integrate finnes; samlet seksfaseflyt og planovergang er ikke ferdig bare fordi kommandoene finnes. E-ark. |
| 26 | `koe/ply-epic-task-eksterne-sakskoblinger.md`<br>`ee9e36e5896811257a001d65a0348bf5cce6ab49e5b092eac9698aca206d478d` | R | G15 | Task/Spec sources er ikke eksterne issuekoblinger; Jira/GitHub sync, unik primærreferanse og commitpolicy er ikke påvist. E-ark. |
| 27 | `koe/ply-execute-malstyrt-leveranse.md`<br>`941e51ac9a1955f81808bb2bbde3e56a3fb1c2ce2b2066cc5c134711af48ae09` | R | G17 | Første execute-skive og timeoutreparasjon er levert; to samtidige featureleveranser er ikke bevist av sekvensielle AL-kjøringer. AL09, AL21, E-ark. |
| 28 | `koe/ply-execute-provideroppdatering-for-taskstart.md`<br>`369edcbcb3df257250cc69de0a5b938fb59313f2315fd73b58e4865715ace3fa` | R | G06 | Oppstartsrestart finnes og nådde ny agent; forebygging av oppdateringsavbrudd og tydelig samlet fortsettelse gjenstår. Ikke bygg restart på nytt. E-provider. |
| 29 | `koe/ply-factory-rig-testoppfolging.md`<br>`74a3b3f60fe67b48f3aa7e5eaee2375b7cedb1fc7fb0712a11a7206dfc450d93` | R | G23 | V4 skiller startup/timeout; gammel v3-årsak er fortsatt ukjent. Samle tids/lastbevis ved nytt relevant funn, ingen retry-til-grønt. E-test. |
| 30 | `koe/ply-filbinding-samtidig-overskriving-under-lesing.md`<br>`ad05154ed6491fcf1ffc67df845608a48937458772a0d943e1a96b5b9f73b1a0` | R | G01 | readFileAs bruker SameFile; AL05 etter-binding-test beviser ikke same-inode/same-size under lesingen. Deterministisk atferdsprøve mangler. AL05, E-test. |
| 31 | `koe/ply-follow-skal-tale-forbigaende-observasjonsfeil.md`<br>`a36a5a1aa302db31b2a172a8aba644772f2ee943080180e6771723c736994bf2` | C | — | AL09 requestbudsjett + AL21 eierpane-race er levert; positive prosessfraværstilfeller er modellert, CLI-unknown er faktisk kjørt. Gamle enkelthendelsers årsaker er fortsatt ukjente. AL09, AL21. |
| 32 | `koe/ply-fortsettelse-gjennom-avklaringer-og-agentbytte.md`<br>`56b055b5da9f977fbf8559f6be7dc8646d0a6fa2d12f75691724c615e690e79b` | R | G07 | NA01 retter neste steg etter svar; generell varig videreføring på tvers av avklaring, agentbytte og tapt historikk er bredere. E-native. |
| 33 | `koe/ply-fravaersbevis-skal-skille-fersk-observasjon-fra-gammel-transport.md`<br>`3e6ae1b9862399508ff75612dd8b49abeb6d2c7ddde3a02bdeb663cd19bd91e8` | C | — | AL20 skiller cached transport fra ferskt inventar med eget tidspunkt; gammel AL06C1-receipt omskrives ikke. AL20. |
| 34 | `koe/ply-gjennomgang-av-alle-godkjenningskrav.md`<br>`59a4865bae30609b519277c40df98b1e951020c79599fa3429f8684988a176ff` | R | G03,G17 | Audit AL02/AL12 er levert; restene er de konkrete gruppene her, ikke en ny total audit eller ny menneskeport. AL12. |
| 35 | `koe/ply-gjenopptakelse-skal-peke-pa-koordinator.md`<br>`b1332d3b93993e320c73e507e46c9a0a6eddb570838e774f349eeb78bfd7f137` | C | — | AL07 navngir koordinator for autorisert recovery; human/headless/unknown beholder egne aktører. AL07. |
| 36 | `koe/ply-kvalifisert-cli-skal-ha-nyttig-buildidentitet.md`<br>`d67c72e120492b7c3668d747b82d240235aaa0fdf72fa75a9d197420eac2950a` | C | — | AL16 gir embedded kjent eller eksplisitt ukjent buildidentitet; native kandidatbinding og faktisk installasjon er egne fakta. AL16. |
| 37 | `koe/ply-lesekontrakt-oppfolging.md`<br>`f072fc4975626fb1beaabc6acea42261198a1a420a5436fa15ebff7a71caf376` | R | G26 | Kjerne-JSON og W1–W9 er levert; appmigrasjon/allowlist/decodere og eventuelle behov utover kjernen eies av appen. Ingen app-kildeinspeksjon i denne avstemmingen. E-app. |
| 38 | `koe/ply-mobilgodkjenning-av-runtimehandlinger.md`<br>`964a7cd533c4c2e54f95f20b00678303b12de652865ae591088489165dd93da4` | R | G08 | Menneskesvar og passkey er ikke faktisk runtime-approval. Providerbundet handling, svarbruk og telefonprøve mangler som samlet bevis. Bevar CL02 og appens scope. E-human. |
| 39 | `koe/ply-native-installasjon-av-kvalifisert-cli.md`<br>`70e116e1811f88ec65ed8e8ab88b864b1c0422b57563f0c82618cb2216311591` | C | — | AL08 erstatter per-jobb installer for self-hosted CLI; AL21 installasjonskvittering og SHA er separat kontrollert. Ingen installasjon i AL22. AL08, AL21. |
| 40 | `koe/ply-nattkjoring-for-qa-og-forbedringsgjeld.md`<br>`e769897fec69f977d82d952105324d5a43f2efe13bc981ac2dea651823d66dad` | R | G27 | Eksisterende worker/week gir utføring; automatisk nattvalg, budsjett, sleep/nett, supervisjon og tillatte rettinger er ikke dermed en levert nattgjeldsreise. E-ark. |
| 41 | `koe/ply-pr-levering-etter-ekstern-merge-og-endret-returvalg.md`<br>`01fae3147c6761962eec0dc0ca57c5ec4403d92650dde3f10bf72279d7b993a0` | R | G20 | PR-flate finnes; ekstern merge/endret returvalg er ikke påvist lukket av PR103s administrative cleanup. Gammel stopped-historikk består. E-PR. |
| 42 | `koe/ply-private-filmodi-uten-nytte-og-upresise-feil.md`<br>`8336726beb3468a3759f39b67a0e3d65667ea54ef54ef2907004564dca88c280` | C | — | AL05 åpner ordinære 0644-inputs og presise rolle/path-feil; autentiserende roller forblir private. Under-read-race ligger rad30. AL05. |
| 43 | `koe/ply-project-repo-medlemskap-fleksible-lokasjoner.md`<br>`ff270bd60040e8a357833939b036d37609bad03f1dae6a451ce7328c321a50d6` | R | G10 | project add oppretter eksplisitt sett; eksisterende medlemsmutasjon og flere repo-anchors per Epic avvises. Bevar registrert PLAN og Sup-07-utkast. E-ark. |
| 44 | `koe/ply-rapportkjoring-skal-tale-ordinare-oppdateringer.md`<br>`af80ab5bd8dc4bb3e26dc8fff96b9338ca69cef62180ec77258eb1607f96ca8e` | C | — | AL04 lar samme rapporteier rapportere etter vanlig CLI/providerdrift med warnings og uendret autoritet. Modellert install/provider, ingen ny live plyapp-prøve. AL04. |
| 45 | `koe/ply-rapportretur-ma-kunne-bevare-uferdige-krav.md`<br>`942ad7169e44b133b0807b8ea8fcb27abcb2818b3605e8db7684a7b41d3144be` | C | — | AL10 leverer ærlig uferdig rapportretur. Akseptert rapport er ikke TaskResult/teknisk pass; senere post-start-avslutning er AL14. AL10, AL14. |
| 46 | `koe/ply-rapportstart-skal-ikke-kreve-en-overflodig-policyfil.md`<br>`2cc7789bb403a8738090cd7366b744401a0020e7440ea77f4555f0987dafac26` | C | — | AL11 ny receipt-rapportstart trenger ingen overflødig policyfil; legacy og ekte kontrollgrenser beholdes. AL11. |
| 47 | `koe/ply-rapportstopp-etter-start-skal-kunne-avsluttes.md`<br>`e26e902197ae91b627f3a3046ffa4cfc836dd5c7a316aea31efd1c994d0a228b` | C | — | AL14 leverer post-start ressursretirement og faktisk AL06C1-avslutning; ikke fullført pilotarbeid eller grønn historikk. AL14, AL20, E-pilot. |
| 48 | `koe/ply-registervalidering-forsinker-oppgavestart.md`<br>`0e4b86458ea203a90b378bd0129101aa94511827db8fe766b6a4064ee90f74db` | C | — | RO01B–RO04B er levert; kildeavgrenset registerkontroll/generasjon finnes. Målte precheck-tall er ikke full provider-startgevinst. E-native. |
| 49 | `koe/ply-review-og-retur-uten-manuell-koordinering.md`<br>`5f62058e4dfe5a145469dee837eca3e5a968eb1e80b12c4bad00b36b660ac67f` | R | G17 | FS01/RV01/RV02 og AL15/17 løser konkrete review/returgrenser; vanlig flyt uten manuell vekking i alle tilstander er ikke samlet bevist. E-native, AL17. |
| 50 | `koe/ply-reviewmandat-skal-skille-tester-fra-native-verifisering.md`<br>`9b87c08f70f9f9d5154419db446cdce0450bb337a4e7f80897967830cb0ce268` | C | — | AL15 og AL17 avgrenser owner-callbacks til faktisk tildelt Run og tillater isolerte fixturetester. Reelt ps-avslag blir ikke et pass. AL15, AL17. |
| 51 | `koe/ply-reviewretur-skal-tale-vanlig-forklaringstekst.md`<br>`c929741220020d29091c9ea7fa886ac3c699783d342b73d00cb8facb92405b30` | C | — | AL13 tillater lang/flerlinjet reviewprosa innen totalbudsjett, med faktisk AL14-reviewretur etter installasjon. AL13. |
| 52 | `koe/ply-samlet-epic-gjennomgang-for-pr.md`<br>`a624d4fbd251e9eed6b42a67f8b59f14d0fe67da58cd630a09954fb6544c2151` | R | G21 | Ingen samlet ny vurdering av dagens Epic mot master er gjort her. Gjenbruk gyldige bevis ved separat gjennomgang; ingen PR/push bestilles. E-PR. |
| 53 | `koe/ply-task-integrasjon-forste-reflog.md`<br>`0244f2ae043c50caaaf42ff6c0a6ca8d5e0305de686a8e1eecf4237cd729d6b7` | C | — | 4c871e1 har faktisk human pass og native første-reflog-integrasjon; gamle journalgap omskrives ikke. E-first. |
| 54 | `koe/ply-task-nytt-repo-og-forste-kode.md`<br>`3aa9df4c9e1f53a4e97f6051a1974252c2d3fa3279e8c2557f10029df5232cf9` | R | G14 | Repo-først gjelder; oppretting av nytt repo/første kode er ikke levert av task create eller project add. E-ark. |
| 55 | `koe/ply-task-prosessjournal-og-observer.md`<br>`93300ca325f0b8e51c039463512fe47fa7fa5480ccb0aee838faaf1ebb7460d3` | R | G25 | Journalens første snitt er menneskeprøvd/integrert; langvarighetsmetning (~292år) og bred race-timeout er fortsatt navngitt lav gjeld. E-journal. |
| 56 | `koe/ply-task-queue-nattko.md`<br>`d80f12fe430e36752711629079115117dfec8a1a77968ecc0218a915e5bcfad4` | R | G24 | Queue/base-refresh finnes; tidsstempelpresisjon, CLI-diagnostikk og sjeldne kombinasjoner er ikke uttømt. E-test. |
| 57 | `koe/ply-task-run-native-avslutning-og-qa-friksjon.md`<br>`cf3f609831c8fca1b5d0a60cfe78b10fe53fc2685349473239940ea39a5e7274` | R | G23 | Native nye AL-closeouts finnes; gammel factory inactive→TaskResult og Ctrl+C/prosessgruppe/siste rapport trenger egen sammenhengende prøve. AL18/21 er ikke samme headless-rigg. E-test. |
| 58 | `koe/ply-task-run-sjeldne-crash-kombinasjoner.md`<br>`ee9f2e21f2e925db1494c7a4d3ce3b923aa7bb0a8462f1beef26b83f25c4ba3c` | R | G24 | F1–F9 og CAS er grunnbevis; full crash/external-write/plattformkombinasjon er ikke bevist. Batch eksisterende returtester. E-test. |
| 59 | `koe/ply-task-spec-oppfolging-og-nattko.md`<br>`75535374676ec4c33ff063de39bb14a233a3d3deb4e28d26d320a7c0c119e0a0` | R | G16,G24 | Problem/Spec/valg/basehistorikk finnes; disposition/alternativer/dekomponering/resultatavhengighet/snapshot-restore og restmatrise må vurderes mot aktuelle kommandoer. Felt alene lukker ikke reise. E-ark. |
| 60 | `koe/ply-trekommando-utviklingsloop.md`<br>`fccfa3174ee47979d38761eb80f7b470b6a86b47c2c7653ded3fa2405c119384` | R | G17,G18 | Task-integrasjon og interaktiv execute finnes; videre mekanisk loop og generell headless-controller er ulike rester, ikke gjenlevering av merge. E-ark. |
| 61 | `koe/ply-workflow-agentrapport-mottak.md`<br>`da045526361d2a99e8de6d563641f7a073fe3cfdc476876f656c1bdedf95f0ff` | C | — | WF01 filbasert start/resultatretur er akseptert og integrert; immutable handoff/return finnes fortsatt. Tail er rad84. E-WF. |
| 62 | `koe/ply-workflow-execute-headless.md`<br>`d19d17fe74d14c16d4170e9e1c7a4aa0250fdc6022055820e87895bb8f31b5d2` | R | G18 | worker/week og factory er eksisterende ikke-interaktive spor; generell execute-headless med samme interaktive kontrakt er ikke påvist. Utred provider/session/supervisjon fremfor parallell motor. E-ark. |
| 63 | `koe/ply-workflow-handoff.md`<br>`7087be98a9360e5dcdaa499b6733b39545fdbe6f91ba112f075c5be02d70c07c` | C | — | WF01 generell filhandoff og redigert inspect/show er levert; samme vertikale levering som rad61. E-WF. |
| 64 | `koe/ply-workflow-merge-taskresultat-til-parent.md`<br>`7e2e9df2214da6aa9a9b5bf792601e5e7f7be2741064c77920ee3cf1358921e9` | C | — | WS04 task integrate leverer resultat→parent; navnet workflow merge trenger ikke gjeninnføres. Kontrollert manuell Git-landing var ikke native live integrasjon. E-WS04. |
| 65 | `koe/ply-workflow-notification-nattgjeld.md`<br>`99fece3618d817ec8e5a1ddde5cc4db0156ec00fc86d5673447040cbdce5eb72` | R | G19 | 64MiB rute-state og arvet help er kjente grenser; gammelt unknown-send må ikke resendes. Senere runtime-route-bevis lukker ikke lagringsskala/diagnostikk automatisk. E-notify. |
| 66 | `koe/ply-workflow-notification-slack-oppsett-og-live-qa.md`<br>`b6256d75149bea36b913fdaa642551568b56089533315c9b61f0a06854930291` | C | — | Første Slack-oppsett, fysisk kanalprøve og kontrollert lokal integrasjon er oppfylt. Ingen ny Slack-probe utført. E-Slack. |
| 67 | `koe/ply-workflow-resume-etter-avbrutt-agent.md`<br>`e02126c4cc297f45bc07d41604dceb32ec369418528968d6a7585e231ae8afe9` | R | G07 | RESUME01 og Codex-sessionfortsettelse finnes; feil next action ved shell-only start og innholdsstopp med levende mottaker er konkrete rester. E-native, E-provider. |
| 68 | `koe/ply-workflow-varsling-automatisering.md`<br>`877a47544254d44e0858b21ad288b04dd9de9b2910ecbf4193b11eca3b4b6f41` | R | G19 | Stående regel og faktisk sending er levert; presis avslutningstiming, videre automatikk og toveis trådsvar er ikke samlet levert. E-notify. |
| 69 | `koe/ply-workspace-epic-create.md`<br>`0a2293810ba37e688526452c60f13e241820376ab83dce8002ae032483b9fcc3` | R | G12 | Epic adopt finnes; design79s nye branch/worktree med intent/retry er ikke samme produktreise. Native kandidat er unbound uten resultater. E-ark. |
| 70 | `koe/ply-workspace-og-workflow-ansvarsgrense.md`<br>`b0a6a4a5c10985520f40b14fe6e5d499a3f0909f6e0e3379f6c061a1a0780f0b` | C | — | Workspace eier register/Git, workflow eier utføring; WS03s valgte vertikale ansvarsgrense støttes av aktuelle cmd/workspace.go og cmd/workflow.go. Ikke løfte om all fremtidig policy. E-WS03. |
| 71 | `koe/ply-workspace-project-init-convenience.md`<br>`c5e02f52f6e6cba1031657e465dc59d2267fdf6bd7a057fde587b52719983b2e` | R | G13 | project add/show/list finnes, project init discovery/preview/defaults finnes ikke som kommando. Gjenbruk eksplisitt medlemssett. E-ark. |
| 72 | `koe/plyapp-autorisert-kontroll-uten-bekreftet-utfall.md`<br>`0f11452ac7848db7b87fbbda3a4e1ca258d83a8a152f4dceb58f2fb8a4e3d621` | U | G09 | Task update er telefonobservert, men opprinnelig autorisert pause/gjenoppta-resultat og komplett statuskort er ikke fysisk bekreftet i bevisene lest. App-eier må observere samme request, ikke gjenta effekt. E-mobile. |
| 73 | `koe/plyapp-home-leseflater-og-livssyklus.md`<br>`30c8fb5e808a4cbfd144d3428283b882463e5d6259c89f14c8755fdbac1a069f` | C | — | W1–W9 og overview b0cfdf5 er agentakseptert/installert/lokalt integrert. Human produkt-QA ble ikke gjort og native human/integration-aksene var unknown; dette er oppgitt i kontrollrapporten. E-app. |
| 74 | `koe/plyapp-lesestotte-teknisk-oppfolging.md`<br>`9c46d1d7f7f75211bc0b2ca9ad0fb12a14ec85b18dcfda1414bb9af11eef04fd` | R | G26 | 30-worktree ~1.20s er måling, ikke sub-sekundsløfte; større inventar og eldre lint gjenstår behovsstyrt. Bevar appansvar og schemas. E-app. |
| 75 | `koe/providerkvote-feilklassifisert-effekt-ma-kunne-avklares.md`<br>`d59c4b9c1bb0b556989f76b9b5f27776000019b78f6bc19e877a199ee8150621` | C | — | SP01 auditert kvote-effektavklaring er levert; faktisk appreise tok to kildeauditer og ordinær retire, med gammel partial bevart. Ikke generell backupautomatikk. E-native. |
| 76 | `koe/providerkvote-for-forste-handling-skal-gi-enkel-backup.md`<br>`6e5575db1790d22220ad88eb45455d1f7675c9c68e7674bd1ab9bd65537d6ad8` | R | G04 | Kvote før første handling og etter nyttig arbeid/review krever fortsatt manuell fersk backupregistrering. Retire/reconcile finnes; ett støttet ende-til-ende backupsteg mangler. E-provider. |
| 77 | `koe/provideroppdateringer-uten-fast-binarkrav.md`<br>`fee69d729ff0e48ebb96311dfc5a26c0ef57bd73402139554f42eb80f600a0df` | R | G06 | PH01 fjerner én hardkodet hash; AL04 viderefører startet rapport. Generell kompatibilitet ved oppdatering for begge providere og faktisk mobilsvar er bredere. E-provider. |
| 78 | `koe/reell-trip-validering-etter-ply-piloten.md`<br>`b3fe0598c15af8f207f51f318d178cb081023ca05ccc4e3140c2c34fc045d532` | R | G29 | Ingen faktisk valgt Trip-feature er kjørt gjennom Ply i undersøkte bevis. Dokument-/fixtureeksempler er ikke live Trip-validering. Avhengig av konkret repo/Spec/mål. E-Trip. |
| 79 | `koe/sl01-capabilities-schema-samsvar.md`<br>`93fc0ad3dd0e1ad53e1957de29b91023bf9fbd85a7c0b9d46daea55123924c1e` | C | — | Offentlig capabilities-schema har recipient_runtime_acceptance; bevart eksisterende SL01-regresjon passerte på levert base. Ikke ny retting eller ny test her. E-schema. |
| 80 | `koe/startup-cli-test-bruker-fjernet-permission-flagg.md`<br>`f04bbe579e34f5a3640a670d2b3c81c6d57c639fa80270ba70d9d12c3d22ab41` | C | — | AL19 bruker gjeldende cause/action resolve-stop for begge modellerte providere; opprinnelig ps-nektet ekstratest er bevart og utelatt med ny review. AL19. |
| 81 | `koe/task-integrasjon-restdekning-nattjobb.md`<br>`36237b6fc654e85c0cdcd2f68f5c27d40d7690981284663e2648c7309c10700e` | R | G24 | WS04 hovedreise er levert; tail-crash/race/format/modekombinasjoner er fortsatt batchbar gjeld. Falsk Unix-mode-dom er ikke produktfeil. E-WS04, E-test. |
| 82 | `koe/tilstandsmaskin-og-overgangskontrakter.md`<br>`1a7cae3d66847c1a7f105c93da512f433eeb2b8d101b75f593226f6287419b46` | C | — | Design06 fullfører bestilt én-node tilstandsmaskin; dette er et dokumentert designutfall, ikke bevis for en AC02-produktkjerne. E-design. |
| 83 | `koe/trip-vedlikehold-cloud-build-og-firebase.md`<br>`b5d9de7bcd976b0f46aa62a023ced8fc54abe6f4fdad13d5b9cedd11ffbb533d` | R | G28 | To faktiske Trip Cloud Build/Firebase-problemer er bevart; korrekte Project/Repo/Epic-bindinger og native registreringsreadback for nettopp disse mangler. E-Trip. |
| 84 | `koe/wf01-restdekning-nattjobb.md`<br>`c53a6ed489dcd0a7d16eee0dfb102ef53fe12a1d1ee332ea865f014fef527b66` | R | G24 | WF01 hovedreise er levert; navngitte schema/store/fault/stream/CAS-varianter er ikke uttømt. Samle med relevant vedlikehold. E-WF, E-test. |
| 85 | `koe/workflow-visualisering.md`<br>`c306775be97803667653b16a4ac88629aff16db7e886753c5bcb5facd23e0bb9` | C | — | f2de0c6 retter faktisk bokser/piler og details/JSON; human pass og lokal integrasjon er kontrollert etter utilstrekkelig første tekstrapport. E-trace. |
| 86 | `koe/ws03-restdekning-nattjobb.md`<br>`b4be15d0f5a33beba3e1b7d1343c40dc571c470ea4e361ea49c3b511fc1a5286` | R | G24 | WS03s startreceipt-gap ble lukket; sjeldne fault/race/power-loss/plattformvarianter består. Ingen retry av gammel integrasjon. E-WS03, E-test. |

## Unike rester og avhengighetsrekkefølge

Rekkefølgen nedenfor er en teknisk avhengighetsanbefaling til eksisterende
koordinator, ikke ny køprioritering. Uavhengige grener kan velges uten å
vente på katalog, headless eller full restmatrise. Ingen nye Tasks er registrert.

| Gruppe | Konkret rest og neste nyttige handling | Avhengighet / eier |
| --- | --- | --- |
| G01 | Deterministisk same-inode/same-size endring **under** `readFileAs`; avgjør faktisk innholdsbinding før eventuell minste retting. | Uavhengig; taskrun/input-eier, rad30. |
| G02 | Gjenskap feil i eksisterende binærscript og vurder bevaring av streams/temp-locator før EXIT-opprydding. | Uavhengig; CLI-test-eier, eksisterende `ply-candidate-cli-blackbox-acceptance`. |
| G03 | Skill fremtidige koordinatorporter fra mottakerkrav; gjør valgfri progressrapport med faktiske testnavn og deklarerte verifier-ID-er forståelig. | AL10/11 levert; taskrun/planner, rad16/34. Ingen ny obligatorisk progressgate. |
| G04 | Ett støttet backupsteg etter faktisk kvote før handling, etter committed arbeid og før review; bevar gammel eier/kilde/stop og ferskt mandat. | SP01/RO04R og retire finnes; provider/runtime-eier. Codex-backup i seg selv lukker ikke gapet. |
| G05 | Samlet legitim Claude-reise og støttede launchsettings/eldre modellflater; inkluder survey-fravær ved neste ordinære start. | AL06B og faktisk AL06C2 er grunnlag, ikke full leverings-QA. Vent på faktisk tilgjengelig provider ved nødvendig live aksept; ingen kvoteprobe. |
| G06 | Forebygg ordinære update-/valgfrie dialogavbrudd og revalider kompatibilitet uten fast binærhash for begge providere. | PH01/AL04/restart gjenbrukes; runtime-eier. Samme-sesjon/mobilprøve er egen faktisk akseptgrense. |
| G07 | Samlet videreføring etter avklaring/agentbytte; rett shell-only/innholdsstopp-next-action fra faktisk observerbar kilde. | RESUME01/02B og NA01 finnes; fortsettelseseier. Ukjent effekt gir observasjon, aldri blind replay. |
| G08 | Native runtimeforespørsel→telefonbeslutning→samme ventende handling. | G05/G06 der provider krever det; Ply-root + app-eier, eksisterende CM/CL-spor. Passkeysvar alene er utilstrekkelig. |
| G09 | Avklar det opprinnelige mobilkontrollutfallet og komplett statuskort fra samme request. | G07/read surfaces; app-eier. Bevar tidligere unknown, ingen gjentatt pause/resume for å få bevis. |
| G10 | Eksisterende Project-medlemsmutasjon og repo-spesifikke Epic-anchors uten privat rebinding. | Workspace-eier; registrert PLAN er relevant nabospec, ikke startet implementasjon her. |
| G11 | Adopter eksisterende Task-worktree og sluttfør historisk `complete/keep`/fjernet kilde/blocked-week via støttet overgang. | Base-refresh og standardcleanup finnes; workspace/lifecycle-eier. Del opp etter faktisk påvist gap, ikke ny generell lockplattform. |
| G12 | Ny Epic-branch/worktree med intent og idempotens, utover `adopt`. | Eksisterende design79 / `ply-candidate-workspace-epic-create`; workspace-eier. G10 kun ved flerrepo-krav. |
| G13 | Project-init discovery/preview/defaults over eksplisitt `project add`. | Uavhengig av G12; eksisterende `ply-candidate-project-init-convenience`. |
| G14 | Problem→nytt repo→første kode med kjent Repo/common-dir/baseline. | Konkret repo-/authorityvalg før effekt; eksisterende `ply-candidate-new-repo-first-code`. G12/13 er mulige hjelpere, ikke automatiske krav. |
| G15 | Eksterne Jira/GitHub-referanser, unik primærkobling og egen commitpolicyconsumer. | Eksisterende `ply-candidate-external-issue-links`; provider/API-undersøkelse mangler. Ingen sync-effekt autorisert av audit. |
| G16 | Avgrens selvhostet Subtask/dekomponering og Task/Spec-disposition/alternativer/resultatavhengighet/restore mot dagens støtte. | Workspace/plan-eier; `ply-ds-01-developer-start`, eksisterende Task/Spec. G10 kun for faktisk flerrepo-gruppering. Ikke gjenoppliv historisk kernel som standard. |
| G17 | Samlet enkel plan→execute→review→ship-flyt, normal vekking/retur og to samtidige features. | AL/RV/FS-støtte finnes; `ply-candidate-agent-loop-automation`. PLAN/PROD forblir registrering-only; ingen ny per-app Epic. |
| G18 | Generell headless execute med samme identitet/resultat/stopp; sammenlign først worker/week/factory og eksisterende providerflater. | G06/G07 og konkret provideravgrensning; runtime-eier. Ikke forutsetning for G01–G17. |
| G19 | Varslingsskala/help/trygg apply-diagnose, presis ferdigtiming, nyttige appinnganger og videre automatikk/toveis samtale. | Notification/launcher + app-eier; skill og faktisk sending finnes. Flere delutfall under samme eksisterende varselsløp; ukjent gammel POST forblir ukjent. |
| G20 | PR-delivery etter ekstern merge/endret returvalg med ærlig historikk og cleanup. | Native delivery-eier; `pull_request` finnes, PR103-administrasjon beviser ikke ny produktstøtte. Ingen remoteeffekt her. |
| G21 | Samlet Epic-mot-master-vurdering og konkret forslag til eventuell videre publisering. | Etter valgt lokalt omfang, separat review-eier. Dette dokumentet dekker koe, ikke samlet integrasjons-QA eller PR-tillatelse. |
| G22 | Revalider headless-pretest og indirekte temp/cleanup mot levert TM01. | TM01 native fix finnes; taskrun/test-eier. Krev konkret rest før kodeendring. |
| G23 | Factory startup-last/timeout ved nytt funn, native inactive→resultat og Ctrl+C/prosessgruppe/siste rapport. | G18 dersom generell headless utvides; eksisterende factory-eier, rad29/57. Gammel v3-årsak kan forbli ukjent. |
| G24 | Én behovsstyrt testbatch: workspace/store/queue-tid, taskretur-crash, Task/Spec-races, taskintegrasjon, WF01-handoff og WS03-fault/plattform. | Seks separate eksisterende testområder (rad56/58/59/81/84/86); relevante eiere. Ikke seks nye dagsporter; ingen påstand om uttømt matrise. |
| G25 | Journalens langvarighetsmetning og målrettet race-løp. | Eksisterende journal-eier; første snitt er levert. Ekstreme intervaller er lav gjeld, ikke blokkering av vanlig bruk. |
| G26 | Behovsstyrt appmigrasjon/read-contract og større worktreeinventar/eldre lint. | Plyapp eier egen implementasjon; Ply eier offentlig lesekontrakt/ytelse. Ikke ny oppstart av registrert PROD her. |
| G27 | Avgrens nattgjeldsvalg/budsjett/supervisjon og tillatte rettinger over eksisterende worker/week. | Etter G18/G23 der headless trengs; `ply-candidate-batched-night-qa`. En scheduler er ikke bevis på trygg produktreise. |
| G28 | Registrer to bevarte Trip vedlikeholdsproblemer med korrekte eksisterende Project/Repo/Epic-bindinger. | Trip-eier kontrollerer faktisk målrepo og readback; `ply-reference-trip-maintenance`. Ingen ny per-app Project/Epic-rebinding. |
| G29 | Én faktisk Trip-Spec/feature gjennom kjeden, med ekte resultat og relevant aksept. | Etter konkret valgt Trip-problem/mål, eventuelt G28. `ply-reference-trip-workflow-validation`; separat Trip-eier. |
| G30 | Revalider AC02s opprinnelige avhengighet/akseptanse mot dagens produkt. Undersøk bevart kilde bare der det gir nødvendig ny informasjon. | AC01/02 er archived/unbound uten resultat; planner + kernel/domeneeier. Manglende gammel dirty-kode er ikke mandat til ny katalogplattform. |

G01–G04 gir konkrete uavhengige neste undersøkelser. G05–G09 har faktiske
provider-/human-/appgrenser. G10–G18 er videre produktkapabiliteter som må
avgrenses mot dagens støtte; G19–G29 er andre produkt-/vedlikeholdsgrener.
G30 er en revalidering før gammel AC02-QA, ikke en blokkering av hele køen.
Ingen timeprognose eller nytt menneskesvar er nødvendig for å lese dette som
rekkefølgegrunnlag; koordinatoren eier neste konkrete sekvensielle valg.

Beskyttede naboer er uendret: `ply-consumer-agent-status-v1@1` (rank1) og
`ply-registered-cross-repo-plan-v1@1` (rank2). Før AL22 var `current=null`;
under AL22 er bare denne oppgavens egen preparation current. De eksakte
Goal-manifestreferansene/order er kontrollert mot bevart før-kø og fersk
readback. Ingen av naboene startet, ingen per-app Epic/Project ble opprettet,
og `sup-07-hendelsesutlost-review-utkast.md` er fortsatt utkast.

## Verifikasjon og grenser

Nytt i AL22: ekstern kontroll av snapshotets egen hash, alle 86 filnavn,
originalbytes/størrelser og SHA256 mot tabellen; filregnskap/gruppekobling;
lokale dokument-/kildelenker; native Task/Run/resultat/integrasjons-/closeout-
referanser; konkrete kildeegenskaper; beskyttede kørefs og diffcheck.
Akseptansen kjører disse kontrollene mot den eksakte dokumentkandidaten.
Private stdout/stderr/exit bevares i bevisroten. Ingen nye produkttester,
bred Go-suite, providerprobe, apprestart, Slack-sending eller CLI-installasjon.

Gjenbrukt: AL12s godkjenningsaudit, AL13–AL21 native review/verification/
delivery/closeout og installert smoke, samt navngitte eldre aksepterte
kontrollrapporter. AL18s positive prosessinspeksjon/Herdr er modellert.
AL19s smalere kvalifiserte testutvalg erstatter ikke den bevarte ps-nektede
ekstratesten. AL20s faktiske gamle AL06C1-readback beholdt immutable receipt;
AL21s CLI-unknown/warn/JSON er faktisk prosesskjørt, positiv absence er modellert.
Dette er ikke en ny fullsuite eller et nytt menneskelig pass.

Separat native review av eksakt clean commit og automatisk native verify
må passere før AL22 kan leveres til avtalt lokal Epic. Teknisk kvalifikasjon,
review, historisk human QA og faktisk AL22-integrasjon er separate fakta.
Selve auditten leverer ingen av de uferdige G-gruppene og avslutter ikke
hele backloggmandatet.
