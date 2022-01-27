package service

import (
	"encoding/json"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/command/server/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/command/utils/log"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
)

type BuildService struct {
	taskService *TaskService
}

func NewBuildService(taskService *TaskService) *BuildService {
	return &BuildService{taskService: taskService}
}

func (s *BuildService) Add(req domain.ReqData) (reply domain.OptResult) {
	build := domain.Build{}

	reqStr, _ := json.Marshal(req.Data)
	err := json.Unmarshal(reqStr, &build)
	if err != nil {
		logUtils.PrintTo(i118Utils.Sprintf("fail_parse_req", err))
		return
	}

	size := s.taskService.GetSize()
	if size == 0 {
		s.taskService.Add(build)
		logUtils.PrintTo(i118Utils.Sprintf("success_add_tak"))
		reply.Success("Success to add task.")
	} else {
		reply.Fail(fmt.Sprintf("Already has %d jobs to be done.", size))
	}

	return
}
