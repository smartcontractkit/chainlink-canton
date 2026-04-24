package authentication

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIssuerEquals(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		a, b  string
		equal bool
	}{
		{
			name:  "trailing_slash",
			a:     "https://idp.example.com",
			b:     "https://idp.example.com/",
			equal: true,
		},
		{
			name:  "trailing_slash_both_sides",
			a:     "https://idp.example.com/oidc",
			b:     "https://idp.example.com/oidc/",
			equal: true,
		},
		{
			name:  "https_default_port_omitted",
			a:     "https://idp.example.com",
			b:     "https://idp.example.com:443",
			equal: true,
		},
		{
			name:  "https_default_port_omitted_with_path",
			a:     "https://idp.example.com/oidc",
			b:     "https://idp.example.com:443/oidc",
			equal: true,
		},
		{
			name:  "http_default_port_omitted",
			a:     "http://idp.example.com",
			b:     "http://idp.example.com:80",
			equal: true,
		},
		{
			name:  "scheme_case",
			a:     "HTTPS://idp.example.com",
			b:     "https://idp.example.com",
			equal: true,
		},
		{
			name:  "host_case",
			a:     "https://Idp.EXAMPLE.com",
			b:     "https://idp.example.com",
			equal: true,
		},
		{
			name:  "surrounding_space",
			a:     "  https://idp.example.com  ",
			b:     "https://idp.example.com/",
			equal: true,
		},
		{
			name:  "ipv6_trailing_and_port",
			a:     "https://[2001:db8::1]",
			b:     "https://[2001:db8::1]:443/",
			equal: true,
		},
		{
			name:  "different_hosts",
			a:     "https://a.example.com",
			b:     "https://b.example.com",
			equal: false,
		},
		{
			name:  "different_path",
			a:     "https://idp.example.com/oidc1",
			b:     "https://idp.example.com/oidc2",
			equal: false,
		},
		{
			name:  "nondefault_https_port",
			a:     "https://idp.example.com:8443",
			b:     "https://idp.example.com",
			equal: false,
		},
		{
			name:  "nondefault_http_port",
			a:     "http://idp.example.com:8080",
			b:     "http://idp.example.com",
			equal: false,
		},
		{
			name:  "path_case_matters",
			a:     "https://idp.example.com/Realm",
			b:     "https://idp.example.com/realm",
			equal: false,
		},
		{
			name:  "query_must_match",
			a:     "https://idp.example.com?k=v",
			b:     "https://idp.example.com",
			equal: false,
		},
		{
			name:  "same_query",
			a:     "https://idp.example.com?k=v",
			b:     "https://idp.example.com/?k=v",
			equal: true,
		},
		{
			name:  "fragment_must_match",
			a:     "https://idp.example.com#frag",
			b:     "https://idp.example.com",
			equal: false,
		},
		{
			name:  "userinfo_equal",
			a:     "https://user:pass@idp.example.com",
			b:     "https://user:pass@idp.example.com/",
			equal: true,
		},
		{
			name:  "userinfo_different",
			a:     "https://a:x@idp.example.com",
			b:     "https://b:x@idp.example.com",
			equal: false,
		},
		{
			name:  "empty",
			a:     "",
			b:     "https://idp.example.com",
			equal: false,
		},
		{
			name:  "invalid",
			a:     "://bad",
			b:     "https://idp.example.com",
			equal: false,
		},
		{
			name:  "relative",
			a:     "/idp",
			b:     "https://idp.example.com",
			equal: false,
		},
		{
			name:  "commutative",
			a:     "https://h/a",
			b:     "https://h:443/a/",
			equal: true,
		},
		// Parsed ports may exceed strconv.Atoi range; must not conflate with -1 sentinel.
		{
			name: "unparseable_port_different",
			//nolint:lll // test literal URLs
			a: "https://idp.example.com:9223372036854775808/x",
			//nolint:lll
			b:     "https://idp.example.com:9223372036854775809/x",
			equal: false,
		},
		{
			name: "unparseable_port_equal",
			//nolint:lll
			a: "https://idp.example.com:9223372036854775808",
			//nolint:lll
			b:     "https://idp.example.com:9223372036854775808/",
			equal: true,
		},
		{
			name: "unparseable_port_not_equal_to_omitted",
			//nolint:lll
			a:     "https://idp.example.com:9223372036854775808",
			b:     "https://idp.example.com",
			equal: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			eq := issuerEquals(tc.a, tc.b)
			require.Equalf(t, tc.equal, eq, "issuerEquals(%q, %q)=%v want %v", tc.a, tc.b, eq, tc.equal)
			eqba := issuerEquals(tc.b, tc.a)
			require.Equalf(t, tc.equal, eqba, "issuerEquals must be symmetric for %q and %q", tc.a, tc.b)
		})
	}
}

func TestParseAbsoluteIssuerURL(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		_, err := parseAbsoluteIssuerURL("")
		require.ErrorIs(t, err, errEmptyIssuerURL)
	})
	t.Run("relative", func(t *testing.T) {
		t.Parallel()
		_, err := parseAbsoluteIssuerURL("/x")
		require.ErrorIs(t, err, errNotAbsoluteIssuerURL)
	})
}
