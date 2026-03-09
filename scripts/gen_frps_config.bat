@echo off
setlocal EnableExtensions EnableDelayedExpansion

REM Generate frps.toml for Windows.
REM Usage:
REM   scripts\gen_frps_config.bat [output] [server_addr] [bind_port] [token] [dashboard]
REM                               [dashboard_addr] [dashboard_port] [dashboard_user] [dashboard_password]
REM Example:
REM   scripts\gen_frps_config.bat frps.toml 0.0.0.0 7000 my_token 1 0.0.0.0 7500 admin admin123

set "OUTPUT=%~1"
if "%OUTPUT%"=="" set "OUTPUT=frps.toml"

set "SERVER_ADDR=%~2"
if "%SERVER_ADDR%"=="" set "SERVER_ADDR=0.0.0.0"

set "BIND_PORT=%~3"
if "%BIND_PORT%"=="" set "BIND_PORT=7000"

set "TOKEN=%~4"

set "DASHBOARD=%~5"
if "%DASHBOARD%"=="" set "DASHBOARD=0"

set "DASHBOARD_ADDR=%~6"
if "%DASHBOARD_ADDR%"=="" set "DASHBOARD_ADDR=0.0.0.0"

set "DASHBOARD_PORT=%~7"
if "%DASHBOARD_PORT%"=="" set "DASHBOARD_PORT=7500"

set "DASHBOARD_USER=%~8"
if "%DASHBOARD_USER%"=="" set "DASHBOARD_USER=admin"

set "DASHBOARD_PASSWORD=%~9"

call :validate_port "%BIND_PORT%" || (
  echo bind_port must be 1-65535
  exit /b 1
)

if "%DASHBOARD%"=="1" (
  call :validate_port "%DASHBOARD_PORT%" || (
    echo dashboard_port must be 1-65535
    exit /b 1
  )
)

if "%TOKEN%"=="" (
  for /f %%i in ('powershell -NoProfile -Command "[guid]::NewGuid().ToString('N')"') do set "TOKEN=%%i"
)

if "%DASHBOARD%"=="1" if "%DASHBOARD_PASSWORD%"=="" (
  for /f %%i in ('powershell -NoProfile -Command "[guid]::NewGuid().ToString('N').Substring(0,20)"') do set "DASHBOARD_PASSWORD=%%i"
)

(
  echo bindAddr = "%SERVER_ADDR%"
  echo bindPort = %BIND_PORT%
  echo.
  echo [auth]
  echo method = "token"
  echo token = "%TOKEN%"
  if "%DASHBOARD%"=="1" (
    echo.
    echo [webServer]
    echo addr = "%DASHBOARD_ADDR%"
    echo port = %DASHBOARD_PORT%
    echo user = "%DASHBOARD_USER%"
    echo password = "%DASHBOARD_PASSWORD%"
  )
)> "%OUTPUT%"

echo Wrote: %OUTPUT%
echo token: %TOKEN%
if "%DASHBOARD%"=="1" (
  echo dashboard: http://%DASHBOARD_ADDR%:%DASHBOARD_PORT%
  echo dashboard user: %DASHBOARD_USER%
  echo dashboard password: %DASHBOARD_PASSWORD%
)
exit /b 0

:validate_port
set "P=%~1"
for /f "delims=0123456789" %%A in ("%P%") do exit /b 1
if "%P%"=="" exit /b 1
if %P% LSS 1 exit /b 1
if %P% GTR 65535 exit /b 1
exit /b 0
