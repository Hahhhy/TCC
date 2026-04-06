package main

import (
	"database/sql"
	"log"
	"net/http"

	"INFRO/handler"
	"INFRO/service"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接MySQL
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/tcc_demo?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(50)

	// 初始化服务
	svc := service.NewPurchaseService(db)
	handler.InitWorkerPool(10) // 10个worker

	// 注册路由
	http.HandleFunc("/purchase", handler.PurchaseHandler(svc))

	log.Println("server start at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
