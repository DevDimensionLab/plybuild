package kibana

import (
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/clock"
	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
)

type recordingKibanaHTTP struct {
	requests  []httpclient.Request
	responses []*http.Response
	err       error
}

func (recording *recordingKibanaHTTP) dependencies() httpclient.Dependencies {
	return httpclient.Dependencies{Client: recording}
}

func (recording *recordingKibanaHTTP) Do(request httpclient.Request) (*http.Response, error) {
	if request.BasicAuth != nil {
		copied := *request.BasicAuth
		request.BasicAuth = &copied
	}
	if request.BearerJSON != nil {
		copied := *request.BearerJSON
		request.BearerJSON = &copied
	}
	if request.POST != nil {
		copied := *request.POST
		copied.Body = append([]byte(nil), request.POST.Body...)
		copied.Header = request.POST.Header.Clone()
		request.POST = &copied
	}
	recording.requests = append(recording.requests, request)
	if recording.err != nil {
		return nil, recording.err
	}
	index := len(recording.requests) - 1
	if index >= len(recording.responses) {
		return nil, errors.New("recorded Kibana response population is empty")
	}
	return recording.responses[index], nil
}

func (recording *recordingKibanaHTTP) assertedRequests() ([]httpclient.Request, error) {
	if len(recording.requests) == 0 {
		return nil, errors.New("recorded Kibana POST request population is empty")
	}
	return recording.requests, nil
}

type recordedKibanaRetryResult struct {
	err      error
	response KibanaResponse
}

type recordingKibanaRetryPOST struct {
	requests []KibanaFetchRequest
	results  []recordedKibanaRetryResult
	sequence *[]string
}

func (recording *recordingKibanaRetryPOST) dependencies(
	dependencies postDependencies,
) postDependencies {
	dependencies.InternalPOST = recording.Post
	return dependencies
}

//nolint:staticcheck // The recorder preserves the existing error-first internalPOST contract.
func (recording *recordingKibanaRetryPOST) Post(request KibanaFetchRequest) (error, KibanaResponse) {
	recording.requests = append(recording.requests, request)
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "post")
	}
	index := len(recording.requests) - 1
	if index >= len(recording.results) {
		return errors.New("recorded Kibana retry result population is empty"), KibanaResponse{}
	}
	return recording.results[index].err, recording.results[index].response
}

func (recording *recordingKibanaRetryPOST) assertedRequests() ([]KibanaFetchRequest, error) {
	if len(recording.requests) == 0 {
		return nil, errors.New("recorded Kibana retry POST population is empty")
	}
	return recording.requests, nil
}

type recordingKibanaRetrySleeper struct {
	durations []time.Duration
	sequence  *[]string
}

func (recording *recordingKibanaRetrySleeper) dependencies(
	dependencies postDependencies,
) postDependencies {
	dependencies.Clock.Sleeper = recording
	return dependencies
}

func (recording *recordingKibanaRetrySleeper) Sleep(duration time.Duration) {
	recording.durations = append(recording.durations, duration)
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "sleep")
	}
}

func (recording *recordingKibanaRetrySleeper) assertedDurations() ([]time.Duration, error) {
	if len(recording.durations) == 0 {
		return nil, errors.New("recorded Kibana retry sleep population is empty")
	}
	return recording.durations, nil
}

type recordingKibanaRetryClock struct {
	now time.Time
}

func (recording *recordingKibanaRetryClock) Now() time.Time {
	return recording.now
}

type recordingKibanaQueryPOST struct {
	requests  []KibanaFetchRequest
	responses []KibanaResponse
	err       error
}

func (recording *recordingKibanaQueryPOST) dependencies() queryDependencies {
	return queryDependencies{Client: recording}
}

func (recording *recordingKibanaQueryPOST) Post(request KibanaFetchRequest) (KibanaResponse, error) {
	recording.requests = append(recording.requests, request)
	if recording.err != nil {
		return KibanaResponse{}, recording.err
	}
	index := len(recording.requests) - 1
	if index >= len(recording.responses) {
		return KibanaResponse{}, errors.New("recorded Kibana query response population is empty")
	}
	return recording.responses[index], nil
}

func (recording *recordingKibanaQueryPOST) assertedRequests() ([]KibanaFetchRequest, error) {
	if len(recording.requests) == 0 {
		return nil, errors.New("recorded Kibana query POST population is empty")
	}
	return recording.requests, nil
}

type trackingKibanaBody struct {
	reader io.Reader
	reads  int
	closed bool
}

func (body *trackingKibanaBody) Read(buffer []byte) (int, error) {
	body.reads++
	return body.reader.Read(buffer)
}

func (body *trackingKibanaBody) Close() error {
	body.closed = true
	return nil
}

type kibanaErrorReader struct {
	err error
}

func (reader kibanaErrorReader) Read([]byte) (int, error) {
	return 0, reader.err
}

