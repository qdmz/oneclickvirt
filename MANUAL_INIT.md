# Quick Manual Initialization Guide

## Problem Found
The API call to /api/v1/public/init returns 500 error because:
1. Database connection is not established when service starts
2. The API expects to connect to the database during initialization

## Easiest Solution: Manual Initialization via Web UI

Please open your browser and go to:
```
http://localhost:8080/#/init
```

Then follow these steps:

### Step 1: Database Configuration
- Type: mysql
- Host: localhost
- Port: 3306
- Database: oneclickvirt
- Username: root
- Password: (leave empty if no password, or enter your root password)
- Click "Test Connection"
- If successful, click "Next"

### Step 2: Admin Account
- Username: admin
- Password: admin123
- Email: admin@example.com

### Step 3: Default User Account
- Username: user
- Password: user123
- Email: user@example.com

### Step 4: Initialize
- Click "Initialize"
- Wait for success message
- You'll be redirected to login page

### Step 5: Login

**Admin Login:**
- URL: http://localhost:8080/#/admin/login
- Username: admin
- Password: admin123

**User Login:**
- URL: http://localhost:8080/#/login
- Username: user
- Password: user123

## After Login - Test New Features

### Admin Menu Check (Sidebar):
- [ ] Agent Management (代理商管理) - NEW
- [ ] KYC Management (实名管理) - NEW
- [ ] Domain Configuration (域名配置) - NEW
- [ ] Domain Management (域名管理) - NEW

### User Menu Check (Sidebar):
- [ ] Domain Management (域名管理)
- [ ] KYC (实名认证)
- [ ] Wallet (钱包)
- [ ] Orders (订单)

### Theme Toggle:
- Click the theme icon in top-right corner
- Switch between dark and light themes
- Check the new Indigo color scheme

## Why Manual is Better

The web UI initialization flow handles:
- Database connection testing
- Error display and retry
- Step-by-step guidance
- Visual feedback

The API method we tried has:
- Less error handling
- Harder to troubleshoot
- Requires exact database credentials

---

If manual initialization also fails, please tell me:
1. What error message do you see?
2. At which step does it fail?
3. What MySQL root password do you have?

Then I can help with the exact setup!
