package service

import (
	"errors"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"strings"
)

type ProjectService struct {
	ProjectRepo *repo.ProjectRepo `inject:""`
}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

func (s *ProjectService) Paginate(req serverDomain.ProjectReqPaginate) (ret domain.PageData, err error) {

	ret, err = s.ProjectRepo.Paginate(req)

	if err != nil {
		return
	}

	return
}

func (s *ProjectService) FindById(id uint) (model.Project, error) {
	return s.ProjectRepo.FindById(id)
}
func (s *ProjectService) FindByPath(projectPath string) (po model.Project, err error) {
	return s.ProjectRepo.FindByPath(projectPath)
}

func (s *ProjectService) Create(project model.Project) (id uint, err error) {
	project.Path = strings.TrimSpace(project.Path)

	if !fileUtils.IsDir(project.Path) {
		err = errors.New(fmt.Sprintf("路径为%s不是目录。", project.Path))
		return
	}

	po, _ := s.ProjectRepo.FindByPath(fileUtils.AddPathSepIfNeeded(project.Path))
	if po.ID != 0 {
		err = errors.New(fmt.Sprintf("路径为%s的项目已存在。", project.Path))
		return
	}

	if project.Name == "" {
		project.Name = fileUtils.GetDirName(project.Path)
	}

	id, err = s.ProjectRepo.Create(project)
	return
}

func (s *ProjectService) Update(id uint, project model.Project) error {
	return s.ProjectRepo.Update(id, project)
}

func (s *ProjectService) DeleteByPath(pth string) (err error) {
	err = s.ProjectRepo.DeleteByPath(pth)
	if err != nil {
		return
	}

	err = s.ProjectRepo.SetCurrProject("")

	return
}

func (s *ProjectService) GetByUser(currProjectPath string) (
	projects []model.Project, currProject model.Project, currProjectConfig commDomain.ProjectConf, scriptTree serverDomain.TestAsset, err error) {
	projects, err = s.ProjectRepo.ListProjectByUser()

	found := false
	for _, p := range projects {
		if p.Path == currProjectPath {
			found = true
			break
		}
	}

	if !found {
		if err != nil {
			logUtils.Errorf("db operation error %s", err.Error())
			return
		}

		name := fileUtils.GetDirName(currProjectPath)
		newLocalProject := model.Project{Path: currProjectPath, Name: name}

		_, err = s.ProjectRepo.Create(newLocalProject)
		if err != nil {
			logUtils.Errorf("db operation error %s", err.Error())
			return
		}

		projects, err = s.ProjectRepo.ListProjectByUser()
	}

	s.ProjectRepo.SetCurrProject(currProjectPath)

	currProject, err = s.ProjectRepo.GetCurrProjectByUser()
	if err != nil {
		logUtils.Errorf("db operation error %s", err.Error())
		return
	}

	if currProject.Type == commConsts.TestFunc {
		scriptTree, err = scriptUtils.LoadScriptTree(currProject.Path)
	}

	currProjectConfig = configUtils.ReadFromFile(currProject.Path)
	currProjectConfig.IsWin = commonUtils.IsWin()

	return
}
