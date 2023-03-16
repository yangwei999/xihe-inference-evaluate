package cloud

import "github.com/opensourceways/xihe-inference-evaluate/domain"

type Cloud interface {
	Create(*domain.CloudPod) error
	NotifyResult(labels map[string]string, status domain.ContainerDetail)
}