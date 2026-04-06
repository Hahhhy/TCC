package handler

import (
	"INFRO/service"
	"encoding/json"
	"net/http"
)

type PurchaseRequest struct {
	OrderNo  string  `json:"order_no"`
	UserId   int64   `json:"user_id"`
	CourseId int64   `json:"course_id"`
	Price    float64 `json:"price"`
}

type PurchaseResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

// 全局 Worker 池
var workerPool *WorkerPool

func InitWorkerPool(size int) {
	workerPool = NewWorkerPool(size)
	workerPool.Start()
}

func PurchaseHandler(s *service.PurchaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req PurchaseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(PurchaseResponse{Success: false, Msg: "bad request"})
			return
		}

		// 幂等性：快速检查订单是否已存在（数据库唯一键会保证）
		// 这里简单交给数据库处理，无需额外逻辑

		// 提交到 Worker 池异步处理
		workerPool.Submit(func() {
			err := s.Purchase(r.Context(), req.OrderNo, req.UserId, req.CourseId, req.Price)
			if err != nil {
				// 记录错误日志，实际生产可写入队列重试
				println("purchase failed:", err.Error())
			}
		})

		// 立即返回接受
		json.NewEncoder(w).Encode(PurchaseResponse{Success: true, Msg: "accepted"})
	}
}