func TestInternalPOSTPassesCompleteRewrittenRequestAndParsesNonSuccessResponse(t *testing.T) {
	body := &trackingKibanaBody{reader: strings.NewReader(strings.Join([]string{
		`{"id":41,"result":{"id":"header-complete","rawResponse":{"took":17,"timed_out":true,"_shards":{"total":5,"successful":4,"skipped":3,"failed":2},"hits":{"total":23}},"isPartial":true,"isRunning":true,"warning":"header-warning","total":29,"loaded":19,"isRestored":true}}`,
		`{"id":73,"result":{"id":"result-complete","rawResponse":{"took":31,"timed_out":true,"_shards":{"total":11,"successful":10,"skipped":9,"failed":8},"hits":{"hits":[{"_index":"complete-index","_type":"complete-type","_id":"complete-hit","_version":7,"fields":{"@timestamp":["2023-01-04T05:06:07Z"],"complete.field":["complete-value"]},"sort":[123,456]}]}},"isPartial":true,"isRunning":true,"warning":"result-warning","total":37,"loaded":33,"isRestored":true}}`,
	}, "\n"))}
	recording := &recordingKibanaHTTP{responses: []*http.Response{{
		StatusCode: http.StatusServiceUnavailable,
		Status:     "503 Must Still Be Parsed",
		Body:       body,
	}}}
	request := completeKibanaFetchRequest()
	request.Body = `{"size":10,"nested":{"size":11},"query":{"complete":true}}`

	err, response := internalPost(recording.dependencies(), request)

	if err != nil {
		t.Fatalf("internal POST returned an error: %v", err)
	}
	requests, assertErr := recording.assertedRequests()
	if assertErr != nil {
		t.Fatal(assertErr)
	}
	want := httpclient.Request{
		URL: request.Url,
		POST: &httpclient.POST{
			Body: []byte(`{"size":500,"nested":{"size":11},"query":{"complete":true}}`),
			Header: http.Header{
				"Accept-Language": []string{request.AcceptLanguage},
				"Authorization":   []string{request.Authorization},
				"Content-Type":    []string{request.ContentType},
				"Kbn-Version":     []string{request.KbnVersion},
			},
		},
	}
	if len(requests) != 1 || !reflect.DeepEqual(requests[0], want) {
		t.Fatalf("recorded Kibana request differs:\n got: %#v\nwant: %#v", requests, want)
	}
	if !body.closed || body.reads == 0 {
		t.Fatalf("Kibana response lifecycle was closed=%t reads=%d", body.closed, body.reads)
	}
	header := response.KibanaResponseHeader
	if header.ID != 41 || header.Result.ID != "header-complete" ||
		header.Result.RawResponse.Took != 17 || !header.Result.RawResponse.TimedOut ||
		header.Result.RawResponse.Shards.Total != 5 || header.Result.RawResponse.Shards.Successful != 4 ||
		header.Result.RawResponse.Shards.Skipped != 3 || header.Result.RawResponse.Shards.Failed != 2 ||
		header.Result.RawResponse.Hits.Total != 23 || !header.Result.IsPartial || !header.Result.IsRunning ||
		header.Result.Warning != "header-warning" || header.Result.Total != 29 || header.Result.Loaded != 19 ||
		!header.Result.IsRestored {
		t.Fatalf("parsed Kibana response header was incomplete: %#v", header)
	}
	result := response.KibanaResult
	if result.ID != 73 || result.Result.ID != "result-complete" || result.Result.RawResponse.Took != 31 ||
		!result.Result.RawResponse.TimedOut || result.Result.RawResponse.Shards.Total != 11 ||
		result.Result.RawResponse.Shards.Successful != 10 || result.Result.RawResponse.Shards.Skipped != 9 ||
		result.Result.RawResponse.Shards.Failed != 8 || len(result.Result.RawResponse.Hits.Hits) != 1 ||
		!result.Result.IsPartial || !result.Result.IsRunning || result.Result.Warning != "result-warning" ||
		result.Result.Total != 37 || result.Result.Loaded != 33 || !result.Result.IsRestored {
		t.Fatalf("parsed Kibana result was incomplete: %#v", result)
	}
	hit := result.Result.RawResponse.Hits.Hits[0]
	if hit.Index != "complete-index" || hit.Type != "complete-type" || hit.ID != "complete-hit" ||
		hit.Version != 7 || !reflect.DeepEqual(hit.Sort, []int64{123, 456}) ||
		!reflect.DeepEqual(hit.Fields["complete.field"], []interface{}{"complete-value"}) {
		t.Fatalf("parsed Kibana hit was incomplete: %#v", hit)
	}
}

