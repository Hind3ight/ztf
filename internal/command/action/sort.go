package action

import (
	scriptService "github.com/aaronchen2k/deeptest/internal/command/service/script"
	assertUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/assert"
)

func Sort(files []string) {
	cases := assertUtils.GetCaseByDirAndFile(files)

	scriptService.Sort(cases)
}
