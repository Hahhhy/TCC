package models

// 订单状态
type OrderStatus string

const (
    OrderStatusInit    OrderStatus = "INIT"      // 初始状态
    OrderStatusPaid    OrderStatus = "PAID"      // 已支付
    OrderStatusCancelled OrderStatus = "CANCELLED" // 已取消

)//各过程状态

// 资源购买信息
type CourseOrder struct {

	
    
}//价格，购买人等

// 用户钱包
type UserWallet struct {
   
}//余额，用户名

// 所有者账户
type OwnerAccount struct {
    
}
// 用户资源权限
type UserCoursePermission struct {
   
}