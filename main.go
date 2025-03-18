// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/mod/modfile"
)

var (
	age           int
	score         int
	createdAtTime time.Time
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

func main() {
	// Read the go.mod file
	fileBytes, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}

	// Parse the go.mod file
	mod, err := modfile.Parse("go.mod", fileBytes, nil)
	if err != nil {
		panic(err)
	}

	// Print the module path
	// fmt.Println("Module Path:", mod.Module.Mod.Path)

	// Print the Go version
	// fmt.Println("Go Version:", mod.Go.Version)

	// Print require statements
	fmt.Println("Require Statements:")
	for _, req := range mod.Require {
		// fmt.Printf("  %s %s\n", req.Mod.Path, req.Mod.Version)
		module := parseModule(req.Mod.Path)
		getRepo(module)
	}

	// Print replace statements
	fmt.Println("Replace Statements:")
	for _, rep := range mod.Replace {
		fmt.Printf("  %s => %s %s\n", rep.Old.Path, rep.New.Path, rep.New.Version)
	}

	// Print exclude statements
	fmt.Println("Exclude Statements:")
	for _, exc := range mod.Exclude {
		fmt.Printf("  %s %s\n", exc.Mod.Path, exc.Mod.Version)
	}
}
func parseModule(module string) (m Module) {
	s := strings.Split(module, "/")

	m = Module{}

	m.Name = module
	m.Host = s[0]
	m.Owner = s[1]
	m.Repo = s[2]

	// fmt.Println(s)
	// fmt.Printf("Host is: %s\n", Host)
	// fmt.Printf("Owwner is: %s\n", Owner)
	// fmt.Printf("Repo is: %s\n", Repo)
	// fmt.Printf("Host is: %s\n", m.Host)
	// fmt.Printf("Owwner is: %s\n", m.Owner)
	// fmt.Printf("Repo is: %s\n", m.Repo)

	return m
}

func getRepo(m Module) (r Repo) {
	client := github.NewClient(nil)
	if m.Host != "github.com" {
		// we don't do anything yet
	} else {
		repo, resp, err := client.Repositories.Get(context.Background(), m.Owner, m.Repo)
		if _, ok := err.(*github.AbuseRateLimitError); ok {
			log.Println("hit rate limit")
		} else if _, ok := err.(*github.AbuseRateLimitError); ok {
			log.Println("high secondary rate limit")
		}
		if err != nil {
			fmt.Println(resp)
			log.Fatal(err)
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

		score := calculate(createdAtTime, r.Fork)

		fmt.Printf("Module is: %s\n", r.Module.Name)
		fmt.Printf("Repo was created at: %s\n", r.CreatedAt.Format("2006-01-02"))
		fmt.Printf("Repo is a fork: %t\n", r.Fork)
		fmt.Printf("Repo last updated at: %s\n", r.UpdatedAt.Format("2006-01-02"))
		fmt.Printf("Module has a score of: %s out of 100.\n", score)
	}
	return r
}

func calculate(created time.Time, fork bool) string {
	now := time.Now()
	difference := now.Sub(created)
	days := int(difference.Hours() / 24) //nolint:mnd

	age = 30
	var score float64
	score = 0

	if fork {
		score += 50
	}

	if days <= 1 {
		score += 50
	}

	if days < age {
		score += (float64(days) / float64(age)) * 50
	}

	fmt.Printf("Module is %d days old.\n", days)
	fn := fmt.Sprintf("%.2f", score)
	return fn
}
