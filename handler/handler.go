package handler

import (
	"INFRO/service"
	"encoding/json"
	"net/http"
)

// 这是请求的DTO和响应的VO
// 购买消息：订单号、用户ID、课程ID、价格
type PurchaseRequest struct {
	OrderNo  string  `json:"order_no"`
	UserId   int64   `json:"user_id"`
	CourseId int64   `json:"course_id"`
	Price    float64 `json:"price"`
}

// 购买响应：成功与否和消息
type PurchaseResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

// 全局 Worker 池
var workerPool *WorkerPool

// 初始化 Worker 池，并且开启Worker 池
func InitWorkerPool(size int) {
	workerPool = NewWorkerPool(size)
	workerPool.Start()
}

// 接收service，返回http.HandlerFunc，就是接受到请求之后怎么实现的那个函数
func PurchaseHandler(s *service.PurchaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// 解析请求体，序列化，失败就直接返回w这个http.ResponseWriter，告诉客户端请求错误
		var req PurchaseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(PurchaseResponse{Success: false, Msg: "bad request"})
			return
		}

		// 幂等性：快速检查订单是否已存在（数据库唯一键会保证）
		// 这里简单交给数据库处理，无需额外逻辑
		//？？？？？？？？？？？

		// 提交到 Worker 池异步处理
		//如果序列化成功了，就把这个请求提交到Worker池里面去处理，真正的购买逻辑在service层实现
		workerPool.Submit(func() {
			err := s.Purchase(r.Context(), req.OrderNo, req.UserId, req.CourseId, req.Price)
			if err != nil {
				// 记录错误日志，实际生产可写入队列重试
				//和第四个一样
				println("purchase failed:", err.Error())
			}
		})

		// 立即返回接受
		//这个会等worker处理完之后才返回吗
		json.NewEncoder(w).Encode(PurchaseResponse{Success: true, Msg: "accepted"})
	}
}
