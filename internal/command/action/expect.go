package action

import (
	commonUtils "github.com/aaronchen2k/deeptest/internal/command/utils/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/command/utils/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/command/utils/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	shellUtils "github.com/aaronchen2k/deeptest/internal/command/utils/shell"
	stringUtils "github.com/aaronchen2k/deeptest/internal/command/utils/string"
	"github.com/aaronchen2k/deeptest/internal/command/utils/vari"
	assertUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/assert"
	configUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/config"
	"path"
	"path/filepath"
	"strings"
)

func GenExpectFiles(files []string) error {
	logUtils.InitLogger()

	cases := assertUtils.GetCaseByDirAndFile(files)

	if len(cases) < 1 {
		logUtils.PrintTo("\n" + i118Utils.Sprintf("no_cases"))
		return nil
	}

	casesToRun, _ := filterCases(cases)
	dryRunScripts(casesToRun)

	return nil
}

func filterCases(cases []string) (casesToRun, casesToIgnore []string) {
	// config interpreter if needed
	if commonUtils.IsWin() {
		conf := configUtils.ReadCurrConfig()
		configChanged := configUtils.InputForScriptInterpreter(cases, &conf, "run")
		if configChanged {
			configUtils.SaveConfig(conf)
		}
	}

	conf := configUtils.ReadCurrConfig()
	for _, cs := range cases {
		if commonUtils.IsWin() {
			if path.Ext(cs) == ".sh" { // filter by os
				continue
			}

			ext := path.Ext(cs)
			if ext != "" {
				ext = ext[1:]
			}
			lang := vari.ScriptExtToNameMap[ext]
			interpreter := commonUtils.GetFieldVal(conf, stringUtils.Ucfirst(lang))
			if interpreter == "-" && vari.Interpreter == "" { // not to ignore if interpreter set
				interpreter = ""

				casesToIgnore = append(casesToIgnore, cs)
			}
			if lang != "bat" && interpreter == "" { // ignore the ones with no interpreter set
				continue
			}
		} else if !commonUtils.IsWin() { // filter by os
			if path.Ext(cs) == ".bat" {
				continue
			}
		}

		casesToRun = append(casesToRun, cs)
	}

	return
}

func dryRunScripts(casesToRun []string) {
	for _, file := range casesToRun {
		dryRunScript(file)
	}
}

func dryRunScript(file string) {
	out, _ := shellUtils.ExecScriptFile(file)
	out = strings.Trim(out, "\n")

	expFile := filepath.Join(filepath.Dir(file), fileUtils.GetFileNameWithoutExt(file)+".exp")
	fileUtils.WriteFile(expFile, out)
}
