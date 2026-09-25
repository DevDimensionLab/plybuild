package logger

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

type logrusLoggerState struct {
	level        logrus.Level
	formatter    logrus.Formatter
	output       io.Writer
	hooks        logrus.LevelHooks
	reportCaller bool
}

type loggerGlobalsSnapshot struct {
	fieldLogger      bool
	collector        *Collector
	collectorEntries []*logrus.Entry
	packageLogger    *logrus.Logger
	packageState     logrusLoggerState
	standardLogger   *logrus.Logger
	standardState    logrusLoggerState
}

func TestCollectorPreservesExactOrderedLevelsAndAppendOnlyEntryPointers(t *testing.T) {
	collector := &Collector{}
	levels := collector.Levels()
	if len(levels) == 0 {
		t.Fatal("Collector level characterization population is empty")
	}
	wantLevels := []logrus.Level{logrus.InfoLevel, logrus.WarnLevel}
	if !reflect.DeepEqual(levels, wantLevels) {
		t.Fatalf("Collector levels were %#v, want exact ordered levels %#v", levels, wantLevels)
	}

	testLogger := logrus.New()
	existing := logrus.NewEntry(testLogger)
	existing.Message = "existing entry"
	info := logrus.NewEntry(testLogger)
	info.Level = logrus.InfoLevel
	info.Message = "complete info entry"
	warn := logrus.NewEntry(testLogger)
	warn.Level = logrus.WarnLevel
	warn.Message = "complete warning entry"
	collector.entries = []*logrus.Entry{existing}

	if err := collector.Fire(info); err != nil {
		t.Fatalf("Collector Fire(info) returned an error: %v", err)
	}
	if err := collector.Fire(warn); err != nil {
		t.Fatalf("Collector Fire(warn) returned an error: %v", err)
	}

	entries := collector.entries
	if len(entries) == 0 {
		t.Fatal("recorded Collector entry population is empty")
	}
	wantEntries := []*logrus.Entry{existing, info, warn}
	if len(entries) != len(wantEntries) {
		t.Fatalf("Collector recorded %d entries, want exactly %d", len(entries), len(wantEntries))
	}
	for index, want := range wantEntries {
		if entries[index] != want {
			t.Fatalf("Collector entry %d identity is %p, want exact pointer %p", index, entries[index], want)
		}
	}

	info.Message = "mutated after Fire"
	if collector.entries[1].Message != info.Message {
		t.Fatalf("Collector copied the fired entry: got message %q, want shared message %q",
			collector.entries[1].Message, info.Message)
	}
}

func TestPackageInitializationInstallsCollectorHookAndCapturesExactInfoWarnEntries(t *testing.T) {
	runWithRestoredLoggerGlobals(t, "initialized package logger", func(t *testing.T) {
		if log == nil || collector == nil {
			t.Fatalf("logger package initialized incomplete globals: log=%p collector=%p", log, collector)
		}
		if log == logrus.StandardLogger() {
			t.Fatal("logger package private logger unexpectedly aliases the standard logrus logger")
		}
		if fieldLogger {
			t.Fatal("logger package initialized field-logger state as true, want false")
		}

		levels := logrus.AllLevels
		if len(levels) == 0 {
			t.Fatal("logrus hook-level characterization population is empty")
		}
		for _, level := range levels {
			hooks := log.Hooks[level]
			wantCount := 0
			if level == logrus.InfoLevel || level == logrus.WarnLevel {
				wantCount = 1
			}
			if len(hooks) != wantCount {
				t.Fatalf("package logger has %d hooks at %s, want exactly %d", len(hooks), level, wantCount)
			}
			if wantCount == 1 && hooks[0] != collector {
				t.Fatalf("package logger hook at %s is %T %p, want exact Collector %p",
					level, hooks[0], hooks[0], collector)
			}
		}

		output := &bytes.Buffer{}
		log.SetOutput(output)
		log.SetLevel(logrus.DebugLevel)
		collector.entries = nil
		log.WithField("kind", "info").Info("complete initialized info")
		log.Error("uncaptured initialized error")
		log.WithField("kind", "warn").Warn("complete initialized warning")
		log.Debug("uncaptured initialized debug")

		entries := collector.entries
		if len(entries) == 0 {
			t.Fatal("initialized Collector recorded-entry population is empty")
		}
		if len(entries) != 2 {
			t.Fatalf("initialized Collector recorded %d entries, want exactly 2", len(entries))
		}
		wantLevels := []logrus.Level{logrus.InfoLevel, logrus.WarnLevel}
		wantMessages := []string{"complete initialized info", "complete initialized warning"}
		wantKinds := []string{"info", "warn"}
		if len(wantLevels) == 0 || len(wantMessages) == 0 || len(wantKinds) == 0 {
			t.Fatal("initialized Collector expectation population is empty")
		}
		for index, entry := range entries {
			if entry == nil {
				t.Fatalf("initialized Collector entry %d is nil", index)
			}
			if entry.Logger != log {
				t.Fatalf("initialized Collector entry %d has logger %p, want exact package logger %p",
					index, entry.Logger, log)
			}
			if entry.Level != wantLevels[index] || entry.Message != wantMessages[index] ||
				entry.Data["kind"] != wantKinds[index] {
				t.Fatalf("initialized Collector entry %d was level=%s message=%q data=%#v",
					index, entry.Level, entry.Message, entry.Data)
			}
		}
	})
}

