package shellUtils

import (
	"bufio"
	"bytes"
	"fmt"
	commonUtils "github.com/aaronchen2k/deeptest/internal/command/utils/common"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	stringUtils "github.com/aaronchen2k/deeptest/internal/command/utils/string"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"

	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"io"
	"os"
	"os/exec"
	"strings"
)

func ExeSysCmd(cmdStr string) (string, error) {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command(cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}
func ExeAppInDir(cmdStr string, dir string) (string, error) {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	if dir != "" {
		cmd.Dir = dir
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}

func ExeAppWithOutput(cmdStr string) []string {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	if consts.ServerProjectDir != "" {
		cmd.Dir = consts.ServerProjectDir
	}

	output := make([]string, 0)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return output
	}

	cmd.Start()

	if err != nil {
		output = append(output, fmt.Sprint(err))
		return output
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		logUtils.Screen(strings.TrimRight(line, "\n"))
		output = append(output, line)
	}

	cmd.Wait()

	return output
}

func ExecScriptFile(filePath string) (string, string) {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		lang := langUtils.GetLangByFile(filePath)

		scriptInterpreter := ""
		if strings.ToLower(lang) != "bat" {
			if consts.Interpreter != "" {
				scriptInterpreter = consts.Interpreter
				fmt.Printf("use interpreter %s for script %s\n", scriptInterpreter, filePath)
			} else {
				scriptInterpreter = commonUtils.GetFieldVal(consts.Config, stringUtils.Ucfirst(lang))
			}
		}
		if scriptInterpreter != "" {
			if strings.Index(strings.ToLower(scriptInterpreter), "autoit") > -1 {
				cmd = exec.Command("cmd", "/C", scriptInterpreter, filePath, "|", "more")
			} else {
				cmd = exec.Command("cmd", "/C", scriptInterpreter, filePath)
			}
		} else if strings.ToLower(lang) == "bat" {
			cmd = exec.Command("cmd", "/C", filePath)
		} else {
			fmt.Printf("use interpreter %s for script %s\n", scriptInterpreter, filePath)
			i118Utils.I118Prt.Printf("no_interpreter_for_run", filePath, lang)
		}
	} else {
		err := os.Chmod(filePath, 0777)
		if err != nil {
			logUtils.Screen("chmod error" + err.Error())
		}

		filePath = "\"" + filePath + "\""
		cmd = exec.Command("/bin/bash", "-c", filePath)
	}

	if consts.ServerWorkDir != "" {
		cmd.Dir = consts.ServerWorkDir
	}

	if cmd == nil {
		msg := "error cmd is nil"
		logUtils.Screen(msg)
		return "", fmt.Sprint(msg)
	}

	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		fmt.Println(err1)
		return "", fmt.Sprint(err1)
	} else if err2 != nil {
		fmt.Println(err2)
		return "", fmt.Sprint(err2)
	}

	cmd.Start()

	reader1 := bufio.NewReader(stdout)
	output1 := make([]string, 0)
	for {
		line, err2 := reader1.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		output1 = append(output1, line)
	}

	reader2 := bufio.NewReader(stderr)
	output2 := make([]string, 0)
	for {
		line, err2 := reader2.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		output2 = append(output2, line)
	}

	cmd.Wait()

	return strings.Join(output1, ""), strings.Join(output2, "")
}
