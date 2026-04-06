-- 订单表
CREATE TABLE course_order (
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_no       VARCHAR(64) NOT NULL UNIQUE,
    user_id        BIGINT NOT NULL,
    course_id      BIGINT NOT NULL,
    price          DECIMAL(10,2) NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'INIT',  -- INIT / TRY / CONFIRMED / CANCELLED
    try_expire_at  DATETIME,
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 用户钱包
CREATE TABLE user_wallet (
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT UNIQUE NOT NULL,
    balance    DECIMAL(10,2) NOT NULL DEFAULT 0,
    frozen     DECIMAL(10,2) NOT NULL DEFAULT 0,
    version    BIGINT NOT NULL DEFAULT 0,          -- 乐观锁
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 所有者账户（简单起见就一条记录）
CREATE TABLE owner_account (
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    account_id VARCHAR(64) UNIQUE NOT NULL,
    balance    DECIMAL(10,2) NOT NULL DEFAULT 0,
    version    BIGINT NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 用户权限
CREATE TABLE user_course_permission (
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    course_id  BIGINT NOT NULL,
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_course (user_id, course_id)
);

-- TCC 事务日志（防悬挂、幂等、空回滚）
CREATE TABLE tcc_transaction_log (
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    tx_id      VARCHAR(64) NOT NULL,        -- 全局事务ID（订单号即可）
    branch_id  VARCHAR(64) NOT NULL,        -- 分支ID（固定"wallet"或"order"等，简单统一用"main"）
    stage      VARCHAR(20) NOT NULL,        -- TRY / CONFIRM / CANCEL
    status     VARCHAR(20) NOT NULL,        -- SUCCESS / FAILED
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_tx_branch_stage (tx_id, branch_id, stage)
);