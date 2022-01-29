package logUtils

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/command/utils/const"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/fatih/color"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

var Logger *logrus.Logger

func GetWholeLine(msg string, char string) string {
	prefixLen := (consts.ScreenWidth - utf8.RuneCountInString(msg) - 2) / 2
	if prefixLen <= 0 { // no width in debug mode
		prefixLen = 6
	}
	postfixLen := consts.ScreenWidth - utf8.RuneCountInString(msg) - 2 - prefixLen - 1
	if postfixLen <= 0 { // no width in debug mode
		postfixLen = 6
	}

	preFixStr := strings.Repeat(char, prefixLen)
	postFixStr := strings.Repeat(char, postfixLen)

	return fmt.Sprintf("%s %s %s", preFixStr, msg, postFixStr)
}

func ColoredStatus(status string) string {
	temp := strings.ToLower(status)

	switch temp {
	case "pass":
		return color.GreenString(i118Utils.Sprintf(temp))
	case "fail":
		return color.RedString(i118Utils.Sprintf(temp))
	case "skip":
		return color.YellowString(i118Utils.Sprintf(temp))
	}

	return status
}

func InitLogger() *logrus.Logger {
	consts.LogDir = fileUtils.GetLogDir()

	if Logger != nil && consts.RunMode != constant.RunModeRequest {
		return Logger
	}

	Logger = logrus.New()
	Logger.Out = ioutil.Discard

	pathMap := lfshook.PathMap{
		logrus.InfoLevel: consts.LogDir + "log.txt",
		logrus.WarnLevel: consts.LogDir + "result.txt",

		logrus.ErrorLevel: consts.LogDir + "err.txt",
	}

	Logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&MyFormatter{},
	))

	Logger.SetFormatter(&MyFormatter{})

	return Logger
}

func Screen(msg string) {
	PrintTo(msg)
}
func Log(msg string) {
	Logger.Infoln(msg)
}
func Result(msg string) {
	Logger.Warnln(msg)
}
func Error(msg string) {
	Logger.Errorln(msg)
}

func ScreenAndResult(msg string) {
	Screen(msg)
	Result(msg)
}

type MyFormatter struct {
	logrus.TextFormatter
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}
