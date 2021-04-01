/*************************************************************************
 # File Name: note/note.go
 # Author: xiezg
 # Mail: xzghyd2008@hotmail.com 
 # Created Time: 2021-03-24 10:22:53
 # Last modified: 2021-04-01 10:05:11
************************************************************************/
package note

import "fmt"
import "daka/db"
import "encoding/json"
import "github.com/xiezg/glog"
import "github.com/xiezg/go-jsonify/jsonify"

func List(ctx interface{}, b []byte) (interface{}, error) {

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

	return jsonify.JsonifyMap(rows)
}

func Update(ctx interface{}, b []byte) (interface{}, error) {

    var req map[string]interface{}

    if err := json.Unmarshal(b, &req); err != nil{
        return nil, err
    }

    stmt, err := db.MyDb.Prepare( "update my_note set title=?, txt=?, key_word=?")
    if err != nil{
        glog.Errorf("prepare fails. err:%v", err)
        return nil, nil
    }

    defer stmt.Close()

    if _, err := stmt.Exec( req["title"], req["txt"], req["key_word"]); err != nil{
		glog.Errorf("update my_note fails. err:%v", err)
        return nil, nil
    }

    return nil, nil
}
