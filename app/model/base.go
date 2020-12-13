package model

import (
	"auto-trading/app/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tableNameSignalEvents = "signal_events"
)

var DbConnection *sql.DB

// キャンドル用のテーブル名を作成して返す
func GetCandleTableName(productCode string, duration time.Duration) string {
	return fmt.Sprintf("%s_%s", productCode, duration)
}

func init() {
	// DBコネクションを行う
	// ※コネクションは随時行う為、errの定義が「err :=」だと、再定義出来ないのでpanicが起こる
	var err error
	DbConnection, err = sql.Open(config.Config.SQLDriver, config.Config.DbName)
	if err != nil {
		log.Fatalln("DBConnection Error", err)
	}
	// テーブル作成のsql文を作成
	sql := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		time DATETIME PRIMARY KEY NOT NULL,
		product_code STRING,
		side STRING,
		price FLOAT,
		size FLOAT)`, tableNameSignalEvents)
	// sqlを実行する
	DbConnection.Exec(sql)

	// 設定したデュレーションの数文のテーブルを作成する
	for _, duration := range config.Config.Durations {
		// 引数を%s_%sで繋げて名前を取得する
		tableName := GetCandleTableName(config.Config.ProductCode, duration)
		sql = fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS %s (
            time DATETIME PRIMARY KEY NOT NULL,
            open FLOAT,
            close FLOAT,
            high FLOAT,
            low FLOAT,
			volume FLOAT)`, tableName)
		DbConnection.Exec(sql)
	}
	fmt.Println("success.")
}