func TestInternalPOSTPreservesDependencyReadAndJSONErrorsAfterRecording(t *testing.T) {
	t.Run("dependency", func(t *testing.T) {
		sentinel := errors.New("Kibana HTTP dependency failed")
		recording := &recordingKibanaHTTP{err: sentinel}

		err, _ := internalPost(recording.dependencies(), completeKibanaFetchRequest())

		if !errors.Is(err, sentinel) {
			t.Fatalf("dependency error was %v, want %v", err, sentinel)
		}
		if _, assertErr := recording.assertedRequests(); assertErr != nil {
			t.Fatal(assertErr)
		}
	})

	t.Run("read", func(t *testing.T) {
		sentinel := errors.New("Kibana response read failed")
		body := &trackingKibanaBody{reader: kibanaErrorReader{err: sentinel}}
		recording := &recordingKibanaHTTP{responses: []*http.Response{{
			StatusCode: http.StatusOK,
			Body:       body,
		}}}

		err, _ := internalPost(recording.dependencies(), completeKibanaFetchRequest())

		if !errors.Is(err, sentinel) {
			t.Fatalf("read error was %v, want %v", err, sentinel)
		}
		if !body.closed || body.reads == 0 {
			t.Fatalf("read-error lifecycle was closed=%t reads=%d", body.closed, body.reads)
		}
	})

	t.Run("header unmarshal", func(t *testing.T) {
		assertKibanaUnmarshalError(t, `{"id":`+"\n"+`{"id":2}`, "header")
	})

	t.Run("result unmarshal", func(t *testing.T) {
		assertKibanaUnmarshalError(t, `{"id":1}`+"\n"+`{"id":`, "result")
	})
}

func TestInternalPOSTDependenciesDefaultToSafeNoRequest(t *testing.T) {
	err, response := internalPost(httpclient.Dependencies{}, completeKibanaFetchRequest())

	if !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("safe Kibana POST default returned %v, want %v", err, httpclient.ErrNoClient)
	}
	if !reflect.DeepEqual(response, KibanaResponse{}) {
		t.Fatalf("safe Kibana POST default returned response %#v", response)
	}
}

func TestKibanaQueryDependenciesDefaultToSafeNoPOST(t *testing.T) {
	request := completeKibanaFetchRequest()
	interval, err := ExtractTimeIntervalFrom(request)
	if err != nil {
		t.Fatalf("extract safe-default interval: %v", err)
	}

	err, response, result, resultExists := executeKibanaQuery(
		queryDependencies{}, request, interval, nil, map[string]bool{}, "",
	)

	if !errors.Is(err, httpclient.ErrNoClient) {
		t.Fatalf("safe Kibana query default returned %v, want %v", err, httpclient.ErrNoClient)
	}
	if !reflect.DeepEqual(response, KibanaResponse{}) || result != nil || len(resultExists) != 0 {
		t.Fatalf("safe Kibana query default returned response=%#v result=%#v exists=%#v", response, result, resultExists)
	}
}

func TestPOSTSelectsExactInternalPOSTAndCompleteSystemClock(t *testing.T) {
	dependencies := systemPostDependencies()
	systemClock := clock.System()

	if dependencies.InternalPOST == nil {
		t.Fatal("Kibana retry selected no internalPOST operation")
	}
	if reflect.ValueOf(dependencies.InternalPOST).Pointer() != reflect.ValueOf(internalPOST).Pointer() {
		t.Fatal("Kibana retry no longer selects the exact existing internalPOST operation")
	}
	if dependencies.Clock.Clock == nil || dependencies.Clock.Sleeper == nil {
		t.Fatalf("Kibana retry selected an incomplete system clock: %#v", dependencies.Clock)
	}
	if reflect.TypeOf(dependencies.Clock.Clock) != reflect.TypeOf(systemClock.Clock) ||
		reflect.TypeOf(dependencies.Clock.Sleeper) != reflect.TypeOf(systemClock.Sleeper) {
		t.Fatalf("Kibana retry clock dependency is %#v, want exact complete system clock %#v",
			dependencies.Clock, systemClock)
	}
	assertPublicPOSTComposition(t)
	assertPrivatePOSTRetryComposition(t)
}

func TestPOSTRecordingDoublesPreserveCompleteCallerDependencies(t *testing.T) {
	callerOperation := internalPOST
	callerClock := &recordingKibanaRetryClock{now: time.Unix(-73, 456789123)}
	callerSleeper := &recordingKibanaRetrySleeper{}
	initial := postDependencies{
		InternalPOST: callerOperation,
		Clock: clock.Dependencies{
			Clock:   callerClock,
			Sleeper: callerSleeper,
		},
	}
	postRecording := &recordingKibanaRetryPOST{}
	sleepRecording := &recordingKibanaRetrySleeper{}

	withPOST := postRecording.dependencies(initial)
	withSleep := sleepRecording.dependencies(initial)

	if withPOST.InternalPOST == nil ||
		reflect.ValueOf(withPOST.InternalPOST).Pointer() != reflect.ValueOf(postRecording.Post).Pointer() ||
		withPOST.Clock.Clock != callerClock || withPOST.Clock.Sleeper != callerSleeper {
		t.Fatalf("Kibana retry POST recorder lost the complete caller-owned dependency: %#v", withPOST)
	}
	if withSleep.InternalPOST == nil ||
		reflect.ValueOf(withSleep.InternalPOST).Pointer() != reflect.ValueOf(callerOperation).Pointer() ||
		withSleep.Clock.Clock != callerClock || withSleep.Clock.Sleeper != sleepRecording {
		t.Fatalf("Kibana retry sleep recorder lost the complete caller-owned dependency: %#v", withSleep)
	}
}

