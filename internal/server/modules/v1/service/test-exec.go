package service

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	analysisUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/analysis"
	"github.com/jinzhu/copier"
	"path/filepath"
)

type TestExecService struct {
	TestExecRepo *repo.TestExecRepo `inject:""`
	ProjectRepo  *repo.ProjectRepo  `inject:""`
}

func NewTestExecService() *TestExecService {
	return &TestExecService{}
}

func (s *TestExecService) List(projectPath string) (ret []serverDomain.TestReportSummary, err error) {
	reportFiles := analysisUtils.ListReport(projectPath)

	for _, seq := range reportFiles {
		var summary serverDomain.TestReportSummary

		report, err1 := analysisUtils.GetReport(projectPath, seq)
		if err1 != nil { // ignore wrong json result
			continue
		}
		copier.Copy(&summary, report)

		summary.Seq = seq
		ret = append(ret, summary)
	}

	return
}

func (s *TestExecService) Get(projectPath string, seq string) (report commDomain.ZtfReport, err error) {
	return analysisUtils.GetReport(projectPath, seq)
}

func (s *TestExecService) Delete(projectPath string, seq string) (err error) {
	dir := filepath.Join(projectPath, commConsts.LogDirName)

	di := filepath.Join(dir, seq)
	err = fileUtils.RmDir(di)

	return
}
