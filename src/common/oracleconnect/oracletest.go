package oracleconnect

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
)

func main1() {
	db, err := sql.Open("oci8", "fcjhost/FCJHOST@54.149.166.130:1521/ORCL")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		return
	}

	rows, err := db.Query("select module from actb_daily_log")
	if err != nil {
		fmt.Println("Error fetching addition")
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var module string
		rows.Scan(&module)
		fmt.Printf("%s\n", module)
	}
}