func TestPOSTReturnsExactFirstResultWithoutSleepOrRetryForNonEmptyHits(t *testing.T) {
	request := completeKibanaFetchRequest()
	firstError := errors.New("complete first Kibana retry result error")
	firstResponse := kibanaResponseForTimestamps(t, []string{"2023-01-04T05:06:07.890123456Z"})
	sequence := []string{}
	postRecording := &recordingKibanaRetryPOST{
		results:  []recordedKibanaRetryResult{{err: firstError, response: firstResponse}},
		sequence: &sequence,
	}
	sleepRecording := &recordingKibanaRetrySleeper{sequence: &sequence}
	dependencies := sleepRecording.dependencies(postRecording.dependencies(postDependencies{}))

	err, response := post(dependencies, request)

	if err != firstError || !reflect.DeepEqual(response, firstResponse) {
		t.Fatalf("Kibana non-empty first result was (%v, %#v), want exact (%v, %#v)",
			err, response, firstError, firstResponse)
	}
	requests, populationErr := postRecording.assertedRequests()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(requests, []KibanaFetchRequest{request}) {
		t.Fatalf("Kibana non-empty retry requests were %#v, want one exact request %#v",
			requests, []KibanaFetchRequest{request})
	}
	if len(sleepRecording.durations) != 0 {
		t.Fatalf("Kibana non-empty result slept with durations %#v, want none", sleepRecording.durations)
	}
	if !reflect.DeepEqual(sequence, []string{"post"}) {
		t.Fatalf("Kibana non-empty result effect order was %#v, want one first request", sequence)
	}
}

func TestPOSTPrintsSleepsAndRetriesOnceAfterEmptyHitsDespiteFirstError(t *testing.T) {
	request := completeKibanaFetchRequest()
	firstError := errors.New("complete discarded first Kibana retry error")
	secondError := errors.New("complete returned second Kibana retry error")
	tests := []struct {
		name           string
		secondResponse KibanaResponse
	}{
		{
			name:           "non-empty second result",
			secondResponse: kibanaResponseForTimestamps(t, []string{"2023-01-04T05:06:07Z"}),
		},
		{
			name:           "empty second result never causes a third request",
			secondResponse: KibanaResponse{},
		},
	}

	if len(tests) == 0 {
		t.Fatal("TestPOSTPrintsSleepsAndRetriesOnceAfterEmptyHitsDespiteFirstError test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sequence := []string{}
			postRecording := &recordingKibanaRetryPOST{
				results: []recordedKibanaRetryResult{
					{err: firstError, response: KibanaResponse{}},
					{err: secondError, response: test.secondResponse},
				},
				sequence: &sequence,
			}
			sleepRecording := &recordingKibanaRetrySleeper{sequence: &sequence}
			dependencies := sleepRecording.dependencies(postRecording.dependencies(postDependencies{}))

			err, response := post(dependencies, request)

			if err != secondError || !reflect.DeepEqual(response, test.secondResponse) {
				t.Fatalf("Kibana retry result was (%v, %#v), want exact second (%v, %#v)",
					err, response, secondError, test.secondResponse)
			}
			requests, populationErr := postRecording.assertedRequests()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			wantRequests := []KibanaFetchRequest{request, request}
			if !reflect.DeepEqual(requests, wantRequests) {
				t.Fatalf("Kibana retry requests were %#v, want two exact requests %#v", requests, wantRequests)
			}
			durations, populationErr := sleepRecording.assertedDurations()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(durations, []time.Duration{15 * time.Second}) {
				t.Fatalf("Kibana retry sleep durations were %#v, want one exact 15-second duration", durations)
			}
			if !reflect.DeepEqual(sequence, []string{"post", "sleep", "post"}) {
				t.Fatalf("Kibana retry effect order was %#v, want request, sleep, request", sequence)
			}
		})
	}
}

func TestPOSTDependenciesDefaultToSafeNoRequestOrRealSleep(t *testing.T) {
	err, response := post(postDependencies{}, completeKibanaFetchRequest())

	if !errors.Is(err, httpclient.ErrNoClient) || !reflect.DeepEqual(response, KibanaResponse{}) {
		t.Fatalf("safe Kibana retry default returned (%v, %#v), want (%v, empty response)",
			err, response, httpclient.ErrNoClient)
	}
}

