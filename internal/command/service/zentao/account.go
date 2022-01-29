package zentaoService

import (
	"github.com/aaronchen2k/deeptest/internal/command/service/client"
	constant "github.com/aaronchen2k/deeptest/internal/command/utils/const"
	"github.com/aaronchen2k/deeptest/internal/command/utils/log"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"

	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/bitly/go-simplejson"
	"github.com/fatih/color"
	"strings"
)

func Login(baseUrl string, account string, password string) bool {
	ok := GetConfig(baseUrl)

	if !ok {
		logUtils.PrintToCmd(i118Utils.Sprintf("fail_to_login"), color.FgRed)
		return false
	}

	uri := ""
	if consts.RequestType == constant.RequestTypePathInfo {
		uri = "user-login.json"
	} else {
		uri = "index.php?m=user&f=login&t=json"
	}
	url := baseUrl + uri

	params := make(map[string]string)
	params["account"] = account
	params["password"] = password

	var body string
	body, ok = client.PostStr(url, params)
	if !ok || (ok && strings.Index(body, "title") > 0) { // use PostObject to login again for new system
		_, ok = client.PostObject(url, params, true)
	}
	if ok {
		if consts.Verbose {
			logUtils.Screen(i118Utils.Sprintf("success_to_login"))
		}
	} else {
		logUtils.PrintToCmd(i118Utils.Sprintf("fail_to_login"), color.FgRed)
	}

	return ok
}

func GetConfig(baseUrl string) bool {
	if consts.RequestType != "" {
		return true
	}

	// get config
	url := baseUrl + "?mode=getconfig"
	body, ok := client.Get(url)
	if !ok {
		return false
	}

	json, _ := simplejson.NewJson([]byte(body))
	consts.ZenTaoVersion, _ = json.Get("version").String()
	consts.SessionId, _ = json.Get("sessionID").String()
	consts.SessionVar, _ = json.Get("sessionVar").String()
	consts.RequestType, _ = json.Get("requestType").String()
	consts.RequestFix, _ = json.Get("requestFix").String()

	// check site path by calling login interface
	uri := ""
	if consts.RequestType == constant.RequestTypePathInfo {
		uri = "user-login.json"
	} else {
		uri = "index.php?m=user&f=login&t=json"
	}
	url = baseUrl + uri
	body, ok = client.Get(url)
	if !ok {
		return false
	}

	return true
}
