package httpUtils

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/ajg/form"
	"github.com/yosssi/gohtml"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func Get(url string) (ret []byte, ok bool) {
	if strings.Index(url, "mode=getconfig") < 0 {
		url = AddToken(url)
	}
	if commConsts.Verbose {
		logUtils.Info(url)
	}

	client := &http.Client{}

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		logUtils.Error(reqErr.Error())
		ok = false
		return
	}

	resp, respErr := client.Do(req)
	if respErr != nil {
		logUtils.Error(respErr.Error())
		ok = false
		return
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if commConsts.Verbose {
		logUtils.PrintUnicode(bodyBytes)
	}
	defer resp.Body.Close()

	var zentaoResp serverDomain.ZentaoResp
	jsonErr := json.Unmarshal(bodyBytes, &zentaoResp)
	if jsonErr != nil {
		if strings.Index(string(bodyBytes), "<html>") > -1 {
			if commConsts.Verbose {
				logUtils.Errorf(i118Utils.Sprintf("server_return") + " HTML - " + gohtml.FormatWithLineNo(string(bodyBytes)))
			}
			return
		} else {
			if commConsts.Verbose {
				logUtils.Infof(jsonErr.Error())
			}
			return
		}
	}

	defer resp.Body.Close()

	status := zentaoResp.Status
	if status == "" { // 非嵌套结构
		ret = bodyBytes
		ok = true
	} else { // 嵌套结构
		ret = []byte(zentaoResp.Data)
		ok = status == "success"
	}

	return
}

func Post(url string, data interface{}, useFormFormat bool) (ret []byte, ok bool) {
	url = AddToken(url)

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("server_address") + url)
	}
	client := &http.Client{}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		logUtils.Error(err.Error())
		return
	}

	dataStr := string(dataBytes)
	if useFormFormat {
		dataStr, _ = form.EncodeToString(data)
		// convert data to post fomat
		re3, _ := regexp.Compile(`([^&]*?)=`)
		dataStr = re3.ReplaceAllStringFunc(dataStr, replacePostData)
	}

	if commConsts.Verbose {
		logUtils.PrintUnicode([]byte(dataStr))
	}

	req, reqErr := http.NewRequest("POST", url, strings.NewReader(dataStr))
	if reqErr != nil {
		logUtils.Error(reqErr.Error())
		return
	}

	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, respErr := client.Do(req)
	if respErr != nil {
		logUtils.Error(respErr.Error())
		return
	}

	bodyBytes, ioErr := ioutil.ReadAll(resp.Body)
	if ioErr != nil {
		logUtils.Error(ioErr.Error())
		return
	}

	if commConsts.Verbose {
		logUtils.PrintUnicode(bodyBytes)
	}
	defer resp.Body.Close()

	ret, ok = GetRespErr(bodyBytes, url)

	return
}

func PostStr(url string, params map[string]string) (ret []byte, ok bool) {
	url = AddToken(url)

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("server_address") + url)
	}
	client := &http.Client{}

	paramStr := ""
	idx := 0
	for pkey, pval := range params {
		if idx > 0 {
			paramStr += "&"
		}
		paramStr = paramStr + pkey + "=" + pval
		idx++
	}

	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("server_params") + paramStr)
	}

	req, reqErr := http.NewRequest("POST", url, strings.NewReader(paramStr))
	if reqErr != nil {
		if commConsts.Verbose {
			logUtils.Infof(reqErr.Error())
		}
		ok = false
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", commConsts.SessionVar+"="+commConsts.SessionId)

	resp, respErr := client.Do(req)
	if respErr != nil {
		if commConsts.Verbose {
			logUtils.Infof(respErr.Error())
		}
		ok = false
		return
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	if commConsts.Verbose {
		logUtils.Infof(i118Utils.Sprintf("server_return") + logUtils.ConvertUnicode(bytes))
	}

	return
}

func GenUrl(server string, path string) string {
	server = UpdateUrl(server)
	url := fmt.Sprintf("%sapi/v1/%s", server, path)
	return url
}

func UpdateUrl(url string) string {
	if strings.LastIndex(url, "/") < len(url)-1 {
		url += "/"
	}
	return url
}

func replacePostData(str string) string {
	str = strings.ToLower(str[:1]) + str[1:]

	if strings.Index(str, ".") > -1 {
		str = strings.Replace(str, ".", "[", -1)
		str = strings.Replace(str, "=", "]=", -1)
	}
	return str
}

func AddToken(url string) (ret string) {
	address := url
	hash := ""

	index := strings.Index(url, "#")
	if index > 0 {
		address = url[:index]
		hash = url[index:]
	}

	paramPir := commConsts.SessionVar + "=" + commConsts.SessionId

	if commConsts.RequestType == commConsts.PathInfo {
		address = address + "?" + paramPir
	} else {
		address = address + "&" + paramPir
	}

	ret = address + "&XDEBUG_SESSION_START=PHPSTORM" + hash

	return
}

func GetRespErr(bytes []byte, url string) (ret []byte, ok bool) {
	ret = bytes

	if len(bytes) == 0 {
		return
	}

	var zentaoResp serverDomain.ZentaoResp
	jsonErr := json.Unmarshal(bytes, &zentaoResp)
	if jsonErr != nil {
		if commConsts.Verbose {
			if strings.Index(url, "login") < 0 { // jsonErr caused by login request return a html
				logUtils.Infof(jsonErr.Error())
			}
		}
		ok = false
		return
	}

	// 嵌套结构，map[status:success, data:{}]
	status := zentaoResp.Status
	if status != "" {
		ret = []byte(zentaoResp.Data)
		ok = status == "success"
		return
	}

	// 非嵌套结构，map[result:success]
	var respData = serverDomain.ZentaoRespData{}
	err := json.Unmarshal(bytes, &respData)

	if err == nil && (respData.Result != "" && respData.Result != "success") {
		ok = false
	}

	return
}
