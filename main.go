package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	FAILEDWITHOUTPUT  = `¯\_(ツ)_/¯` + " \n %s \n"
	SUCCESSWITHOUTPUT = `ᕙ_(⇀‸↼)_ᕗ` + " \n %s \n"
)

var jiraProjects = []string{"SHOP-", "FEED-"}

func main() {
	app := cli.NewApp()
	app.Name = "gits"
	app.Usage = "Git with convenient + short commands"
	app.Version = "1.0"

	app.Commands = []cli.Command{
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Git status",
			Action:  Status,
		},
		{
			Name:    "pull",
			Aliases: []string{"pl"},
			Usage:   "Git pull",
			Action:  Pull,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "rebase, r",
					Usage: "git pull rebase",
				},
			},
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "f",
					Usage: "force push",
					},
			},
		},
		{
			Name:    "commitpush",
			Aliases: []string{"cp"},
			Usage:   "Git commit and push",
			Action:  CommitAndPush,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "f",
					Usage: "force push",
				},
			},
		},
		{
			Name:    "gitignore",
			Usage:   "Add default .gitignore file",
			Action:  AddGitIgnore,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Status(c *cli.Context) error {
	cmd := exec.Command("git", "status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf(FAILEDWITHOUTPUT, out)
		return err
	}
	fmt.Printf(SUCCESSWITHOUTPUT, out)
	return nil
}

func Pull(c *cli.Context) error {
	args := []string{"pull"}
	isRebase := c.Bool("r")
	if isRebase {
		args = append(args, "--rebase")
	}
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf(FAILEDWITHOUTPUT, out)
		return err
	}
	fmt.Printf(SUCCESSWITHOUTPUT, out)
	return nil
}

func Add(c *cli.Context) error {
	cmd := exec.Command("git", "add", ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf(FAILEDWITHOUTPUT, out)
		return err
	}
	fmt.Printf(SUCCESSWITHOUTPUT, out)
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
		fmt.Printf(FAILEDWITHOUTPUT, out)
		return err
	}
	fmt.Printf(SUCCESSWITHOUTPUT, out)
	return nil
}

func Push(c *cli.Context) error {
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return cli.NewExitError("Could not get current branch", 1)
	}
	args := []string{"push", "origin", currentBranch}
	isForce := c.Bool("f")
	if isForce {
		args = append(args, "-f")
	}
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf(FAILEDWITHOUTPUT, out)
		return err
	}
	fmt.Printf(SUCCESSWITHOUTPUT, out)
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

func AddGitIgnore(c *cli.Context) error {
	file, err := os.OpenFile(".gitignore", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf(FAILEDWITHOUTPUT, err)
		return err
	}
	file.WriteString(".idea")
	file.WriteString("\n")
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
