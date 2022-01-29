package zentaoService

import (
	testingService "github.com/aaronchen2k/deeptest/internal/command/service/testing"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"

	configUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/config"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"strconv"
)

func CommitZTFTestResult(resultDir string, productId string, taskId string, noNeedConfirm bool) {
	conf := configUtils.ReadCurrConfig()
	ok := Login(conf.Url, conf.Account, conf.Password)
	if !ok {
		return
	}

	report := testingService.GetZTFTestReportForSubmit(resultDir)

	if consts.ProductId == "" && productId != "" {
		consts.ProductId = productId
	}

	if taskId == "" && !noNeedConfirm {
		taskId = stdinUtils.GetInput("\\d*", "",
			i118Utils.Sprintf("pls_enter")+" "+i118Utils.Sprintf("task_id")+
				i118Utils.Sprintf("task_id_empty_to_create"))
	}

	taskIdInt, _ := strconv.Atoi(taskId)
	CommitTestResult(report, taskIdInt)
}
