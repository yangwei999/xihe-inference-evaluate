package evaluateimpl

import "github.com/opensourceways/xihe-inference-evaluate/infrastructure/inferenceimpl"

type Config struct {
	Image        string                  `json:"image"  required:"true"`
	OBS          inferenceimpl.OBSConfig `json:"-"`
	CrdNamespace string                  `json:"crd_namespace"    required:"true"`
	CrdCpu       string                  `json:"crd_cpu"          required:"true"`
	CrdMemory    string                  `json:"crd_memory"       required:"true"`
}
