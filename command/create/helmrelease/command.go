package helmrelease

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
		Use:   "helm-release",
		Short: "Create helm-release CR",
		Long: `Creation of a helm-release CR to use in a k8s cluster running helm-operator.
Given a values yaml file is merged with a yaml compliance with helm-release kubernetes.`,
		Run: c.Execute,
	}

	c.cobraCommand.PersistentFlags().StringVarP(&flags.name, "name", "n", "", "Name of the project that will be created")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.namespace, "namespace", "", "default", "Namespace field for CR")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.gitRepoURL, "git-repo-url", "", "", "Git repo where is placed Helm Chart")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.gitRepoRef, "git-repo-ref", "", "", "Git repo reference for the chart, can be a branch or a tag")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.chartName, "chart-name", "c", "", "Helm Chart name inside the git repo url")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.valuesFile, "values-file", "v", "values.yaml", "Path to values file")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.registry, "registry", "R", "", "Registry where docker image is hosted")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.repository, "repository", "r", "", "Repository where docker image is hosted")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.outputFile, "output-file", "o", "", "Output file name without extension")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.tag, "tag", "t", "", "Docker image tag")

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
	generateHelmRelease(flags.name, flags.namespace, flags.gitRepoURL, flags.chartName, flags.gitRepoRef, flags.valuesFile, flags.registry, flags.repository, flags.outputFile, flags.tag)

	return nil
}

type Image struct {
	Registry   string `yaml:"registry"`
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

type ImageSpec struct {
	Image Image
}

type Metadata struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Annotations map[string]string `yaml:"annotations"`
}

type Chart struct {
	Git  string `yaml:"git"`
	Path string `yaml:"path"`
	Ref  string `yaml:"ref"`
}

type Spec struct {
	ReleaseName string `yaml:"releaseName"`
	Chart       Chart
	Values      map[string]interface{} `yaml:"values"`
}

type HelmReleaseCR struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   Metadata
	Spec       Spec
}

func generateHelmRelease(name string, namespace string, gitRepoURL string, chartName string, gitRef string, valuesFile string, registry string, repository string, ouputFile string, tag string) {

	var values map[string]interface{}

	imageValues := &ImageSpec{
		Image: Image{
			Registry:   registry,
			Repository: repository,
			Tag:        tag,
		},
	}

	var imageInterface map[string]interface{}
	inrec, _ := yaml.Marshal(imageValues)
	yaml.Unmarshal(inrec, &imageInterface)

	//Read values files
	bs, err := ioutil.ReadFile(valuesFile)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(bs, &values); err != nil {
		panic(err)
	}

	//Override image values
	for k, v := range imageInterface {
		values[k] = v
	}

	t := HelmReleaseCR{
		ApiVersion: "helm.fluxcd.io/v1",
		Kind:       "HelmRelease",
		Metadata: Metadata{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{"fluxcd.io/automated": "false"},
		},
		Spec: Spec{
			ReleaseName: name,
			Chart: Chart{
				Git:  gitRepoURL,
				Path: chartName,
				Ref:  gitRef,
			},
			Values: values,
		},
	}

	output, err := yaml.Marshal(t)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(ouputFile+".yaml", output, 755)
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(output))
}
