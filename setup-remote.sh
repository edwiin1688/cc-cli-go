#!/bin/bash
# Git Remote 設定腳本
# 請將下方的 URL 替換成您的 GitHub repository URL

# HTTPS 方式（推薦給初學者）
git remote add origin https://github.com/YOUR_USERNAME/cc-cli-go.git

# 或是使用 SSH 方式（需要設定 SSH key）
# git remote add origin git@github.com:YOUR_USERNAME/cc-cli-go.git

# 驗證 remote 設定
git remote -v

# 推送到 GitHub
git push -u origin master