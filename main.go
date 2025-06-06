// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/mod/modfile"
	"golang.org/x/oauth2"

	"github.com/esacteksab/httpcache"
	"github.com/esacteksab/httpcache/diskcache"
)

var (
	age           int
	score         float64
	createdAtTime time.Time
	client        *github.Client
)

type Repo struct {
	Module    Module
	CreatedAt time.Time
	Fork      bool
	UpdatedAt *github.Timestamp
}

type Module struct {
	Module string
	Name   string
	Host   string
	Owner  string
	Repo   string
}

type CachingTransport struct {
	Transport http.RoundTripper
}

func (t *CachingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Check if response came from cache
	if resp.Header.Get(httpcache.XFromCache) == "1" {
		fmt.Println("✅ RESPONSE SERVED FROM CACHE")
	} else {
		fmt.Println("❌ RESPONSE NOT FROM CACHE")
	}

	// Check for auth header (don't print the actual token)
	if req.Header.Get("Authorization") != "" {
		fmt.Println("🔑 Request contains Authorization header")
	} else {
		fmt.Println("⚠️  No Authorization header found!")
	}

	return resp, nil
}

func main() {
	// Read the go.mod file
	fileBytes, err := os.ReadFile("go.mod")
	if err != nil {
		log.Fatal(err)
	}

	// Parse the go.mod file
	mod, err := modfile.Parse("go.mod", fileBytes, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Print require statements
	fmt.Println("Require Statements:")
	for _, req := range mod.Require {
		// fmt.Printf("  %s %s\n", req.Mod.Path, req.Mod.Version)
		module := parseModule(req.Mod.Path)
		getRepo(module)
	}
}

func parseModule(module string) (m Module) {
	s := strings.Split(module, "/")
	if len(s) > 2 { //nolint:mnd
		m = Module{}

		m.Name = module
		m.Host = s[0]
		m.Owner = s[1]
		m.Repo = s[2]
	}
	return m
}

func getRepo(m Module) (r Repo) { //nolint:funlen
	projectCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("Cache Directory: %s not found", projectCacheDir)
	}

	clydeDir := "clyde"

	// Need to create clydeDir in $XDG_CACHE_HOME
	if err := os.MkdirAll(projectCacheDir+"/"+clydeDir, 0o750); err != nil { //nolint:mnd
		log.Fatalf("Could not create clyde directory: %v", err)
	}

	if err := os.MkdirAll(projectCacheDir, 0o750); err != nil { //nolint:mnd
		log.Fatalf("failed to create cache directory: %v\n", err)
	}

	cache := diskcache.New(projectCacheDir + "/" + clydeDir)

	if m.Host == "github.com" {
		// Get the GitHub token
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			fmt.Println(
				"⚠️ No GITHUB_TOKEN found in environment. Using unauthenticated client with lower rate limits.",
			)
		} else {
			fmt.Println("🔑 Found GITHUB_TOKEN in environment.")
		}

		ctx := context.Background()

		if token != "" {
			cacheTransport := httpcache.NewTransport(cache)
			ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
			authTransport := &oauth2.Transport{
				Base:   cacheTransport,
				Source: ts,
			}

			cachingTransport := &CachingTransport{Transport: authTransport}
			httpClient := &http.Client{Transport: cachingTransport}
			client = github.NewClient(httpClient)

			fmt.Println("🔐 Using authenticated GitHub client")
		} else {
			// For unauthenticated requests, we can use a simpler chain
			// Debug -> Cache -> HTTP
			cacheTransport := httpcache.NewTransport(cache)
			debugTransport := &CachingTransport{Transport: cacheTransport}
			httpClient := &http.Client{Transport: debugTransport}
			client = github.NewClient(httpClient)

			fmt.Println("🔓 Using unauthenticated GitHub client")
		}

		repo, resp, err := client.Repositories.Get(ctx, m.Owner, m.Repo)
		if _, ok := err.(*github.RateLimitError); ok {
			fmt.Println("hit rate limit")
		} else if _, ok := err.(*github.AbuseRateLimitError); ok {
			fmt.Println("high secondary rate limit")
		}

		rate := resp.Rate
		fmt.Printf("Rate limit: %d/%d, resets at %v\n",
			rate.Remaining,
			rate.Limit,
			(rate.Reset.Time).Local().Format("2006-01-02-15:04:05"))

		// Check if rate limit is for authenticated or unauthenticated requests
		if rate.Limit >= 5000 { //nolint:mnd
			fmt.Println("✅ Using authenticated rate limits (5000+/hour)")
		} else if rate.Limit <= 60 { //nolint:mnd
			fmt.Println("❌ Using unauthenticated rate limits (60/hour)")
		}

		r = Repo{}

		// If repo.Fork is a *bool, dereference it first
		if repo.Fork != nil {
			r.Fork = *repo.Fork
		} else {
			r.Fork = false // maybe?
		}

		// Converting *github.Timestamp to time.Time so I can manipulate it later with .Sub()
		if repo.CreatedAt != nil {
			createdAtTime = repo.CreatedAt.Time
		} else {
			createdAtTime = time.Time{}
		}

		r.Module = m
		r.CreatedAt = createdAtTime
		r.UpdatedAt = repo.UpdatedAt

		fmt.Printf("Module is: %s\n", r.Module.Name)

		if r.Fork {
			fmt.Println("🍴 Repo is a fork")
		} else if !r.Fork {
			fmt.Println("🍰 Repo is not a fork")
		}

		fmt.Printf("Repo was created at: %s\n", r.CreatedAt.Format("2006-01-02"))
		fmt.Printf("Repo last updated at: %s\n", r.UpdatedAt.Format("2006-01-02"))

		calculate(createdAtTime, r.Fork)
	} else {
		fmt.Printf("❓ Module %s not hosted on GitHub.", m.Name)
	}
	fmt.Println("\nB===================================================D")
	// Because go vet doesn't like the redundant \n on the line above. Heere's an
	// empty line. Check. Mate.
	fmt.Println("")
	return r
}

func calculate(created time.Time, fork bool) float64 {
	now := time.Now()
	difference := now.Sub(created)
	days := int(difference.Hours() / 24) //nolint:mnd
	fmt.Printf("Module is %d days old.\n", days)

	age = 30
	score = 0

	if fork {
		score += 50
	}

	if days <= 1 {
		score += 50
	}

	if days < age {
		score += (float64(days) / float64(age)) * 50 //nolint:mnd
	}

	if score >= 100 { //nolint:mnd
		fmt.Printf("⛔ Module has a score of: %.2f out of 100.\n", score)
	} else if score >= 50 && score < 100 {
		fmt.Printf("💩 Module has a score of: %.2f out of 100.\n", score)
	} else if score < 50 { //nolint:mnd
		fmt.Printf("✨ Module has a score of: %.2f out of 100.\n", score)
	}

	fns := fmt.Sprintf("%.2f", score)
	fn, _ := strconv.ParseFloat(fns, 64)
	return fn
}
