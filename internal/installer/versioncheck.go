package installer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

const (
	githubRepoOwner = "tldr-it-stepankutaj"
	githubRepoName  = "setup-mac"
	githubAPIURL    = "https://api.github.com/repos/%s/%s/releases/latest"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	PublishedAt string         `json:"published_at"`
	HTMLURL     string         `json:"html_url"`
	Assets      []ReleaseAsset `json:"assets"`
}

// ReleaseAsset represents a release asset
type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// VersionChecker checks for new versions
type VersionChecker struct {
	currentVersion string
	timeout        time.Duration
}

// NewVersionChecker creates a new version checker
func NewVersionChecker(currentVersion string) *VersionChecker {
	return &VersionChecker{
		currentVersion: currentVersion,
		timeout:        5 * time.Second,
	}
}

// CheckForUpdate checks if a newer version is available
func (v *VersionChecker) CheckForUpdate(ctx context.Context) (*GitHubRelease, bool, error) {
	url := fmt.Sprintf(githubAPIURL, githubRepoOwner, githubRepoName)

	client := &http.Client{Timeout: v.timeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "setup-mac/"+v.currentVersion)

	resp, err := client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, err
	}

	// Compare versions
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(v.currentVersion, "v")

	// Simple comparison - if they're different and latest is "greater"
	isNewer := v.isNewerVersion(latestVersion, currentVersion)

	return &release, isNewer, nil
}

// isNewerVersion compares semantic versions
func (v *VersionChecker) isNewerVersion(latest, current string) bool {
	// Handle development versions
	if current == "" || current == "dev" || strings.Contains(current, "dirty") {
		return false // Don't suggest updates for dev builds
	}

	// Simple string comparison for now
	// For proper semver, we'd use a library
	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	for i := 0; i < len(latestParts) && i < len(currentParts); i++ {
		if latestParts[i] > currentParts[i] {
			return true
		} else if latestParts[i] < currentParts[i] {
			return false
		}
	}

	return len(latestParts) > len(currentParts)
}

// GetDownloadURL returns the download URL for the current platform
func (v *VersionChecker) GetDownloadURL(release *GitHubRelease) string {
	arch := runtime.GOARCH
	expectedName := fmt.Sprintf("setup-mac-darwin-%s.tar.gz", arch)

	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			return asset.BrowserDownloadURL
		}
	}

	// Fallback to release page
	return release.HTMLURL
}

// CheckAndPrompt checks for updates and prompts user if available
func (v *VersionChecker) CheckAndPrompt(ctx context.Context, prompt *ui.Prompt) error {
	release, isNewer, err := v.CheckForUpdate(ctx)
	if err != nil {
		// Silently ignore errors - update check is not critical
		return nil
	}

	if !isNewer {
		return nil
	}

	// Show update notification
	fmt.Println()
	ui.PrintInfo(fmt.Sprintf("New version available: %s (current: %s)", release.TagName, v.currentVersion))

	downloadURL := v.GetDownloadURL(release)

	// Ask user if they want to update
	if prompt != nil && prompt.Interactive {
		update, err := prompt.Confirm("Would you like to download the update?", false)
		if err != nil {
			return nil
		}

		if update {
			ui.PrintInfo(fmt.Sprintf("Download from: %s", downloadURL))
			ui.PrintInfo("After downloading, extract and run: sudo cp setup-mac /usr/local/bin/")
		}
	} else {
		ui.PrintInfo(fmt.Sprintf("Download: %s", downloadURL))
	}
	fmt.Println()

	return nil
}
