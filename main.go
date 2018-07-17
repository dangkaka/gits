package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"errors"
	"fmt"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitfmt"
	app.Description = "Git format"
	app.Version = "1.0"

	app.Commands = []cli.Command{
		{
			Name:   "pull",
			Usage:  "Git pull",
			Action: Pull,
		},
		{
			Name:   "add",
			Usage:  "Git add",
			Action: Add,
		},
		{
			Name:   "commit",
			Usage:  "Git commit",
			Action: Commit,
		},
		{
			Name:   "push",
			Usage:  "Git push",
			Action: Push,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Pull(c *cli.Context) error {
	cmd := exec.Command("git", "pull")
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("*** DONE *** \n")
	return nil
}

func Add(c *cli.Context) error {
	cmd := exec.Command("git", "add", ".")
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("*** DONE *** \n")
	return nil
}

func Commit(c *cli.Context) error {
	msg := c.Args().Get(0)
	if msg == "" {
		return errors.New("Empty commit message")
	}
	cmd := exec.Command("git", "commit", "-m", msg)
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("*** DONE *** \n %s \n", out)
	return nil
}

func Push(c *cli.Context) error {
	cmd := exec.Command("git", "push")
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("*** DONE *** \n")
	return nil
}
