package service

import (
	"encoding/json"
	"github.com/aaronchen2k/deeptest/internal/command/model"
	"github.com/aaronchen2k/deeptest/internal/command/server/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"

	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	configUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/config"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
)

var ()

type ConfigService struct {
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (s *ConfigService) Update(req domain.ReqData) {
	conf := model.Config{}

	reqStr, _ := json.Marshal(req.Data)
	err := json.Unmarshal(reqStr, &conf)
	if err != nil {
		logUtils.PrintTo(i118Utils.Sprintf("fail_parse_req", err))
		return
	}

	if conf.Version != 0 {
		consts.Config.Version = conf.Version
	}
	if conf.Language != "" {
		consts.Config.Language = conf.Language
	}
	if conf.Url != "" {
		consts.Config.Url = conf.Url
	}
	if conf.Account != "" {
		consts.Config.Account = conf.Account
	}
	if conf.Password != "" {
		consts.Config.Password = conf.Password
	}
	if conf.Javascript != "" {
		consts.Config.Javascript = conf.Javascript
	}
	if conf.Lua != "" {
		consts.Config.Lua = conf.Lua
	}
	if conf.Perl != "" {
		consts.Config.Perl = conf.Perl
	}
	if conf.Php != "" {
		consts.Config.Php = conf.Php
	}
	if conf.Python != "" {
		consts.Config.Python = conf.Python
	}
	if conf.Ruby != "" {
		consts.Config.Ruby = conf.Ruby
	}
	if conf.Tcl != "" {
		consts.Config.Tcl = conf.Tcl
	}
	if conf.Lua != "" {
		consts.Config.Autoit = conf.Autoit
	}

	configUtils.SaveConfig(consts.Config)
}
