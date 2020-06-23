package project

import (
	"context"
	"fmt"
	"os"

	"github.com/ferrandinand/create-release/service/azuredevops"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
)

var flags = &Flags{}

// Config represents the configuration used to create a new draughtsman
// command.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger
}

type Command struct {
	// Dependencies.
	logger micrologger.Logger

	// Internals.
	cobraCommand *cobra.Command
}

// New creates a new configured draughtsman command.
func New(config Config) (*Command, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	c := &Command{
		logger: config.Logger,

		// Internals.
		cobraCommand: nil,
	}

	c.cobraCommand = &cobra.Command{
		Use:   "project",
		Short: "Managing projects with CLI",
		Long:  `Creation project command is a wrapper that triggers a pipeline to created different resources.`,
		Run:   c.Execute,
	}

	flags.AZDVOToken = os.Getenv("AZDVO_TOKEN")
	flags.AZDVOUser = os.Getenv("AZDVO_USER")

	c.cobraCommand.PersistentFlags().StringVarP(&flags.projectName, "name", "n", "", "Name of the project that will be created")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.projectType, "type", "t", "container-project", "Template project to be used")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.projectPlatform, "platform", "p", "azuredevops", "In which platform will be created the project")

	return c, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	err := flags.Validate()
	if err != nil {
		c.logger.Log("level", "error", "message", err.Error(), "stack", fmt.Sprintf("%#v", err), "verbosity", 0)
		os.Exit(1)
	}

	err = c.execute()
	if err != nil {
		c.logger.Log("level", "error", "message", err.Error(), "stack", fmt.Sprintf("%#v", err), "verbosity", 0)
		os.Exit(1)
	}
}

func (c *Command) execute() error {
	fmt.Printf("Running creation project %s.\n", flags.projectName)

	var err error

	ctx := context.Background()

	var azdvoService *azuredevops.Service
	{
		config := azuredevops.Config{
			ADVOUser:  flags.AZDVOUser,
			AZDOToken: flags.AZDVOToken,
		}

		azdvoService, err = azuredevops.New(config)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	if flags.projectPlatform == "azuredevops" {
		buildId := 1
		buildOrganization := "unicc-opd"
		buildMasterProject := "opd_dli"
		buildParameters := "{\"definition\": {\"id\": %d }, \"parameters\": \"{\\\"project_name\\\": \\\"%s\\\",\\\"project_type\\\": \\\"%s\\\"}\"}"

		params := fmt.Sprintf(buildParameters, buildId, flags.projectName, flags.projectType)

		err = azdvoService.RunPipeline(ctx, buildOrganization, buildMasterProject, params)
		if err != nil {
			fmt.Printf("Creation failed. Reason: %s\n", err.Error())
			return err
		}
		fmt.Printf("Creating project %s in %s organization", flags.projectName, buildOrganization)
	} else {
		fmt.Print("Only azuredevops project creation implemented")
	}

	return nil
}
