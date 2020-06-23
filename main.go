package main

import (
	"github.com/ferrandinand/create-release/command"

	"github.com/giantswarm/micrologger"
)

func main() {
	var err error

	// Create a new logger which is used by all packages.
	var newLogger micrologger.Logger
	{
		newLogger, err = micrologger.New(micrologger.Config{})
		if err != nil {
			panic(err)
		}
	}

	var newCommand *command.Command
	{
		c := command.Config{
			Logger: newLogger,
		}

		newCommand, err = command.New(c)
		if err != nil {
			panic(err)
		}
	}

	newCommand.CobraCommand().Execute()
}
