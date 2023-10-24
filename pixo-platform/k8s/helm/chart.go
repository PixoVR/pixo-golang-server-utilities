package helm

type Chart struct {
	Name        string
	RepoURL     string
	Version     string
	ReleaseName string
	Namespace   string
}
