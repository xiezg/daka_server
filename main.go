/*************************************************************************
# File Name: main.go
# Author: xiezg
# Mail: xzghyd2008@hotmail.com
# Created Time: 2020-03-08 10:45:57
# Last modified: 2020-03-22 20:43:59
************************************************************************/

package main

import "time"
import "daka/db"
import "net/http"
import "encoding/json"
import "github.com/gorilla/mux"
import "github.com/xiezg/muggle/auth"

func commit_action(b []byte) (interface{}, error) {

	msg := struct {
		ActionType int    `json:"ActionType"`
		CommitTime string `json:"CommitTime"`
		Remarks    string `json:"Remarks"`
	}{}

	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, err
	}

	err := db.TaskCommit(msg.ActionType, msg.CommitTime, msg.Remarks)

	return nil, err
}

func query_action_list(b []byte) (interface{}, error) {

	msg := struct {
		CurTime int64 `json:"time"`
	}{}

	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, err
	}

	return db.QueryActionList(time.Unix(msg.CurTime, 0))
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/daka/api/query_action", auth.Auth(query_action_list))
	r.HandleFunc("/daka/api/login", auth.Login(db.QueryAccountInfo, "/daka/index.html"))
	r.HandleFunc("/daka/api/commit_action", auth.Auth(commit_action))

	http.ListenAndServe("127.0.0.1:8081", r)
}
