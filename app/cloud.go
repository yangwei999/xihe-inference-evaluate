package app

import (
	"errors"

	"github.com/opensourceways/xihe-inference-evaluate/domain"
	"github.com/opensourceways/xihe-inference-evaluate/domain/cloud"
)

type CloudPodCreateCmd domain.CloudPod

func (cmd *CloudPodCreateCmd) Validate() error {
	b := cmd.PodId == "" ||
		cmd.SurvivalTime.SurvivalTime() < 0

	if b {
		return errors.New("invalid cmd")
	}

	return nil
}

type CloudService interface {
	Create(*CloudPodCreateCmd) error
}

func NewCloudService(
	manager cloud.Cloud,
) CloudService {
	return &cloudService{
		manager: manager,
	}
}

type cloudService struct {
	manager cloud.Cloud
}

func (s *cloudService) Create(cmd *CloudPodCreateCmd) error {
	return s.manager.Create((*domain.CloudPod)(cmd))
}
