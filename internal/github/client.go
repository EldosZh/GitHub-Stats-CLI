package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.github.com"

type User struct {
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	Bio         string    `json:"bio"`
	Location    string    `json:"location"`
	Company     string    `json:"company"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	PublicRepos int       `json:"public_repos"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	HTMLURL     string    `json:"html_url"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	Watchers    int       `json:"watchers_count"`
	OpenIssues  int       `json:"open_issues_count"`
	Fork        bool      `json:"fork"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func NewClient(baseURL, token string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), token: token, httpClient: &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) GetUser(ctx context.Context, username string) (User, error) {
	var user User
	err := c.getJSON(ctx, "/users/"+url.PathEscape(username), &user)
	return user, err
}

func (c *Client) GetRepositories(ctx context.Context, username string) ([]Repository, error) {
	var all []Repository
	for page := 1; ; page++ {
		endpoint := fmt.Sprintf("/users/%s/repos?per_page=100&page=%d&sort=updated", url.PathEscape(username), page)
		var repos []Repository
		if err := c.getJSON(ctx, endpoint, &repos); err != nil {
			return nil, err
		}
		all = append(all, repos...)
		if len(repos) < 100 {
			return all, nil
		}
	}
}

func (c *Client) getJSON(ctx context.Context, endpoint string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("create GitHub request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "github-stats-cli")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request GitHub API: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("GitHub user or resource not found")
		}
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			if resp.Header.Get("X-RateLimit-Remaining") == "0" || resp.StatusCode == http.StatusTooManyRequests || strings.Contains(strings.ToLower(string(body)), "rate limit") {
				return fmt.Errorf("GitHub API rate limit exceeded; set GITHUB_TOKEN to increase the limit")
			}
		}
		return fmt.Errorf("GitHub API returned %s", resp.Status)
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 16<<20)).Decode(target); err != nil {
		return fmt.Errorf("decode GitHub response: %w", err)
	}
	return nil
}
