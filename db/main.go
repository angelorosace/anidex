package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitializeDB() (*sql.DB, error) {

	dbUrl := buildDBUrl()

	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func buildDBUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("MYSQLUSER"), os.Getenv("MYSQLPASSWORD"), os.Getenv("MYSQLHOST"), os.Getenv("MYSQLPORT"), os.Getenv("MYSQLDATABASE"))
}
