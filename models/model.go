package models

import "time"

//订单状态
type OrderStatus string

const (
	OrderInit      OrderStatus = "INIT"
	OrderTry       OrderStatus = "TRY"
	OrderConfirmed OrderStatus = "CONFIRMED"
	OrderCancelled OrderStatus = "CANCELLED"
)

//资源购买信息：订单号、用户ID、课程ID、价格、状态、冻结时间（冻结是发生在Try阶段的，所以放在订单里）
type CourseOrder struct {
	Id          int64
	OrderNo     string
	UserId      int64
	CourseId    int64
	Price       float64
	Status      OrderStatus
	TryExpireAt *time.Time
}

//用户钱包信息：用户ID、余额、冻结金额、版本号（乐观锁）
//乐观锁就是：每次更新时检查版本号，只有版本号匹配才更新成功，否则说明数据被修改过，需要重试。
type UserWallet struct {
	UserId  int64
	Balance float64
	Frozen  float64
	Version int64
}

//平台方账户信息：账户ID、余额、版本号（乐观锁）
type OwnerAccount struct {
	AccountId string
	Balance   float64
	Version   int64
}

//用户资源权限
type UserCoursePermission struct {
	UserId    int64
	CourseId  int64
	Status    int
	Version   int64
	GrantedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt *time.Time
}
