package controller

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type TestExecCtrl struct {
	TestExecService *service.TestExecService `inject:""`
	BaseCtrl
}

func NewTestExecCtrl() *TestExecCtrl {
	return &TestExecCtrl{}
}

// List 分页列表
func (c *TestExecCtrl) List(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	data, err := c.TestExecService.List(projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.RequestErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

// Get 详情
func (c *TestExecCtrl) Get(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	seq := ctx.Params().Get("seq")

	if seq == "" {
		logUtils.Errorf("参数解析失败")
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}

	exec, err := c.TestExecService.Get(projectPath, seq)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.RequestErr.Code, Data: nil,
			Msg: fmt.Sprintf("获取编号为%s的日志失败。", seq)})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: exec, Msg: domain.NoErr.Msg})
}

// Delete 删除
func (c *TestExecCtrl) Delete(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	seq := ctx.Params().Get("seq")

	if projectPath == "" || seq == "" {
		logUtils.Errorf("参数解析失败")
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}

	err := c.TestExecService.Delete(projectPath, seq)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.RequestErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}
