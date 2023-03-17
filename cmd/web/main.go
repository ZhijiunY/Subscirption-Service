package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {

	// connect to the database
	db := initDB()
	db.Ping()

	// create sessions

	// able to login to the account

	// create channels

	// create waitgroup

	// set up the application config

	// set up mail

	// listen for web connections
}

// *sql.DB 返回資料庫
func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to database")
	}
	return conn
}

func connectToDB() *sql.DB {
	// 嘗試連接到資料庫固定次數，如果連接不到，就會死

	// 從0開始的計次
	counts := 0

	dsn := os.Getenv("DSN") // 設定環境變數，獲取dsn字串，來自os.Getenv調用環境變數

	for {
		// connection 和 err 來自調用尚不存在的開放資料庫
		connection, err := openDB(dsn)

		// 確認連接
		if err != nil {
			log.Println("postgres not yet ready...")
		} else {
			log.Print("connected to database!")
			return connection
		}

		// 如果遇到錯誤，再適十次
		if counts > 10 {
			return nil
		}

		// 否則
		log.Print("Backing off for 1 seconds")
		time.Sleep(1 * time.Second)
		// 增加 counts++
		counts++
		continue
	}

}

// openDB function
// 連接dsn 是一個 string, 返回 sql.DB 和 error
func openDB(dsn string) (*sql.DB, error) {
	// 確保多次嘗試連接到資料庫

	// 檢查錯誤時聲明一個變數DB。
	db, err := sql.Open("pgx", dsn)
	// 確認連接
	if err != nil {
		return nil, err
	}

	// 為了安全起見
	// db.Ping()檢查數據庫連接狀態的函數
	// 數據庫的連接可能會出現斷開、網絡問題等問題，
	// 因此在進行數據庫操作前，我們需要先確認數據庫是否正常連接。
	// 而這時候就可以使用 db.Ping() 方法來進行檢查。
	// 如果數據庫連接正常，該方法返回 nil；否則返回相應的錯誤信息，可以通過該錯誤信息進行相應的處理。
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}
