-- データベース作成（存在しなければ）
CREATE DATABASE IF NOT EXISTS saikyo_db
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_0900_ai_ci;

-- DBを選択
USE saikyo_db;

-- アプリケーション用ユーザー作成（存在しなければ）
CREATE USER IF NOT EXISTS 'saikyo_user'@'%' IDENTIFIED BY 'saikyo_pass';

-- 権限付与
GRANT ALL PRIVILEGES ON saikyo_db.* TO 'saikyo_user'@'%';

-- 忘れずに権限を適用
FLUSH PRIVILEGES;
