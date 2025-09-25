package util

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ResolveHostFromURL parses rawurl, returns:
//   - scheme://host (without port)
//   - host (without port)
//   - resolved IPs
// If rawurl has no scheme, it defaults to "http://".
func ResolveHostFromURL(rawurl string) (string, string, []string, error) {
	// Ensure URL has a scheme so url.Parse works predictably
	if !strings.Contains(rawurl, "://") {
		rawurl = "http://" + rawurl
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return "", "", nil, err
	}

	// Extract host (might include port)
	host := u.Host
	if host == "" {
		host = u.Path
	}

	// If host contains path (rare), trim it
	if idx := strings.Index(host, "/"); idx != -1 {
		host = host[:idx]
	}

	// Remove port if exists
	if h, _, err2 := net.SplitHostPort(host); err2 == nil {
		host = h
	}

	if host == "" {
		return "", "", nil, fmt.Errorf("no host found in url: %s", rawurl)
	}

	// Resolve IPs
	ips, err := net.LookupHost(host)
	if err != nil {
		// return host and scheme but with DNS error
		return fmt.Sprintf("%s://%s", u.Scheme, host), host, nil, err
	}

	return fmt.Sprintf("%s://%s", u.Scheme, host), host, ips, nil
}
