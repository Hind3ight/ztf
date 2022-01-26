package cron

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/command/server/service"
	serverUtils "github.com/aaronchen2k/deeptest/internal/command/server/utils/common"
	serverConst "github.com/aaronchen2k/deeptest/internal/command/server/utils/const"
	cronUtils "github.com/aaronchen2k/deeptest/internal/command/server/utils/cron"
)

type CronService struct {
	heartBeatService *service.HeartBeatService

	buildService *service.BuildService
	taskService  *service.TaskService
	execService  *service.ExecService
}

func NewCronService(heartBeatService *service.HeartBeatService,
	buildService *service.BuildService, taskService *service.TaskService,
	execService *service.ExecService) *CronService {
	return &CronService{heartBeatService: heartBeatService,
		buildService: buildService, taskService: taskService, execService: execService}
}

func (s *CronService) Init() {
	cronUtils.AddTaskFuc(
		"HeartBeat",
		fmt.Sprintf("@every %ds", serverConst.HeartBeatInterval),
		func() { s.heartBeat() },
	)

	cronUtils.AddTaskFuc(
		"CheckRunning",
		fmt.Sprintf("@every %ds", serverConst.CheckUpgradeInterval),
		func() {
			if s.taskService.CheckRunning() { // ignore if task is running
				return
			}
		},
	)
}

func (s *CronService) heartBeat() {
	if serverUtils.IsVmAgent() { // vm
		// is running，register busy
		if s.taskService.CheckRunning() {
			s.heartBeatService.HeartBeat(true)
			return
		}

		// no task to run, register idle
		if s.taskService.GetSize() == 0 {
			s.heartBeatService.HeartBeat(false)
			return
		}

		// has task to run，register busy, then run
		build := s.taskService.Peek()
		s.heartBeatService.HeartBeat(true)

		// run
		s.taskService.Start()

		s.execService.Exec(build)

		s.taskService.Remove()
		s.taskService.End()
	}
}
