package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	log.SetFlags(10)
	app := cli.NewApp()
	app.Commands = []cli.Command{
		GenerateKeyCommand,
		ForwardCommand,
		QueryCommand,
		ReceiveCommand,
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