func TestDebugLoggerAndContextPreservePrivateLoggerEmptyFieldsAndStateSplit(t *testing.T) {
	runWithRestoredLoggerGlobals(t, "private logger helpers", func(t *testing.T) {
		packageLogger := logrus.New()
		packageOutput := &bytes.Buffer{}
		packageLogger.SetLevel(logrus.ErrorLevel)
		packageLogger.SetOutput(packageOutput)
		packageCollector := &Collector{}
		packageLogger.AddHook(packageCollector)
		log = packageLogger
		collector = packageCollector
		fieldLogger = true
		standard := logrus.StandardLogger()
		standardLevel := standard.GetLevel()

		debugFieldLogger := DebugLogger()
		debugEntry, ok := debugFieldLogger.(*logrus.Entry)
		if !ok {
			t.Fatalf("DebugLogger returned %T, want *logrus.Entry", debugFieldLogger)
		}
		assertEmptyPrivateEntry(t, "DebugLogger", debugEntry, packageLogger)
		callerFields := []string{"caller", "line", "func"}
		if len(callerFields) == 0 {
			t.Fatal("DebugLogger caller-field characterization population is empty")
		}
		for _, field := range callerFields {
			if _, found := debugEntry.Data[field]; found {
				t.Fatalf("DebugLogger unexpectedly populated caller-derived field %q: %#v", field, debugEntry.Data)
			}
		}
		if packageLogger.GetLevel() != logrus.DebugLevel {
			t.Fatalf("DebugLogger set private level to %s, want debug", packageLogger.GetLevel())
		}
		if standard.GetLevel() != standardLevel {
			t.Fatalf("DebugLogger changed standard logrus level from %s to %s",
				standardLevel, standard.GetLevel())
		}

		contextFieldLogger := Context()
		contextEntry, ok := contextFieldLogger.(*logrus.Entry)
		if !ok {
			t.Fatalf("Context returned %T, want *logrus.Entry", contextFieldLogger)
		}
		assertEmptyPrivateEntry(t, "Context", contextEntry, packageLogger)
		if packageLogger.GetLevel() != logrus.DebugLevel {
			t.Fatalf("Context changed private logger level to %s, want retained debug", packageLogger.GetLevel())
		}
		if !fieldLogger {
			t.Fatal("DebugLogger or Context changed the independent field-logger state")
		}
		if packageOutput.Len() != 0 || len(packageCollector.entries) != 0 {
			t.Fatalf("logger helpers emitted output or hook entries: output=%q entries=%#v",
				packageOutput.String(), packageCollector.entries)
		}
	})
}

