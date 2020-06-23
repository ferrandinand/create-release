package create

import (
	"github.com/ferrandinand/create-release/command/create/helmrelease"
	"github.com/ferrandinand/create-release/command/create/project"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
)

// Config represents the configuration used to update a new update command.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger
}

type Command struct {
	// Internals.
	cobraCommand *cobra.Command
}

func New(config Config) (*Command, error) {
	var err error

	newCommand := &Command{
		// Internals.
		cobraCommand: nil,
	}

	newCommand.cobraCommand = &cobra.Command{
		Use:   "create",
		Short: "Create option to manage operations",
		Long:  "Manage operation creation from a centralized CLI",
		Run:   newCommand.Execute,
	}

	var projectCommand *project.Command
	{
		c := project.Config{
			Logger: config.Logger,
		}

		projectCommand, err = project.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var helmReleaseCommand *helmrelease.Command
	{
		c := helmrelease.Config{
			Logger: config.Logger,
		}

		helmReleaseCommand, err = helmrelease.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newCommand.cobraCommand.AddCommand(projectCommand.CobraCommand())
	newCommand.cobraCommand.AddCommand(helmReleaseCommand.CobraCommand())

	return newCommand, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, nil)
}
