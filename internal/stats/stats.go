package stats

import (
	"sort"
	"strings"
	"time"

	"github-stats-cli/internal/github"
)

type Summary struct {
	Repositories int `json:"repositories"`
	Original     int `json:"original_repositories"`
	Forked       int `json:"forked_repositories"`
	Stars        int `json:"stars"`
	Forks        int `json:"forks"`
	Watchers     int `json:"watchers"`
	OpenIssues   int `json:"open_issues"`
}
type Language struct {
	Name         string  `json:"name"`
	Repositories int     `json:"repositories"`
	Percentage   float64 `json:"percentage"`
}

func Summarize(repos []github.Repository) Summary {
	s := Summary{Repositories: len(repos)}
	for _, r := range repos {
		s.Stars += r.Stars
		s.Forks += r.Forks
		s.Watchers += r.Watchers
		s.OpenIssues += r.OpenIssues
		if r.Fork {
			s.Forked++
		} else {
			s.Original++
		}
	}
	return s
}

func Languages(repos []github.Repository) []Language {
	counts := make(map[string]int)
	for _, r := range repos {
		if name := strings.TrimSpace(r.Language); name != "" {
			counts[name]++
		}
	}
	result := make([]Language, 0, len(counts))
	total := 0
	for _, n := range counts {
		total += n
	}
	for name, n := range counts {
		result = append(result, Language{Name: name, Repositories: n, Percentage: float64(n) * 100 / float64(total)})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Repositories == result[j].Repositories {
			return result[i].Name < result[j].Name
		}
		return result[i].Repositories > result[j].Repositories
	})
	return result
}

func Sort(repos []github.Repository, by string) {
	sort.SliceStable(repos, func(i, j int) bool {
		switch by {
		case "forks":
			return repos[i].Forks > repos[j].Forks
		case "updated":
			return repos[i].UpdatedAt.After(repos[j].UpdatedAt)
		default:
			return repos[i].Stars > repos[j].Stars
		}
	})
}

func Top(repos []github.Repository, n int) []github.Repository {
	if n > len(repos) {
		n = len(repos)
	}
	if n < 0 {
		n = 0
	}
	return repos[:n]
}

func FormatDate(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("2006-01-02")
}
