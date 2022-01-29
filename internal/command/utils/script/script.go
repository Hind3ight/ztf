package scriptUtils

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/command/model"
	constant "github.com/aaronchen2k/deeptest/internal/command/utils/const"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/zentao"
	"github.com/emirpasic/gods/maps"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"regexp"
	"strconv"
	"strings"
)

func GetStepAndExpectMap(file string) (stepMap, stepTypeMap, expectMap maps.Map, isOldFormat bool) {
	if !fileUtils.FileExist(file) {
		return
	}

	lang := langUtils.GetLangByFile(file)
	txt := fileUtils.ReadFile(file)

	isOldFormat = strings.Index(txt, "[esac]") > -1
	_, checkpoints := zentaoUtils.ReadCaseInfo(txt, lang, isOldFormat)
	lines := strings.Split(checkpoints, "\n")

	if isOldFormat {
		groupBlockArr := getGroupBlockArr(lines)
		groupArr := getStepNestedArrObsolete(groupBlockArr)
		_, stepMap, stepTypeMap, expectMap = getSortedTextFromNestedStepsObsolete(groupArr)
	} else {
		groupArr := getStepNestedArr(lines)
		_, stepMap, stepTypeMap, expectMap = getSortedTextFromNestedSteps(groupArr)
	}

	return
}

func SortFile(file string) {
	stepsTxt := ""

	if fileUtils.FileExist(file) {
		txt := fileUtils.ReadFile(file)
		lang := langUtils.GetLangByFile(file)
		isOldFormat := strings.Index(txt, "[esac]") > -1
		info, content := zentaoUtils.ReadCaseInfo(txt, lang, isOldFormat)
		lines := strings.Split(content, "\n")

		groupBlockArr := getGroupBlockArr(lines)
		groupArr := getStepNestedArrObsolete(groupBlockArr)
		stepsTxt, _, _, _ = getSortedTextFromNestedStepsObsolete(groupArr)

		// replace info
		from := ""
		to := ""
		if isOldFormat {
			from = `(?s)\[case\].*\[esac\]`
			to = "[case]\n" + info + "\n" + stepsTxt + "\n\n[esac]"
		} else {
			from = fmt.Sprintf(`(?s)%s.*%s`, constant.LangCommentsRegxMap[lang][0], constant.LangCommentsRegxMap[lang][1])
			to = fmt.Sprintf("%s\n"+info+"\n"+stepsTxt+"\n\n%s",
				constant.LangCommentsRegxMap[lang][0], constant.LangCommentsRegxMap[lang][1])
		}
		re, _ := regexp.Compile(from)
		script := re.ReplaceAllString(txt, to)

		fileUtils.WriteFile(file, script)
	}
}

func getGroupBlockArr(lines []string) [][]string {
	groupBlockArr := make([][]string, 0)

	idx := 0
	for true {
		if idx >= len(lines) {
			break
		}

		var groupContent []string
		line := strings.TrimSpace(lines[idx])
		if isGroup(line) { // must match a group
			groupContent = make([]string, 0)
			groupContent = append(groupContent, line)

			idx++

			for true {
				if idx >= len(lines) {
					groupBlockArr = append(groupBlockArr, groupContent)
					break
				}

				line = strings.TrimSpace(lines[idx])
				if isGroup(line) {
					groupBlockArr = append(groupBlockArr, groupContent)

					break
				} else if line != "" && !isGroup(line) {
					groupContent = append(groupContent, line)
				}

				idx++
			}
		} else {
			idx++
		}
	}

	return groupBlockArr
}

func getStepNestedArrObsolete(blocks [][]string) []model.TestStep {
	ret := make([]model.TestStep, 0)

	for _, block := range blocks {
		name := block[0]
		group := model.TestStep{Desc: name}

		if isStepsIdent(block[1]) { // muti line
			group.MutiLine = true
			childs := loadMutiLineSteps(block[1:])

			group.Children = append(group.Children, childs...)
		} else {
			childs := loadSingleLineSteps(block[1:])

			group.Children = append(group.Children, childs...)
		}

		ret = append(ret, group)
	}

	return ret
}

