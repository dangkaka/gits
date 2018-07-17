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
			Name:   "add",
			Usage:  "Git add",
			Action: Add,
		},
		{
			Name:   "commit",
			Usage:  "Git commit",
			Action: Commit,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Add(c *cli.Context) error {
	cmd := exec.Command("git", "add", ".")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("\n*** Success *** %s\n", out)
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
	fmt.Printf("\n*** Success *** %s\n", out)
	return nil
}
