package dingtalk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var tokenMap = make(map[string]string)

func execDingCommand(msg GitlabWebhookModel, token string) string {

	kind := msg.Object_kind
	status := msg.Object_attributes.Status
	projectName := msg.Project.Name

	if kind == "pipeline" && status == "failed" {
		strFormat := "%s Pipeline Failed"
		return fmt.Sprintf(strFormat, projectName)
	}

	oldStatus, ok := tokenMap[token]
	if kind == "pipeline" && status == "success" && ok && oldStatus == "failed" {
		strFormat := "%s Pipeline Success"
		return fmt.Sprintf(strFormat, projectName)
	}
	return ""
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
	//b, err := json.Marshal(obj)

	msg := execDingCommand(obj, token)

	tokenMap[token] = obj.Object_attributes.Status

	if msg == "" {
		return
	}

	var dingToken = []string{token}
	cli := InitDingTalk(dingToken, ".")
	println("send msg %s, %s", token, msg)
	cli.SendTextMessage(msg)
	_, _ = w.Write([]byte(""))
}