func TestExternalErrorPreservesExactLegacyTextWithoutWrapping(t *testing.T) {
	tests := []struct {
		name    string
		source  error
		message string
		want    string
	}{
		{
			name:    "multiline external output",
			source:  errors.New("dependency failed: ø"),
			message: "STDOUT:\nfirst line\n\nSTDERR:\nsecond\\line",
			want: "dependency failed: ø\n\n##### EXTERNAL ERROR MESSAGE #####\n\n" +
				"  STDOUT:\n  first line\n  \n  STDERR:\n  second\\line\n" +
				"##### ENDS EXTERNAL ERROR MESSAGE #####\n\n",
		},
		{
			name:    "nil error and empty message",
			message: "",
			want: "<nil>\n\n##### EXTERNAL ERROR MESSAGE #####\n\n" +
				"  \n##### ENDS EXTERNAL ERROR MESSAGE #####\n\n",
		},
	}
	if len(tests) == 0 {
		t.Fatal("ExternalError characterization population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := ExternalError(test.source, test.message)
			if actual == nil || actual.Error() != test.want {
				t.Fatalf("ExternalError text was %q, want exact bytes %q", actual, test.want)
			}
			if test.source != nil && errors.Is(actual, test.source) {
				t.Fatal("ExternalError unexpectedly wrapped the source error")
			}
		})
	}
}

func TestLoggerConfigurationAndStdOutPreservePrivateStandardStateSplit(t *testing.T) {
	runWithRestoredLoggerGlobals(t, "configuration and standard output", func(t *testing.T) {
		packageLogger := logrus.New()
		packageOutput := &bytes.Buffer{}
		packageFormatter := &logrus.TextFormatter{DisableColors: true, DisableTimestamp: true}
		packageLogger.SetLevel(logrus.WarnLevel)
		packageLogger.SetOutput(packageOutput)
		packageLogger.SetFormatter(packageFormatter)
		log = packageLogger
		fieldLogger = false
		standard := logrus.StandardLogger()
		standardBeforeJSON := snapshotLogrusLoggerState(standard)

		if IsFieldLogger() {
			t.Fatal("field-logger state was true before SetFieldLogger")
		}
		SetFieldLogger()
		if !IsFieldLogger() {
			t.Fatal("SetFieldLogger did not transition field-logger state to true")
		}
		SetFieldLogger()
		if !IsFieldLogger() {
			t.Fatal("second SetFieldLogger call changed true state")
		}

		SetJsonLogging()
		jsonFormatter, ok := packageLogger.Formatter.(*logrus.JSONFormatter)
		if !ok {
			t.Fatalf("SetJsonLogging selected %T, want *logrus.JSONFormatter", packageLogger.Formatter)
		}
		if !reflect.DeepEqual(jsonFormatter, &logrus.JSONFormatter{}) {
			t.Fatalf("SetJsonLogging selected non-default JSON formatter: %#v", jsonFormatter)
		}
		assertLogrusLoggerState(t, "standard logger after SetJsonLogging", standard, standardBeforeJSON)
		if packageLogger.GetLevel() != logrus.WarnLevel || packageLogger.Out != packageOutput {
			t.Fatalf("SetJsonLogging changed private level/output: level=%s output=%T",
				packageLogger.GetLevel(), packageLogger.Out)
		}

		levels := []struct {
			level      logrus.Level
			wantStdout bool
		}{
			{level: logrus.PanicLevel},
			{level: logrus.FatalLevel},
			{level: logrus.ErrorLevel},
			{level: logrus.WarnLevel},
			{level: logrus.InfoLevel},
			{level: logrus.DebugLevel, wantStdout: true},
			{level: logrus.TraceLevel, wantStdout: true},
		}
		if len(levels) == 0 {
			t.Fatal("StdOut level characterization population is empty")
		}
		for _, test := range levels {
			logrus.SetLevel(test.level)
			actual := StdOut()
			var want *os.File
			if test.wantStdout {
				want = os.Stdout
			}
			if actual != want {
				t.Fatalf("StdOut at standard level %s returned %p, want exact pointer %p",
					test.level, actual, want)
			}
			if packageLogger.GetLevel() != logrus.WarnLevel {
				t.Fatalf("StdOut standard-level probe changed private level to %s", packageLogger.GetLevel())
			}
		}
		if packageOutput.Len() != 0 {
			t.Fatalf("configuration helpers wrote unexpected private log bytes %q", packageOutput.String())
		}
	})
}

