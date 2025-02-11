@echo off
setlocal enabledelayedexpansion

rem 设置日志文件
set LOG_FILE=init_data.log

rem 切换到脚本所在目录
cd /d "%~dp0"

rem 设置编码
chcp 65001 > nul

rem 记录开始时间
echo [%date% %time%] 开始初始化数据... > "%LOG_FILE%"

rem 设置MongoDB Shell路径
set MONGOSH_PATH=D:\Program Files\MongoDB\mongosh-2.3.9-win32-x64\bin\mongosh.exe

rem 检查MongoDB Shell程序
if not exist "%MONGOSH_PATH%" (
    echo [%date% %time%] 错误：找不到MongoDB Shell程序 >> %LOG_FILE%
    goto error
)

rem 检查数据库连接
echo [%date% %time%] 检查数据库连接... >> %LOG_FILE%
"%MONGOSH_PATH%" --eval "db.runCommand({ping:1})" >> %LOG_FILE% 2>&1
if !errorlevel! neq 0 (
    echo [%date% %time%] 错误：无法连接到数据库 >> %LOG_FILE%
    goto error
)

rem 初始化基础数据
echo [%date% %time%] 初始化基础数据... >> %LOG_FILE%
echo 当前脚本路径: %~dp0 >> %LOG_FILE%
set "JS_PATH=%~dp0init_data.js"
set "JS_PATH=%JS_PATH:\=/%"
echo 完整JS文件路径: %JS_PATH% >> %LOG_FILE%
"%MONGOSH_PATH%" --eval "load('%JS_PATH%')" >> %LOG_FILE% 2>&1
if !errorlevel! neq 0 (
    echo [%date% %time%] 错误：数据初始化失败 >> %LOG_FILE%
    goto error
)

rem 成功完成
echo [%date% %time%] 数据初始化完成 >> %LOG_FILE%
goto end

:error
echo [%date% %time%] 初始化过程中出现错误，请查看 %LOG_FILE% >> %LOG_FILE%
exit /b 1

:end
endlocal
pause 