package upstream

import (
	"context"
	"testing"
)

func TestCheckUpdate(t *testing.T) {
	info, err := CheckUpdate(context.Background(), "1.0.0")
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if info.LatestVersion != "1.0.0" { t.Fatalf("version mismatch") }
	if info.DownloadURL == "" { t.Fatal("expected download url") }
}
