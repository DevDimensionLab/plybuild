package workspace

import (
	"crypto/rand"
	"io"
	"regexp"
	"time"
)

type IntegrationAuthorityID string
type IntegrationIntentID string
type IntegrationAttemptID string
type IntegrationResultID string

type TaskLifecycleIDSource interface {
	NewIntegrationAuthorityID() (IntegrationAuthorityID, error)
	NewIntegrationIntentID() (IntegrationIntentID, error)
	NewIntegrationAttemptID() (IntegrationAttemptID, error)
	NewIntegrationResultID() (IntegrationResultID, error)
}

type WorkClock interface{ Now() time.Time }

type cryptoTaskLifecycleIDSource struct{ Reader io.Reader }

func (source cryptoTaskLifecycleIDSource) reader() io.Reader {
	if source.Reader != nil {
		return source.Reader
	}
	return rand.Reader
}
func (source cryptoTaskLifecycleIDSource) next(prefix string) (string, error) {
	value, err := randomWorkItemID(source.reader(), prefix)
	return value, err
}
func (source cryptoTaskLifecycleIDSource) NewIntegrationAuthorityID() (IntegrationAuthorityID, error) {
	v, e := source.next("iauth_")
	return IntegrationAuthorityID(v), e
}
func (source cryptoTaskLifecycleIDSource) NewIntegrationIntentID() (IntegrationIntentID, error) {
	v, e := source.next("iint_")
	return IntegrationIntentID(v), e
}
func (source cryptoTaskLifecycleIDSource) NewIntegrationAttemptID() (IntegrationAttemptID, error) {
	v, e := source.next("iat_")
	return IntegrationAttemptID(v), e
}
func (source cryptoTaskLifecycleIDSource) NewIntegrationResultID() (IntegrationResultID, error) {
	v, e := source.next("ires_")
	return IntegrationResultID(v), e
}

type systemWorkClock struct{}

func (systemWorkClock) Now() time.Time { return time.Now() }

var lifecycleIDPatterns = map[string]*regexp.Regexp{
	"integration authority": regexp.MustCompile(`^iauth_[0-9a-f]{32}$`),
	"integration intent":    regexp.MustCompile(`^iint_[0-9a-f]{32}$`),
	"integration attempt":   regexp.MustCompile(`^iat_[0-9a-f]{32}$`),
	"integration result":    regexp.MustCompile(`^ires_[0-9a-f]{32}$`),
}
