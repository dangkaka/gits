package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"strings"
)

var jiraProjects = []string{"SHOP-", "FEED-"}

func main() {
	app := cli.NewApp()
	app.Name = "gits"
	app.Usage = "Git with convenient + short commands"
	app.Version = "1.0"

	app.Commands = []cli.Command{
		{
			Name:    "pull",
			Aliases: []string{"pl"},
			Usage:   "Git pull",
			Action:  Pull,
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Git add",
			Action:  Add,
		},
		{
			Name:    "commit",
			Aliases: []string{"c"},
			Usage:   "Git commit",
			Action:  Commit,
		},
		{
			Name:    "push",
			Aliases: []string{"p"},
			Usage:   "Git push",
			Action:  Push,
		},
		{
			Name:    "commitpush",
			Aliases: []string{"cp"},
			Usage:   "Git commit and push",
			Action:  CommitAndPush,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Pull(c *cli.Context) error {
	cmd := exec.Command("git", "pull")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("*** FAILED ***")
		return err
	}
	fmt.Printf("*** DONE *** \n")
	return nil
}

func Add(c *cli.Context) error {
	cmd := exec.Command("git", "add", ".")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("*** FAILED ***")
		return err
	}
	fmt.Printf("*** DONE *** \n")
	return nil
}

func Commit(c *cli.Context) error {
	msg := c.Args().Get(0)
	if msg == "" {
		return cli.NewExitError("Empty commit message", 1)
	}
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return cli.NewExitError("Could not get current branch", 1)
	}
	if isJiraProject(currentBranch) {
		msg = currentBranch + ": " + msg
	}
	cmd := exec.Command("git", "commit", "-m", msg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("*** FAILED *** \n %s \n", out)
		return err
	}
	fmt.Printf("*** DONE *** \n %s \n", out)
	return nil
}

func Push(c *cli.Context) error {
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return cli.NewExitError("Could not get current branch", 1)
	}
	cmd := exec.Command("git", "push", "origin", currentBranch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("*** FAILED ***")
		return err
	}
	fmt.Printf("*** DONE *** \n %s \n", out)
	return nil
}

func CommitAndPush(c *cli.Context) error {
	err := Commit(c)
	if err != nil {
		return err
	}
	err = Push(c)
	if err != nil {
		return err
	}
	return nil
}

func isJiraProject(branch string) bool {
	for _, prefix := range jiraProjects {
		if strings.HasPrefix(branch, prefix) {
			return true
		}
	}
	return false
}

func getCurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
