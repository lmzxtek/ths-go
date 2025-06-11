<#
.SYNOPSIS
多平台 Go 项目构建脚本 (PowerShell 版本)

.DESCRIPTION
自动编译 Windows、Linux 和 macOS 的可执行文件
#>

# 配置参数
$ProjectName = "ths"
$OutputDir = "bin"
# $env:GIN_MODE="release"

# 目标平台列表
$Platforms = @(
    [PSCustomObject]@{ GOOS = "windows"; GOARCH = "amd64"; Ext = ".exe" }
    # [PSCustomObject]@{ GOOS = "linux";   GOARCH = "amd64"; Ext = "" },
    # [PSCustomObject]@{ GOOS = "darwin";  GOARCH = "amd64"; Ext = "" },
    # [PSCustomObject]@{ GOOS = "darwin";  GOARCH = "arm64"; Ext = "" }
)

# 清理旧构建
if (Test-Path $OutputDir) {
    Remove-Item -Path $OutputDir -Recurse -Force
}
New-Item -ItemType Directory -Path $OutputDir | Out-Null


# 主构建循环
foreach ($Platform in $Platforms) {
    $env:GOOS = $Platform.GOOS
    $env:GOARCH = $Platform.GOARCH

    # 生成输出文件名
    $OutputFile = "${OutputDir}/${ProjectName}-$($env:GOOS)-$($env:GOARCH)$($Platform.Ext)"
    
    # 执行编译
    Write-Host "Building → $($env:GOOS)/$($env:GOARCH)" -ForegroundColor Cyan
    go build -ldflags="-s -w" -trimpath -o $OutputFile main.go

    # 编译结果检查
    if ($LASTEXITCODE -ne 0) {
        Write-Host "构建失败: $($env:GOOS)/$($env:GOARCH)" -ForegroundColor Red
        exit 1
    }
}

Write-Host "`n所有平台构建完成！" -ForegroundColor Green
Write-Host "输出目录: $((Resolve-Path $OutputDir).Path)`n" -ForegroundColor Yellow