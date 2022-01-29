package service

import (
	"github.com/aaronchen2k/deeptest/internal/command/action"
	"github.com/aaronchen2k/deeptest/internal/command/server/domain"
	serverUtils "github.com/aaronchen2k/deeptest/internal/command/server/utils/common"
	constant "github.com/aaronchen2k/deeptest/internal/command/utils/const"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"

	"strings"
)

type ExecService struct {
}

func NewExecService() *ExecService {
	return &ExecService{}
}

func (s *ExecService) Exec(build domain.Build) (reply domain.OptResult) {
	serverVerbose := consts.Verbose
	consts.Verbose = build.Debug
	consts.RunMode = constant.RunModeRequest
	defer rollback(serverVerbose)

	s.prepareCodes(&build)
	s.prepareDir(&build)

	resultDir := ""
	if stringUtils.FindInArr(build.UnitTestType, constant.UnitTestTypes) { // unit test
		consts.ProductId = build.ProductId

		consts.UnitTestType = build.UnitTestType
		consts.UnitTestTool = build.UnitTestTool

		resultDir = action.RunUnitTest(build.UnitTestCmd)

	} else { // ztf functional test
		consts.ProductId = build.ProductId

		action.RunZTFTest(build.Files, build.SuiteId, build.TaskId)
		resultDir = consts.LogDir
	}

	serverUtils.BakLog(resultDir)
	return
}

func (s *ExecService) prepareCodes(build *domain.Build) {
	if build.WorkDir != "" {
		build.WorkDir = fileUtils.AddPathSepIfNeeded(build.WorkDir)
	}

	if build.ScmAddress != "" { // git
		serverUtils.CheckoutCodes(build)

	} else if strings.Index(build.ScriptUrl, "http") == 0 { // zip
		serverUtils.DownloadCodes(build)

	} else { // folder
		if build.ScriptUrl != "" {
			build.ScriptUrl = fileUtils.AddPathSepIfNeeded(build.ScriptUrl)
		}
		build.ProjectDir = build.ScriptUrl
	}
}

func (s *ExecService) prepareDir(build *domain.Build) {
	consts.ServerWorkDir = build.WorkDir
	consts.ServerProjectDir = build.ProjectDir

	if consts.ServerProjectDir == "" && consts.ServerWorkDir != "" {
		consts.ServerProjectDir = consts.ServerWorkDir
	} else if consts.ServerProjectDir != "" && consts.ServerWorkDir == "" {
		consts.ServerWorkDir = consts.ServerProjectDir
	} else if consts.ServerProjectDir == "" && consts.ServerWorkDir == "" {
		consts.ServerWorkDir = fileUtils.AbsolutePath(".")
		consts.ServerProjectDir = consts.ServerWorkDir
	}

	if consts.ServerWorkDir != "" {
		consts.ServerWorkDir = fileUtils.AddPathSepIfNeeded(consts.ServerWorkDir)
	}
	if consts.ServerProjectDir != "" {
		consts.ServerProjectDir = fileUtils.AddPathSepIfNeeded(consts.ServerProjectDir)
	}
}

func rollback(serverVerbose bool) {
	consts.Verbose = serverVerbose
	consts.RunMode = constant.RunModeCommon
}
