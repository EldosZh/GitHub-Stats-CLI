package stats

import (
	"github-stats-cli/internal/github"
	"testing"
	"time"
)

func TestSummaryLanguagesAndSort(t *testing.T) {
	repos := []github.Repository{{Name: "a", Stars: 5, Forks: 1, Watchers: 2, OpenIssues: 3, Language: "Go", UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}, {Name: "b", Stars: 2, Forks: 8, Language: "Go", Fork: true, UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}, {Name: "c", Stars: 0, Language: "Rust"}}
	s := Summarize(repos)
	if s.Repositories != 3 || s.Original != 2 || s.Forked != 1 || s.Stars != 7 || s.Forks != 9 || s.Watchers != 2 || s.OpenIssues != 3 {
		t.Fatalf("summary=%+v", s)
	}
	langs := Languages(repos)
	if len(langs) != 2 || langs[0].Name != "Go" || langs[0].Repositories != 2 || langs[0].Percentage != 100*2.0/3 {
		t.Fatalf("languages=%+v", langs)
	}
	Sort(repos, "forks")
	if repos[0].Name != "b" {
		t.Fatal("fork sort")
	}
	Sort(repos, "updated")
	if repos[0].Name != "b" {
		t.Fatal("updated sort")
	}
	Sort(repos, "stars")
	if repos[0].Name != "a" {
		t.Fatal("star sort")
	}
	if len(Top(repos, 10)) != 3 || len(Top(nil, 5)) != 0 {
		t.Fatal("top or empty list")
	}
}

func TestEmptyLanguages(t *testing.T) {
	if len(Languages(nil)) != 0 || Summarize(nil).Repositories != 0 {
		t.Fatal("expected empty statistics")
	}
}
