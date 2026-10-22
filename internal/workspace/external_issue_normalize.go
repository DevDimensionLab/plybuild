package workspace

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var externalDNS = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)
var externalSegment = regexp.MustCompile(`^[A-Za-z0-9_~.-]+$`)
var externalJiraKey = regexp.MustCompile(`^[A-Z][A-Z0-9_]*-[1-9][0-9]*$`)
var externalOwner = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,99}$`)
var externalRepo = regexp.MustCompile(`^[a-z0-9_.-]{1,100}$`)

// NormalizeExternalIssue validates supplied offline metadata. It performs no I/O.
// Scheme is retained in the URL but not issue identity: HTTP/HTTPS aliases cannot
// claim the same issue for different work items.
func NormalizeExternalIssue(in ExternalIssue) (ExternalIssue, error) {
	fail := func(detail string) (ExternalIssue, error) {
		return ExternalIssue{}, WorkInvalidArguments("external issue: " + detail)
	}
	if !validTaskText(in.URL, 1, 2048) {
		return fail("URL must contain 1..2048 characters without controls or surrounding whitespace")
	}
	u, err := url.Parse(in.URL)
	if err != nil || u == nil {
		return fail("invalid URL")
	}
	u.Scheme = strings.ToLower(u.Scheme)
	if (u.Scheme != "https" && u.Scheme != "http") || u.Opaque != "" || u.User != nil || u.Host == "" || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || strings.Contains(in.URL, "#") {
		return fail("use an absolute http or https URL without credentials, query or fragment")
	}
	host := strings.ToLower(strings.TrimSuffix(u.Hostname(), "."))
	if host == "" || len(host) > 253 {
		return fail("invalid DNS host")
	}
	for _, label := range strings.Split(host, ".") {
		if !externalDNS.MatchString(label) {
			return fail("host must be an ASCII DNS name or IPv4 address")
		}
	}
	port := u.Port()
	if strings.HasSuffix(u.Host, ":") {
		return fail("empty port")
	}
	if port != "" {
		n, e := strconv.Atoi(port)
		if e != nil || n < 1 || n > 65535 {
			return fail("invalid port")
		}
		port = strconv.Itoa(n)
	}
	// Only each scheme's default port is an alias of an omitted port.
	if u.Scheme == "https" && port == "443" || u.Scheme == "http" && port == "80" {
		port = ""
	}
	if port != "" {
		host = net.JoinHostPort(host, port)
	}
	if strings.Contains(u.Path, "\\") || !strings.HasPrefix(u.Path, "/") {
		return fail("invalid path")
	}
	path := strings.TrimSuffix(u.Path, "/")
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	for _, part := range parts {
		if len(part) == 0 || len(part) > 128 || part == "." || part == ".." || !externalSegment.MatchString(part) {
			return fail("path must contain bounded unambiguous ASCII segments")
		}
	}
	// Encoded slashes and backslashes are separators with ambiguous router meaning.
	escaped := strings.ToLower(u.EscapedPath())
	if strings.Contains(escaped, "%2f") || strings.Contains(escaped, "%5c") {
		return fail("encoded path separators are not supported")
	}
	out := ExternalIssue{Provider: in.Provider, Host: host}
	switch in.Provider {
	case ExternalJira:
		if len(parts) < 2 || parts[len(parts)-2] != "browse" {
			return fail("Jira requires /<optional-context>/browse/KEY routing")
		}
		out.JiraContextPath = ""
		if len(parts) > 2 {
			out.JiraContextPath = "/" + strings.Join(parts[:len(parts)-2], "/")
		}
		if in.JiraContextPath != "" && in.JiraContextPath != out.JiraContextPath {
			return fail("supplied Jira context disagrees with URL")
		}
		key := strings.ToUpper(in.Key)
		if len(key) > 128 || !externalJiraKey.MatchString(key) || strings.ToUpper(parts[len(parts)-1]) != key {
			return fail("Jira key must match the last decoded URL path segment (for example PLY-123)")
		}
		if in.Owner != "" || in.Repository != "" || in.Number != 0 {
			return fail("Jira does not take GitHub identity fields")
		}
		out.Key = key
		parts[len(parts)-1] = key
	case ExternalGitHub:
		if in.JiraContextPath != "" {
			return fail("GitHub does not take a Jira context")
		}
		if in.Host == "" {
			return fail("GitHub requires an explicit host agreeing with the URL authority")
		}
		owner, repo := strings.ToLower(in.Owner), strings.ToLower(in.Repository)
		if !externalOwner.MatchString(owner) || !externalRepo.MatchString(repo) || repo == "." || repo == ".." || in.Number < 1 || in.Key != "" {
			return fail("GitHub requires owner, repository and a positive int64 issue number; no Jira key")
		}
		if len(parts) != 4 || strings.ToLower(parts[0]) != owner || strings.ToLower(parts[1]) != repo || parts[2] != "issues" {
			return fail("GitHub URL must be /owner/repository/issues/number agreeing with supplied identity")
		}
		number, e := strconv.ParseInt(parts[3], 10, 64)
		if e != nil || number != in.Number || strings.HasPrefix(parts[3], "+") {
			return fail("GitHub URL issue number disagrees")
		}
		out.Owner, out.Repository, out.Number = owner, repo, number
		parts = []string{owner, repo, "issues", strconv.FormatInt(number, 10)}
	default:
		return fail("provider must be jira or github")
	}
	if in.Host != "" && strings.ToLower(in.Host) != host {
		return fail("supplied host must match normalized URL authority: " + host)
	}
	out.URL = u.Scheme + "://" + host + "/" + strings.Join(parts, "/")
	identity := string(out.Provider) + "/" + host + out.JiraContextPath + "/" + out.CommitKey()
	out.ID = "ext_" + strings.TrimPrefix(digestTaskBytes([]byte(identity)), "sha256:")
	if in.ID != "" && in.ID != out.ID {
		return fail("supplied binding ID disagrees with normalized identity")
	}
	return out, nil
}
func (issue ExternalIssue) CommitKey() string {
	if issue.Provider == ExternalJira {
		return issue.Key
	}
	return fmt.Sprintf("%s/%s#%d", issue.Owner, issue.Repository, issue.Number)
}
