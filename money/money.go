/*************************************************************************
# File Name: money/money.go
# Author: xiezg
# Mail: xzghyd2008@hotmail.com
# Created Time: 2021-01-31 15:09:39
# Last modified: 2021-02-23 21:50:29
************************************************************************/
package money

import "daka/db"
import "github.com/xiezg/glog"
import "github.com/xiezg/go-jsonify/jsonify"

type expense_item_info struct {
	Time   int `json:"time"`
	Price  int `json:"price"`
	Number int `json:"number"`
	Name   int `json:"name"`
}

func init() {

	create_table_sql := "create table if not exists  expense_list( id int auto_increment primary key, name char(255) not null, number float not null, price float not null, time timestamp default CURRENT_TIMESTAMP );"

	if _, err := db.MyDb.Exec(create_table_sql); err != nil {
		glog.Fatalf("create table fails. err:%v", err)
	}
}

func DailyExpenseList(ctx interface{}, b []byte) (interface{}, error) {

	sql := "select name, number, price, time from expense_list order by time"

	rows, err := db.MyDb.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return jsonify.JsonifyMap(rows)
}

func DailyExpenseAdd(ctx interface{}, b []byte) (interface{}, error) {
	return nil, nil
}

//func ExpenseAdd( b []byte)(interface[], error){
//
//    create table if not exists  expense_list( id int auto_increment primary key, name char(255) not null, number float not null, price float not null, time timestamp default CURRENT_TIMESTAMP );
//
//    var info expense_item_info
//    if err := json.Unmarshal(b, &info); err != nil{
//        return nil, err
//    }
//
//    return nil, nil
//}
//
//func Temp_charge_tcard_rules_abnormal(communityId string, obj interface{}) error {
//
//	mm := obj.(map[string]interface{})
//
//	if _, err := MyDb.Exec("insert into `temp_charge_tcard_rules` ( "+
//		"`uniq_id`, "+
//		"`exception_price` )"+
//		" values( 1, ? ) on DUPLICATE key update "+
//		"`exception_price`=values(`exception_price`)",
//		mm["price"]); err != nil {
//		return err
//	}
//
//	return nil
//}
