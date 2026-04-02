# Git Remote 設定指南

## 選項 1：建立新的 GitHub Repository

### 步驟 1：在 GitHub 建立 Repository

1. 前往 https://github.com/new
2. 填寫資訊：
   - Repository name: `cc-cli-go`
   - Description: `CC-CLI-Go - Claude Code CLI Implementation`
   - 選擇 Public 或 Private
   - **不要勾選** "Add a README file"（我們已經有了）
   - **不要勾選** "Add .gitignore"（我們已經有了）
3. 點擊 "Create repository"

### 步驟 2：複製 Repository URL

建立後，GitHub 會顯示 URL，例如：

- HTTPS: `https://github.com/your-username/cc-cli-go.git`
- SSH: `git@github.com:your-username/cc-cli-go.git`

### 步驟 3：設定 Remote 並推送

```bash
# 進入專案目錄
cd /Users/user-name/github/cc-cli-go

# 設定 remote（替換成您的 URL）
git remote add origin https://github.com/YOUR_USERNAME/cc-cli-go.git

# 驗證設定
git remote -v

# 推送到 GitHub
git push -u origin master
```

---

## 選項 2：使用現有的 GitHub Repository

如果您已有 GitHub repository：

```bash
# 直接設定 remote
git remote add origin YOUR_REPOSITORY_URL

# 推送
git push -u origin master
```

---

## HTTPS vs SSH

### HTTPS（推薦給初學者）

- URL 格式: `https://github.com/username/repo.git`
- 每次推送需要輸入 GitHub 帳號密碼
- 或設定 Personal Access Token

### SSH（推薦給進階使用者）

- URL 格式: `git@github.com:username/repo.git`
- 需要先設定 SSH key
- 設定後不用輸入密碼

---

## 常用指令

```bash
# 查看 remote
git remote -v

# 修改 remote URL
git remote set-url origin NEW_URL

# 刪除 remote
git remote remove origin

# 推送
git push -u origin master

# 之後的推送只需要
git push

# 拉取更新
git pull origin master
```

---

## 需要幫忙？

如果您：

1. 需要我幫您建立 GitHub repository（我可以用 gh CLI）
2. 已有 repository URL，需要我幫您設定 remote
3. 需要設定 SSH key

請告訴我您的 GitHub username 或 repository URL！
