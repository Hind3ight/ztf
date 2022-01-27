package scriptUtils

import (
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/command/utils/script"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
)

func Sort(cases []string) {
	for _, file := range cases {
		scriptUtils.SortFile(file)
	}

	logUtils.PrintTo(i118Utils.Sprintf("success_sort_steps", len(cases)))
}
