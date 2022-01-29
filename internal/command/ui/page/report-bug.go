package page

import (
	"fmt"
	zentaoService "github.com/aaronchen2k/deeptest/internal/command/service/zentao"
	"github.com/aaronchen2k/deeptest/internal/command/ui"
	"github.com/aaronchen2k/deeptest/internal/command/ui/widget"
	constant "github.com/aaronchen2k/deeptest/internal/command/utils/const"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"

	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/awesome-gocui/gocui"
	"github.com/fatih/color"
	"strings"
)

var filedValMap map[string]int

func InitReportBugPage(resultDir string, caseId string) error {
	DestoryReportBugPage()

	consts.CurrBug, consts.CurrBugStepIds = zentaoService.PrepareBug(resultDir, caseId)
	bug := consts.CurrBug

	w, h := consts.Cui.Size()
	x := 1
	y := 1

	var bugVersion string
	for _, val := range bug.OpenedBuild { // 取字符串值显示
		bugVersion = val
	}

	// title
	left := x
	right := left + widget.TextWidthFull - 5
	titleInput := widget.NewTextWidget("titleInput", left, y, widget.TextWidthFull-5, bug.Title)
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], titleInput.Name())

	// steps
	left = right + ui.Space
	stepsWidth := w - left - 3
	stepsInput := widget.NewTextareaWidget("stepsInput", left, y, stepsWidth, h-constant.CmdViewHeight-2, bug.Steps)
	stepsInput.Title = i118Utils.Sprintf("steps")
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], stepsInput.Name())

	// module
	y += 3
	left = x
	right = left + widget.SelectWidth
	moduleInput := widget.NewSelectWidgetWithDefault("module", left, y, widget.SelectWidth, 6,
		i118Utils.Sprintf("module"),
		consts.ZenTaoBugFields.Modules, zentaoService.GetNameById(bug.Module, consts.ZenTaoBugFields.Modules),
		bugSelectFieldCheckEvent())
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], moduleInput.Name())

	// type
	left = right + ui.Space
	right = left + widget.SelectWidth
	typeInput := widget.NewSelectWidgetWithDefault("type", left, y, widget.SelectWidth, 6,
		i118Utils.Sprintf("category"),
		consts.ZenTaoBugFields.Categories, zentaoService.GetNameById(bug.Type, consts.ZenTaoBugFields.Categories),
		bugSelectFieldCheckEvent())
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], typeInput.Name())

	// version
	left = right + ui.Space
	right = left + widget.SelectWidth
	versionInput := widget.NewSelectWidgetWithDefault("version", left, y, widget.SelectWidth, 6,
		i118Utils.Sprintf("version"),
		consts.ZenTaoBugFields.Versions, zentaoService.GetNameById(bugVersion, consts.ZenTaoBugFields.Versions),
		bugSelectFieldCheckEvent())
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], versionInput.Name())

	// severity
	y += 7
	left = x
	right = left + widget.SelectWidth
	severityInput := widget.NewSelectWidgetWithDefault("severity", left, y, widget.SelectWidth, 6,
		i118Utils.Sprintf("severity"),
		consts.ZenTaoBugFields.Severities, zentaoService.GetNameById(bug.Severity, consts.ZenTaoBugFields.Severities),
		bugSelectFieldCheckEvent())
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], severityInput.Name())

	// priority
	left = right + ui.Space
	right = left + widget.SelectWidth
	priorityInput := widget.NewSelectWidgetWithDefault("priority", left, y, widget.SelectWidth, 6,
		i118Utils.Sprintf("priority"),
		consts.ZenTaoBugFields.Priorities, zentaoService.GetNameById(bug.Pri, consts.ZenTaoBugFields.Priorities),
		bugSelectFieldCheckEvent())
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], priorityInput.Name())

	// msg
	y += 7
	left = x
	reportBugMsg := widget.NewPanelWidget("reportBugMsg", left, y, widget.TextWidthFull-5, 2, "")
	reportBugMsg.Frame = false
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], reportBugMsg.Name())

	// buttons
	y += 5
	buttonX := x + widget.SelectWidth + ui.Space
	submitInput := widget.NewButtonWidgetAutoWidth("submitInput", buttonX, y,
		i118Utils.Sprintf("submit"), reportBug)
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], submitInput.Name())

	cancelReportBugInput := widget.NewButtonWidgetAutoWidth("cancelReportBugInput",
		buttonX+11, y, i118Utils.Sprintf("cancel"), cancelReportBug)
	ui.ViewMap["reportBug"] = append(ui.ViewMap["reportBug"], cancelReportBugInput.Name())

	ui.BindEventForInputWidgets(ui.ViewMap["reportBug"])

	consts.Cui.SetCurrentView("titleInput")

	return nil
}

