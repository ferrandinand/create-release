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

const (
	defaultAZDVOURL = "https://dev.azure.com/unicc-opd/opd_dli/"
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
		// Dependencies.
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

	c.cobraCommand.PersistentFlags().StringP("project-name", "n", "", "Registry where docker image is hosted")
	c.cobraCommand.PersistentFlags().StringP("project-type", "p", "", "Repository where docker image is hosted")
	c.cobraCommand.PersistentFlags().StringP("project-platform", "", "azuredevops", "Output file name without extension")
	c.cobraCommand.PersistentFlags().StringP("organization", "", "", "Azure Devops organization where the project will be created.")

	return c, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	err := flags.Validate()
	if err != nil {
		fmt.Print("pedo")
		c.logger.Log("level", "error", "message", err.Error(), "stack", fmt.Sprintf("%#v", err), "verbosity", 0)
		os.Exit(1)
	}

	err = c.execute()
	if err != nil {
		fmt.Print("pedo2")
		c.logger.Log("level", "error", "message", err.Error(), "stack", fmt.Sprintf("%#v", err), "verbosity", 0)
		os.Exit(1)
	}
}

func (c *Command) execute() error {
	fmt.Print("WARNING: This command is deprecated for updating draughtsman configuration in a running installation.\n")
	fmt.Print("Use 'opsctl update draughtsman' instead.\n\n")

	var err error

	ctx := context.Background()

	var azuredevops *azuredevops.Service
	{
		config := azuredevops.Config{
			Logger: c.logger,

			ADVOURL:   defaultAZDVOURL,
			ADVOUser:  flags.AZDVOUser,
			AZDOToken: flags.AZDVOToken,
		}

		azuredevops, err = azuredevops.New(config)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	buildParameters := "{\"definition\": {\"id\": %s }, \"parameters\": \"{\\\"project_name\\\": \\\"%s\\\",\\\"project_type\\\": \\\"%s\\\"}\"}"
	params := fmt.Sprintf(buildParameters, "test", "test2")

	err = azuredevops.runDevopsPipeline(ctx, "unicc-opd", "opd_dli", params)
	if err != nil {
		panic(err)
	}

	return nil
}

//func runDevopsPipeline(projectName string, organization string) ADVOResponse {
//
//	var jsonFormat = "{\"definition\": {\"id\": 1 }, \"parameters\": \"{\\\"project_name\\\": \\\"" + projectName + "\\\" }\"}"
//
//	var jsonStr = []byte(jsonFormat)
//
//	fmt.Println(string(jsonStr))
//	var baseUrl = "https://dev.azure.com/" + organization + "/opd_dli/_apis/build/builds?api-version=5.1"
//	fmt.Println(string(baseUrl))
//	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonStr))
//
//	pat := os.Getenv("AZDVO_TOKEN")
//	user := os.Getenv("AZDVO_USER")
//
//	req.SetBasicAuth(user, pat)
//	req.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//
//	resp, err := client.Do(req)
//
//	if err != nil {
//		panic(err)
//	}
//
//	responseData, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(string(responseData))
//
//	data := ADVOResponse{}
//	json.Unmarshal([]byte(responseData), &data)
//
//	return data
//
//}
