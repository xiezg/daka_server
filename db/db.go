/*************************************************************************
# File Name: db/db.go
# Author: xiezg
# Mail: xzghyd2008@hotmail.com
# Created Time: 2020-03-08 11:29:42
# Last modified: 2020-04-02 07:47:39
************************************************************************/

package db

import "time"
import "flag"
import "strconv"
import "database/sql"
import "github.com/golang/glog"
import "github.com/xiezg/muggle/db"

var MyDb *sql.DB

const (
	DAY_SECONDS = 24 * 60 * 60
)

func init() {
	flag.Parse()
}

func init() {

	db, err := db.InitMysql("rm-2zeg3thu1693r3609.mysql.rds.aliyuncs.com", 3306, "daka", "Daka@123")

	if err != nil {
		glog.Fatalf("mysql open fails. err:%v", err)
	}

	db.Exec("use daka")

	glog.Info("database connect success")
	MyDb = db
}

func QueryAccountInfo(name string, pwd string) (int, error) {
	return 0, nil
}

//LEFT JOIN 关键字会从左表 (Persons) 那里返回所有的行，即使在右表 (Orders) 中没有匹配的行
func QueryActionList(date time.Time) ([]interface{}, error) {

	beginUnixTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).Unix()
	endUnixTime := beginUnixTime + DAY_SECONDS

	sql := `SELECT
    a.id,
    a.action,
    a.action_time,
    a.take_time,
    a.warning,
    b.commit_time,
    b.remarks,
    b.result 
FROM
    plana a
    LEFT JOIN task_status b ON a.id = b.action_type
    AND UNIX_TIMESTAMP( b.commit_time ) >= ` + strconv.Itoa(int(beginUnixTime)) + `
    AND UNIX_TIMESTAMP( b.commit_time ) < ` + strconv.Itoa(int(endUnixTime)) + `
ORDER BY
    a.action_time;`

	rows, err := MyDb.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []interface{}

	for rows.Next() {
		var id int
		var action string
		var action_time string
		var commit_time *string
		var take_time *string
		var warning *string
		var remarks *string
		var status *int

		if err := rows.Scan(&id, &action, &action_time, &take_time, &warning, &commit_time, &remarks, &status); err != nil {
			return nil, err
		}

		result = append(result, struct {
			Action_type int
			Action      string
			Action_time string
			Take_time   *string
			Warning     *string
			Commit_time *string
			Remarks     *string
			Status      *int
		}{Action_type: id, Action: action, Action_time: action_time, Take_time: take_time, Warning: warning, Commit_time: commit_time, Remarks: remarks, Status: status})
	}

	return result, nil
}

func TaskCommit( uid int, action_type int, commit_time string, remarks string) error {

	if commit_time == "" {
		commit_time = time.Now().Format( "2006-01-02T15:04:05" )
	}

	sql := "INSERT INTO task_status ( uid, action_type, commit_time, remarks, result ) VALUES ( ?,?,?,?,1 ) ON DUPLICATE KEY UPDATE remarks = VALUES(remarks)"

	if _, err := MyDb.Exec(sql, uid, action_type, commit_time, remarks ); err != nil {
		return err
	}

	return nil
}