func TestLogEntriesExposesExactCollectorSliceAndEntryPointers(t *testing.T) {
	runWithRestoredLoggerGlobals(t, "collector entry exposure", func(t *testing.T) {
		testLogger := logrus.New()
		first := logrus.NewEntry(testLogger)
		first.Level = logrus.InfoLevel
		first.Message = "complete first recorded entry"
		second := logrus.NewEntry(testLogger)
		second.Level = logrus.WarnLevel
		second.Message = "complete second recorded entry"
		recordingCollector := &Collector{entries: []*logrus.Entry{first, second}}
		collector = recordingCollector
		log = testLogger

		entries := LogEntries()
		if len(entries) == 0 {
			t.Fatal("LogEntries recorded-entry population is empty")
		}
		if len(entries) != 2 || cap(entries) != cap(recordingCollector.entries) {
			t.Fatalf("LogEntries returned len=%d cap=%d, want exact len=%d cap=%d",
				len(entries), cap(entries), len(recordingCollector.entries), cap(recordingCollector.entries))
		}
		want := []*logrus.Entry{first, second}
		for index, entry := range entries {
			if entry != want[index] {
				t.Fatalf("LogEntries entry %d identity is %p, want exact pointer %p", index, entry, want[index])
			}
		}
		if &entries[0] != &recordingCollector.entries[0] {
			t.Fatal("LogEntries copied the collector slice instead of exposing its exact backing storage")
		}

		replacement := logrus.NewEntry(testLogger)
		replacement.Message = "replacement through exposed slice"
		entries[0] = replacement
		if recordingCollector.entries[0] != replacement || LogEntries()[0] != replacement {
			t.Fatal("LogEntries result no longer shares the collector's exact entry-pointer storage")
		}
	})
}

func assertEmptyPrivateEntry(t *testing.T, name string, entry *logrus.Entry, wantLogger *logrus.Logger) {
	t.Helper()
	if entry == nil {
		t.Fatalf("%s returned a nil *logrus.Entry", name)
	}
	if entry.Logger != wantLogger {
		t.Fatalf("%s entry logger identity is %p, want exact private logger %p", name, entry.Logger, wantLogger)
	}
	if entry.Data == nil || len(entry.Data) != 0 {
		t.Fatalf("%s entry data was %#v, want an allocated empty field map", name, entry.Data)
	}
}

func runWithRestoredLoggerGlobals(t *testing.T, name string, action func(*testing.T)) {
	t.Helper()
	wantRestored := snapshotLoggerGlobals()
	t.Run(name, func(t *testing.T) {
		before := snapshotLoggerGlobals()
		t.Cleanup(func() { restoreLoggerGlobals(before) })
		action(t)
	})
	assertLoggerGlobals(t, wantRestored)
}

func snapshotLoggerGlobals() loggerGlobalsSnapshot {
	standard := logrus.StandardLogger()
	snapshot := loggerGlobalsSnapshot{
		fieldLogger:    fieldLogger,
		collector:      collector,
		packageLogger:  log,
		standardLogger: standard,
		standardState:  snapshotLogrusLoggerState(standard),
	}
	if collector != nil {
		snapshot.collectorEntries = collector.entries
	}
	if log != nil {
		snapshot.packageState = snapshotLogrusLoggerState(log)
	}
	return snapshot
}

func restoreLoggerGlobals(snapshot loggerGlobalsSnapshot) {
	fieldLogger = snapshot.fieldLogger
	collector = snapshot.collector
	if collector != nil {
		collector.entries = snapshot.collectorEntries
	}
	log = snapshot.packageLogger
	if log != nil {
		restoreLogrusLoggerState(log, snapshot.packageState)
	}
	restoreLogrusLoggerState(snapshot.standardLogger, snapshot.standardState)
}

func snapshotLogrusLoggerState(logger *logrus.Logger) logrusLoggerState {
	return logrusLoggerState{
		level:        logger.GetLevel(),
		formatter:    logger.Formatter,
		output:       logger.Out,
		hooks:        logger.Hooks,
		reportCaller: logger.ReportCaller,
	}
}

func restoreLogrusLoggerState(logger *logrus.Logger, state logrusLoggerState) {
	logger.SetLevel(state.level)
	logger.SetFormatter(state.formatter)
	logger.SetOutput(state.output)
	logger.ReplaceHooks(state.hooks)
	logger.SetReportCaller(state.reportCaller)
}

