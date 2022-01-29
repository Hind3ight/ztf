package configUtils

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/command/model"
	commonUtils "github.com/aaronchen2k/deeptest/internal/command/utils/common"
	"github.com/aaronchen2k/deeptest/internal/command/utils/const"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"

	assertUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/assert"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/display"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/fatih/color"
	"gopkg.in/ini.v1"
	"os"
	"path"
	"reflect"
	"strings"
)

func InitConfig() {
	consts.ExeDir, consts.IsDebug = fileUtils.GetZTFDir()
	CheckConfigPermission()

	consts.ConfigPath = consts.ExeDir + constant.ConfigFile
	consts.Config = getInst()

	// screen size
	InitScreenSize()

	// internationalization
	i118Utils.InitI118(consts.Config.Language)

	consts.ScriptExtToNameMap = langUtils.GetExtToNameMap()
}

func InitScreenSize() {
	w, h := display.GetScreenSize()
	consts.ScreenWidth = w
	consts.ScreenHeight = h
}

func PrintCurrConfig() {
	logUtils.PrintToWithColor("\n"+i118Utils.Sprintf("current_config"), color.FgCyan)

	val := reflect.ValueOf(consts.Config)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(consts.Config).NumField(); i++ {
		if !commonUtils.IsWin() && i > 4 {
			break
		}

		val := val.Field(i)
		name := typeOfS.Field(i).Name

		fmt.Printf("  %s: %v \n", name, val.Interface())
	}
}

func ReadCurrConfig() model.Config {
	config := model.Config{}

	configPath := consts.ConfigPath
	if consts.ServerWorkDir != "" {
		configPath = consts.ServerWorkDir + constant.ConfigFile
	}

	if !fileUtils.FileExist(configPath) {
		config.Language = "en"
		i118Utils.InitI118("en")

		return config
	}

	ini.MapTo(&config, consts.ConfigPath)

	config.Url = commonUtils.AddSlashForUrl(config.Url)

	return config
}
func SaveConfig(conf model.Config) error {
	configPath := consts.ConfigPath
	if consts.ServerWorkDir != "" {
		configPath = consts.ServerWorkDir + constant.ConfigFile
	}

	fileUtils.MkDirIfNeeded(path.Dir(configPath))

	if conf.Version == 0 {
		conf.Version = constant.ConfigVer
	}

	cfg := ini.Empty()
	cfg.ReflectFrom(&conf)

	cfg.SaveTo(configPath)
	if i118Utils.I118Prt == nil { // first time, i118 may not be init.
		logUtils.PrintToWithColor(fmt.Sprintf("Successfully update config file %s.", configPath), color.FgCyan)
	} else {
		logUtils.PrintToWithColor(i118Utils.Sprintf("success_update_config", configPath), color.FgCyan)
	}

	consts.Config = ReadCurrConfig()
	return nil
}

func getInst() model.Config {
	isSetAction := len(os.Args) > 1 && (os.Args[1] == "set" || os.Args[1] == "-set")
	if !isSetAction {
		CheckConfigReady()
	}

	ini.MapTo(&consts.Config, consts.ConfigPath)

	if consts.Config.Version < constant.ConfigVer { // old config file, re-init
		if consts.Config.Language != "en" && consts.Config.Language != "zh" {
			consts.Config.Language = "en"
		}

		SaveConfig(consts.Config)
	}

	return consts.Config
}

func CheckConfigPermission() {
	//err := syscall.Access(consts.ExeDir, syscall.O_RDWR)

	err := fileUtils.MkDirIfNeeded(consts.ExeDir + "conf")
	if err != nil {
		logUtils.PrintToWithColor(
			i118Utils.Sprintf("perm_deny", consts.ExeDir), color.FgRed)
		os.Exit(0)
	}
}

func CheckConfigReady() {
	if !fileUtils.FileExist(consts.ConfigPath) {
		InputForSet()
	}
}

func InputForSet() {
	conf := ReadCurrConfig()

	var configSite bool

	logUtils.PrintToWithColor(i118Utils.Sprintf("begin_config"), color.FgCyan)

	enCheck := ""
	var numb string
	if conf.Language == "en" {
		enCheck = "*"
		numb = "1"
	}
	zhCheck := ""
	if conf.Language == "zh" {
		zhCheck = "*"
		numb = "2"
	}

	numbSelected := stdinUtils.GetInput("(1|2)", numb, "enter_language", enCheck, zhCheck)

	if numbSelected == "1" {
		conf.Language = "en"
	} else {
		conf.Language = "zh"
	}

	stdinUtils.InputForBool(&configSite, true, "config_zentao_site")
	if configSite {
		conf.Url = stdinUtils.GetInput("((http|https)://.*)", conf.Url, "enter_url", conf.Url)
		conf.Url = getZenTaoBaseUrl(conf.Url)

		conf.Account = stdinUtils.GetInput("(.{2,})", conf.Account, "enter_account", conf.Account)
		conf.Password = stdinUtils.GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)
	}

	if commonUtils.IsWin() {
		var configInterpreter bool
		stdinUtils.InputForBool(&configInterpreter, true, "config_script_interpreter")
		if configInterpreter {
			scripts := assertUtils.GetCaseByDirAndFile([]string{"."})
			InputForScriptInterpreter(scripts, &conf, "set")
		}
	}

	SaveConfig(conf)
	PrintCurrConfig()
}

func CheckRequestConfig() {
	conf := ReadCurrConfig()
	if conf.Url == "" || conf.Account == "" || conf.Password == "" {
		InputForRequest()
	}
}

func InputForRequest() {
	conf := ReadCurrConfig()

	logUtils.PrintToWithColor(i118Utils.Sprintf("need_config"), color.FgCyan)

	conf.Url = stdinUtils.GetInput("(http://.*)", conf.Url, "enter_url", conf.Url)
	conf.Account = stdinUtils.GetInput("(.{2,})", conf.Account, "enter_account", conf.Account)
	conf.Password = stdinUtils.GetInput("(.{2,})", conf.Password, "enter_password", conf.Password)

	SaveConfig(conf)
}

func InputForScriptInterpreter(scripts []string, config *model.Config, from string) bool {
	configChanged := false

	langs := assertUtils.GetScriptType(scripts)

	for _, lang := range langs {
		if lang == "bat" || lang == "shell" {
			continue
		}

		deflt := commonUtils.GetFieldVal(*config, lang)
		if from == "run" && deflt != "" { // already set when run, "-" means ignore
			continue
		}

		if deflt == "-" {
			deflt = ""
		}
		sampleOrDefaultTips := ""
		if deflt == "" {
			sampleOrDefaultTips = i118Utils.Sprintf("for_example", langUtils.LangMap[lang]["interpreter"]) + " " +
				i118Utils.Sprintf("empty_to_ignore")
		} else {
			sampleOrDefaultTips = deflt
		}

		configChanged = true

		inter := stdinUtils.GetInputForScriptInterpreter(deflt, "set_script_interpreter", lang, sampleOrDefaultTips)
		commonUtils.SetFieldVal(config, lang, inter)
	}

	return configChanged
}

func getZenTaoBaseUrl(url string) string {
	arr := strings.Split(url, "/")

	base := url
	last := arr[len(arr)-1]
	if strings.Index(last, ".php") > -1 || strings.Index(last, ".html") > -1 ||
		strings.Index(last, "user-login") > -1 || strings.Index(last, "?") == 0 {
		base = base[:strings.LastIndex(base, "/")]
	}

	if strings.Index(base, "?") > -1 {
		base = base[:strings.LastIndex(base, "?")]
	}

	return base
}
