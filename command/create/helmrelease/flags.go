package helmrelease

// Flags represents flags used by command.
type Flags struct {
	name       string
	namespace  string
	gitRepoURL string
	gitRepoRef string
	chartName  string
	valuesFile string
	registry   string
	repository string
	outputFile string
	tag        string
}

// Validate will check flags correctness.
func (f *Flags) Validate() error {
	return nil
}
