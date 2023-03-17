package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "8888"

func main() {

	// connect to the database
	db := initDB()

	// create sessions to connect to redis
	session := initSession()

	// create loggers
	// 訊息日誌，紀錄：日誌、時間
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// 錯誤日誌，紀錄： 日治、時間、位置
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// able to login to the account

	// create channels

	// create waitgroup
	wg := sync.WaitGroup{}

	// set up the application config
	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Wait:     &wg,
	}

	// set up mail

	// listen for web connections
	app.serve()
}

// ---------- ---------- ----------
// << connect to the database >>
// *sql.DB 返回資料庫
func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to database")
	}
	return conn
}

// from initDB
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

// from connectToDB
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

// ---------- ---------- ----------
// << create sessions to connect to redis >>

// Session
func initSession() *scs.SessionManager { // scs from "github.com/alexedwards/scs/v2"

	// set up session
	// 創建新的Session管理器的函數
	session := scs.New()

	// store 存儲庫
	session.Store = redisstore.New(initRedis())
	// session.Lifetime 可以用來控制會話的有效期限，以防止資源的過度浪費和系統的不穩定運行。
	session.Lifetime = 24 * time.Hour
	// 控制會話（session）的 Cookie 是否應該在用戶端的瀏覽器上持續存在，即使瀏覽器被關閉或重新啟動。
	// 希望用戶端的瀏覽器在關閉後能夠自動重新登錄，可以將 session.Cookie.Persist 設定為 true
	session.Cookie.Persist = true
	// 是一個枚舉值，表示同站點請求時瀏覽器應該如何處理 Cookie 的 SameSite 屬性。
	session.Cookie.SameSite = http.SameSiteDefaultMode
	// 是一個用於控制瀏覽器是否只在使用安全協議（HTTPS）時發送會話（session）的 Cookie 的屬性。
	session.Cookie.Secure = true

	return session
}

// Redis
func initRedis() *redis.Pool { // redis 已經從 docker compose 連線了

	// 創建變數，指向 redis.Pool
	// redis.Pool是一個用於管理Redis連接池的Go語言庫。
	// Redis是一個快速、高效的內存鍵值數據庫，
	// 許多應用程序都使用Redis來存儲和查詢數據。
	// 然而，每次與Redis建立連接都需要額外的開銷，包括網絡開銷和認證開銷。
	// 這導致了一些性能問題，因此通常需要對Redis連接進行池化以減少開銷。
	redisPool := &redis.Pool{
		// 設置MaxIdle參數時需要根據應用程序的實際情況進行調整，
		// 以確保系統在運行期間可以維持穩定的性能和資源使用率。
		MaxIdle: 10, // default 預設值
		// 使用Redis Pool庫時，Dial是一個重要的函數，用於創建和初始化Redis連接。
		// Dial函數需要指定Redis服務器的地址和端口，並返回一個Redis連接和錯誤信息。
		Dial: func() (redis.Conn, error) { // 撥打讀取伺服器
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}

	return redisPool
}

// ---------- ---------- ----------
// << set up the application config >>

// ---------- ---------- ----------
// << create loggers >>

// ---------- ---------- ----------
// << listen for web connections >>
func (app *Config) serve() {
	// start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting web server ...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