func assertLoggerGlobals(t *testing.T, want loggerGlobalsSnapshot) {
	t.Helper()
	if fieldLogger != want.fieldLogger {
		t.Fatalf("fieldLogger restored as %t, want %t", fieldLogger, want.fieldLogger)
	}
	if collector != want.collector {
		t.Fatalf("collector global identity is %p, want exact pointer %p", collector, want.collector)
	}
	assertEntrySliceIdentity(t, "collector entries", collector.entries, want.collectorEntries)
	if log != want.packageLogger {
		t.Fatalf("private log global identity is %p, want exact pointer %p", log, want.packageLogger)
	}
	assertLogrusLoggerState(t, "private logger", log, want.packageState)
	if logrus.StandardLogger() != want.standardLogger {
		t.Fatalf("standard logrus logger identity is %p, want exact pointer %p",
			logrus.StandardLogger(), want.standardLogger)
	}
	assertLogrusLoggerState(t, "standard logger", want.standardLogger, want.standardState)
}

func assertLogrusLoggerState(t *testing.T, name string, logger *logrus.Logger, want logrusLoggerState) {
	t.Helper()
	if logger.GetLevel() != want.level {
		t.Fatalf("%s level is %s, want %s", name, logger.GetLevel(), want.level)
	}
	if !sameInterfaceIdentity(logger.Formatter, want.formatter) {
		t.Fatalf("%s formatter identity is %T %#x, want %T %#x",
			name, logger.Formatter, interfacePointer(logger.Formatter), want.formatter, interfacePointer(want.formatter))
	}
	if !sameInterfaceIdentity(logger.Out, want.output) {
		t.Fatalf("%s output identity is %T %#x, want %T %#x",
			name, logger.Out, interfacePointer(logger.Out), want.output, interfacePointer(want.output))
	}
	assertHookIdentity(t, name, logger.Hooks, want.hooks)
	if logger.ReportCaller != want.reportCaller {
		t.Fatalf("%s ReportCaller is %t, want %t", name, logger.ReportCaller, want.reportCaller)
	}
}

func assertHookIdentity(t *testing.T, name string, actual, want logrus.LevelHooks) {
	t.Helper()
	if (actual == nil) != (want == nil) || len(actual) != len(want) {
		t.Fatalf("%s hooks were %#v, want exact topology %#v", name, actual, want)
	}
	for level, wantHooks := range want {
		actualHooks, found := actual[level]
		if !found || (actualHooks == nil) != (wantHooks == nil) || len(actualHooks) != len(wantHooks) {
			t.Fatalf("%s hooks at %s were %#v, want %#v", name, level, actualHooks, wantHooks)
		}
		for index, wantHook := range wantHooks {
			if !sameInterfaceIdentity(actualHooks[index], wantHook) {
				t.Fatalf("%s hook %d at %s is %T %#x, want %T %#x", name, index, level,
					actualHooks[index], interfacePointer(actualHooks[index]), wantHook, interfacePointer(wantHook))
			}
		}
	}
}

func assertEntrySliceIdentity(t *testing.T, name string, actual, want []*logrus.Entry) {
	t.Helper()
	if (actual == nil) != (want == nil) || len(actual) != len(want) || cap(actual) != cap(want) {
		t.Fatalf("%s were len=%d cap=%d nil=%t, want len=%d cap=%d nil=%t",
			name, len(actual), cap(actual), actual == nil, len(want), cap(want), want == nil)
	}
	if len(actual) > 0 && &actual[0] != &want[0] {
		t.Fatalf("%s backing identity is %p, want exact pointer %p", name, &actual[0], &want[0])
	}
	for index, entry := range want {
		if actual[index] != entry {
			t.Fatalf("%s entry %d identity is %p, want exact pointer %p", name, index, actual[index], entry)
		}
	}
}

func sameInterfaceIdentity(actual, want interface{}) bool {
	if actual == nil || want == nil {
		return actual == nil && want == nil
	}
	actualValue := reflect.ValueOf(actual)
	wantValue := reflect.ValueOf(want)
	if actualValue.Type() != wantValue.Type() {
		return false
	}
	return actualValue.Pointer() == wantValue.Pointer()
}

func interfacePointer(value interface{}) uintptr {
	if value == nil {
		return 0
	}
	return reflect.ValueOf(value).Pointer()
}
