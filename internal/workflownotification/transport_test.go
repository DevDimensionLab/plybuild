package workflownotification

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestHTTPClassificationAndBoundedAdapter(t *testing.T) {
	cases := []struct {
		status      int
		body        string
		headers     []string
		state, code string
	}{
		{200, "ok", nil, "transport_acknowledged", "ok"}, {200, "ok\n", nil, "transport_acknowledged", "ok"}, {200, "ok\n\n", nil, "unknown", "transport_unconfirmed"},
		{400, "invalid_payload", nil, "rejected", "invalid_payload"}, {400, "user_not_found\n", nil, "rejected", "user_not_found"}, {403, "action_prohibited", nil, "rejected", "action_prohibited"}, {404, "channel_not_found", nil, "rejected", "channel_not_found"}, {410, "channel_is_archived", nil, "rejected", "channel_is_archived"},
		{403, "other", nil, "unknown", "transport_unconfirmed"}, {500, "rollup_error", nil, "unknown", "transport_unconfirmed"}, {302, "ok", nil, "unknown", "transport_unconfirmed"},
		{429, "ignored", []string{"60"}, "rate_limited", "rate_limited"}, {429, "ignored", []string{"60", "60"}, "unknown", "transport_unconfirmed"}, {429, "ignored", []string{"0"}, "unknown", "transport_unconfirmed"}, {429, "ignored", []string{"86401"}, "unknown", "transport_unconfirmed"}, {429, "ignored", []string{"+60"}, "unknown", "transport_unconfirmed"}, {200, strings.Repeat("x", 4097), nil, "unknown", "transport_unconfirmed"},
	}
	now := time.Date(2026, 10, 3, 19, 0, 0, 0, time.UTC)
	t.Setenv("HTTPS_PROXY", "http://must-not-be-used.invalid")
	for _, c := range cases {
		calls := 0
		o := postHTTP(context.Background(), syntheticSecret, []byte(`{"text":"fixture"}`), roundTripFunc(func(r *http.Request) (*http.Response, error) {
			calls++
			if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" || r.GetBody != nil {
				t.Fatal("unexpected transport request")
			}
			deadline, ok := r.Context().Deadline()
			if !ok || time.Until(deadline) > 10*time.Second {
				t.Fatal("missing total deadline")
			}
			headers := http.Header{}
			for _, v := range c.headers {
				headers.Add("Retry-After", v)
			}
			headers.Set("Location", "https://must-not-follow.invalid/")
			return &http.Response{StatusCode: c.status, Body: io.NopCloser(strings.NewReader(c.body)), Header: headers, Request: r}, nil
		}), nil)
		a := classify(o, now)
		if calls != 1 || a.State != c.state || *a.ResponseCode != c.code {
			t.Errorf("status %d body length %d: %+v calls %d", c.status, len(c.body), a, calls)
		}
		if len(o.Body) > 4097 {
			t.Fatal("unbounded response read")
		}
	}
	for _, zero := range []bool{false, true} {
		calls := 0
		o := postHTTP(context.Background(), syntheticSecret, nil, roundTripFunc(func(*http.Request) (*http.Response, error) { calls++; return nil, context.DeadlineExceeded }), func() bool { return zero })
		a := classify(o, now)
		want := "unknown"
		if zero {
			want = "not_sent"
		}
		if a.State != want || calls != 1 {
			t.Fatal("unsafe error/retry classification", a)
		}
	}
}

type failingBody struct{}

func (failingBody) Read([]byte) (int, error) {
	return 0, errors.New("raw response secret must not escape")
}
func (failingBody) Close() error { return nil }
func TestTruncatedBodyIsUnknown(t *testing.T) {
	o := postHTTP(context.Background(), syntheticSecret, nil, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: failingBody{}, Header: http.Header{}, Request: r}, nil
	}), nil)
	if classify(o, time.Now()).State != "unknown" {
		t.Fatal("body read failure accepted")
	}
}
