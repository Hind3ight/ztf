package vari

import commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"

var (
	IsDebug bool
	Config  = commDomain.ProjectConf{}
	//Cui            *gocui.Gui
	MainViewHeight int

	ConfigPath       string
	ExeDir           string
	ServerWorkDir    string
	ServerProjectDir string
	LogDir           string
	UnitTestType     string
	UnitTestTool     string
	UnitTestResult   string
	UnitTestResults  = "results"
	ProductId        string

	ZenTaoVersion string
	SessionVar    string
	SessionId     string
	RequestType   string
	RequestFix    string

	ScriptExtToNameMap map[string]string
	CurrScriptFile     string // scripts/tc-001.py
	CurrResultDate     string // 2019-08-15T173802
	CurrCaseId         int    // 2019-08-15T173802

	ScreenWidth  int
	ScreenHeight int
	//ZenTaoBugFields model.ZentaoBugFields
	//
	//CurrBug        model.Bug
	CurrBugStepIds string

	Verbose     bool
	Interpreter string

	// server
	RunMode     string
	IP          string
	MAC         string
	Port        int
	Platform    string
	AgentLogDir string
)
