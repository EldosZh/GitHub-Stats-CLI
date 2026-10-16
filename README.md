# GitHub Stats CLI

A fast and simple command-line tool written in Go for viewing GitHub profile and repository statistics directly from your terminal.

## Features

- GitHub profile and repository statistics
- Top repositories, sorted by stars, forks, or update date
- Language statistics, JSON output, and pagination
- Optional authenticated requests and GitHub API rate-limit handling

## Tech Stack

- Go
- GitHub REST API
- Go standard library

## Installation

```sh
git clone https://github.com/YOUR_USERNAME/github-stats-cli.git
cd github-stats-cli
go build -o github-stats .
```

On Windows, the executable may be named `github-stats.exe`.

## Usage

```sh
github-stats torvalds
github-stats torvalds --top 10
github-stats torvalds --sort forks
github-stats torvalds --json
github-stats --help
github-stats --version
```

Languages are shown by default. Language counts represent the number and share of repositories with each primary language; they do not measure lines of code. Repositories without a primary language are omitted from this section.

## Authentication

`GITHUB_TOKEN` is optional. Setting it allows authenticated API requests, which have higher rate limits.

Linux/macOS:

```sh
export GITHUB_TOKEN="your_token_here"
```

Windows PowerShell:

```powershell
$env:GITHUB_TOKEN="your_token_here"
```

## Example Output

Illustrative values only:

```text
GitHub Stats
────────────────────────────────────────
User:          octocat
Name:          The Octocat
Followers:     12000
Public Repos:  8

Repository Stats
Repositories:  8
Original Repos: 6
Forked Repos:  2
Total Stars:   420
Total Forks:   95
```

## How It Works

The CLI requests public data from the GitHub REST API, decodes the JSON into Go structs, calculates repository statistics, then formats the results for the terminal or as JSON.

## Project Structure

```text
main.go                  CLI parsing and output
internal/github/         GitHub API client and response types
internal/stats/           Aggregation, sorting, and language counts
```

## Testing

```sh
go test ./...
```

Client tests use `httptest`, so the test suite works offline.

## What I Learned

- Consuming REST APIs with `net/http` and decoding JSON
- Handling pagination, optional authentication, and API rate limits
- Aggregating and sorting data
- Building a CLI with useful error handling
- Testing HTTP clients with `httptest`

## Future Improvements

Contribution statistics, repository activity, caching, terminal colors, or comparing two users.
