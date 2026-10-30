# Adopt an existing Task source

Denne guiden beskriver de gjenværende direkte `worktree adopt`- og
`release-adoption`-kommandoene. Den eldre
[adopsjonskontrakten](task-worktree-adoption-spec.md) er historisk; dens
execute-, leverings- og automatiske oppryddingsløp finnes ikke lenger.

Adopsjon bruker fortsatt en ubrukt aktiv Task med sin første, gjeldende Goal
registrert i køen. Denne varianten krever Goal-typen med automatisk lokal
Epic-levering og registry-format 7. Det er en begrensning ved adopsjonskommandoen;
vanlig [Task-registrering og worktree-oppretting](../README.md#workspace) trenger
ikke denne dokumentkjeden.

Kjør fra registrert Epic-returrot. Kilden må være en fysisk, ren linked worktree
på en full lokal branch, med egne commits etter gjeldende Epic-base. Adopsjonen
kontrollerer også registrert eierskap og tilgjengelige writer-observasjoner;
ukjent writer-status blir ikke tolket som fravær.

```sh
ply workspace task worktree adopt TASK_ID \
  --project PROJECT --repo REPO --epic EPIC \
  --path /physical/source --ref refs/heads/useful \
  --expected-oid SOURCE_OID --format json
# Gjenta identiske flagg med --apply --confirmation sha256:PREVIEW_DIGEST.

ply workspace task show TASK_ID --format json
```

Preview skriver ingenting. Apply binder den eksisterende kilden og bevarer
source-, parent- og Goal-identitet. Den flytter ingen Git-ref, starter ingen agent
og utfører ingen review, test, integrasjon eller sletting. Velg senere operasjoner
uttrykkelig; [kommandooversikten](../README.md#human-owned-process) beskriver
avgrensningen.

En ubenyttet binding kan frigis med en egen forhåndsvisning og bekreftelse:

```sh
ply workspace task worktree release-adoption ADOPTION_ID \
  --reason 'Den valgte kildebindingen skal frigis' --format json
# Gjenta med --apply --confirmation sha256:RELEASE_PREVIEW_DIGEST.
```

Release endrer bare den eksakte metadata-bindingen og bevarer Git, filer og
historikk. Selection, preparation eller TaskResult som allerede bruker bindingen
hindrer denne begrensede frigivelsen. Ukjent publiseringsutfall krever lesing av
samme operasjon før videre handling; det gir ikke grunnlag for å slette kilden.

## Continue committed work through a fresh Task

Nyttig arbeid kan beskrives som eksplisitte input til et nytt oppdrag. Bevar
opprinnelig base `B`, kildecommit `S`, tree og patch-hash slik at mottakeren kan
vurdere det faktiske bidraget. Disse kommandoene leser den gamle Git-kilden:

```sh
git -C /physical/old-source diff --binary --full-index B S > /private/draft.patch
git -C /physical/old-source rev-parse S 'S^{tree}'
shasum -a 256 /private/draft.patch
```

Registrer den nye Tasken og opprett dens worktree med de direkte kommandoene i
[README](../README.md#workspace). Velg selv hvordan patchen skal vurderes og
brukes. Ucommittede filer i den gamle kilden er ikke med i denne eksporten.
Ingen backup-goal-, execute- eller automatisk erstatningsstart følger av eksporten.
