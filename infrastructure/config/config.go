package config

import (
	"fmt"
	"os"
	"strconv"
)

type CRDConfig struct {
	CRDImage     string `json:"crd_image"        required:"true"`
	CRDInitImage string `json:"crd_init_image"        required:"true"`
	TemplateFile string `json:"crd_template"     required:"true"`

	// Specifies the terminal container port for connection
	ContainerPort int `json:"container_port"   required:"true"`

	// CrdCpu specifies the number of cpu
	CRDCpu float32 `json:"crd_cpu"               required:"true"`

	// CrdMemory specifies the memory in megabyte.
	CRDMemory int `json:"crd_memory"             required:"true"`
}

func (cfg *CRDConfig) CRDCpuString() string {
	return fmt.Sprintf("%v", cfg.CRDCpu)
}

func (cfg *CRDConfig) CRDMemoryString() string {
	return fmt.Sprintf("%vMi", cfg.CRDMemory)
}

func (cfg *CRDConfig) CRDContainerPortString() string {
	return strconv.Itoa(cfg.ContainerPort)
}

func (cfg *CRDConfig) SetDefault() {
	if cfg.ContainerPort <= 0 {
		cfg.ContainerPort = 8080
	}

	// get image tag from environment variable
	if image, ok := os.LookupEnv("CRD_IMAGE"); ok {
		cfg.CRDImage = image
	}

	if initImage, ok := os.LookupEnv("CRD_INIT_IMAGE"); ok {
		cfg.CRDInitImage = initImage
	}
}

func (cfg *CRDConfig) Validate() error {
	if cfg.CRDImage == "" {
		return fmt.Errorf("crd image must be set")
	}

	if cfg.CRDInitImage == "" {
		return fmt.Errorf("crd init image must be set")
	}

	return nil
}
