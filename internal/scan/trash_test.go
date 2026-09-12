package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTrashContainsDirectoryWithUntrustedID(t *testing.T) {
	root := t.TempDir()
	store := filepath.Join(root, "store")
	project := filepath.Join(store, "sessions", "%2Ftmp")
	source := filepath.Join(project, "actual-session")
	if err := os.MkdirAll(source, 0o700); err != nil {
		t.Fatal(err)
	}

	target, err := Trash(store, Session{
		ID:   "../../../escaped",
		Path: source,
	})
	if err != nil {
		t.Fatal(err)
	}

	trashRoot := filepath.Join(store, TrashDir)
	if filepath.Dir(filepath.Dir(target)) != trashRoot {
		t.Fatalf("trash target escaped its dated directory: %s", target)
	}
	wantName := filepath.Base(project) + "__" + filepath.Base(source)
	if got := filepath.Base(target); got != wantName {
		t.Fatalf("trash name = %q, want filesystem-derived %q", got, wantName)
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("trashed session missing at %s: %v", target, err)
	}
	if _, err := os.Stat(source); !os.IsNotExist(err) {
		t.Fatalf("source still exists after trash: %v", err)
	}
}

func TestTrashContainsFileBackedSidecarWithUntrustedID(t *testing.T) {
	root := t.TempDir()
	store := filepath.Join(root, "store")
	project := filepath.Join(store, "sessions", "project")
	transcript := filepath.Join(project, "rollout-X.jsonl")
	sidecar := strings.TrimSuffix(transcript, ".jsonl")
	if err := os.MkdirAll(project, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(transcript, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(sidecar, 0o700); err != nil {
		t.Fatal(err)
	}

	transcriptTarget, err := Trash(store, Session{ID: "../../../../../pwned", Path: transcript})
	if err != nil {
		t.Fatal(err)
	}
	trashRoot := filepath.Join(store, TrashDir)
	if filepath.Dir(filepath.Dir(transcriptTarget)) != trashRoot {
		t.Fatalf("transcript target escaped its dated directory: %s", transcriptTarget)
	}
	wantSidecar := filepath.Join(filepath.Dir(transcriptTarget), filepath.Base(project)+"__"+filepath.Base(sidecar))
	if _, err := os.Stat(wantSidecar); err != nil {
		t.Fatalf("sidecar was not contained in trash at %s: %v", wantSidecar, err)
	}
	if _, err := os.Stat(sidecar); !os.IsNotExist(err) {
		t.Fatalf("sidecar still exists at source: %v", err)
	}
}
