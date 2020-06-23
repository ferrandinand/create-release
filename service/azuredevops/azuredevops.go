package azuredevops

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/giantswarm/microerror"
)

const (
	baseURL    = "https://dev.azure.com"
	buildsPath = "%s/%s/_apis/build/builds?api-version=5.1"
)

// Config represents the configuration used to create a new azdvo service.
type Config struct {
	// Settings.
	ADVOUser  string
	AZDOToken string
}

// Service represents a azdvo service.
type Service struct {
	// Settings.
	ADVOUser  string
	AZDOToken string
}

// DefaultConfig provides a default configuration to create a azdvo service.
func DefaultConfig() Config {
	newConfig := Config{
		// Settings.
		ADVOUser:  "",
		AZDOToken: "",
	}

	return newConfig
}

// New creates a new configured azdvo service.
func New(config Config) (*Service, error) {
	// Settings.
	if config.ADVOUser == "" || config.AZDOToken == "" {
		return nil, microerror.Maskf(
			invalidConfigError,
			"ADVOUser and AZDOToken must not be empty",
		)
	}

	// Create service.
	newService := &Service{
		// Settings.
		ADVOUser:  config.ADVOUser,
		AZDOToken: config.AZDOToken,
	}

	return newService, nil
}

// runPipeline will trigger Azure Devops pipeline.
func (s *Service) RunPipeline(ctx context.Context, organization string, project string, params string) error {
	client := http.DefaultClient

	// Prepare URL.
	buildPath := fmt.Sprintf(buildsPath, organization, project)
	buildURL, err := url.Parse(baseURL + "/" + buildPath)
	if err != nil {
		return microerror.Mask(err)
	}

	// Prepare build parameters.
	var formatParams = []byte(params)

	req, err := http.NewRequest("POST", buildURL.String(), bytes.NewBuffer(formatParams))
	if err != nil {
		return microerror.Mask(err)
	}

	// Use azdvo credentials to access.
	req.SetBasicAuth(s.ADVOUser, s.AZDOToken)

	req.Header.Set("Content-Type", "application/json")

	// Perform request.
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else if resp.StatusCode != http.StatusOK {
		return microerror.Maskf(
			unexpectedResponseError,
			fmt.Sprint(resp.StatusCode),
		)
	}
	defer resp.Body.Close()

	return nil
}
