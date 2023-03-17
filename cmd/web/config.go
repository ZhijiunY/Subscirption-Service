package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
)

// 用於app上的接收器
type Config struct {
	Session  *scs.SessionManager // Session管理
	DB       *sql.DB             // 資料庫
	InfoLog  *log.Logger         // 資訊日誌，將內容寫入主控台或日誌欓
	ErrorLog *log.Logger         // 錯誤日誌記錄器
	Wait     *sync.WaitGroup
}
