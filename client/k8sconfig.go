package client

type Config struct {
	Group          string `json:"group"            required:"true"`
	Version        string `json:"version"          required:"true"`
	Kind           string `json:"kind"             required:"true"`
	KubeConfigFile string `json:"kube_config_file" required:"true"`
}
