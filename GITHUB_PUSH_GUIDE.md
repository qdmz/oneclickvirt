# GitHub推送指南

## 当前状态
✅ 代码修复完成
✅ Git提交完成
⏳ GitHub推送待完成（需要认证）

## 推送方法

### 方法1：使用Personal Access Token（推荐）

#### 1. 生成GitHub Personal Access Token
1. 访问：https://github.com/settings/tokens
2. 点击"Generate new token" → "Generate new token (classic)"
3. 设置token名称：`oneclickvirt-fix-push`
4. 选择权限：
   - ✅ `repo` (完整仓库访问权限)
   - ✅ `workflow` (工作流权限)
5. 点击"Generate token"
6. **重要**：复制生成的token（只显示一次）

#### 2. 配置Git使用Token
```bash
cd /root/.openclaw/workspace/oneclickvirt

# 方法A：直接在URL中包含token
git remote set-url origin https://YOUR_TOKEN@github.com/qdmz/oneclickvirt.git

# 方法B：使用Git凭据存储
git config --global credential.helper store
git push origin main
# 然后输入用户名和token（密码位置）
```

#### 3. 推送代码
```bash
git push origin main
```

### 方法2：使用SSH密钥

#### 1. 生成SSH密钥
```bash
ssh-keygen -t ed25519 -C "xiaomeili@openclaw.com"
# 按回车使用默认路径
# 可以设置密码或留空
```

#### 2. 查看公钥
```bash
cat ~/.ssh/id_ed25519.pub
```

#### 3. 添加SSH密钥到GitHub
1. 访问：https://github.com/settings/keys
2. 点击"New SSH key"
3. 粘贴公钥内容
4. 点击"Add SSH key"

#### 4. 测试SSH连接
```bash
ssh -T git@github.com
```

#### 5. 推送代码
```bash
cd /root/.openclaw/workspace/oneclickvirt
git remote set-url origin git@github.com:qdmz/oneclickvirt.git
git push origin main
```

### 方法3：使用GitHub CLI

#### 1. 安装GitHub CLI
```bash
# Linux
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update
sudo apt install gh
```

#### 2. 登录GitHub
```bash
gh auth login
# 选择：GitHub.com
# 选择：SSH
# 按照提示完成认证
```

#### 3. 推送代码
```bash
cd /root/.openclaw/workspace/oneclickvirt
git push origin main
```

## 推送验证

### 检查推送状态
```bash
cd /root/.openclaw/workspace/oneclickvirt
git status
git log --oneline -5
```

### 查看远程分支
```bash
git branch -r
git log origin/main --oneline -5
```

## 常见问题

### 问题1：认证失败
**错误**：`fatal: Authentication failed`
**解决**：
- 检查token是否正确
- 确认token有足够的权限
- 尝试重新生成token

### 问题2：权限不足
**错误**：`remote: Permission denied`
**解决**：
- 确认你是仓库的协作者
- 检查token权限设置
- 联系仓库管理员

### 问题3：推送冲突
**错误**：`! [rejected] main -> main (non-fast-forward)`
**解决**：
```bash
# 先拉取远程更新
git pull origin main --rebase

# 解决冲突后推送
git push origin main
```

## 推送后操作

### 1. 创建Pull Request（如果是fork）
```bash
# 如果是从fork推送，需要创建PR
gh pr create --title "fix: 修复MySQL事务隔离级别设置错误" --body "修复了创建实例时的数据库事务错误"
```

### 2. 验证GitHub上的提交
访问：https://github.com/qdmz/oneclickvirt/commits/main

### 3. 检查CI/CD状态
如果有GitHub Actions，检查构建状态。

## 自动化脚本

### 推送脚本
```bash
#!/bin/bash

# GitHub推送脚本
cd /root/.openclaw/workspace/oneclickvirt

echo "开始推送到GitHub..."

# 检查是否有未提交的更改
if [ -n "$(git status --porcelain)" ]; then
    echo "警告：有未提交的更改"
    git status
    exit 1
fi

# 推送到GitHub
echo "推送到origin/main..."
git push origin main

if [ $? -eq 0 ]; then
    echo "✅ 推送成功！"
    echo "查看提交：https://github.com/qdmz/oneclickvirt/commits/main"
else
    echo "❌ 推送失败，请检查认证信息"
    exit 1
fi
```

### 使用方法
```bash
chmod +x push_to_github.sh
./push_to_github.sh
```

## 安全建议

1. **不要在代码中硬编码token**
2. **定期轮换Personal Access Token**
3. **使用最小权限原则**
4. **启用双因素认证**
5. **定期审查已授权的应用**

## 后续维护

### 定期同步
```bash
# 拉取远程更新
git pull origin main

# 推送本地更改
git push origin main
```

### 分支管理
```bash
# 创建新分支
git checkout -b feature/new-feature

# 推送新分支
git push origin feature/new-feature
```

## 联系支持

如果遇到问题：
- GitHub文档：https://docs.github.com
- Git文档：https://git-scm.com/docs
- 项目Issues：https://github.com/qdmz/oneclickvirt/issues

---

**最后更新**：2026-04-24
**维护人员**：小美丽
