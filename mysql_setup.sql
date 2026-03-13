-- MySQL Setup Script for OneClickVirt

-- 1. Create database
CREATE DATABASE IF NOT EXISTS oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 2. Create user (if not exists)
-- Replace 'your_root_password' with your actual MySQL root password
CREATE USER IF NOT EXISTS 'oneclickvirt'@'%' IDENTIFIED BY '123456';
CREATE USER IF NOT EXISTS 'oneclickvirt'@'localhost' IDENTIFIED BY '123456';

-- 3. Grant all privileges
GRANT ALL PRIVILEGES ON oneclickvirt.* TO 'oneclickvirt'@'%';
GRANT ALL PRIVILEGES ON oneclickvirt.* TO 'oneclickvirt'@'localhost';

-- 4. Flush privileges
FLUSH PRIVILEGES;

-- 5. Show users and databases
SELECT User, Host FROM mysql.user WHERE User = 'oneclickvirt';
SHOW DATABASES LIKE 'oneclickvirt';
