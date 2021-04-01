/*************************************************************************
# File Name: daily/daily.go
# Author: xiezg
# Mail: xzghyd2008@hotmail.com
# Created Time: 2021-02-23 21:41:57
# Last modified: 2021-04-01 11:04:30
************************************************************************/
package daily

import "fmt"
import "time"
import "daka/db"
import "encoding/json"

import "github.com/xiezg/glog"
import "github.com/xiezg/go-jsonify/jsonify"

func init() {

	today_zero := time.Now().Unix()
	today_zero -= today_zero % 86400 //取整为UTC零点
	today_zero -= 3600 * 8           //调整为北京时区，北京时区比UTC快8个小时

	time.AfterFunc(time.Duration(today_zero+86400-time.Now().Unix())*time.Second, TaskCreateToday)
}

func TaskCreateToday() {

	glog.Info("start create task")

	time.AfterFunc(24*time.Hour, TaskCreateToday)

	todayDate := time.Now().Format("2006-01-02")

	rows, err := db.MyDb.Query("select * from daily_task where task_date=?", todayDate)
	if err != nil {
		return
	}

	defer rows.Close()
	if rows.Next() {
		return
	}

	default_msg := `gantt
dateformat HH:mm
axisFormat %H:%M
title ` + todayDate + `
whole day : crit,done,09:30,22:00
section 上午
上午 : crit,done,am,10:00,12:00
section 中午
中午 : crit,done,rest,after am,120m
section 下午
下午 : crit,done, 14:00,22:00`

	if _, err := db.MyDb.Exec("insert into daily_task (task_msg,task_date)values(?,?)", default_msg, todayDate); err != nil {
		return
	}

	return
}

func TaskList(ctx interface{}, b []byte) (interface{}, error) {

	sql := "select * from daily_task order by task_date desc"

	var req map[string]int64

	if b != nil && json.Unmarshal(b, &req) == nil {
		sql = fmt.Sprintf("select * from daily_task where id=%v", req["task_id"])
	}

	rows, err := db.MyDb.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return jsonify.JsonifyMap(rows)
}

func TaskSet(ctx interface{}, b []byte) (interface{}, error) {

	var req map[string]interface{}

	if b == nil || json.Unmarshal(b, &req) != nil {
		return nil, fmt.Errorf("invalid request")
	}

	sql := "update daily_task set task_msg =? where id=?"

	if _, err := db.MyDb.Exec(sql, req["task_msg"], req["task_id"]); err != nil {
		return nil, err
	}

	return nil, nil
}
