package authentication

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// issuerEquals reports whether the urls a and b refer to the same authorization server issuer
// for AS metadata checks:
//
// * http(s) scheme and host are compared case-insensitively,
// * default ports (80 and 443) are treated as equal to the explicit :80 / :443 form,
// * a trailing "/" on the path is ignored,
// * userinfo, raw query, and fragment must match when present.
func issuerEquals(a, b string) bool {
	ua, err1 := parseAbsoluteIssuerURL(strings.TrimSpace(a))
	ub, err2 := parseAbsoluteIssuerURL(strings.TrimSpace(b))
	if err1 != nil || err2 != nil {
		return false
	}
	if !strings.EqualFold(ua.Scheme, ub.Scheme) {
		return false
	}
	if !strings.EqualFold(ua.Hostname(), ub.Hostname()) {
		return false
	}
	if !portsEqual(ua, ub) {
		return false
	}
	if !userInfoEqual(ua, ub) {
		return false
	}
	// Sanity check:
	// Hierarchical http(s) URLs; opaque form should not appear for issuers in practice.
	if ua.Opaque != "" || ub.Opaque != "" {
		return ua.Opaque == ub.Opaque
	}
	if normalizeIssuerPath(ua) != normalizeIssuerPath(ub) {
		return false
	}
	if ua.RawQuery != ub.RawQuery {
		return false
	}
	if ua.Fragment != ub.Fragment {
		return false
	}

	return true
}

func parseAbsoluteIssuerURL(s string) (*url.URL, error) {
	if s == "" {
		return nil, errEmptyIssuerURL
	}
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errNotAbsoluteIssuerURL
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
	default:
		return nil, errInvalidIssuerScheme
	}

	return u, nil
}

var errEmptyIssuerURL = errors.New("empty issuer URL")

// errNotAbsoluteIssuerURL is returned for relative URLs or without host.
var errNotAbsoluteIssuerURL = errors.New("issuer must be absolute with host")

// errInvalidIssuerScheme is returned when the scheme is not http or https.
var errInvalidIssuerScheme = errors.New("issuer URL scheme must be http or https")

// portsEqual compares URL ports, treating default http(s) ports as equal when the port
// is omitted or explicitly 80/443, and using string equality when a port is present but
// cannot be parsed (e.g. Atoi out-of-range), so two different unparseable ports are never
// conflated.
func portsEqual(ua, ub *url.URL) bool {
	pa, pb := ua.Port(), ub.Port()
	na, okA := portAsInt(ua.Scheme, pa)
	nb, okB := portAsInt(ub.Scheme, pb)
	if okA && okB {
		return na == nb
	}
	if !okA && !okB {
		// e.g. two :9223372036854775808 must match; two different unparseable strings must not.
		return pa == pb
	}

	// e.g. default 443 vs unparseable explicit port, or 8443 vs unparseable.
	return false
}

// portAsInt returns the port as an integer, using default 80/443 for http(s) when empty;
// the second return is false when a non-empty port cannot be parsed to int.
// parseAbsoluteIssuerURL only permits http and https, so the empty-port default is unreachable
// in issuerEquals but kept defensively for other callers of portAsInt.
func portAsInt(scheme, port string) (n int, ok bool) {
	if port == "" {
		switch strings.ToLower(scheme) {
		case "https":
			return 443, true
		case "http":
			return 80, true
		default:
			return 0, false
		}
	}
	n, err := strconv.Atoi(port)
	if err != nil {
		return 0, false
	}

	return n, true
}

// normalizeIssuerPath trims a single trailing slash on the path so "" and "/" match for the root.
func normalizeIssuerPath(u *url.URL) string {
	// Use EscapedPath for stable bytes; issuer is typically ASCII.
	return strings.TrimSuffix(u.EscapedPath(), "/")
}

func userInfoEqual(a, b *url.URL) bool {
	switch {
	case a.User == nil && b.User == nil:
		return true
	case a.User == nil || b.User == nil:
		return false
	default:
		// String() is stable for comparison of equivalent userinfo.
		return a.User.String() == b.User.String()
	}
}
