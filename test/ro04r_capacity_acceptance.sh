#!/bin/sh
set -eu
: "${PLY_CANDIDATE_BINARY:?exact candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate commit required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
RO04R_GO=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go
export GOTOOLCHAIN=local
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
"$PLY_CANDIDATE_BINARY" version --format json | python3 -c 'import json,os,sys; x=json.load(sys.stdin); assert x["source_commit"]==os.environ["PLY_CANDIDATE_OID"] and not x["dirty"]'
"$RO04R_GO" test -mod=readonly ./internal/taskrun -run '^(TestAR03RetirementDoesNotReplayDeniedHerdrInspection|TestCapacityReconcileBindsResolvedHistory|TestCapacityReconcileCompletionAfterOverload|TestCapacityReconcileEarlierRunningProcessMustComplete|TestCapacityReconcileLiteralExecutorInventory|TestCapacityReconcileNativeRetirement|TestCapacityReconcileNativeSafetyGates|TestCapacityReconcilePreCapacityTerminalFailureIsHistory|TestCapacityReconcileRejectedResolutionOnly|TestCapacityReconcileRejectsOmittedRunningResult|TestCapacityReconcileRejectsUnfinishedProcesses|TestCapacityReconcileRejectsUnknownEvidence|TestDeliveryReceiptDoesNotDowngradeLegacyContract|TestDeliveryReceiptDoesNotRequireBlanketPermission|TestDeliveryReceiptLaterRuntimeRefusalPreservesStop|TestDeliveryReceiptReportedRootsNeverAttestSandbox|TestDeliveryReceiptResolutionPreservesUncertainExternalEffects|TestDeliveryReceiptResolutionSelectsFrozenStartupEntrypoint|TestDeliveryReceiptReuseRejectsChangedInputsWithoutRunningAcceptance|TestDeliveryRetireEndsSelectedOwnerAndFreesQueueWithoutOwnerCallback|TestDeliveryRetireStopsOnActiveWriterOrUnknownEffect|TestQuotaReconcileNeverReopensOwnerAndRecoversLostState|TestQuotaReconcileRetainsSafetyGates|TestQuotaReconcileSourceChangeDuringPaneClose|TestQuotaReconcileTranscriptBoundaries|TestQuotaReconcileWindowPreservesHistoricalEffects)$' -skip '^TestDeliveryReceiptCLIJourney$' -timeout=20m -count=1
"$RO04R_GO" vet -mod=readonly ./internal/taskrun ./cmd
git diff --check
printf 'Retained concrete regression tests passed; removed process journeys are excluded.\n'
