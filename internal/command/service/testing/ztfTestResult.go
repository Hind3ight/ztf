package testingService

import (
	"encoding/json"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/command/model"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"strings"
)

func GetZTFTestReportForSubmit(resultDir string) model.TestReport {
	resultPath := resultDir + "result.json"

	content := fileUtils.ReadFile(resultPath)
	content = strings.Replace(content, "\n", "", -1)

	var report model.TestReport
	json.Unmarshal([]byte(content), &report)

	return report
}

func GetStepHtml(step model.StepLog) string {
	stepResults := make([]string, 0)

	stepStatus := stringUtils.BoolToPass(step.Status)

	stepTxt := fmt.Sprintf(
		"<p><b>%s %s</b></p>",
		step.Id, stepStatus)

	for _, checkpoint := range step.CheckPoints {
		checkpointStatus := stringUtils.BoolToPass(checkpoint.Status)

		text := fmt.Sprintf(
			"<p>&nbsp;Checkpoint: %s</p>"+
				"<p>&nbsp;&nbsp;Expect</p>"+
				"&nbsp;&nbsp;&nbsp;%s"+
				"<p>&nbsp;&nbsp;Actual<p/>"+
				"&nbsp;&nbsp;&nbsp;%s",
			checkpointStatus, checkpoint.Expect, checkpoint.Actual)

		stepResults = append(stepResults, text)
	}

	return stepTxt + strings.Join(stepResults, "<br/>") + "<br/>"
}

func GetStepText(step model.StepLog) string {
	stepResults := make([]string, 0)

	stepStatus := stringUtils.BoolToPass(step.Status)

	stepTxt := fmt.Sprintf(
		"%s %s\n",
		step.Id, stepStatus)

	for _, checkpoint := range step.CheckPoints {
		checkpointStatus := stringUtils.BoolToPass(checkpoint.Status)

		text := fmt.Sprintf(
			" Checkpoint: %s\n"+
				"  Expect\n"+
				"   %s\n"+
				"  Actual\n"+
				"   %s",
			checkpointStatus, checkpoint.Expect, checkpoint.Actual)

		stepResults = append(stepResults, text)
	}

	return stepTxt + strings.Join(stepResults, "\n") + "\n"
}
