// Package command implements the root command for the command line tool.
package command

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"

	"github.com/ferrandinand/create-release/command/create"
)

// Config represents the configuration used to create a new root command.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger
}

type Command struct {
	// Internals.
	cobraCommand *cobra.Command
}

// New creates a new root command.
func New(config Config) (*Command, error) {
	var err error

	newCommand := &Command{
		cobraCommand: nil,
	}

	// apply settings from environment variables



	newCommand.cobraCommand = &cobra.Command{
		Use:   config.Name,
		Short: config.Description,
		Long:  config.Description,
		Run:   newCommand.Execute,
	}


	var createCommand *create.Command
	{
		c := create.Config{
			Logger: config.Logger,
		}

		createCommand, err = create.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newCommand.cobraCommand.AddCommand(createCommand.CobraCommand())

	return newCommand, nil
}

// CobraCommand returns the spf13/cobra command
func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

// Execute is called to actuall run the main command
func (c *Command) Execute(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, nil)
}
