package writer

import (
	"path/filepath"
	"testing"

	"github.com/thomas/gotrunk/profile"
)

func TestWriteAndLoadTOML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "config.toml")

	values := map[string]string{
		"api_url": "https://api.example.com",
		"api_key": "secret-1234",
		"tenant":  "acme",
	}
	if err := Write(path, profile.FormatTOML, profile.ModeOverwrite, values); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := Load(path, profile.FormatTOML)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	for k, v := range values {
		if got[k] != v {
			t.Errorf("clé %q: got %v, want %v", k, got[k], v)
		}
	}
}

func TestMergePreservesUntouchedKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	first := map[string]string{
		"api_url":   "https://old.example.com",
		"api_key":   "old-key",
		"tenant":    "acme",
		"extra_key": "preserved",
	}
	if err := Write(path, profile.FormatTOML, profile.ModeOverwrite, first); err != nil {
		t.Fatalf("first write: %v", err)
	}

	update := map[string]string{
		"api_url": "https://new.example.com",
	}
	if err := Write(path, profile.FormatTOML, profile.ModeMerge, update); err != nil {
		t.Fatalf("merge write: %v", err)
	}

	got, err := Load(path, profile.FormatTOML)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got["api_url"] != "https://new.example.com" {
		t.Errorf("api_url not updated: %v", got["api_url"])
	}
	if got["extra_key"] != "preserved" {
		t.Errorf("extra_key lost in merge: %v", got["extra_key"])
	}
	if got["api_key"] != "old-key" {
		t.Errorf("api_key lost in merge: %v", got["api_key"])
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	got, err := Load(filepath.Join(t.TempDir(), "nope.toml"), profile.FormatTOML)
	if err != nil {
		t.Fatalf("load missing: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestRoundTripJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	values := map[string]string{"a": "1", "b": "two"}
	if err := Write(path, profile.FormatJSON, profile.ModeOverwrite, values); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := Load(path, profile.FormatJSON)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got["a"] != "1" || got["b"] != "two" {
		t.Errorf("unexpected: %v", got)
	}
}

func TestBackupCreatesCopy(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	values := map[string]string{"a": "original"}
	if err := Write(path, profile.FormatJSON, profile.ModeOverwrite, values); err != nil {
		t.Fatalf("write: %v", err)
	}

	backupPath, err := Backup(path)
	if err != nil {
		t.Fatalf("backup: %v", err)
	}
	if backupPath == "" {
		t.Fatal("expected backup path, got empty string")
	}

	got, err := Load(backupPath, profile.FormatJSON)
	if err != nil {
		t.Fatalf("load backup: %v", err)
	}
	if got["a"] != "original" {
		t.Errorf("backup content unexpected: %v", got)
	}

	if _, err := Load(path, profile.FormatJSON); err != nil {
		t.Errorf("original file should still exist: %v", err)
	}
}

func TestBackupMissingFileReturnsEmpty(t *testing.T) {
	backupPath, err := Backup(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("backup of missing: %v", err)
	}
	if backupPath != "" {
		t.Errorf("expected empty backup path for missing file, got %q", backupPath)
	}
}

func TestRoundTripYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	values := map[string]string{"a": "1", "b": "two"}
	if err := Write(path, profile.FormatYAML, profile.ModeOverwrite, values); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := Load(path, profile.FormatYAML)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got["a"] != "1" || got["b"] != "two" {
		t.Errorf("unexpected: %v", got)
	}
}
