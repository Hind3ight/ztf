package logUtils

import (
	"encoding/json"
	"fmt"
	constant "github.com/aaronchen2k/deeptest/internal/command/utils/const"
	fileUtils "github.com/aaronchen2k/deeptest/internal/command/utils/file"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

var LoggerStandard *zap.Logger
var LoggerExecConsole *zap.Logger

var LoggerExecFile *zap.Logger
var LoggerExecResult *zap.Logger

var (
	usageFile  = fmt.Sprintf("res%sdoc%susage.txt", string(os.PathSeparator), string(os.PathSeparator))
	sampleFile = fmt.Sprintf("res%sdoc%ssample.txt", string(os.PathSeparator), string(os.PathSeparator))
)

func Info(str string) {
	LoggerStandard.Info(str)
}
func Infof(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerStandard.Info(msg)
}
func Warn(str string) {
	LoggerStandard.Warn(str)
}
func Warnf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerStandard.Warn(msg)
}
func Error(str string) {
	LoggerStandard.Error(str)
}
func Errorf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerStandard.Error(msg)
}

func ExecConsole(attr color.Attribute, str string) {
	msg := color.New(attr).Sprint(str)
	LoggerExecConsole.Info(msg)
}
func ExecConsolef(clr color.Attribute, str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	msg = color.New(clr).Sprint(msg)

	LoggerExecConsole.Info(msg)
}

func ExecFile(str string) {
	LoggerExecFile.Info(str)
}
func ExecFilef(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerExecFile.Info(msg)
}

func ExecResult(str string) {
	LoggerExecResult.Info(str)
}
func ExecResultf(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	LoggerExecResult.Info(msg)
}

func PrintUnicode(str []byte) {
	msg := ConvertUnicode(str)
	LoggerStandard.Info(msg)
}

func ConvertUnicode(str []byte) string {
	var a interface{}

	temp := strings.Replace(string(str), "\\\\", "\\", -1)

	err := json.Unmarshal([]byte(temp), &a)

	var msg string
	if err == nil {
		msg = fmt.Sprint(a)
	} else {
		msg = temp
	}

	return msg
}

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

func PrintUsage() {
	PrintToWithColor("Usage: ", color.FgCyan)

	usage := fileUtils.ReadResData(usageFile)
	exeFile := constant.AppName
	if commonUtils.IsWin() {
		exeFile += ".exe"
	}
	usage = fmt.Sprintf(usage, exeFile)
	fmt.Printf("%s\n", usage)

	PrintToWithColor("\nExample: ", color.FgCyan)
	sample := fileUtils.ReadResData(sampleFile)
	if !commonUtils.IsWin() {
		regx, _ := regexp.Compile(`\\`)
		sample = regx.ReplaceAllString(sample, "/")

		regx, _ = regexp.Compile(constant.AppName + `.exe`)
		sample = regx.ReplaceAllString(sample, constant.AppName)

		regx, _ = regexp.Compile(`/bat/`)
		sample = regx.ReplaceAllString(sample, "/shell/")

		regx, _ = regexp.Compile(`\.bat\s{4}`)
		sample = regx.ReplaceAllString(sample, ".shell")
	}
	fmt.Printf("%s\n", sample)
}

func PrintTo(str string) {
	output := color.Output
	fmt.Fprint(output, str+"\n")
}
func PrintTof(format string, params ...interface{}) {
	output := color.Output
	fmt.Fprintf(output, format+"\n", params...)
}

func PrintToWithColor(msg string, attr color.Attribute) {
	output := color.Output

	if attr == -1 {
		fmt.Fprint(output, msg+"\n")
	} else {
		color.New(attr).Fprintf(output, msg+"\n")
	}
}

func PrintToCmd(msg string, attr color.Attribute) {
	output := color.Output

	if attr == -1 {
		fmt.Fprint(output, msg+"\n")
	} else {
		clr := color.New(attr)
		clr.Fprint(output, msg+"\n")
	}
}
