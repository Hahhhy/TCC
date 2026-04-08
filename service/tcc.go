package service

import "database/sql"

// 记录TCC阶段成功
func RecordTCCStage(tx *sql.Tx, txId, branchId, stage string) error {
	_, err := tx.Exec(
		"INSERT INTO tcc_transaction_log (tx_id, branch_id, stage, status) VALUES (?, ?, ?, 'SUCCESS')",
		txId, branchId, stage,
	)
	return err
}

// 检查某个阶段是否已经成功过（幂等）
func IsStageDone(tx *sql.Tx, txId, branchId, stage string) (bool, error) {
	var status string
	err := tx.QueryRow(
		"SELECT status FROM tcc_transaction_log WHERE tx_id = ? AND branch_id = ? AND stage = ? FOR UPDATE",
		txId, branchId, stage,
	).Scan(&status)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return status == "SUCCESS", nil
}
