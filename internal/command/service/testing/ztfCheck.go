package testingService

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/command/model"
	"github.com/aaronchen2k/deeptest/internal/command/utils/const"
	"github.com/aaronchen2k/deeptest/internal/command/utils/lang"
	"github.com/aaronchen2k/deeptest/internal/command/utils/log"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/command/utils/script"
	stringUtils "github.com/aaronchen2k/deeptest/internal/command/utils/string"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/zentao"
	"github.com/emirpasic/gods/maps"
	"github.com/mattn/go-runewidth"
	"strconv"
	"strings"
)

func CheckCaseResult(file string, logs string, report *model.TestReport, idx int, total int, secs string, pathMaxWidth int, numbMaxWidth int) {
	_, _, expectMap, isOldFormat := scriptUtils.GetStepAndExpectMap(file)

	isIndependent, expectIndependentContent := zentaoUtils.GetDependentExpect(file)
	if isIndependent {
		if isOldFormat {
			expectMap = scriptUtils.GetExpectMapFromIndependentFileObsolete(expectMap, expectIndependentContent, false)
		} else {
			expectMap = scriptUtils.GetExpectMapFromIndependentFile(expectMap, expectIndependentContent, false)
		}
	}

	skip := false
	actualArr := make([][]string, 0)
	if isOldFormat {
		skip, actualArr = zentaoUtils.ReadLogArrObsolete(logs)
	} else {
		skip, actualArr = zentaoUtils.ReadLogArr(logs)
	}

	language := langUtils.GetLangByFile(file)
	ValidateCaseResult(file, language, expectMap, skip, actualArr, report,
		idx, total, secs, pathMaxWidth, numbMaxWidth)
}

func ValidateCaseResult(scriptFile string, langType string,
	expectMap maps.Map, skip bool, actualArr [][]string, report *model.TestReport,
	idx int, total int, secs string, pathMaxWidth int, numbMaxWidth int) {

	_, caseId, productId, title := zentaoUtils.GetCaseInfo(scriptFile)

	stepLogs := make([]model.StepLog, 0)
	caseResult := constant.PASS.String()
	noExpects := true

	if skip {
		caseResult = constant.SKIP.String()
	} else {
		idx := 0

		for _, numbInterf := range expectMap.Keys() { // iterate by checkpoints
			expectInterf, _ := expectMap.Get(numbInterf)

			numb := strings.TrimSpace(numbInterf.(string))
			expect := strings.TrimSpace(expectInterf.(string))

			if expect == "" {
				continue
			}

			noExpects = false

			expectLines := strings.Split(expect, "\n")
			var actualLines []string
			if len(actualArr) > idx {
				actualLines = actualArr[idx]
			}

			stepResult, checkpointLogs := ValidateStepResult(langType, expectLines, actualLines)
			stepLog := model.StepLog{Id: numb, Status: stepResult, CheckPoints: checkpointLogs}
			stepLogs = append(stepLogs, stepLog)
			if !stepResult {
				caseResult = constant.FAIL.String()
			}

			idx++
		}
	}

	if noExpects {
		caseResult = constant.SKIP.String()
	}

	if caseResult == constant.FAIL.String() {
		report.Fail = report.Fail + 1
	} else if caseResult == constant.PASS.String() {
		report.Pass = report.Pass + 1
	} else if caseResult == constant.SKIP.String() {
		report.Skip = report.Skip + 1
	}
	report.Total = report.Total + 1

	cs := model.FuncResult{Id: caseId, ProductId: productId, Title: title,
		Path: scriptFile, Status: caseResult, Steps: stepLogs}
	report.FuncResult = append(report.FuncResult, cs)

	// print case result to console
	statusColor := logUtils.ColoredStatus(cs.Status)
	width := strconv.Itoa(len(strconv.Itoa(total)))
	numbWidth := strconv.Itoa(numbMaxWidth)

	path := cs.Path
	lent := runewidth.StringWidth(path)

	if pathMaxWidth > lent {
		postFix := strings.Repeat(" ", pathMaxWidth-lent)
		path += postFix
	}

	format := "(%" + width + "d/%d) %s [%s] [%" + numbWidth + "d. %s] (%ss)"
	logUtils.Screen(fmt.Sprintf(format, idx+1, total, statusColor, path, cs.Id, cs.Title, secs))
	logUtils.Result(fmt.Sprintf(format, idx+1, total, i118Utils.Sprintf(cs.Status), path, cs.Id, cs.Title, secs))
}

func ValidateStepResult(langType string, expectLines []string, actualLines []string) (bool, []model.CheckPointLog) {
	stepResult := true

	checkpointLogs := make([]model.CheckPointLog, 0)

	indx2 := 0
	for _, expect := range expectLines {
		log := "N/A"
		if len(actualLines) > indx2 {
			log = actualLines[indx2]
		}

		expect = strings.TrimSpace(expect)
		var pass bool
		if expect[:1] == "`" && expect[len(expect)-1:] == "`" {
			expect = expect[1 : len(expect)-1]
			pass = stringUtils.MatchString(expect, log, langType)
		} else {
			pass = strings.Contains(log, expect)
		}

		if !pass {
			stepResult = false
		}

		cp := model.CheckPointLog{Numb: indx2 + 1, Status: pass, Expect: expect, Actual: log}
		checkpointLogs = append(checkpointLogs, cp)

		indx2++
	}

	return stepResult, checkpointLogs

}
