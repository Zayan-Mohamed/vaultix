package main

import (
	"fmt"
	"os"

	"github.com/Zayan-Mohamed/vaultix/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error

	switch command {
	case "init":
		err = cli.Init(args)
	case "add":
		err = cli.Add(args)
	case "list":
		err = cli.List(args)
	case "extract":
		err = cli.Extract(args)
	case "drop":
		err = cli.Drop(args)
	case "remove":
		err = cli.Remove(args)
	case "clear":
		err = cli.Clear(args)
	case "recover":
		err = cli.Recover(args)
	case "help", "-h", "--help":
		cli.PrintUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n\n", command)
		cli.PrintUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
