package action

import (
	zentaoService "github.com/aaronchen2k/deeptest/internal/command/service/zentao"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/command/utils/script"
	assertUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/assert"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/zentao"
	"log"
)

func CommitCases(files []string) {
	cases := assertUtils.GetCaseByDirAndFile(files)

	for _, cs := range cases {
		pass, id, _, title := zentaoUtils.GetCaseInfo(cs)

		if pass {
			stepMap, stepTypeMap, expectMap, isOldFormat := scriptUtils.GetStepAndExpectMap(cs)
			log.Println(isOldFormat)

			isIndependent, expectIndependentContent := zentaoUtils.GetDependentExpect(cs)
			if isIndependent {
				expectMap = scriptUtils.GetExpectMapFromIndependentFileObsolete(expectMap, expectIndependentContent, true)
			}

			zentaoService.CommitCase(id, title, stepMap, stepTypeMap, expectMap)
		}
	}
}
