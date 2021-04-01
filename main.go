/*************************************************************************
# File Name: main.go
# Author: xiezg
# Mail: xzghyd2008@hotmail.com
# Created Time: 2020-03-08 10:45:57
# Last modified: 2021-03-25 11:11:59
************************************************************************/

package main

import "time"
import "net/http"
import "daka/db"
import "daka/daily"
import "daka/note"
import "daka/money"
import "encoding/json"
import "github.com/gorilla/mux"
import "github.com/xiezg/muggle/auth"
import sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"

func commit_action(ctx interface{}, b []byte) (interface{}, error) {

	uid := 1

	msg := struct {
		ActionType int    `json:"ActionType"`
		CommitTime string `json:"CommitTime"`
		Remarks    string `json:"Remarks"`
	}{}

	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, err
	}

	err := db.TaskCommit(uid, msg.ActionType, msg.CommitTime, msg.Remarks)

	return nil, err
}

func query_action_list(ctx interface{}, b []byte) (interface{}, error) {

	msg := struct {
		CurTime int64 `json:"time"`
	}{}

	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, err
	}

	return db.QueryActionList(time.Unix(msg.CurTime, 0))
}

func send_sms(actionName string, actionTime time.Time) error {

	SecretId := "AKIDMcdS3SAdOtHn49cb4KKBnDWbVAZjMbCe"
	SecretKey := "MVO2RhiU5UeYCbsugEWZMtMgkYBcolTv"
	phoneNum := "+8618710166030"
	templateId := "555667"
	appId := "1400342302"
	sign := "谢振国模板素材"
	param1 := actionName
	param2 := actionTime.Format("15:04:05")

	req := sms.NewSendSmsRequest()
	req.PhoneNumberSet = []*string{&phoneNum}
	req.TemplateID = &templateId
	req.SmsSdkAppid = &appId
	req.Sign = &sign
	req.TemplateParamSet = []*string{&param1, &param2}

	client, err := sms.NewClientWithSecretId(SecretId, SecretKey, "")
	if err != nil {
		return err
	}

	_, err = client.SendSms(req)

	return err
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/daka/api/note/list", auth.Auth(note.List))
	r.HandleFunc("/daka/api/note/update", auth.Auth(note.Update))

	r.HandleFunc("/daka/api/daily/task/list", auth.Auth(daily.TaskList))
	r.HandleFunc("/daka/api/daily/task/set", auth.Auth(daily.TaskSet))
	r.HandleFunc("/money/api/daily/expense/add", auth.Auth(money.DailyExpenseAdd))
	r.HandleFunc("/money/api/daily/expense/list", auth.Auth(money.DailyExpenseList))
	r.HandleFunc("/daka/api/query_action", auth.Auth(query_action_list))
	r.HandleFunc("/daka/api/login", auth.Login(db.QueryAccountInfo, ""))
	r.HandleFunc("/daka/api/commit_action", auth.Auth(commit_action))

	http.ListenAndServe("127.0.0.1:8081", r)
}