func getStepNestedArr(lines []string) (ret []model.TestStep) {
	parent := model.TestStep{}
	increase := 0
	for index := 0; index < len(lines); index++ {
		line := lines[index]
		lineTrim := strings.TrimSpace(line)
		if lineTrim == "" || lineTrim == ">>" {
			continue
		}

		if strings.Index(line, " ") != 0 {
			parent, increase = parserNextLines(line, lines[index+1:])
			index += increase

			if strings.TrimSpace(parent.Expect) == "" && strings.Index(line, ">>") > -1 {
				parent.Expect = constant.ExpectResultPass
			}
			ret = append(ret, parent)
		} else { // 有缩进
			child := model.TestStep{}
			child, increase = parserNextLines(line, lines[index+1:])
			index += increase

			if parent.Desc != "" {
				if strings.TrimSpace(child.Expect) == "" && strings.Index(line, ">>") > -1 {
					child.Expect = constant.ExpectResultPass
				}

				ret[len(ret)-1].Children = append(ret[len(ret)-1].Children, child)
			}
		}
	}

	return
}
func parserNextLines(str string, nextLines []string) (ret model.TestStep, increase int) {
	arr := strings.Split(str, ">>")
	desc := strings.TrimSpace(arr[0])

	expect := ""
	if len(arr) > 1 {
		expect = strings.TrimSpace(arr[1])
	}

	if strings.Index(str, ">>") < 0 || expect != "" { // no >> or single line expect
		ret = model.TestStep{Desc: desc, Expect: expect}
		return
	}

	if strings.Index(str, ">>") > -1 { // will test if it has multi-line expect
		for index, line := range nextLines {
			if strings.TrimSpace(line) == ">>" {
				increase = index
				break
			}

			if strings.Index(line, ">>") > -1 {
				expect = ""
				break
			}

			if len(expect) > 0 {
				expect += " | "
			}
			expect += strings.TrimSpace(line)
		}

		if increase == 0 { // multi-line
			expect = ""
		}
	}

	ret = model.TestStep{Desc: desc, Expect: expect}
	return
}

func loadMutiLineSteps(arr []string) []model.TestStep {
	childs := make([]model.TestStep, 0)

	child := model.TestStep{}
	idx := 0
	for true {
		if idx >= len(arr) {
			if child.Desc != "" {
				childs = append(childs, child)
			}

			break
		}

		line := arr[idx]
		line = strings.TrimSpace(line)

		if isStepsIdent(line) {
			if idx > 0 {
				childs = append(childs, child)
			}

			child = model.TestStep{}
			idx++

			stp := ""
			for true { // retrieve next lines
				if idx >= len(arr) || hasBrackets(arr[idx]) {
					child.Desc = stp
					break
				}

				stp += arr[idx] + "\n"
				idx++
			}
		}

		if isExpectsIdent(line) {
			idx++

			exp := ""
			for true { // retrieve next lines
				if idx >= len(arr) || hasBrackets(arr[idx]) {
					child.Expect = exp
					break
				}

				temp := strings.TrimSpace(arr[idx])
				if temp == ">>" {
					temp = ""
				}
				exp += temp + "\n"
				idx++
			}
		}

	}

	return childs
}

func loadSingleLineSteps(arr []string) []model.TestStep {
	childs := make([]model.TestStep, 0)

	for _, line := range arr {
		line = strings.TrimSpace(line)

		sections := strings.Split(line, ">>")
		expect := ""
		if len(sections) > 1 { // has expect
			expect = strings.TrimSpace(sections[1])
		}

		child := model.TestStep{Desc: sections[0], Expect: expect}

		childs = append(childs, child)
	}

	return childs
}

func isGroupIdent(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[\s*group\s*\]`, str)
	return pass
}

func isStepsIdent(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[.*steps\.*\]`, str)
	return pass
}

