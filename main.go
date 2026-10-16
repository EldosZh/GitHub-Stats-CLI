package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github-stats-cli/internal/github"
	"github-stats-cli/internal/stats"
)

const version = "github-stats-cli v1.0.0"

type output struct {
	User            github.User         `json:"user"`
	RepositoryStats stats.Summary       `json:"repository_stats"`
	TopRepositories []github.Repository `json:"top_repositories"`
	Languages       []stats.Language    `json:"languages"`
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		usage(stdout)
		return nil
	}
	if len(args) == 1 && args[0] == "--version" {
		fmt.Fprintln(stdout, version)
		return nil
	}
	if len(args) == 0 {
		usage(stdout)
		return nil
	}
	username := args[0]
	if strings.HasPrefix(username, "-") {
		return fmt.Errorf("username is required\n\nUsage: github-stats <username> [flags]")
	}
	flags := flag.NewFlagSet("github-stats", flag.ContinueOnError)
	flags.SetOutput(stderr)
	top := flags.Int("top", 5, "number of top repositories to display")
	sortBy := flags.String("sort", "stars", "repository sort order: stars, forks, updated")
	jsonOutput := flags.Bool("json", false, "output JSON")
	flags.Usage = func() { usage(stderr) }
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	if *top < 1 {
		return fmt.Errorf("--top must be a positive integer")
	}
	if *sortBy != "stars" && *sortBy != "forks" && *sortBy != "updated" {
		return fmt.Errorf("invalid sort value %q. Use stars, forks, or updated.", *sortBy)
	}
	client := github.NewClient(os.Getenv("GITHUB_API_URL"), os.Getenv("GITHUB_TOKEN"))
	ctx := context.Background()
	user, err := client.GetUser(ctx, username)
	if err != nil {
		return err
	}
	repos, err := client.GetRepositories(ctx, username)
	if err != nil {
		return err
	}
	stats.Sort(repos, *sortBy)
	result := output{User: user, RepositoryStats: stats.Summarize(repos), TopRepositories: stats.Top(repos, *top), Languages: stats.Languages(repos)}
	if *jsonOutput {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}
	printHuman(stdout, result)
	return nil
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "GitHub Stats CLI\n\nUsage:\n  github-stats <username> [flags]\n\nExample:\n  github-stats torvalds\n\nFlags:\n  --top N       number of top repositories (default 5)\n  --sort VALUE  stars, forks, or updated (default stars)\n  --json        output JSON\n  --help        show this help\n  --version     show version")
}

func value(s string) string {
	if strings.TrimSpace(s) == "" {
		return "N/A"
	}
	return s
}
func printHuman(w io.Writer, o output) {
	u := o.User
	fmt.Fprintln(w, "GitHub Stats\n────────────────────────────────────────")
	fmt.Fprintf(w, "User:          %s\nName:          %s\nBio:           %s\nLocation:      %s\nCompany:       %s\nFollowers:     %d\nFollowing:     %d\nPublic Repos:  %d\nCreated:       %s\nProfile:        https://github.com/%s\n", u.Login, value(u.Name), value(u.Bio), value(u.Location), value(u.Company), u.Followers, u.Following, u.PublicRepos, stats.FormatDate(u.CreatedAt), u.Login)
	s := o.RepositoryStats
	fmt.Fprintf(w, "\nRepository Stats\n────────────────────────────────────────\nRepositories:  %d\nOriginal Repos: %d\nForked Repos:  %d\nTotal Stars:   %d\nTotal Forks:   %d\nWatchers:      %d\nOpen Issues:   %d\n", s.Repositories, s.Original, s.Forked, s.Stars, s.Forks, s.Watchers, s.OpenIssues)
	fmt.Fprintln(w, "\nTop Repositories\n────────────────────────────────────────")
	if len(o.TopRepositories) == 0 {
		fmt.Fprintln(w, "No public repositories.")
	}
	for i, r := range o.TopRepositories {
		desc := strings.TrimSpace(r.Description)
		if desc == "" {
			desc = "No description"
		}
		if len(desc) > 100 {
			desc = desc[:97] + "..."
		}
		lang := value(r.Language)
		fmt.Fprintf(w, "%d. %s\n   ★ %d   Forks %d   %s\n   %s\n   %s\n", i+1, r.Name, r.Stars, r.Forks, lang, r.HTMLURL, desc)
	}
	fmt.Fprintln(w, "\nLanguages (repositories by primary language)\n────────────────────────────────────────")
	if len(o.Languages) == 0 {
		fmt.Fprintln(w, "No language data available.")
	}
	for _, l := range o.Languages {
		fmt.Fprintf(w, "%-18s %3d  %5.1f%%\n", l.Name, l.Repositories, l.Percentage)
	}
}
