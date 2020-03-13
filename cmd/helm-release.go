package cmd

import (
	"fmt"

	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// helmReleaseCmd represents the helmRelease command
var helmReleaseCmd = &cobra.Command{
	Use:   "helm-release",
	Short: "Generate Helm Release custom resource",
	Long:  `create-release generates a valid Helm Release custom resource to be commited to a git repo to be used by Helm Operator.`,

	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		namespace, _ := cmd.Flags().GetString("namespace")
		helmURL, _ := cmd.Flags().GetString("helm-url")
		valuesFile, _ := cmd.Flags().GetString("values-file")

		generateHelmRelease(name, namespace, helmURL, valuesFile)
	},
}

func init() {
	rootCmd.AddCommand(helmReleaseCmd)

	helmReleaseCmd.Flags().StringP("name", "n", "", "name for the application")
	helmReleaseCmd.Flags().StringP("namespace", "", "default", "namespace field for CR")
	helmReleaseCmd.Flags().StringP("helm-url", "", "", "where is placed Helm Chart")
	helmReleaseCmd.Flags().StringP("values-file", "v", "values.yaml", "path to values file")
}

type Metadata struct {
	Name        string   `yaml:"name"`
	Namespace   string   `yaml:"namespace"`
	Annotations []string `yaml:"annotations"`
}

type Chart struct {
	Repository string `yaml:"repository"`
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
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

func generateHelmRelease(name string, namespace string, repository string, valuesFile string) {

	var values map[string]interface{}
	bs, err := ioutil.ReadFile(valuesFile)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(bs, &values); err != nil {
		panic(err)
	}

	t := HelmReleaseCR{
		ApiVersion: "helm.fluxcd.io/v1",
		Kind:       "HelmRelease",
		Metadata: Metadata{
			Name:        name,
			Namespace:   namespace,
			Annotations: []string{"fluxcd.io/automated: true"},
		},
		Spec: Spec{
			ReleaseName: name,
			Chart: Chart{
				Repository: repository,
				Name:       name,
				Version:    "0.0.1",
			},
			Values: values,
		},
	}

	output, err := yaml.Marshal(t)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(name+".yaml", output, 755)
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(output))
}
