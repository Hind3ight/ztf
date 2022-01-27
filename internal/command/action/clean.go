package action

import (
	constant "github.com/aaronchen2k/deeptest/internal/command/utils/const"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	"github.com/aaronchen2k/deeptest/internal/command/utils/vari"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"os"
)

func Clean() {
	path := vari.ExeDir + constant.LogDir
	bak := path[:len(path)-1] + "-bak" + string(os.PathSeparator) + path[len(path):]

	os.RemoveAll(path)
	os.RemoveAll(bak)

	logUtils.PrintTo(i118Utils.Sprintf("success_to_clean_logs"))
}
