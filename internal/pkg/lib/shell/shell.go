package shellUtils

import (
	"bufio"
	"bytes"
	"fmt"
	_commonUtils "github.com/aaronchen2k/deeptest/internal/command/utils/common"
	_logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"github.com/kataras/iris/v12/websocket"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func ExeSysCmd(cmdStr string) (string, error) {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	output := out.String()

	return output, err
}

func ExeShell(cmdStr string) (string, error) {
	return ExeShellInDir(cmdStr, "")
}

func ExeShellInDir(cmdStr string, dir string) (ret string, err error) {
	ret, err, _ = ExeShellInDirWithPid(cmdStr, dir)
	return
}

func ExeShellWithPid(cmdStr string) (string, error, int) {
	return ExeShellInDirWithPid(cmdStr, "")
}

func ExeShellInDirWithPid(cmdStr string, dir string) (ret string, err error, pid int) {
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

	err = cmd.Run()
	if err != nil {
		logUtils.Error(i118Utils.Sprintf("fail_to_exec_command", cmdStr, cmd.Dir, err))
	}

	pid = cmd.Process.Pid
	ret = stringUtils.TrimAll(out.String())
	return
}

func ExeShellWithOutput(cmdStr string) ([]string, error) {
	return ExeShellWithOutputInDir(cmdStr, "")
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
		_logUtils.Screen(strings.TrimRight(line, "\n"))
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
				scriptInterpreter = _commonUtils.GetFieldVal(consts.Config, stringUtils.UcFirst(lang))
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
			_logUtils.Screen("chmod error" + err.Error())
		}

		filePath = "\"" + filePath + "\""
		cmd = exec.Command("/bin/bash", "-c", filePath)
	}

	if consts.ServerWorkDir != "" {
		cmd.Dir = consts.ServerWorkDir
	}

	if cmd == nil {
		msg := "error cmd is nil"
		_logUtils.Screen(msg)
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

func ExeShellWithOutputInDir(cmdStr string, dir string) ([]string, error) {
	return ExeShellWithEnvVarsAndOutputInDir(cmdStr, dir, nil)
}

func ExeShellWithEnvVarsAndOutputInDir(cmdStr, dir string, envVars []string) ([]string, error) {
	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	if dir != "" {
		cmd.Dir = dir
	}
	if envVars != nil && len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}

	output := make([]string, 0)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return output, err
	}

	cmd.Start()

	if err != nil {
		return output, err
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		logUtils.Info(strings.TrimRight(line, "\n"))
		output = append(output, line)
	}

	cmd.Wait()

	return output, nil
}

func ExeShellCallback(ch chan int, cmdStr, dir string,
	fun func(info string, msg websocket.Message), msg websocket.Message) (err error) {

	var cmd *exec.Cmd
	if commonUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	if dir != "" {
		cmd.Dir = dir
	}

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return
	}

	cmd.Start()

	if err != nil {
		return
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}

		line = strings.Trim(line, "\n")
		fun(line, msg)

		select {
		case <-ch:
			fmt.Println("exiting...")
			ch <- 1
			return
		default:
			fmt.Println("continue...")
		}
	}

	cmd.Wait()
	return
}

func GetProcess(app string) (string, error) {
	var cmd *exec.Cmd

	tmpl := ""
	cmdStr := ""
	if commonUtils.IsWin() {
		tmpl = `tasklist`
		cmdStr = fmt.Sprintf(tmpl)

		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		tmpl = `ps -ef | grep "%s" | grep -v "grep" | awk '{print $2}'`
		cmdStr = fmt.Sprintf(tmpl, app)

		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	output := ""
	if commonUtils.IsWin() {
		arr := strings.Split(out.String(), "\n")
		for _, line := range arr {
			if strings.Index(line, app+".exe") > -1 {
				arr2 := regexp.MustCompile(`\s+`).Split(line, -1)
				output = arr2[1]
				break
			}
		}
	} else {
		output = out.String()
	}

	return output, err
}

func KillProcess(app string) (string, error) {
	var cmd *exec.Cmd

	tmpl := ""
	cmdStr := ""
	if commonUtils.IsWin() {
		// tasklist | findstr ztf.exe
		tmpl = `taskkill.exe /f /im %s.exe`
		cmdStr = fmt.Sprintf(tmpl, app)

		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		tmpl = `ps -ef | grep '%s' | grep -v "grep" | awk '{print $2}' | xargs kill -9`
		cmdStr = fmt.Sprintf(tmpl, app)

		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	output := out.String()

	return output, err
}

func KillProcessById(pid int) {
	cmdStr := fmt.Sprintf("kill -9 %d", pid)
	ExeShell(cmdStr)
}
