package dao

import (
	"INFRO/models"
	"database/sql"
)

// 查询钱包（带乐观锁版本号）
func GetWalletForUpdate(tx *sql.Tx, userId int64) (*models.UserWallet, error) {
	var w models.UserWallet
	err := tx.QueryRow(
		"SELECT user_id, balance, frozen, version FROM user_wallet WHERE user_id = ? FOR UPDATE",
		userId,
	).Scan(&w.UserId, &w.Balance, &w.Frozen, &w.Version)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// Try阶段：冻结余额（乐观锁）
func FreezeBalance(tx *sql.Tx, userId int64, amount float64, version int64) (bool, error) {
	result, err := tx.Exec(
		"UPDATE user_wallet SET frozen = frozen + ?, version = version + 1 WHERE user_id = ? AND balance >= ? AND version = ?",
		amount, userId, amount, version,
	)
	if err != nil {
		return false, err
	}
	rows, _ := result.RowsAffected()
	return rows == 1, nil
}

// Confirm阶段：冻结转实际扣款
func ConfirmDeduct(tx *sql.Tx, userId int64, amount float64) error {
	_, err := tx.Exec(
		"UPDATE user_wallet SET balance = balance - frozen, frozen = 0 WHERE user_id = ? AND frozen >= ?",
		userId, amount,
	)
	return err
}

// Cancel阶段：释放冻结
func CancelFreeze(tx *sql.Tx, userId int64, amount float64) error {
	_, err := tx.Exec(
		"UPDATE user_wallet SET frozen = frozen - ? WHERE user_id = ? AND frozen >= ?",
		amount, userId, amount,
	)
	return err
}

// 增加所有者余额
func AddOwnerBalance(tx *sql.Tx, amount float64) error {
	_, err := tx.Exec(
		"UPDATE owner_account SET balance = balance + ?, version = version + 1 WHERE account_id = 'owner'",
		amount,
	)
	return err
}

// 授予权限（幂等）
func GrantPermission(tx *sql.Tx, userId, courseId int64) error {
	_, err := tx.Exec(
		"INSERT IGNORE INTO user_course_permission (user_id, course_id) VALUES (?, ?)",
		userId, courseId,
	)
	return err
}

// 更新订单状态
func UpdateOrderStatus(tx *sql.Tx, orderNo string, status models.OrderStatus) error {
	_, err := tx.Exec(
		"UPDATE course_order SET status = ? WHERE order_no = ?",
		status, orderNo,
	)
	return err
}

// 创建订单（INIT状态）
func CreateOrder(tx *sql.Tx, orderNo string, userId, courseId int64, price float64) error {
	_, err := tx.Exec(
		"INSERT INTO course_order (order_no, user_id, course_id, price, status) VALUES (?, ?, ?, ?, ?)",
		orderNo, userId, courseId, price, models.OrderInit,
	)
	return err
}
