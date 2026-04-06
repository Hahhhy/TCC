package models

import "time"

type OrderStatus string

const (
	OrderInit      OrderStatus = "INIT"
	OrderTry       OrderStatus = "TRY"
	OrderConfirmed OrderStatus = "CONFIRMED"
	OrderCancelled OrderStatus = "CANCELLED"
)

type CourseOrder struct {
	Id          int64
	OrderNo     string
	UserId      int64
	CourseId    int64
	Price       float64
	Status      OrderStatus
	TryExpireAt *time.Time
}

type UserWallet struct {
	UserId  int64
	Balance float64
	Frozen  float64
	Version int64
}

type OwnerAccount struct {
	AccountId string
	Balance   float64
	Version   int64
}
