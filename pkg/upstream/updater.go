package upstream

import (
	"context"
	"fmt"
)

type UpdateInfo struct {
	LatestVersion string `json:"latest_version"`
	DownloadURL   string `json:"download_url"`
	ReleaseNotes  string `json:"release_notes"`
}

func CheckUpdate(ctx context.Context, currentVersion string) (*UpdateInfo, error) {
	return &UpdateInfo{
		LatestVersion: currentVersion,
		DownloadURL:   fmt.Sprintf("https://github.com/HP-jh/TarsDify/releases/tag/v%s", currentVersion),
		ReleaseNotes:  "stub",
	}, nil
}

func (c *Client) CheckUpdate(ctx context.Context, currentVersion string) (*UpdateInfo, error) {
	return CheckUpdate(ctx, currentVersion)
}
