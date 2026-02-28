# PowerShell构建脚本 - 用于Windows
# 构建前端并嵌入到Go二进制文件中

Write-Host "开始构建..." -ForegroundColor Green

# 检查Node.js是否安装
if (-not (Get-Command node -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到Node.js，请先安装Node.js" -ForegroundColor Red
    exit 1
}

# 检查Go是否安装
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到Go，请先安装Go" -ForegroundColor Red
    exit 1
}

# 进入前端目录
Set-Location front

# 安装依赖（如果需要）
if (-not (Test-Path "node_modules")) {
    Write-Host "安装前端依赖..." -ForegroundColor Yellow
    if (Get-Command pnpm -ErrorAction SilentlyContinue) {
        pnpm install
    } elseif (Get-Command npm -ErrorAction SilentlyContinue) {
        npm install
    } else {
        Write-Host "错误: 未找到pnpm或npm" -ForegroundColor Red
        exit 1
    }
}

# 构建前端
Write-Host "构建前端..." -ForegroundColor Yellow
if (Get-Command pnpm -ErrorAction SilentlyContinue) {
    pnpm build
} else {
    npm run build
}

if ($LASTEXITCODE -ne 0) {
    Write-Host "前端构建失败" -ForegroundColor Red
    exit 1
}

# 返回项目根目录
Set-Location ..

# 检查dist目录是否存在
if (-not (Test-Path "front/dist")) {
    Write-Host "错误: 前端构建输出目录 front/dist 不存在" -ForegroundColor Red
    exit 1
}

# 将dist目录复制到ginServer/static/dist（embed需要）
Write-Host "复制前端构建文件到embed目录..." -ForegroundColor Yellow
if (Test-Path "ginServer/static/dist") {
    Remove-Item -Recurse -Force "ginServer/static/dist"
}
Copy-Item -Recurse "front/dist" "ginServer/static/dist"

# 构建Go二进制文件
Write-Host "构建Go二进制文件..." -ForegroundColor Yellow
$buildTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
$gitCommit = git rev-parse --short HEAD 2>$null
if (-not $gitCommit) {
    $gitCommit = "unknown"
}

go build -ldflags "-X 'main.BuildTime=$buildTime' -X 'main.GitCommit=$gitCommit'" -o trashcan.exe main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "Go构建失败" -ForegroundColor Red
    exit 1
}

Write-Host "构建完成！" -ForegroundColor Green
Write-Host "可执行文件: trashcan.exe" -ForegroundColor Cyan