func TestPOSTSystemPreservesFreshClientRedirects(t *testing.T) {
	previous := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = previous })
	redirectBody := &trackingKibanaBody{reader: strings.NewReader("redirect")}
	finalBody := &trackingKibanaBody{reader: strings.NewReader(kibanaResponseWithTimestamps(t, []string{"2023-01-04T05:06:07Z"}))}
	requests := []*http.Request{}
	requestBodies := []string{}
	http.DefaultTransport = kibanaRoundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request)
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		requestBodies = append(requestBodies, string(body))
		if len(requests) == 1 {
			return &http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Status:     "307 Temporary Redirect",
				Header:     http.Header{"Location": []string{"https://kibana.example.invalid/final"}},
				Body:       redirectBody,
			}, nil
		}
		return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Body: finalBody}, nil
	})
	request := completeKibanaFetchRequest()
	request.Url = "https://kibana.example.invalid/start"

	err, response := POST(request)

	if err != nil {
		t.Fatalf("redirected Kibana POST returned an error: %v", err)
	}
	if len(response.KibanaResult.Result.RawResponse.Hits.Hits) != 1 {
		t.Fatalf("redirected Kibana result was incomplete: %#v", response)
	}
	if len(requests) != 2 {
		t.Fatalf("Kibana redirect request population was %d, want 2", len(requests))
	}
	for index, recorded := range requests {
		if recorded.Method != http.MethodPost || !strings.Contains(requestBodies[index], `"size":500`) {
			t.Fatalf("Kibana redirect request %d was %s with body %q", index, recorded.Method, requestBodies[index])
		}
		assertHeader(t, recorded.Header, "accept-language", testKibanaAcceptLanguage)
		assertHeader(t, recorded.Header, "authorization", testKibanaAuthorization)
		assertHeader(t, recorded.Header, "content-type", testKibanaContentType)
		assertHeader(t, recorded.Header, "kbn-version", testKibanaVersion)
	}
	if !redirectBody.closed || !finalBody.closed {
		t.Fatalf("Kibana redirect bodies were closed=(%t, %t)", redirectBody.closed, finalBody.closed)
	}
}

func TestExecuteKibanaQueryKeepsPOSTDependencyAcrossRecursiveIntervals(t *testing.T) {
	firstTimestamps := make([]string, kibanaMaxResult)
	for index := range firstTimestamps {
		firstTimestamps[index] = "2023-01-06T10:00:00Z"
	}
	firstTimestamps[len(firstTimestamps)-1] = "2023-01-03T09:00:00Z"
	secondTimestamp := "2023-01-02T12:00:00Z"
	recording := &recordingKibanaQueryPOST{responses: []KibanaResponse{
		kibanaResponseForTimestamps(t, firstTimestamps),
		kibanaResponseForTimestamps(t, []string{secondTimestamp}),
	}}
	request := completeKibanaFetchRequest()
	requestedInterval := TimeInterval{
		Gte: mustParseTime(t, "2023-01-02T03:04:05Z"),
		Lte: mustParseTime(t, "2023-01-07T08:09:10Z"),
	}

	err, _, _, _ := executeKibanaQuery(
		recording.dependencies(), request, requestedInterval, nil, map[string]bool{}, "",
	)

	if err != nil {
		t.Fatalf("recursive Kibana query returned an error: %v", err)
	}
	requests, assertErr := recording.assertedRequests()
	if assertErr != nil {
		t.Fatal(assertErr)
	}
	if len(requests) != 2 {
		t.Fatalf("recursive Kibana request population was %d, want 2", len(requests))
	}
	firstInterval, err := ExtractTimeIntervalFrom(requests[0])
	if err != nil {
		t.Fatalf("extract first recursive request interval: %v", err)
	}
	secondInterval, err := ExtractTimeIntervalFrom(requests[1])
	if err != nil {
		t.Fatalf("extract second recursive request interval: %v", err)
	}
	wantSecondLte := mustParseTime(t, firstTimestamps[len(firstTimestamps)-1])
	if !firstInterval.Gte.Equal(requestedInterval.Gte) || !firstInterval.Lte.Equal(requestedInterval.Lte) ||
		!secondInterval.Gte.Equal(requestedInterval.Gte) || !secondInterval.Lte.Equal(wantSecondLte) {
		t.Fatalf("recursive intervals were (%s, %s) then (%s, %s)",
			firstInterval.Gte, firstInterval.Lte, secondInterval.Gte, secondInterval.Lte)
	}
	for index, recorded := range requests {
		if recorded.Url != request.Url || recorded.AcceptLanguage != request.AcceptLanguage ||
			recorded.Authorization != request.Authorization || recorded.ContentType != request.ContentType ||
			recorded.KbnVersion != request.KbnVersion {
			t.Fatalf("recursive Kibana caller request %d was incomplete: %#v", index, recorded)
		}
	}
}

