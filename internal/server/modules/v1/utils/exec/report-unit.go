package scriptUtils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	dateUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/date"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GenUnitTestReport(req serverDomain.WsReq, startTime, endTime int64,
	ch chan int, sendOutputMsg, sendExecMsg func(info, isRunning string, wsMsg websocket.Message), wsMsg websocket.Message) (
	report commDomain.ZtfReport) {

	testSuites := RetrieveUnitResult(req.ProjectPath, startTime, req.Framework, req.Tool)
	cases, classNameMaxWidth, duration := ParserUnitTestResult(testSuites)

	if duration == 0 {
		duration = float32(endTime - startTime)
	}

	report = commDomain.ZtfReport{
		TestEnv:       commonUtils.GetOs(),
		TestType:      commConsts.TestUnit,
		TestFramework: req.Framework,
		TestTool:      req.Tool,
		Pass:          0, Fail: 0, Total: 0}

	failedCount := 0
	failedCaseLines := make([]string, 0)
	failedCaseLinesDesc := make([]string, 0)

	for idx, cs := range cases {
		if cs.Failure != nil {
			report.Fail++

			if failedCount > 0 { // 换行
				failedCaseLinesDesc = append(failedCaseLinesDesc, "")
			}
			className := cases[idx].TestSuite

			line := fmt.Sprintf("[%s] %d.%s", className, cs.Id, cs.Title)
			failedCaseLines = append(failedCaseLines, line)

			failedCaseLinesDesc = append(failedCaseLinesDesc, line)
			failDesc := fmt.Sprintf("   %s - %s", cs.Failure.Type, cs.Failure.Desc)
			failedCaseLinesDesc = append(failedCaseLinesDesc, failDesc)
		} else {
			report.Pass++
		}
		report.Total++

		if startTime == 0 {
			if report.StartTime == 0 || cs.StartTime < report.StartTime {
				report.StartTime = cs.StartTime
			}
			if cs.EndTime > report.EndTime {
				report.EndTime = cs.EndTime
			}
		}
	}
	report.UnitResult = cases
	if duration == 0 {
		report.Duration = report.EndTime - report.StartTime
	} else {
		report.Duration = int64(duration)
	}

	postFix := ":"
	if len(cases) == 0 {
		postFix = "."
	}

	temp := i118Utils.Sprintf("found_scripts", strconv.Itoa(len(cases))) + postFix
	sendExecMsg(temp, "", wsMsg)
	logUtils.ExecConsolef(color.FgCyan, temp)
	logUtils.ExecResult(temp)

	width := strconv.Itoa(len(strconv.Itoa(report.Total)))
	for idx, cs := range cases {
		testSuite := stringUtils.AddPostfix(cs.TestSuite, classNameMaxWidth, " ")

		format := "(%" + width + "d/%d) %s [%s] [%" + width + "d. %s] (%.3fs)"
		msg := fmt.Sprintf(format, idx+1, report.Total, cs.Status, testSuite, cs.Id, cs.Title, cs.Duration)
		sendExecMsg(msg, "", wsMsg)
		logUtils.ExecConsolef(color.FgCyan, temp)
		logUtils.ExecResult(msg)
	}

	if report.Fail > 0 {
		msg := "\n" + i118Utils.Sprintf("failed_scripts")
		msg += strings.Join(failedCaseLines, "\n")
		msg += strings.Join(failedCaseLinesDesc, "\n")

		sendExecMsg(msg, "", wsMsg)
		logUtils.ExecConsolef(color.FgCyan, temp)
		logUtils.ExecResult(msg)
	}

	// 生成统计行
	secTag := ""
	if commConsts.Language == "en" && report.Duration > 1 {
		secTag = "s"
	}

	fmtStr := "%d(%.1f%%) %s"
	passRate := 0
	failRate := 0
	skipRate := 0
	if report.Total > 0 {
		passRate = report.Pass * 100 / report.Total
		failRate = report.Fail * 100 / report.Total
		skipRate = report.Skip * 100 / report.Total
	}

	passStr := fmt.Sprintf(fmtStr, report.Pass, float32(passRate), i118Utils.Sprintf("pass"))
	failStr := fmt.Sprintf(fmtStr, report.Fail, float32(failRate), i118Utils.Sprintf("fail"))
	skipStr := fmt.Sprintf(fmtStr, report.Skip, float32(skipRate), i118Utils.Sprintf("skip"))

	// 执行%d个用例，耗时%d秒%s。%s，%s，%s。报告%s。
	msg := dateUtils.DateTimeStr(time.Now()) + " " +
		i118Utils.Sprintf("run_result",
			report.Total, report.Duration, secTag,
			passStr, failStr, skipStr,
		)
	sendExecMsg(msg, "", wsMsg)
	logUtils.ExecConsole(color.FgCyan, msg)
	logUtils.ExecResult(msg)

	resultPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultText)
	msg = "                    " + i118Utils.Sprintf("run_report", resultPath) + "\n"

	sendExecMsg(msg, "false", wsMsg)
	logUtils.ExecConsole(color.FgCyan, msg)
	logUtils.ExecResult(msg)

	//report.ProductId, _ = strconv.Atoi(vari.ProductId)
	json, _ := json.MarshalIndent(report, "", "\t")
	jsonPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultJson)
	fileUtils.WriteFile(jsonPath, string(json))

	return
}

