package controller

import (
	"github.com/opensourceways/xihe-inference-evaluate/app"
	"github.com/opensourceways/xihe-inference-evaluate/domain"
)

type CloudPodCreateRequest struct {
	PodId        string `json:"pod_id"`
	SurvivalTime int64  `json:"survival_time"`
	User         string `json:"user"`
}

func (req *CloudPodCreateRequest) toCmd() (cmd app.CloudPodCreateCmd, err error) {
	cmd.PodId = req.PodId
	cmd.User = req.User

	if cmd.SurvivalTime, err = domain.NewSurvivalTime(req.SurvivalTime); err != nil {
		return
	}

	return
}
