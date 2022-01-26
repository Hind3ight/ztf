package zentaoService

import (
	testingService "github.com/aaronchen2k/deeptest/internal/command/service/testing"
	configUtils "github.com/aaronchen2k/deeptest/internal/command/utils/config"
	i118Utils "github.com/aaronchen2k/deeptest/internal/command/utils/i118"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/command/utils/stdin"
	"github.com/aaronchen2k/deeptest/internal/command/utils/vari"
	"strconv"
)

func CommitZTFTestResult(resultDir string, productId string, taskId string, noNeedConfirm bool) {
	conf := configUtils.ReadCurrConfig()
	ok := Login(conf.Url, conf.Account, conf.Password)
	if !ok {
		return
	}

	report := testingService.GetZTFTestReportForSubmit(resultDir)

	if vari.ProductId == "" && productId != "" {
		vari.ProductId = productId
	}

	if taskId == "" && !noNeedConfirm {
		taskId = stdinUtils.GetInput("\\d*", "",
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("task_id")+
				i118Utils.Sprintf("task_id_empty_to_create"))
	}

	taskIdInt, _ := strconv.Atoi(taskId)
	CommitTestResult(report, taskIdInt)
}
