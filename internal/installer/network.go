package installer

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

// NetworkChecker handles network connectivity checks
type NetworkChecker struct {
	timeout time.Duration
}

// NewNetworkChecker creates a new network checker
func NewNetworkChecker() *NetworkChecker {
	return &NetworkChecker{
		timeout: 10 * time.Second,
	}
}

// CheckConnectivity verifies internet connectivity
func (n *NetworkChecker) CheckConnectivity(ctx context.Context) error {
	spinner := ui.NewSpinner("Checking network connectivity...")
	spinner.Start()
	defer spinner.Stop()

	// List of hosts to check (in order of preference)
	hosts := []string{
		"https://raw.githubusercontent.com", // GitHub raw (used by Homebrew installer)
		"https://github.com",                // GitHub
		"https://brew.sh",                   // Homebrew
		"https://apple.com",                 // Apple (for Xcode CLT)
	}

	var lastErr error
	for _, host := range hosts {
		if err := n.checkHost(ctx, host); err == nil {
			spinner.Success("Network connectivity OK")
			return nil
		} else {
			lastErr = err
		}
	}

	spinner.Fail("No network connectivity")
	return fmt.Errorf("network connectivity check failed: %w\n\nPlease check your internet connection and try again", lastErr)
}

// checkHost attempts to connect to a specific host
func (n *NetworkChecker) checkHost(ctx context.Context, url string) error {
	client := &http.Client{
		Timeout: n.timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Any response means we have connectivity
	return nil
}

// CheckConnectivityQuick does a quick DNS check without full HTTP request
func (n *NetworkChecker) CheckConnectivityQuick(ctx context.Context) bool {
	_, err := net.LookupHost("github.com")
	return err == nil
}