func TestExecuteKibanaQueryReturnsPOSTErrorAfterCompleteCallerDelivery(t *testing.T) {
	sentinel := errors.New("Kibana query POST dependency failed")
	recording := &recordingKibanaQueryPOST{err: sentinel}
	request := completeKibanaFetchRequest()
	interval := TimeInterval{
		Gte: mustParseTime(t, "2023-02-03T04:05:06Z"),
		Lte: mustParseTime(t, "2023-02-07T08:09:10Z"),
	}

	err, _, _, _ := executeKibanaQuery(
		recording.dependencies(), request, interval, map[string]string{"complete": "field"}, map[string]bool{}, "",
	)

	if !errors.Is(err, sentinel) {
		t.Fatalf("Kibana query POST error was %v, want %v", err, sentinel)
	}
	requests, assertErr := recording.assertedRequests()
	if assertErr != nil {
		t.Fatal(assertErr)
	}
	want := CreateRequestForInterval(interval, request)
	if len(requests) != 1 || !reflect.DeepEqual(requests[0], want) {
		t.Fatalf("Kibana caller dependency received an incomplete request:\n got: %#v\nwant: %#v", requests, want)
	}
}

func assertPublicPOSTComposition(t *testing.T) {
	t.Helper()
	function := parsedKibanaFunction(t, "POST")
	if function.Type.Params == nil || len(function.Type.Params.List) != 1 ||
		function.Type.Results == nil || len(function.Type.Results.List) != 2 ||
		function.Body == nil || len(function.Body.List) != 1 {
		t.Fatal("public Kibana POST signature or single-return composition changed")
	}
	returned, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returned.Results) != 1 {
		t.Fatal("public Kibana POST no longer directly returns one private composition")
	}
	call, ok := returned.Results[0].(*ast.CallExpr)
	if !ok || len(call.Args) != 2 {
		t.Fatal("public Kibana POST no longer returns post with dependencies and exact request")
	}
	callee, calleeOK := call.Fun.(*ast.Ident)
	selector, selectorOK := call.Args[0].(*ast.CallExpr)
	var selectorName *ast.Ident
	selectorNameOK := false
	if selectorOK {
		selectorName, selectorNameOK = selector.Fun.(*ast.Ident)
	}
	request, requestOK := call.Args[1].(*ast.Ident)
	if !calleeOK || callee.Name != "post" || !selectorOK || !selectorNameOK ||
		selectorName.Name != "systemPostDependencies" || len(selector.Args) != 0 ||
		!requestOK || request.Name != "reguest" {
		t.Fatal("public Kibana POST no longer selects exact private production dependencies beside its request")
	}
}