func reportBug(g *gocui.Gui, v *gocui.View) error {
	bug := consts.CurrBug

	titleView, _ := g.View("titleInput")
	stepsView, _ := g.View("stepsInput")
	moduleView, _ := g.View("module")
	typeView, _ := g.View("type")
	versionView, _ := g.View("version")
	severityView, _ := g.View("severity")
	priorityView, _ := g.View("priority")

	title := strings.TrimSpace(titleView.Buffer())
	stepsStr := strings.TrimSpace(stepsView.Buffer())

	moduleStr := strings.TrimSpace(ui.GetSelectedRowVal(moduleView))
	typeStr := strings.TrimSpace(ui.GetSelectedRowVal(typeView))
	versionStr := strings.TrimSpace(ui.GetSelectedRowVal(versionView))
	severityStr := strings.TrimSpace(ui.GetSelectedRowVal(severityView))
	priorityStr := strings.TrimSpace(ui.GetSelectedRowVal(priorityView))

	if title == "" {
		v, _ := consts.Cui.View("reportBugMsg")
		color.New(color.FgMagenta).Fprintf(v, i118Utils.ReadI18nJson("title_cannot_be_empty"))
		return nil
	}

	bug.Title = title
	bug.Steps = strings.Replace(stepsStr, "\n", "<br/>", -1)

	bug.Type = zentaoService.GetIdByName(typeStr, consts.ZenTaoBugFields.Categories)

	bug.Module = zentaoService.GetIdByName(moduleStr, consts.ZenTaoBugFields.Modules)

	versionKey := zentaoService.GetIdByName(versionStr, consts.ZenTaoBugFields.Versions)
	build := make(map[string]string)
	if versionKey == "trunk" {
		build["0"] = "trunk"
	} else {
		build[versionKey] = versionStr
	}
	bug.OpenedBuild = build

	bug.Severity = zentaoService.GetIdByName(severityStr, consts.ZenTaoBugFields.Severities)
	bug.Pri = zentaoService.GetIdByName(priorityStr, consts.ZenTaoBugFields.Priorities)

	consts.CurrBug = bug
	ok, msg := zentaoService.CommitBug()

	msgView, _ := consts.Cui.View("reportBugMsg")
	msgView.Clear()

	if ok {
		color.New(color.FgGreen).Fprintf(msgView, msg)

		consts.Cui.DeleteView("submitInput")

		cancelReportBugInput, _ := consts.Cui.View("cancelReportBugInput")
		cancelReportBugInput.Clear()
		fmt.Fprint(cancelReportBugInput, " "+i118Utils.Sprintf("close"))
	} else {
		color.New(color.FgMagenta).Fprintf(msgView, msg)
	}

	return nil
}

func bugSelectFieldCheckEvent() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		name := v.Name()

		g.SetCurrentView(name)

		return nil
	}
}

func init() {
	filedValMap = make(map[string]int)
}

func cancelReportBug(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func DestoryReportBugPage() {
	for _, v := range ui.ViewMap["reportBug"] {
		consts.Cui.DeleteView(v)
		consts.Cui.DeleteKeybindings(v)
	}
}
