package action

import (
	scriptService "github.com/aaronchen2k/deeptest/internal/command/service/script"
	assertUtils "github.com/aaronchen2k/deeptest/internal/command/utils/assert"
)

func View(files []string, keywords string) {
	cases := assertUtils.GetCaseByDirAndFile(files)

	scriptService.View(cases, keywords)
}
