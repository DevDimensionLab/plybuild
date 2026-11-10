package workflownotification

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func SystemDependencies() Dependencies {
	return Dependencies{Now: time.Now, LookupEnv: os.LookupEnv, Transport: sendHTTP}
}

// A fresh transport/connection per attempt prevents net/http's stale pooled
// connection retry path. No proxy lookup, redirect following or GetBody replay.
func sendHTTP(ctx context.Context, endpoint string, payload []byte) Observation {
	var connected atomic.Bool
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	tr := &http.Transport{Proxy: nil, DisableKeepAlives: true, ForceAttemptHTTP2: false, TLSHandshakeTimeout: 10 * time.Second, ResponseHeaderTimeout: 10 * time.Second}
	tr.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		c, e := dialer.DialContext(ctx, network, address)
		if e == nil {
			connected.Store(true)
		}
		return c, e
	}
	defer tr.CloseIdleConnections()
	return postHTTP(ctx, endpoint, payload, tr, func() bool { return !connected.Load() })
}
func postHTTP(ctx context.Context, endpoint string, payload []byte, rt http.RoundTripper, zeroBytes func() bool) Observation {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, e := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if e != nil {
		return Observation{Dispatch: "not_started", Failed: true}
	}
	req.GetBody = nil
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: rt, Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, e := client.Do(req)
	if e != nil {
		dispatch := "unknown"
		if zeroBytes != nil && zeroBytes() {
			dispatch = "not_started"
		}
		return Observation{Failed: true, Dispatch: dispatch}
	}
	defer resp.Body.Close()
	body, e := io.ReadAll(io.LimitReader(resp.Body, 4097))
	return Observation{Status: resp.StatusCode, Body: body, RetryAfter: resp.Header.Values("Retry-After"), Dispatch: "started", Failed: e != nil, Oversize: len(body) > 4096}
}
func classify(o Observation, now time.Time) Attempt {
	completed := stamp(now)
	code := "transport_unconfirmed"
	a := Attempt{State: "unknown", CompletedAt: &completed, Dispatch: o.Dispatch, ResponseCode: &code}
	if a.Dispatch != "started" && a.Dispatch != "not_started" {
		a.Dispatch = "unknown"
	}
	if o.Status > 0 {
		s := o.Status
		a.HTTPStatus = &s
	}
	if o.Failed {
		if o.Dispatch == "not_started" {
			a.State = "not_sent"
			code = "before_dispatch"
		}
		return a
	}
	if o.Oversize || len(o.Body) > 4096 {
		return a
	}
	body := string(o.Body)
	if strings.HasSuffix(body, "\n") {
		body = strings.TrimSuffix(body, "\n")
	}
	switch {
	case o.Status == 200 && body == "ok":
		a.State = "transport_acknowledged"
		code = "ok"
	case o.Status == 400 && (body == "invalid_payload" || body == "user_not_found"), o.Status == 403 && body == "action_prohibited", o.Status == 404 && body == "channel_not_found", o.Status == 410 && body == "channel_is_archived":
		a.State = "rejected"
		code = body
	case o.Status == 429 && len(o.RetryAfter) == 1:
		s := o.RetryAfter[0]
		valid := s != ""
		for _, c := range s {
			if c < '0' || c > '9' {
				valid = false
			}
		}
		n, e := strconv.Atoi(s)
		if valid && e == nil && n >= 1 && n <= 86400 {
			a.State = "rate_limited"
			code = "rate_limited"
			deadline := stamp(now.Add(time.Duration(n) * time.Second))
			a.RetryNotBefore = &deadline
		}
	}
	return a
}
