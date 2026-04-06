## 问题一

### 实现普通的本地数据库事务购买逻辑
- models：订单状态是INIT、PAID、FAILED
- purchase过程：开启事务
```
// 开启事务
tx, err := s.db.Begin()
if err != nil {
    return err
}
// 确保事务最终被 Rollback 或 Commit
defer func() {
    if err != nil {
        tx.Rollback()
    }
}()
```
锁定用户钱包——检查余额——锁定平台账户——创建订单——扣除用户钱——增加平台账户余额——授予用户权限——更新订单状态——提交事务
- testing：