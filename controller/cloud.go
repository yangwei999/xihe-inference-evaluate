package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-inference-evaluate/app"
	"github.com/opensourceways/xihe-inference-evaluate/domain/cloud"
)

func AddRouterForCloudController(
	rg *gin.RouterGroup,
	manager cloud.Cloud,
) {
	ctl := CloudController{
		s: app.NewCloudService(manager),
	}

	rg.POST("/v1/cloud/pod", ctl.Create)
}

type CloudController struct {
	baseController

	s app.CloudService
}

// @Summary Create
// @Description create cloud pod
// @Tags  Cloud
// @Accept json
// @Success 201
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/cloud/pod [post]
func (ctl *CloudController) Create(ctx *gin.Context) {
	req := CloudPodCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if err = ctl.s.Create(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData("successfully"))
}