func assertPrivatePOSTRetryComposition(t *testing.T) {
	t.Helper()
	retryFunction := parsedKibanaFunction(t, "post")
	if retryFunction.Body == nil || len(retryFunction.Body.List) != 3 {
		t.Fatal("private Kibana POST retry composition no longer has request, selection, and first return")
	}
	assertKibanaRetryPOSTAssignment(t, retryFunction.Body.List[0])

	retry, ok := retryFunction.Body.List[1].(*ast.IfStmt)
	if !ok || retry.Body == nil || len(retry.Body.List) != 4 {
		t.Fatal("private Kibana POST retry selection no longer has print, sleep, request, and return")
	}
	zeroHitSelection := false
	ast.Inspect(retry.Cond, func(node ast.Node) bool {
		comparison, ok := node.(*ast.BinaryExpr)
		if !ok || comparison.Op != token.EQL {
			return true
		}
		length, lengthOK := comparison.X.(*ast.CallExpr)
		zero, zeroOK := comparison.Y.(*ast.BasicLit)
		var identifier *ast.Ident
		identifierOK := false
		if lengthOK {
			identifier, identifierOK = length.Fun.(*ast.Ident)
		}
		zeroHitSelection = lengthOK && zeroOK && identifierOK &&
			identifier.Name == "len" && len(length.Args) == 1 && zero.Value == "0"
		return !zeroHitSelection
	})
	if !zeroHitSelection {
		t.Fatal("private Kibana POST retry selection is no longer based only on zero hit-list length")
	}

	printStatement, ok := retry.Body.List[0].(*ast.ExprStmt)
	if !ok {
		t.Fatal("Kibana retry print is no longer the first branch operation")
	}
	printCall, ok := printStatement.X.(*ast.CallExpr)
	var printName *ast.Ident
	printNameOK := false
	var printText *ast.BasicLit
	printTextOK := false
	if ok {
		printName, printNameOK = printCall.Fun.(*ast.Ident)
		if len(printCall.Args) == 1 {
			printText, printTextOK = printCall.Args[0].(*ast.BasicLit)
		}
	}
	if !ok || !printNameOK || printName.Name != "println" || len(printCall.Args) != 1 ||
		!printTextOK || printText.Value != strconv.Quote("sleep and retry") {
		t.Fatal("Kibana retry no longer prints the exact sleep and retry text first")
	}

	sleepStatement, ok := retry.Body.List[1].(*ast.ExprStmt)
	if !ok {
		t.Fatal("Kibana retry sleep is no longer the second branch operation")
	}
	sleepCall, ok := sleepStatement.X.(*ast.CallExpr)
	var sleepSelector *ast.SelectorExpr
	sleepSelectorOK := false
	if ok {
		sleepSelector, sleepSelectorOK = sleepCall.Fun.(*ast.SelectorExpr)
	}
	var sleepPackage *ast.Ident
	sleepPackageOK := false
	if sleepSelectorOK {
		sleepPackage, sleepPackageOK = sleepSelector.X.(*ast.Ident)
	}
	var clockValue *ast.SelectorExpr
	clockValueOK := false
	var duration *ast.BinaryExpr
	durationOK := false
	if ok && len(sleepCall.Args) == 2 {
		clockValue, clockValueOK = sleepCall.Args[0].(*ast.SelectorExpr)
		duration, durationOK = sleepCall.Args[1].(*ast.BinaryExpr)
	}
	var dependencyName *ast.Ident
	dependencyNameOK := false
	if clockValueOK {
		dependencyName, dependencyNameOK = clockValue.X.(*ast.Ident)
	}
	var fifteen *ast.BasicLit
	fifteenOK := false
	var second *ast.SelectorExpr
	secondOK := false
	if durationOK {
		fifteen, fifteenOK = duration.X.(*ast.BasicLit)
		second, secondOK = duration.Y.(*ast.SelectorExpr)
	}
	var timePackage *ast.Ident
	timePackageOK := false
	if secondOK {
		timePackage, timePackageOK = second.X.(*ast.Ident)
	}
	if !ok || !sleepSelectorOK || !sleepPackageOK || sleepPackage.Name != "clock" ||
		sleepSelector.Sel.Name != "Sleep" || len(sleepCall.Args) != 2 || !clockValueOK ||
		!dependencyNameOK || dependencyName.Name != "dependencies" || clockValue.Sel.Name != "Clock" ||
		!durationOK || duration.Op != token.MUL || !fifteenOK || fifteen.Value != "15" ||
		!secondOK || !timePackageOK || timePackage.Name != "time" || second.Sel.Name != "Second" {
		t.Fatal("Kibana retry no longer delegates the exact 15*time.Second duration through its clock")
	}

	assertKibanaRetryPOSTAssignment(t, retry.Body.List[2])
	assertKibanaErrorResponseReturn(t, retry.Body.List[3])
	assertKibanaErrorResponseReturn(t, retryFunction.Body.List[2])

	calls := map[string]int{}
	ast.Inspect(retryFunction.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch function := call.Fun.(type) {
		case *ast.Ident:
			calls[function.Name]++
		case *ast.SelectorExpr:
			if receiver, ok := function.X.(*ast.Ident); ok {
				calls[receiver.Name+"."+function.Sel.Name]++
			}
		}
		return true
	})
	wantCalls := map[string]int{"dependencies.Post": 2, "println": 1, "clock.Sleep": 1, "len": 1}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("private Kibana retry operations were %#v, want only %#v", calls, wantCalls)
	}
}

func assertKibanaRetryPOSTAssignment(t *testing.T, statement ast.Stmt) {
	t.Helper()
	assignment, ok := statement.(*ast.AssignStmt)
	if !ok || assignment.Tok != token.DEFINE || len(assignment.Lhs) != 2 || len(assignment.Rhs) != 1 {
		t.Fatal("Kibana retry request no longer assigns exact error and response values")
	}
	errName, errOK := assignment.Lhs[0].(*ast.Ident)
	responseName, responseOK := assignment.Lhs[1].(*ast.Ident)
	call, callOK := assignment.Rhs[0].(*ast.CallExpr)
	var selector *ast.SelectorExpr
	selectorOK := false
	if callOK {
		selector, selectorOK = call.Fun.(*ast.SelectorExpr)
	}
	var dependency *ast.Ident
	dependencyOK := false
	if selectorOK {
		dependency, dependencyOK = selector.X.(*ast.Ident)
	}
	var request *ast.Ident
	requestOK := false
	if callOK && len(call.Args) == 1 {
		request, requestOK = call.Args[0].(*ast.Ident)
	}
	if !errOK || errName.Name != "err" || !responseOK || responseName.Name != "response" ||
		!callOK || !selectorOK || !dependencyOK || dependency.Name != "dependencies" ||
		selector.Sel.Name != "Post" || len(call.Args) != 1 || !requestOK || request.Name != "reguest" {
		t.Fatal("Kibana retry request no longer calls the dependency once with the exact request")
	}
}

