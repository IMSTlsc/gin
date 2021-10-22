package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	//auto_qc_violation_code_tab
	db, err := sql.Open("mysql", "root:88888888@tcp(127.0.0.1:3306)/sp")
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(15)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return
	}
	sql := fmt.Sprintf("select id,code_description from auto_qc_violation_code_tab") //查询出所有数据
	rows, err := db.Query(sql)
	defer rows.Close()

	var codes []*autoQcCode // 解析所有数据
	for rows.Next() {
		code := &autoQcCode{}
		rows.Scan(&code.id, &code.codeDes)
		codes = append(codes, code)
	}
	codeId2NewDes := make(map[int]string) //记录需要更新的数据id，以及更新的内容
	for _, code := range codes {
		var codeDes codeDescription
		json.Unmarshal([]byte(code.codeDes), &codeDes)
		if codeDes.Age > 12 { //筛选复合更新条件的数据
			newCodeDes := newCodeDescription{
				Name: codeDes.Name,
				Age:  codeDes.Age,
				ExtraM: extraMes{
					"this is extra",
				},
			}
			marMes, _ := json.Marshal(newCodeDes) // 解析、编码后保存到本地map
			codeId2NewDes[code.id] = string(marMes)
		}
	}
	for id, newDes := range codeId2NewDes { // 遍历本地map，更新数据库中数据
		sql = fmt.Sprintf("update auto_qc_violation_code_tab set code_description = '%v' where id = '%v'", newDes, id)
		fmt.Println(sql)
		effected, err := db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(effected)
	}

}

type autoQcCode struct {
	id      int
	codeDes string
}
type codeDescription struct {
	Name string
	Age  int
}

type newCodeDescription struct {
	Name   string
	Age    int
	ExtraM extraMes
}

type extraMes struct {
	Mes string
}

func insert(db sql.DB) {
	code1 := &codeDescription{
		Name: "this is a code",
		Age:  13,
	}
	marMessage, err := json.Marshal(&code1)
	if err != nil {
		log.Fatal(err)
	}
	sql := fmt.Sprintf("insert into auto_qc_violation_code_tab (code_name,code_description) values ('%s','%s')", "default code", string(marMessage))
	rows, err := db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(rows.RowsAffected())
}

func query(db sql.DB) {
	var (
		id   int
		Name string
		des  string
	)

	rows, err := db.Query("select id,code_name,code_description from auto_qc_violation_code_tab where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() { //遍历数据行
		err := rows.Scan(&id, &Name, &des) //获取每行中数据，赋值的类型应当和数据库中类型一致
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, Name, des)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
