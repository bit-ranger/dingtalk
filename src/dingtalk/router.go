package dingtalk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func execDingCommand(msg GitlabWebhookModel) string {
	kind := msg.Object_kind
	status := msg.Object_attributes.Detailed_status
	if kind == "pipeline" && status == "failed" {
		strFormat :=
			`
    Pipeline %s
`
		return fmt.Sprintf(strFormat, status)
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

	msg := execDingCommand(obj)
	if msg == "" {
		return
	}

	var dingToken = []string{token}
	cli := InitDingTalk(dingToken, ".")
	println("send msg %s, %s", token, msg)
	cli.SendTextMessage(msg)
	_, _ = w.Write([]byte(""))
}
