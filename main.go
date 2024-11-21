package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-github/v66/github"
)

var client *github.Client

func addRunner(owner, repo string) {
	r, _, err := client.Actions.CreateRegistrationToken(context.Background(), repo, owner)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully created registration token")

	token := r.GetToken()
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("bash",
		"./config.sh",
		"--url", fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		"--token", token,
		"--name", fmt.Sprintf("runner-%s-%d", host, time.Now().Unix()))
	runnerLoc := os.Getenv("RUNNER_LOCATION")
	cmd.Env = os.Environ()
	if runnerLoc != "" {
		cmd.Dir = runnerLoc
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func run() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("No command provided")
	}
	bin := args[0]
	left := make([]string, 0, len(args)-1)
	if len(args) != 1 {
		left = args[1:]
	}

	cmd := exec.Command(bin, left...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}

func removeRunner(owner, repo string) {
	r, _, err := client.Actions.CreateRemoveToken(context.Background(), repo, owner)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully created remove token")

	token := r.GetToken()
	cmd := exec.Command("bash", "./config.sh", "remove", "--token", token)
	runnerLoc := os.Getenv("RUNNER_LOCATION")
	cmd.Env = os.Environ()
	if runnerLoc != "" {
		cmd.Dir = runnerLoc
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func prepareRemove(owner, repo string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		removeRunner(owner, repo)
		os.Exit(0)
	}()
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("GITHUB_TOKEN is required")
	}

	client = github.NewClient(nil).WithAuthToken(token)
	_, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		panic(err)
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		panic("GITHUB_REPOSITORY is required")
	}

	owner, repo, found := strings.Cut(repo, "/")
	if !found {
		panic("Invalid GITHUB_REPOSITORY " + repo)
	}

	addRunner(owner, repo)
	prepareRemove(owner, repo)

	run()
}
