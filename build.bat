@echo off
REM build.bat - 多平台自动化构建脚本（Windows 版本）

REM 清理旧构建并创建输出目录
rd /s /q bin 2>nul
mkdir bin

REM 定义目标平台列表
set PLATFORMS=windows/amd64,linux/amd64,darwin/amd64,darwin/arm64

REM 遍历所有平台
for %%p in (%PLATFORMS%) do (
    REM 分割平台参数
    for /f "tokens=1,2 delims=/" %%a in ("%%p") do (
        set GOOS=%%a
        set GOARCH=%%b

        REM 设置输出文件名
        set OUTPUT=bin/ths-%%a-%%b
        if "%%a" == "windows" set OUTPUT=bin/ths-%%a-%%b.exe

        REM 执行编译命令
        echo [BUILDING] %%a/%%b → %OUTPUT%
        set GOOS=%%a
        set GOARCH=%%b
        go build -o %OUTPUT% main.go
    )
)

echo ----------------------------------
echo 构建完成！文件保存在 bin 目录
pause