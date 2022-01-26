package action

import (
	zentaoService "github.com/aaronchen2k/deeptest/internal/command/service/zentao"
	fileUtils "github.com/aaronchen2k/deeptest/internal/command/utils/file"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/command/utils/stdin"
)

func CommitZTFTestResult(files []string, productId string, taskId string, noNeedConfirm bool) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		stdinUtils.InputForDir(&resultDir, "", "result")
	}

	resultDir = fileUtils.AddPathSepIfNeeded(resultDir)
	zentaoService.CommitZTFTestResult(resultDir, productId, taskId, noNeedConfirm)
}
