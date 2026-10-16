package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserAndAuth(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing auth header")
		}
		fmt.Fprint(w, `{"login":"octocat","followers":7}`)
	}))
	defer s.Close()
	u, err := NewClient(s.URL, "test-token").GetUser(context.Background(), "octocat")
	if err != nil || u.Login != "octocat" || u.Followers != 7 {
		t.Fatalf("user=%+v err=%v", u, err)
	}
}

func TestAPIErrorResponses(t *testing.T) {
	for _, tc := range []struct {
		code int
		body string
		want string
	}{{404, "{}", "not found"}, {500, "{}", "500"}, {200, "{", "decode"}, {403, `{"message":"rate limit"}`, "rate limit"}} {
		t.Run(tc.want, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(tc.code); fmt.Fprint(w, tc.body) }))
			defer s.Close()
			_, err := NewClient(s.URL, "").GetUser(context.Background(), "nobody")
			if err == nil || !contains(err.Error(), tc.want) {
				t.Fatalf("err=%v want %q", err, tc.want)
			}
		})
	}
}

func TestPagination(t *testing.T) {
	calls := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Query().Get("per_page") != "100" {
			t.Errorf("per_page=%s", r.URL.Query().Get("per_page"))
		}
		if calls == 1 {
			w.Write([]byte("["))
			for i := 0; i < 100; i++ {
				if i > 0 {
					w.Write([]byte(","))
				}
				fmt.Fprintf(w, `{"name":"r%d"}`, i)
			}
			w.Write([]byte("]"))
			return
		}
		fmt.Fprint(w, `[{"name":"last"}]`)
	}))
	defer s.Close()
	repos, err := NewClient(s.URL, "").GetRepositories(context.Background(), "user")
	if err != nil || len(repos) != 101 || calls != 2 {
		t.Fatalf("len=%d calls=%d err=%v", len(repos), calls, err)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
