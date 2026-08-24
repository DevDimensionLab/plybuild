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

func TestPOSTKeepsZeroHitRetrySelectionAndOutputWithoutWaiting(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate Kibana retry contract test")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(filename), "kibana.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse Kibana implementation: %v", err)
	}
	var retryFunction *ast.FuncDecl
	for _, declaration := range parsed.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && function.Recv == nil && function.Name.Name == "POST" {
			retryFunction = function
			break
		}
	}
	if retryFunction == nil {
		t.Fatal("Kibana POST retry function is missing")
	}
	internalCalls := 0
	zeroHitSelection := false
	retryOutput := false
	retrySleep := false
	ast.Inspect(retryFunction.Body, func(node ast.Node) bool {
		switch value := node.(type) {
		case *ast.CallExpr:
			if identifier, ok := value.Fun.(*ast.Ident); ok {
				if identifier.Name == "internalPOST" {
					internalCalls++
				}
				if identifier.Name == "println" && len(value.Args) == 1 {
					literal, ok := value.Args[0].(*ast.BasicLit)
					retryOutput = ok && literal.Value == strconv.Quote("sleep and retry")
				}
			}
			if selector, ok := value.Fun.(*ast.SelectorExpr); ok {
				packageName, packageOK := selector.X.(*ast.Ident)
				retrySleep = retrySleep || (packageOK && packageName.Name == "time" && selector.Sel.Name == "Sleep")
			}
		case *ast.BinaryExpr:
			literal, ok := value.Y.(*ast.BasicLit)
			call, callOK := value.X.(*ast.CallExpr)
			if ok && callOK && value.Op == token.EQL {
				identifier, identifierOK := call.Fun.(*ast.Ident)
				zeroHitSelection = zeroHitSelection ||
					(identifierOK && literal.Value == "0" && identifier.Name == "len")
			}
		}
		return true
	})
	if internalCalls != 2 || !zeroHitSelection || !retryOutput || !retrySleep {
		t.Fatalf("retry contract was calls=%d zero-hit=%t output=%t sleep=%t",
			internalCalls, zeroHitSelection, retryOutput, retrySleep)
	}
}

func TestRecordedKibanaPOSTRequestsRejectEmptyPopulations(t *testing.T) {
	httpRecording := &recordingKibanaHTTP{}
	queryRecording := &recordingKibanaQueryPOST{}

	if _, err := httpRecording.assertedRequests(); err == nil {
		t.Fatal("empty recorded Kibana HTTP POST request population passed")
	}
	if _, err := queryRecording.assertedRequests(); err == nil {
		t.Fatal("empty recorded Kibana query POST population passed")
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
