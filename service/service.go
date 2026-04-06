package service

import (
	"INFRO/dao"
	"INFRO/models"
	"context"
	"database/sql"
	"errors"
)

type PurchaseService struct {
	db *sql.DB
}

// 对外的接口一个构造PurchaseService实例的函数
func NewPurchaseService(db *sql.DB) *PurchaseService {
	return &PurchaseService{db: db}
}

// Try：预留资源
func (s *PurchaseService) Try(ctx context.Context, orderNo string, userId, courseId int64, price float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 防悬挂：检查是否已经 Cancel
	cancelDone, err := IsCancelDone(tx, orderNo, "main")
	if err != nil {
		return err
	}
	if cancelDone {
		return errors.New("transaction already cancelled, refuse try")
	}

	// 2. 幂等：检查 Try 是否已成功
	done, err := IsStageDone(tx, orderNo, "main", "TRY")
	if err != nil {
		return err
	}
	if done {
		return tx.Commit() // 已经 Try 过，直接返回
	}

	// 3. 获取钱包（带行锁）和版本号
	wallet, err := dao.GetWalletForUpdate(tx, userId)
	if err != nil {
		return err
	}
	if wallet.Balance < price {
		return errors.New("insufficient balance")
	}

	// 4. 冻结余额（乐观锁）
	ok, err := dao.FreezeBalance(tx, userId, price, wallet.Version)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("concurrency conflict, please retry")
	}

	// 5. 更新订单状态为 TRY
	err = dao.UpdateOrderStatus(tx, orderNo, models.OrderTry)
	if err != nil {
		return err
	}

	// 6. 记录 Try 成功日志
	err = RecordTCCStage(tx, orderNo, "main", "TRY")
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Confirm：正式提交
func (s *PurchaseService) Confirm(ctx context.Context, orderNo string, userId, courseId int64, price float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 幂等
	done, err := IsStageDone(tx, orderNo, "main", "CONFIRM")
	if err != nil {
		return err
	}
	if done {
		return tx.Commit()
	}

	// 检查 Try 是否成功
	tryDone, err := IsStageDone(tx, orderNo, "main", "TRY")
	if err != nil {
		return err
	}
	if !tryDone {
		return errors.New("try not done, cannot confirm")
	}

	// 正式扣款（冻结转实际）
	err = dao.ConfirmDeduct(tx, userId, price)
	if err != nil {
		return err
	}

	// 给所有者加钱
	err = dao.AddOwnerBalance(tx, price)
	if err != nil {
		return err
	}

	// 授权
	err = dao.GrantPermission(tx, userId, courseId)
	if err != nil {
		return err
	}

	// 更新订单状态
	err = dao.UpdateOrderStatus(tx, orderNo, models.OrderConfirmed)
	if err != nil {
		return err
	}

	// 记录 Confirm 日志
	err = RecordTCCStage(tx, orderNo, "main", "CONFIRM")
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Cancel：回滚
func (s *PurchaseService) Cancel(ctx context.Context, orderNo string, userId int64, price float64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 幂等
	done, err := IsStageDone(tx, orderNo, "main", "CANCEL")
	if err != nil {
		return err
	}
	if done {
		return tx.Commit()
	}

	// 空回滚：如果 Try 没成功，直接记录 Cancel 成功即可
	tryDone, _ := IsStageDone(tx, orderNo, "main", "TRY")
	if tryDone {
		// 释放冻结
		err = dao.CancelFreeze(tx, userId, price)
		if err != nil {
			return err
		}
	}

	// 更新订单状态
	err = dao.UpdateOrderStatus(tx, orderNo, models.OrderCancelled)
	if err != nil {
		return err
	}

	// 记录 Cancel 日志
	err = RecordTCCStage(tx, orderNo, "main", "CANCEL")
	if err != nil {
		return err
	}

	return tx.Commit()
}

// 完整购买流程：Try -> Confirm
func (s *PurchaseService) Purchase(ctx context.Context, orderNo string, userId, courseId int64, price float64) error {
	// 1. Try
	if err := s.Try(ctx, orderNo, userId, courseId, price); err != nil {
		// Try失败，不需要Cancel（没有预留任何资源）
		return err
	}
	// 2. Confirm（这里可以改成异步，但简单点同步）
	if err := s.Confirm(ctx, orderNo, userId, courseId, price); err != nil {
		// Confirm失败，触发Cancel
		_ = s.Cancel(ctx, orderNo, userId, price)
		return err
	}
	return nil
}
