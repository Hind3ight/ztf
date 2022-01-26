package scriptUtils

import (
	i118Utils "github.com/aaronchen2k/deeptest/internal/command/utils/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/command/utils/script"
)

func Sort(cases []string) {
	for _, file := range cases {
		scriptUtils.SortFile(file)
	}

	logUtils.PrintTo(i118Utils.Sprintf("success_sort_steps", len(cases)))
}
