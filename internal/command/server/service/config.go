package service

import (
	"encoding/json"
	"github.com/aaronchen2k/deeptest/internal/command/model"
	"github.com/aaronchen2k/deeptest/internal/command/server/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	"github.com/aaronchen2k/deeptest/internal/command/utils/vari"
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
		vari.Config.Version = conf.Version
	}
	if conf.Language != "" {
		vari.Config.Language = conf.Language
	}
	if conf.Url != "" {
		vari.Config.Url = conf.Url
	}
	if conf.Account != "" {
		vari.Config.Account = conf.Account
	}
	if conf.Password != "" {
		vari.Config.Password = conf.Password
	}
	if conf.Javascript != "" {
		vari.Config.Javascript = conf.Javascript
	}
	if conf.Lua != "" {
		vari.Config.Lua = conf.Lua
	}
	if conf.Perl != "" {
		vari.Config.Perl = conf.Perl
	}
	if conf.Php != "" {
		vari.Config.Php = conf.Php
	}
	if conf.Python != "" {
		vari.Config.Python = conf.Python
	}
	if conf.Ruby != "" {
		vari.Config.Ruby = conf.Ruby
	}
	if conf.Tcl != "" {
		vari.Config.Tcl = conf.Tcl
	}
	if conf.Lua != "" {
		vari.Config.Autoit = conf.Autoit
	}

	configUtils.SaveConfig(vari.Config)
}
