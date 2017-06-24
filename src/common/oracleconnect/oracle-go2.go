package oracleconnect

import (
	"context"
	"database/sql"
	"time"
	//"log"

	"fmt"
	ora "gopkg.in/rana/ora.v4"
)

func main2() {

	db, err := sql.Open("ora", "fcjhost/FCJHOST@54.149.166.130:1521/ORCL")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	// Set timeout (Go 1.8)
	ctx, err1 := context.WithTimeout(context.Background(), 5*time.Second)
	if err1 != nil {
		fmt.Println(err.Error())
	}

	// Set prefetch count (Go 1.8)
	ctx = ora.WithStmtCfg(ctx, ora.Cfg().StmtCfg.SetPrefetchRowCount(50000))
	rows, _ := db.QueryContext(ctx, "SELECT * FROM user_objects")
	defer rows.Close()
}
