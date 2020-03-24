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
		gitRepoURL, _ := cmd.Flags().GetString("git-repo-url")
		chartName, _ := cmd.Flags().GetString("chart-name")
		valuesFile, _ := cmd.Flags().GetString("values-file")
		registry, _ := cmd.Flags().GetString("registry")
		repository, _ := cmd.Flags().GetString("repository")
		tag, _ := cmd.Flags().GetString("tag")

		generateHelmRelease(name, namespace, gitRepoURL, chartName, valuesFile, registry, repository, tag)
	},
}

func init() {
	rootCmd.AddCommand(helmReleaseCmd)

	helmReleaseCmd.Flags().StringP("name", "n", "", "name for the application")
	helmReleaseCmd.Flags().StringP("namespace", "", "default", "namespace field for CR")
	helmReleaseCmd.Flags().StringP("git-repo-url", "", "", "Git repo where is placed Helm Chart")
	helmReleaseCmd.Flags().StringP("chart-name", "c", "", "Helm Chart name inside the git repo url")
	helmReleaseCmd.Flags().StringP("values-file", "v", "values.yaml", "path to values file")
	helmReleaseCmd.Flags().StringP("registry", "R", "", "Registry where docker image is hosted")
	helmReleaseCmd.Flags().StringP("repository", "r", "", "Repository where docker image is hosted")
	helmReleaseCmd.Flags().StringP("tag", "t", "", "Docker image tag")
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
	Name        string   `yaml:"name"`
	Namespace   string   `yaml:"namespace"`
	Annotations []string `yaml:"annotations"`
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

func generateHelmRelease(name string, namespace string, gitRepoURL string, chartName string, valuesFile string, registry string, repository string, tag string) {

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
			Annotations: []string{"fluxcd.io/automated: true"},
		},
		Spec: Spec{
			ReleaseName: name,
			Chart: Chart{
				Git:  gitRepoURL,
				Path: chartName,
				Ref:  "master",
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