func assertKibanaErrorResponseReturn(t *testing.T, statement ast.Stmt) {
	t.Helper()
	returned, ok := statement.(*ast.ReturnStmt)
	if !ok || len(returned.Results) != 2 {
		t.Fatal("Kibana retry no longer returns exact error and response values")
	}
	errName, errOK := returned.Results[0].(*ast.Ident)
	responseName, responseOK := returned.Results[1].(*ast.Ident)
	if !errOK || errName.Name != "err" || !responseOK || responseName.Name != "response" {
		t.Fatal("Kibana retry return order is no longer error before response")
	}
}

func parsedKibanaFunction(t *testing.T, name string) *ast.FuncDecl {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate Kibana retry contract test")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(filename), "kibana.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse Kibana implementation: %v", err)
	}
	for _, declaration := range parsed.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && function.Recv == nil && function.Name.Name == name {
			return function
		}
	}
	t.Fatalf("private Kibana function %q is missing", name)
	return nil
}

func TestRecordedKibanaPOSTRequestsRejectEmptyPopulations(t *testing.T) {
	httpRecording := &recordingKibanaHTTP{}
	queryRecording := &recordingKibanaQueryPOST{}
	retryPOSTRecording := &recordingKibanaRetryPOST{}
	retrySleepRecording := &recordingKibanaRetrySleeper{}

	if _, err := httpRecording.assertedRequests(); err == nil {
		t.Fatal("empty recorded Kibana HTTP POST request population passed")
	}
	if _, err := queryRecording.assertedRequests(); err == nil {
		t.Fatal("empty recorded Kibana query POST population passed")
	}
	if _, err := retryPOSTRecording.assertedRequests(); err == nil {
		t.Fatal("empty recorded Kibana retry POST population passed")
	}
	if _, err := retrySleepRecording.assertedDurations(); err == nil {
		t.Fatal("empty recorded Kibana retry sleep population passed")
	}
}

func assertKibanaUnmarshalError(t *testing.T, responseBody, label string) {
	t.Helper()
	body := &trackingKibanaBody{reader: strings.NewReader(responseBody)}
	recording := &recordingKibanaHTTP{responses: []*http.Response{{
		StatusCode: http.StatusUnauthorized,
		Status:     "401 Must Not Be Rejected Before Unmarshal",
		Body:       body,
	}}}

	err, _ := internalPost(recording.dependencies(), completeKibanaFetchRequest())

	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Fatalf("%s unmarshal error was %T %v, want *json.SyntaxError", label, err, err)
	}
	if !body.closed || body.reads == 0 {
		t.Fatalf("%s unmarshal lifecycle was closed=%t reads=%d", label, body.closed, body.reads)
	}
}

func completeKibanaFetchRequest() KibanaFetchRequest {
	return KibanaFetchRequest{
		Url:            "https://kibana.example.invalid/internal/bsearch?compress=false",
		AcceptLanguage: testKibanaAcceptLanguage,
		Authorization:  testKibanaAuthorization,
		ContentType:    testKibanaContentType,
		KbnVersion:     testKibanaVersion,
		Body:           testKibanaBody(),
	}
}

func kibanaResponseWithTimestamps(t *testing.T, timestamps []string) string {
	t.Helper()
	hits := make([]map[string]interface{}, len(timestamps))
	for index, timestamp := range timestamps {
		hits[index] = map[string]interface{}{
			"_id": "hit-" + strconv.Itoa(index),
			"fields": map[string]interface{}{
				"@timestamp": []string{timestamp},
			},
		}
	}
	header, err := json.Marshal(map[string]interface{}{
		"id": 1,
		"result": map[string]interface{}{
			"rawResponse": map[string]interface{}{
				"hits": map[string]interface{}{"total": len(hits)},
			},
		},
	})
	if err != nil {
		t.Fatalf("marshal Kibana response header: %v", err)
	}
	result, err := json.Marshal(map[string]interface{}{
		"id": 2,
		"result": map[string]interface{}{
			"rawResponse": map[string]interface{}{
				"hits": map[string]interface{}{"hits": hits},
			},
		},
	})
	if err != nil {
		t.Fatalf("marshal Kibana result: %v", err)
	}
	return string(header) + "\n" + string(result)
}

func kibanaResponseForTimestamps(t *testing.T, timestamps []string) KibanaResponse {
	t.Helper()
	records := strings.Split(kibanaResponseWithTimestamps(t, timestamps), "\n")
	response := KibanaResponse{}
	if err := json.Unmarshal([]byte(records[0]), &response.KibanaResponseHeader); err != nil {
		t.Fatalf("unmarshal generated Kibana response header: %v", err)
	}
	if err := json.Unmarshal([]byte(records[1]), &response.KibanaResult); err != nil {
		t.Fatalf("unmarshal generated Kibana result: %v", err)
	}
	return response
}
