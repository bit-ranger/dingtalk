package dingding

import (
	"dingtalk/log"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var tokenMap = make(map[string]string)

func execDingCommand(msg GitlabWebhookModel, token string) *dingMap {
	kind := msg.Object_kind
	status := msg.Object_attributes.Status
	pipelineId := msg.Object_attributes.Id
	webUrl := msg.Project.WebUrl
	projectName := msg.Project.Name

	if kind == "pipeline" && status == "failed" {
		dm := DingMap()
		dm.Set(projectName, H2)
		dm.Set(fmt.Sprintf("任务: #%d", pipelineId), N)
		dm.Set(fmt.Sprintf("状态: $$%s$$", status), RED)
		dm.Set(fmt.Sprintf("地址: $$%s/-/pipelines/%d$$", webUrl, pipelineId), BLUE)
		return dm
	}

	oldStatus, ok := tokenMap[token]
	if kind == "pipeline" && status == "success" && ok && oldStatus == "failed" {
		dm := DingMap()
		dm.Set(projectName, H2)
		dm.Set(fmt.Sprintf("任务: #%d", pipelineId), N)
		dm.Set(fmt.Sprintf("状态: $$%s$$", status), GREEN)
		dm.Set(fmt.Sprintf("地址: $$%s/-/pipelines/%d$$", webUrl, pipelineId), BLUE)
		return dm
	}
	return nil
}

// MyHandler实现Handler接口
type GitlabWebhookHandler struct {
}

func (h *GitlabWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("access_token")

	body := r.Body
	var err error
	var buf []byte
	obj := GitlabWebhookModel{}

	buf, err = ioutil.ReadAll(body)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &obj)
	if err != nil {
		print("unmarshal err %s", err.Error())
		return
	}
	
	status := obj.Object_attributes.Status
	if status == "success" || status== "failed" {
		//只处理success和failed两种状态的请求
	} else {
		return
	}
	
	pipelineId := obj.Object_attributes.Id

	log.Logger.WithFields(logrus.Fields{
		"pipelineId": pipelineId,
		"token":  token,
		"status": tokenMap[token],
	}).Info("dingding pipline old status")

	dingMap := execDingCommand(obj, token)

	tokenMap[token] = obj.Object_attributes.Status
	
	log.Logger.WithFields(logrus.Fields{
		"pipelineId": pipelineId,
		"token":  token,
		"status": tokenMap[token],
	}).Info("dingding pipline new status")

	if dingMap == nil {
		return
	}

	var dingToken = []string{token}
	cli := InitDingTalk(dingToken, ".")

	log.Logger.WithFields(logrus.Fields{
		"pipelineId": pipelineId,
		"token":   token,
		"message": dingMap.l,
	}).Info("send msg...")

	// 发送钉钉消息
	err = cli.SendMarkDownMessageBySlice(tokenMap[token], dingMap.Slice())
	if err != nil {
		log.Logger.Error(err.Error())
	}
	_, _ = w.Write([]byte(""))
}
