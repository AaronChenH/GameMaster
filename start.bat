@echo off
chcp 65001
echo 正在启动 Galaxy Empire Manager...

echo 检查并更新依赖...
go mod tidy

set DB_PATH=data\db
set MONGO_PATH=C:\Program Files\MongoDB\Server\7.0\bin\mongod.exe

echo 检查 MongoDB 服务...
tasklist | findstr mongod.exe > nul
if errorlevel 1 (
    echo MongoDB 服务未运行，正在启动...
    start "MongoDB" "%MONGO_PATH%" --dbpath="%DB_PATH%"
    ping 127.0.0.1 -n 4 > nul
)

echo 启动后台服务...
go run main.go

pause 