func RetrieveUnitResult(projectPath string, startTime int64,
	testFramework commConsts.UnitTestFramework, testTool commConsts.UnitTestTool) (suites []commDomain.UnitTestSuite) {
	resultFiles := make([]string, 0)

	resultDir := ""

	if testFramework == commConsts.JUnit && testTool == commConsts.Maven {
		resultDir = filepath.Join("target", "surefire-reports")
	} else if testFramework == commConsts.TestNG && testTool == commConsts.Maven {
		resultDir = filepath.Join("target", "surefire-reports", "junitreports")
	} else if testFramework == commConsts.RobotFramework || testFramework == commConsts.Cypress {
		resultDir = "results"
	} else {
		resultDir = "results"
	}

	resultDir = filepath.Join(projectPath, resultDir)
	resultFiles, _ = GetSuiteFiles(resultDir, startTime)

	for _, file := range resultFiles {
		testSuite, err := GetTestSuite(file, testFramework)

		if err == nil {
			suites = append(suites, testSuite)
		}
	}

	return
}

func GetSuiteFiles(resultDir string, startTime int64) (resultFiles []string, err error) {
	if fileUtils.IsDir(resultDir) {
		dir, err := ioutil.ReadDir(resultDir)
		if err == nil {
			for _, fi := range dir {
				name := fi.Name()
				ext := path.Ext(name)
				if ext == ".xml" && fi.ModTime().Unix() >= startTime {
					pth := filepath.Join(resultDir, name)
					resultFiles = append(resultFiles, pth)
				}
			}
		}
	} else {
		resultFiles = append(resultFiles, resultDir)
	}

	return
}

func GetTestSuite(xmlFile string, testFramework commConsts.UnitTestFramework) (
	testSuite commDomain.UnitTestSuite, err error) {

	content := fileUtils.ReadFile(xmlFile)

	if testFramework == commConsts.JUnit || testFramework == commConsts.TestNG {
		testSuite = commDomain.UnitTestSuite{}
		err = xml.Unmarshal([]byte(content), &testSuite)

	} else if testFramework == commConsts.PHPUnit {
		phpTestSuite := commDomain.PhpUnitSuites{}
		err = xml.Unmarshal([]byte(content), &phpTestSuite)
		if err == nil {
			testSuite = ConvertPhpUnitResult(phpTestSuite)
		}
	} else if testFramework == commConsts.PyTest {
		pyTestSuite := commDomain.PyTestSuites{}
		err = xml.Unmarshal([]byte(content), &pyTestSuite)
		if err == nil {
			testSuite = ConvertPyTestResult(pyTestSuite)
		}
	} else if testFramework == commConsts.Jest {
		jestSuite := commDomain.JestSuites{}
		err = xml.Unmarshal([]byte(content), &jestSuite)
		if err == nil {
			testSuite = ConvertJestResult(jestSuite)
		}
	} else if testFramework == commConsts.GTest {
		gTestSuite := commDomain.GTestSuites{}
		err = xml.Unmarshal([]byte(content), &gTestSuite)
		if err == nil {
			testSuite = ConvertGTestResult(gTestSuite)
		}
	} else if testFramework == commConsts.QTest {
		qTestSuite := commDomain.QTestSuites{}
		err = xml.Unmarshal([]byte(content), &qTestSuite)
		if err == nil {
			testSuite = ConvertQTestResult(qTestSuite)
		}
	} else if testFramework == commConsts.CppUnit {
		content = strings.Replace(content, "ISO-8859-1", "UTF-8", -1)

		cppUnitSuites := commDomain.CppUnitSuites{}
		err = xml.Unmarshal([]byte(content), &cppUnitSuites)
		if err == nil {
			testSuite = ConvertCppUnitResult(cppUnitSuites)
		}
	} else if testFramework == commConsts.RobotFramework {
		robotResult := commDomain.RobotResult{}
		err = xml.Unmarshal([]byte(content), &robotResult)
		if err == nil {
			testSuite = ConvertRobotResult(robotResult)
		}
	} else if testFramework == commConsts.Cypress {
		cyResult := commDomain.CypressTestsuites{}
		err = xml.Unmarshal([]byte(content), &cyResult)
		if err == nil {
			testSuite = ConvertCyResult(cyResult)
		}
	}

	return
}

