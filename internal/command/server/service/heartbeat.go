package service

import (
	serverModel "github.com/aaronchen2k/deeptest/internal/command/server/domain"
	"github.com/aaronchen2k/deeptest/internal/command/server/utils/common"
	serverConst "github.com/aaronchen2k/deeptest/internal/command/server/utils/const"
	"github.com/aaronchen2k/deeptest/internal/command/service/client"
	zentaoService "github.com/aaronchen2k/deeptest/internal/command/service/zentao"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"

	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/zentao"
)

var (
	sysInfo serverModel.SysInfo
)

type HeartBeatService struct {
}

func NewHeartBeatService() *HeartBeatService {
	return &HeartBeatService{}
}

func (s *HeartBeatService) HeartBeat(isBusy bool) {
	if sysInfo.OsName == "" {
		sysInfo = serverUtils.GetSysInfo()
	}

	// send request
	zentaoService.GetConfig(consts.Config.Url)

	url := consts.Config.Url + zentaoUtils.GenApiUri("agent", "heartbeat", "")
	data := map[string]interface{}{"type": consts.Platform, "sys": sysInfo}

	status := serverConst.VmActive
	if isBusy {
		status = serverConst.VmBusy
	}
	data["status"] = status

	_, ok := client.PostObject(url, data, false)
	if ok {
		logUtils.PrintTo(i118Utils.Sprintf("success_heart_beat"))
	} else {
		logUtils.PrintTo(i118Utils.Sprintf("fail_heart_beat"))
	}

	return
}
