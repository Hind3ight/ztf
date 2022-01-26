package zentaoService

import (
	"github.com/aaronchen2k/deeptest/internal/command/model"
	"github.com/aaronchen2k/deeptest/internal/command/service/client"
	configUtils "github.com/aaronchen2k/deeptest/internal/command/utils/config"
	i118Utils "github.com/aaronchen2k/deeptest/internal/command/utils/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	"github.com/aaronchen2k/deeptest/internal/command/utils/vari"
	"github.com/aaronchen2k/deeptest/internal/command/utils/zentao"
	"github.com/bitly/go-simplejson"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

func CommitTestResult(report model.TestReport, testTaskId int) {
	if vari.ProductId == "" {
		if len(report.FuncResult) > 0 {
			vari.ProductId = strconv.Itoa(report.FuncResult[0].ProductId)
		} else if len(report.UnitResult) > 0 {
			vari.ProductId = strconv.Itoa(report.ProductId)
		}
	}

	if report.Total == 0 {
		logUtils.Screen(color.CyanString(i118Utils.Sprintf("ignore_to_submit_result_no_result_empty")))
		return
	}
	if vari.ProductId == "" {
		logUtils.Screen(color.CyanString(i118Utils.Sprintf("ignore_to_submit_result_no_product_id")))
		return
	}

	conf := configUtils.ReadCurrConfig()
	ok := Login(conf.Url, conf.Account, conf.Password)
	if !ok {
		return
	}

	report.ZentaoData = os.Getenv("ZENTAO_DATA")
	report.BuildUrl = os.Getenv("BUILD_URL")
	report.ProductId, _ = strconv.Atoi(vari.ProductId)
	report.TaskId = testTaskId

	if len(report.FuncResult) > 0 {
		report.ProductId = report.FuncResult[0].ProductId
	}

	url := conf.Url + zentaoUtils.GenApiUri("ci", "commitResult", "")
	resp, ok := client.PostObject(url, report, false)

	if ok {
		json, err1 := simplejson.NewJson([]byte(resp))
		if err1 == nil {
			result, err2 := json.Get("result").String()
			if err2 != nil || result != "success" {
				ok = false
			}
		} else {
			ok = false
		}
	}

	msg := "\n"
	if ok {
		msg += color.GreenString(i118Utils.Sprintf("success_to_submit_test_result"))
	} else {
		msg = i118Utils.Sprintf("fail_to_submit_test_result")
		if strings.Index(resp, "login") > -1 {
			msg = i118Utils.Sprintf("fail_to_login")
		}
		msg = color.RedString(msg)
	}

	logUtils.Screen(msg)
	logUtils.Screen(logUtils.GetWholeLine("=", "=") + "\n")

	if report.Fail > 0 || !ok {
		os.Exit(1)
	}
}
