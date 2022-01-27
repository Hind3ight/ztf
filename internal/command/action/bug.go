package action

import (
	"fmt"
	testingService "github.com/aaronchen2k/deeptest/internal/command/service/testing"
	"github.com/aaronchen2k/deeptest/internal/command/ui/page"
	constant "github.com/aaronchen2k/deeptest/internal/command/utils/const"
	fileUtils "github.com/aaronchen2k/deeptest/internal/command/utils/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/command/utils/stdin"
	stringUtils "github.com/aaronchen2k/deeptest/internal/command/utils/string"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

func CommitBug(files []string) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		stdinUtils.InputForDir(&resultDir, "", "result")
	}
	resultDir = fileUtils.AddPathSepIfNeeded(resultDir)

	report := testingService.GetZTFTestReportForSubmit(resultDir)

	ids := make([]string, 0)
	lines := make([]string, 0)
	for _, cs := range report.FuncResult {
		if cs.Status != constant.PASS.String() {
			lines = append(lines, fmt.Sprintf("%d. %s %s", cs.Id, cs.Title, logUtils.ColoredStatus(cs.Status)))
			ids = append(ids, strconv.Itoa(cs.Id))
		}
	}

	for {
		logUtils.PrintToWithColor("\n"+i118Utils.Sprintf("enter_case_id_for_report_bug"), color.FgCyan)
		logUtils.PrintToWithColor(strings.Join(lines, "\n"), -1)

		var caseId string
		fmt.Scanln(&caseId)

		if caseId == "exit" {
			color.Unset()
			os.Exit(0)
		} else {
			if stringUtils.FindInArr(caseId, ids) {
				page.CuiReportBug(resultDir, caseId)
			} else {
				logUtils.PrintToWithColor(i118Utils.Sprintf("invalid_input"), color.FgRed)
			}
		}
	}
}
