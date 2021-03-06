package project

import (
	"github.com/giantswarm/microerror"
)

// Flags represents flags used by command.
type Flags struct {
	AZDVOToken      string
	AZDVOUser       string
	projectName     string
	projectType     string
	projectPlatform string
}

// Validate will check flags correctness.
func (f *Flags) Validate() error {
	if f.AZDVOToken == "" {
		return microerror.Maskf(invalidFlagsError, "Azure Devops PAT token must not be empty")
	}
	if f.AZDVOUser == "" {
		return microerror.Maskf(invalidFlagsError, "Azure Devops user must not be empty")
	}
	if f.projectName == "" {
		return microerror.Maskf(invalidFlagsError, "Azure Devops project name must not be empty")
	}

	return nil
}
