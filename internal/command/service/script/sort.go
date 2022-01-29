package scriptUtils

import (
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/script"
)

func Sort(cases []string) {
	for _, file := range cases {
		scriptUtils.SortFile(file)
	}

	logUtils.PrintTo(i118Utils.Sprintf("success_sort_steps", len(cases)))
}
