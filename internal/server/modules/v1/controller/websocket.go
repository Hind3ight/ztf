package controller

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	execHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/exec"
	websocketUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/websocket"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
)

var (
	ch chan int
)

type WebSocketCtrl struct {
	Namespace         string
	WorkspaceService  *service.WorkspaceService `inject:""`
	*websocket.NSConn `stateless:"true"`
}

func NewWebSocketCtrl() *WebSocketCtrl {
	inst := &WebSocketCtrl{Namespace: serverConfig.WsDefaultNameSpace}
	return inst
}

func (c *WebSocketCtrl) OnNamespaceConnected(msg websocket.Message) error {
	websocketUtils.SetConn(c.Conn)

	logUtils.Infof(i118Utils.Sprintf("ws_namespace_connected", c.Conn.ID(), msg.Room))

	data := serverDomain.WsResp{Msg: "from server: connected to websocket"}
	websocketUtils.Broadcast(msg.Namespace, "", "OnVisit", data)
	return nil
}

// OnNamespaceDisconnect
// This will call the "OnVisit" event on all clients, except the current one,
// it can't because it's left but for any case use this type of design.
func (c *WebSocketCtrl) OnNamespaceDisconnect(msg websocket.Message) error {
	logUtils.Infof(i118Utils.Sprintf("ws_namespace_disconnected", c.Conn.ID()))

	data := serverDomain.WsResp{Msg: fmt.Sprintf("ws_connected")}
	websocketUtils.Broadcast(msg.Namespace, "", "OnVisit", data)
	return nil
}

// OnChat This will call the "OnVisit" event on all clients,
// including the current one, with the 'newCount' variable.
func (c *WebSocketCtrl) OnChat(wsMsg websocket.Message) (err error) {
	ctx := websocket.GetContext(c.Conn)
	logUtils.Infof(i118Utils.Sprintf("ws_onchat", ctx.RemoteAddr(), wsMsg.Room, string(wsMsg.Body)))

	req := serverDomain.WsReq{}
	err = json.Unmarshal(wsMsg.Body, &req)
	if err != nil {
		msg := i118Utils.Sprintf("wrong_req_params", err.Error())
		websocketUtils.SendExecMsg(msg, "", commConsts.Error, &wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)
		return
	}

	act := req.Act

	if act == commConsts.ExecInit {
		msg := i118Utils.Sprintf("success_to_conn")
		//websocketUtils.SendExecMsg(msg, strconv.FormatBool(execHelper.GetRunning()), wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)
		return
	}

	if act == commConsts.ExecStop {
		if ch != nil {
			if !execHelper.GetRunning() {
				ch = nil
			} else {
				ch <- 1
				ch = nil
			}
		}

		execHelper.SetRunning(false)

		msg := i118Utils.Sprintf("end_task")
		websocketUtils.SendExecMsg(msg, "false", commConsts.Run, &wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)
		return
	}

	if execHelper.GetRunning() && (act == commConsts.ExecCase || act == commConsts.ExecModule ||
		act == commConsts.ExecSuite || act == commConsts.ExecTask || act == commConsts.ExecUnit) {
		msg := i118Utils.Sprintf("pls_stop_previous")
		websocketUtils.SendExecMsg(msg, "true", commConsts.Run, &wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)

		return
	}

	// populate test set's props with parent
	execHelper.PopulateTestSetProps(&req)
	for idx, _ := range req.TestSets {
		testSet := &req.TestSets[idx]

		if testSet.WorkspaceId != 0 {
			po, _ := c.WorkspaceService.Get(uint(testSet.WorkspaceId))
			testSet.WorkspacePath = po.Path
		}
	}

	ch = make(chan int, 1)
	go func() {
		execHelper.Exec(ch, req, &wsMsg)
		execHelper.SetRunning(false)
	}()

	execHelper.SetRunning(true)

	msg := i118Utils.Sprintf("start_task")
	websocketUtils.SendExecMsg(msg, "true", commConsts.Run, &wsMsg)
	logUtils.ExecConsole(color.FgCyan, msg)

	return
}
