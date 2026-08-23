package kibana

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/testutil"
)

const (
	testKibanaURL            = "https://kibana.example.test/internal/bsearch?compress=true"
	testKibanaGte            = "2022-11-27T21:37:48.801Z"
	testKibanaLte            = "2022-12-02T11:43:33.929Z"
	testKibanaAcceptLanguage = "nb-NO"
	testKibanaAuthorization  = "Bearer token-123"
	testKibanaContentType    = "application/json"
	testKibanaVersion        = "8.9.0"
)

func TestLoadFromFetchRequest(t *testing.T) {
	request, err := LoadFromFetchRequest(writeFetchFixture(t))
	if err != nil {
		t.Fatalf("load fetch request fixture: %v", err)
	}

	if request.Url != strings.Replace(testKibanaURL, "compress=true", "compress=false", 1) {
		t.Errorf("URL was %q", request.Url)
	}
	if request.AcceptLanguage != testKibanaAcceptLanguage {
		t.Errorf("accept-language was %q", request.AcceptLanguage)
	}
	if request.Authorization != testKibanaAuthorization {
		t.Errorf("authorization was %q", request.Authorization)
	}
	if request.ContentType != testKibanaContentType {
		t.Errorf("content-type was %q", request.ContentType)
	}
	if request.KbnVersion != testKibanaVersion {
		t.Errorf("kbn-version was %q", request.KbnVersion)
	}
	if request.Body != testKibanaBody() {
		t.Errorf("body was %q, want %q", request.Body, testKibanaBody())
	}

	interval, err := ExtractTimeIntervalFrom(request)
	if err != nil {
		t.Fatalf("extract interval from parsed request: %v", err)
	}
	if interval.Gte.Format(time.RFC3339Nano) != testKibanaGte || interval.Lte.Format(time.RFC3339Nano) != testKibanaLte {
		t.Errorf("parsed interval was %s to %s", interval.Gte.Format(time.RFC3339Nano), interval.Lte.Format(time.RFC3339Nano))
	}
}

func TestExecuteKibanaQueryRecordsRequestAndConvertsHits(t *testing.T) {
	request, err := LoadFromFetchRequest(writeFetchFixture(t))
	if err != nil {
		t.Fatalf("load fetch request fixture: %v", err)
	}

	requestedInterval := TimeInterval{
		Gte: mustParseTime(t, "2023-01-02T03:04:05Z"),
		Lte: mustParseTime(t, "2023-01-07T08:09:10Z"),
	}
	var recordedMethod string
	var recordedURL string
	var recordedHeader http.Header
	var recordedBody string
	stubKibanaTransport(t, func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("read request body: %w", err)
		}
		recordedMethod = req.Method
		recordedURL = req.URL.String()
		recordedHeader = req.Header.Clone()
		recordedBody = string(body)

		response := strings.Join([]string{
			`{"id":1,"result":{"rawResponse":{"hits":{"total":2}}}}`,
			`{"id":2,"result":{"rawResponse":{"took":7,"hits":{"hits":[{"fields":{"@timestamp":["2023-01-06T10:00:00Z"],"service.name":["orders-api"],"owner.team":["payments"]}},{"fields":{"@timestamp":["2023-01-03T09:00:00Z"],"service.name":["billing-api"],"owner.team":["finance"]}}]}}}}`,
		}, "\n")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(response)),
		}, nil
	})

	filter := map[string]string{
		"application": "service.name",
		"team":        "owner.team",
	}
	err, _, result, resultExists := ExecuteKibanaQuery(request, requestedInterval, filter, map[string]bool{}, "")
	if err != nil {
		t.Fatalf("execute Kibana query: %v", err)
	}

	if recordedMethod != http.MethodPost {
		t.Errorf("request method was %q, want POST", recordedMethod)
	}
	wantURL := strings.Replace(testKibanaURL, "compress=true", "compress=false", 1)
	if recordedURL != wantURL {
		t.Errorf("request URL was %q, want %q", recordedURL, wantURL)
	}
	assertHeader(t, recordedHeader, "accept-language", testKibanaAcceptLanguage)
	assertHeader(t, recordedHeader, "authorization", testKibanaAuthorization)
	assertHeader(t, recordedHeader, "content-type", testKibanaContentType)
	assertHeader(t, recordedHeader, "kbn-version", testKibanaVersion)

	if !strings.Contains(recordedBody, `"size":500`) {
		t.Errorf("request body did not contain rewritten result size: %s", recordedBody)
	}
	if strings.Contains(recordedBody, `"size":10`) {
		t.Errorf("request body retained original result size: %s", recordedBody)
	}
	recordedInterval, err := ExtractTimeIntervalFrom(KibanaFetchRequest{Body: recordedBody})
	if err != nil {
		t.Fatalf("extract recorded request interval: %v", err)
	}
	if !recordedInterval.Gte.Equal(requestedInterval.Gte) || !recordedInterval.Lte.Equal(requestedInterval.Lte) {
		t.Errorf("recorded interval was %s to %s, want %s to %s", recordedInterval.Gte, recordedInterval.Lte, requestedInterval.Gte, requestedInterval.Lte)
	}

	wantResults := map[string]bool{
		`{"application":"orders-api","team":"payments"}`: true,
		`{"application":"billing-api","team":"finance"}`: true,
	}
	if len(result) != len(wantResults) {
		t.Fatalf("converted %d hits, want %d: %v", len(result), len(wantResults), result)
	}
	for _, hit := range result {
		if !wantResults[hit] {
			t.Errorf("unexpected converted hit %s", hit)
		}
	}
	if len(resultExists) != len(wantResults) {
		t.Errorf("result existence set contained %d entries, want %d", len(resultExists), len(wantResults))
	}
}

type kibanaRoundTripFunc func(*http.Request) (*http.Response, error)

func (f kibanaRoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func stubKibanaTransport(t *testing.T, transport kibanaRoundTripFunc) {
	t.Helper()
	previous := http.DefaultTransport
	http.DefaultTransport = transport
	t.Cleanup(func() {
		http.DefaultTransport = previous
	})
}

func writeFetchFixture(t *testing.T) string {
	t.Helper()
	fixture := fmt.Sprintf(`fetch(%q, {
  "headers": {
    "accept-language": %q,
    "authorization": %q,
    "content-type": %q,
    "kbn-version": %q
  },
  "body": %s,
})`, testKibanaURL, testKibanaAcceptLanguage, testKibanaAuthorization, testKibanaContentType, testKibanaVersion, strconv.Quote(testKibanaBody()))
	path := filepath.Join(t.TempDir(), "request.fetch")
	if err := testutil.WriteFileOutsideWorkingTree(path, []byte(fixture), 0o600); err != nil {
		t.Fatalf("write fetch request fixture: %v", err)
	}
	return path
}

func testKibanaBody() string {
	return `{"size":10,"query":{"range":{"@timestamp":{"gte":"` + testKibanaGte + `","lte":"` + testKibanaLte + `"}}}}`
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		t.Fatalf("parse test timestamp %q: %v", value, err)
	}
	return parsed
}

func assertHeader(t *testing.T, header http.Header, name, expected string) {
	t.Helper()
	if actual := header.Get(name); actual != expected {
		t.Errorf("header %s was %q, want %q", name, actual, expected)
	}
}