func isExpectsIdent(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[.*expects\.*\]`, str)
	return pass
}

func hasBrackets(str string) bool {
	pass, _ := regexp.MatchString(`(?i)()\[.*\]`, str)
	return pass
}

func isGroup(str string) bool {
	ret := strings.Index(str, ">>") < 0 && hasBrackets(str) && !isStepsIdent(str) && !isExpectsIdent(str)

	return ret
}

func getSortedTextFromNestedStepsObsolete(groups []model.TestStep) (string, maps.Map, maps.Map, maps.Map) {
	ret := make([]string, 0)
	stepMap := linkedhashmap.New()
	stepTypeMap := linkedhashmap.New()
	expectMap := linkedhashmap.New()

	groupNumb := 1
	for _, group := range groups {
		desc := group.Desc

		if desc == "[group]" {
			ret = append(ret, "\n"+desc)

			for idx, child := range group.Children { // level 1 item
				numbStr := getNumbStr(groupNumb, -1)
				stepTypeMap.Put(numbStr, "item")

				if group.MutiLine {
					// steps
					tag := replaceNumb("[steps]", groupNumb, -1, true)
					ret = append(ret, "  "+tag)

					stepTxt := printMutiStepOrExpect(child.Desc)
					ret = append(ret, stepTxt)
					stepMap.Put(numbStr, stepTxt)

					// expects
					tag = replaceNumb("[expects]", groupNumb, -1, true)
					ret = append(ret, "  "+tag)

					expectTxt := printMutiStepOrExpect(child.Expect)
					ret = append(ret, expectTxt)
					if idx < len(group.Children)-1 {
						ret = append(ret, "")
					}
					expectMap.Put(numbStr, expectTxt)
				} else {
					stepTxt := strings.TrimSpace(child.Desc)
					stepTxtWithNumb := replaceNumb(stepTxt, groupNumb, -1, false)
					stepMap.Put(numbStr, stepTxt)

					expectTxt := child.Expect
					expectTxt = strings.TrimSpace(expectTxt)
					expectMap.Put(numbStr, expectTxt)

					if expectTxt != "" {
						expectTxt = ">> " + expectTxt
					}

					ret = append(ret, fmt.Sprintf("  %s %s", stepTxtWithNumb, expectTxt))
				}

				groupNumb++
			}
		} else {
			desc = replaceNumb(group.Desc, groupNumb, -1, true)
			ret = append(ret, "\n"+desc)

			numbStr := getNumbStr(groupNumb, -1)
			stepMap.Put(numbStr, getGroupName(group.Desc))
			stepTypeMap.Put(numbStr, "group")
			expectMap.Put(numbStr, "")

			childNumb := 1
			for _, child := range group.Children {
				numbStr := getNumbStr(groupNumb, childNumb)
				stepTypeMap.Put(numbStr, "item")

				if group.MutiLine {
					// steps
					tag := replaceNumb("[steps]", groupNumb, childNumb, true)
					ret = append(ret, "  "+tag)

					stepTxt := printMutiStepOrExpect(child.Desc)
					ret = append(ret, stepTxt)
					stepMap.Put(numbStr, stepTxt)

					// expects
					tag = replaceNumb("[expects]", groupNumb, childNumb, true)
					ret = append(ret, "  "+tag)

					expectTxt := printMutiStepOrExpect(child.Expect)
					ret = append(ret, expectTxt)
					expectMap.Put(numbStr, expectTxt)
				} else {
					stepTxt := strings.TrimSpace(child.Desc)
					stepMap.Put(numbStr, stepTxt)

					expectTxt := child.Expect
					expectTxt = strings.TrimSpace(expectTxt)
					expectMap.Put(numbStr, expectTxt)

					if expectTxt != "" {
						expectTxt = ">> " + expectTxt
					}

					ret = append(ret, fmt.Sprintf("  %s %s", stepTxt, expectTxt))
				}

				childNumb++
			}

			groupNumb++
		}
	}

	return strings.Join(ret, "\n"), stepMap, stepTypeMap, expectMap
}

func getSortedTextFromNestedSteps(groups []model.TestStep) (ret string, stepMap, stepTypeMap, expectMap maps.Map) {
	arr := make([]string, 0)
	stepMap = linkedhashmap.New()
	stepTypeMap = linkedhashmap.New()
	expectMap = linkedhashmap.New()

	for idx1, group := range groups {
		numbStr := getNumbStr(idx1+1, -1)
		stepType := "step"
		if len(group.Children) > 0 {
			stepType = "group"
		}
		stepTypeMap.Put(numbStr, stepType)

		stepTxt := strings.TrimSpace(group.Desc)
		stepMap.Put(numbStr, stepTxt)

		expectTxt := strings.TrimSpace(group.Expect)
		expectTxt = strings.TrimRight(expectTxt, "]]")
		expectTxt = strings.TrimSpace(expectTxt)

		expectMap.Put(numbStr, expectTxt)

		if expectTxt != "" {
			expectTxt = ">> " + expectTxt
		}
		arr = append(arr, fmt.Sprintf("  %s %s", stepTxt, expectTxt))

		for idx2, child := range group.Children {
			numbStr := getNumbStr(idx1+1, idx2+1)
			stepTypeMap.Put(numbStr, "item")

			stepTxt := strings.TrimSpace(child.Desc)
			stepMap.Put(numbStr, stepTxt)

			expectTxt := strings.TrimSpace(child.Expect)
			expectMap.Put(numbStr, expectTxt)

			if expectTxt != "" {
				expectTxt = ">> " + expectTxt
			}

			arr = append(arr, fmt.Sprintf("  %s %s", stepTxt, expectTxt))
		}
	}

	ret = strings.Join(arr, "\n")
	return
}

func replaceNumb(str string, groupNumb int, childNumb int, withBrackets bool) string {
	numb := getNumbStr(groupNumb, childNumb)

	reg := `[\d\.\s]*(.*)`
	repl := numb + " ${1}"
	if withBrackets {
		reg = `\[` + reg + `\]`
		repl = `[` + repl + `]`
	}

	regx, _ := regexp.Compile(reg)
	str = regx.ReplaceAllString(str, repl)

	return str
}
func getNumbStr(groupNumb int, childNumb int) string {
	numb := strconv.Itoa(groupNumb) + "."
	if childNumb != -1 {
		numb += strconv.Itoa(childNumb) + "."
	}

	return numb
}
func getGroupName(str string) string {
	reg := `\[\d\.\s]*(.*)\]`
	repl := "${1}"

	regx, _ := regexp.Compile(reg)
	str = regx.ReplaceAllString(str, repl)

	return str
}

func printMutiStepOrExpect(str string) string {
	str = strings.TrimSpace(str)

	ret := make([]string, 0)

	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimSpace(line)

		ret = append(ret, fmt.Sprintf("%s%s", strings.Repeat(" ", 4), line))
	}

	return strings.Join(ret, "\r\n")
}

func GetExpectMapFromIndependentFileObsolete(expectMap maps.Map, content string, withEmptyExpect bool) maps.Map {
	retMap := linkedhashmap.New()

	expectArr := zentaoUtils.ReadExpectIndependentArrObsolete(content)

	idx := 0
	for _, keyIfs := range expectMap.Keys() {
		valueIfs, _ := expectMap.Get(keyIfs)

		key := strings.TrimSpace(keyIfs.(string))
		value := strings.TrimSpace(valueIfs.(string))

		if value == "" && len(expectArr) > idx {
			retMap.Put(key, strings.Join(expectArr[idx], "\r\n"))
			idx++
		} else {
			if withEmptyExpect {
				retMap.Put(key, "")
			}
		}
	}

	return retMap
}

func GetExpectMapFromIndependentFile(expectMap maps.Map, content string, withEmptyExpect bool) maps.Map {
	retMap := linkedhashmap.New()

	expectArr := zentaoUtils.ReadExpectIndependentArr(content)

	idx := 0
	for _, keyIfs := range expectMap.Keys() {
		valueIfs, _ := expectMap.Get(keyIfs)

		key := strings.TrimSpace(keyIfs.(string))
		value := strings.TrimSpace(valueIfs.(string))

		if value == "" && len(expectArr) > idx {
			retMap.Put(key, strings.Join(expectArr[idx], " | "))
			idx++
		} else {
			if withEmptyExpect {
				retMap.Put(key, "")
			}
		}
	}

	return retMap
}
