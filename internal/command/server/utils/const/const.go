package serverConst

import constant "github.com/aaronchen2k/deeptest/internal/command/utils/const"

const (
	HeartBeatInterval    = 60
	CheckUpgradeInterval = 30

	AgentRunTime = 30 * 60
	AgentLogDir  = "log-agent"

	QiNiuURL         = "https://dl.cnezsoft.com/" + constant.AppName + "/"
	AgentVersionURL  = QiNiuURL + "version.txt"
	AgentDownloadURL = QiNiuURL + "%s/%s/" + constant.AppName + ".zip"
)
