package mydb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Job struct {
	ID              int
	Name, Time, Cmd string
	STime, ETime    int
	Status          uint8
	Running         uint8
}

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", "root", "wida", "127.0.0.1", 3306, "mycron"))
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)
	db.Ping()
}

func GetCronList() (jobss []Job, e error) {
	ut := int64(time.Now().Unix())
	rows, err := db.Query("SELECT * FROM cron where status = 1 and sTime < ? and eTime > ?", ut, ut)
	if err != nil {
		panic(err.Error())
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))
	jobs := make([]Job, len(values))
	i := 0
	// Fetch rows
	for rows.Next() {
		err = rows.Scan(&jobs[i].ID, &jobs[i].Name, &jobs[i].Time, &jobs[i].Cmd, &jobs[i].STime, &jobs[i].ETime, &jobs[i].Status, &jobs[i].Running)
		if err != nil {
			panic(err.Error())
		}
		i++
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
	return jobs, nil
}

func (job Job) ChangeRunningStatu(status int) (int64, error) {
	stmtIns, err := db.Prepare("update cron set isrunning = ? where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(status, job.ID)
	if err != nil {
		panic(err.Error())
	}
	return result.RowsAffected()
}