func ParserUnitTestResult(testSuites []commDomain.UnitTestSuite) (
	cases []commDomain.UnitResult, classNameMaxWidth int, dur float32) {

	idx := 1
	for _, suite := range testSuites {
		if suite.Time != 0 { // for junit, there is a time on suite level
			dur += suite.Time
		}

		for _, cs := range suite.Cases {
			cs.Id = idx

			if cs.Failure != nil {
				cs.Status = "fail"

				cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "<![CDATA[", "", -1)
				cs.Failure.Desc = strings.Replace(cs.Failure.Desc, "]]>", "", -1)
			} else {
				cs.Status = "pass"
			}

			lent2 := runewidth.StringWidth(cs.TestSuite)
			if lent2 > classNameMaxWidth {
				classNameMaxWidth = lent2
			}

			cases = append(cases, cs)
			idx++
		}
	}

	return
}

func ConvertJestResult(jestSuite commDomain.JestSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}
	testSuite.Time = jestSuite.Time

	for _, suite := range jestSuite.TestSuites {
		for _, cs := range testSuite.Cases {
			caseResult := commDomain.UnitResult{}
			caseResult.Title = cs.Title
			caseResult.Duration = cs.Duration

			if suite.Title != "" && suite.Title != "undefined" {
				caseResult.TestSuite = suite.Title
			} else {
				caseResult.TestSuite = jestSuite.Title
			}

			caseResult.Failure = cs.Failure

			testSuite.Cases = append(testSuite.Cases, caseResult)
		}
	}

	return testSuite
}

func ConvertPhpUnitResult(phpUnitSuite commDomain.PhpUnitSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	var total float32 = 0
	for _, cs := range phpUnitSuite.Cases {
		caseResult := commDomain.UnitResult{}
		caseResult.Title = cs.Title
		caseResult.Duration = cs.Time

		total += cs.Time

		if cs.Groups != "" && cs.Groups != "default" {
			caseResult.TestSuite = cs.Groups
		} else {
			caseResult.TestSuite = cs.TestSuite
		}

		if cs.Status != 0 {
			fail := commDomain.Failure{}
			fail.Desc = cs.Fail
			caseResult.Failure = &fail
		}

		testSuite.Cases = append(testSuite.Cases, caseResult)
	}
	testSuite.Duration = int64(total)
	testSuite.Time = total

	return testSuite
}

func ConvertPyTestResult(pytestSuites commDomain.PyTestSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	var total float32 = 0
	for _, suite := range pytestSuites.TestSuites {
		total += suite.Time

		for _, cs := range suite.Cases {
			caseResult := commDomain.UnitResult{}
			caseResult.Title = cs.Title
			caseResult.Duration = cs.Duration

			if suite.Title != "" && suite.Title != "pytest" {
				caseResult.TestSuite = suite.Title
			} else {
				caseResult.TestSuite = cs.TestSuite
			}

			if cs.Failure != nil {
				fail := commDomain.Failure{}
				fail.Type = cs.Failure.Type
				fail.Desc = cs.Failure.Desc
				caseResult.Failure = &fail
			} else if cs.Error != nil {
				fail := commDomain.Failure{}
				fail.Type = cs.Error.Message
				fail.Desc = cs.Error.Text
				caseResult.Failure = &fail
			}

			testSuite.Cases = append(testSuite.Cases, caseResult)

		}
	}

	testSuite.Duration = int64(total)
	testSuite.Time = total

	return testSuite
}

func ConvertGTestResult(gTestSuite commDomain.GTestSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}
	testSuite.Time = gTestSuite.Time

	for _, suite := range gTestSuite.TestSuites {
		for _, cs := range suite.Cases {
			caseResult := commDomain.UnitResult{}
			caseResult.Title = cs.Title
			caseResult.Duration = cs.Duration
			caseResult.Status = cs.Status

			if suite.Title != "" && suite.Title != "pytest" {
				caseResult.TestSuite = suite.Title
			}

			if cs.Failure != nil {
				fail := commDomain.Failure{}
				fail.Type = cs.Failure.Type
				fail.Desc = cs.Failure.Desc
				caseResult.Failure = &fail
			}

			testSuite.Cases = append(testSuite.Cases, caseResult)

		}
	}

	return testSuite
}

