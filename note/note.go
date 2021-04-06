/*************************************************************************
 # File Name: note/note.go
 # Author: xiezg
 # Mail: xzghyd2008@hotmail.com
 # Created Time: 2021-03-24 10:22:53
 # Last modified: 2021-04-06 14:13:09
************************************************************************/
package note

import "fmt"
import "daka/db"
import "encoding/json"
import "github.com/xiezg/glog"
import "github.com/gomarkdown/markdown"
import "github.com/xiezg/go-jsonify/jsonify"

func List(ctx interface{}, b []byte) (rst interface{}, err error) {

	sql := "select id,title,txt,key_word, create_time, last_update from my_note order by create_time desc"

	var req map[string]int64

	if b != nil && json.Unmarshal(b, &req) == nil {
		sql = fmt.Sprintf("select * from my_note where id=%v", req["note_id"])
	}

	rows, err := db.MyDb.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	defer func() {

		for _, item := range rst.([]map[string]interface{}) {
			item["html"] = string(markdown.ToHTML([]byte(item["txt"].(string)), nil, nil))
		}
	}()
	return jsonify.JsonifyMap(rows)
}

func Update(ctx interface{}, b []byte) (interface{}, error) {

	var req map[string]interface{}

	if err := json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	args_str := ""
	args := make([]interface{}, 0)

	if req["title"] != nil {
		args = append(args, req["title"])

		if len(args_str) != 0 {
			args_str = args_str + ","
		}
		args_str = args_str + " title=?"
	}

	if req["txt"] != nil {
		args = append(args, req["txt"])

		if len(args_str) != 0 {
			args_str = args_str + ","
		}
		args_str += " txt=?"
	}

	if req["key_word"] != nil {
		args = append(args, req["key_word"])

		if len(args_str) != 0 {
			args_str = args_str + ","
		}
		args_str += " key_word=?"
	}

	args = append(args, req["id"])
	args_str += " where id=?"

	stmt, err := db.MyDb.Prepare("update my_note set" + args_str)
	if err != nil {
		glog.Errorf("prepare fails. err:%v", err)
		return nil, nil
	}

	defer stmt.Close()

	if _, err := stmt.Exec(args...); err != nil {
		glog.Errorf("update my_note fails. err:%v", err)
		return nil, nil
	}

	return nil, nil
}
