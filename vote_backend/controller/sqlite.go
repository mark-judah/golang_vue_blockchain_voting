package controller

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitSqlite() {
	database, err := sql.Open("sqlite3", "./nodeDB.db")
	if err != nil {
		panic(err)
	}
	statement, err2 := database.Prepare("CREATE TABLE IF NOT EXISTS transactions (transaction_id TEXT PRIMARY KEY, desktop_id TEXT, candidate_id TEXT, hash TEXT,created_by TEXT)")
	if err2 != nil {
		panic(err2)
	}
	statement.Exec()
}