func ConvertCppUnitResult(cppunitSuite commDomain.CppUnitSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	for _, cs := range cppunitSuite.FailedTests.Cases {
		caseResult := commDomain.UnitResult{}
		caseResult.Id = cs.Id
		caseResult.Title = cs.Title

		fail := commDomain.Failure{}
		fail.Type = cs.FailureType
		fail.Desc = cs.Message
		caseResult.Failure = &fail

		testSuite.Cases = append(testSuite.Cases, caseResult)
	}

	for _, cs := range cppunitSuite.SuccessfulTests.Cases {
		caseResult := commDomain.UnitResult{}
		caseResult.Id = cs.Id
		caseResult.Title = cs.Title

		testSuite.Cases = append(testSuite.Cases, caseResult)
	}

	return testSuite
}

func ConvertQTestResult(qTestSuite commDomain.QTestSuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	for _, cs := range qTestSuite.Cases {
		caseResult := commDomain.UnitResult{}
		caseResult.TestSuite = qTestSuite.Name
		caseResult.Title = cs.Title
		caseResult.Status = cs.Result

		if cs.Failure != nil {
			fail := commDomain.Failure{}
			fail.Type = cs.Failure.Type
			fail.Desc = cs.Failure.Desc
			caseResult.Failure = &fail
		}

		testSuite.Cases = append(testSuite.Cases, caseResult)
	}

	return testSuite
}

func ConvertRobotResult(result commDomain.RobotResult) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	suiteMap := map[string]string{}
	for _, state := range result.Statistics.Suite.States {
		suiteMap[state.ID] = state.Text
	}

	tests := make([]commDomain.RobotTest, 0)
	for _, suite := range result.Suites {
		RetrieveRobotTests(suite, &tests)
	}

	for _, cs := range tests {
		caseResult := commDomain.UnitResult{}
		caseResult.Title = cs.Name
		caseResult.Status = strings.ToLower(cs.Status.Status)

		suiteId := cs.ID[0:strings.LastIndex(cs.ID, "-")]
		caseResult.TestSuite = suiteMap[suiteId]

		templ := "20060102 15:04:05.000"
		startTime, _ := time.ParseInLocation(templ, cs.Status.StartTime, time.Local)
		endTime, _ := time.ParseInLocation(templ, cs.Status.EndTime, time.Local)

		caseResult.StartTime = startTime.Unix()
		caseResult.EndTime = endTime.Unix()
		caseResult.Duration = float32(caseResult.EndTime - caseResult.StartTime)

		if caseResult.Status != "pass" {
			fail := commDomain.Failure{}
			fail.Type = ""
			fail.Desc = cs.Status.Text
			caseResult.Failure = &fail
		}

		testSuite.Cases = append(testSuite.Cases, caseResult)
	}

	return testSuite
}

func RetrieveRobotTests(suite commDomain.RobotSuite, tests *[]commDomain.RobotTest) {
	for _, suite := range suite.Suites {
		RetrieveRobotTests(suite, tests)
	}

	for _, test := range suite.Tests {
		*tests = append(*tests, test)
	}
}

func ConvertCyResult(result commDomain.CypressTestsuites) commDomain.UnitTestSuite {
	testSuite := commDomain.UnitTestSuite{}

	for _, suite := range result.Testsuites {
		if suite.Name == "Root Suite" {
			continue
		}

		templ := "20060102 15:04:05.000"
		duration := suite.Time
		startTime, _ := time.ParseInLocation(templ, suite.Timestamp, time.Local)
		//endTime := time.Unix(startTime.Unix() + int64(duration), 0)

		testSuite.Duration = int64(duration)
		testSuite.Time = float32(startTime.Unix())

		for _, cs := range suite.Testcases {
			caseResult := commDomain.UnitResult{}
			caseResult.TestSuite = suite.Name
			caseResult.Title = cs.Name
			caseResult.Duration = float32(cs.Time)

			if len(cs.Failures) > 0 {
				caseResult.Status = "fail"

				fail := commDomain.Failure{}
				fail.Type = cs.Failures[0].Type
				fail.Desc = cs.Failures[0].Message
				caseResult.Failure = &fail
			} else {
				caseResult.Status = "pass"
			}

			testSuite.Cases = append(testSuite.Cases, caseResult)
		}
	}

	return testSuite
}
