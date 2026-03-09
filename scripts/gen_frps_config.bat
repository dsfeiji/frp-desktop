@echo off
setlocal EnableExtensions EnableDelayedExpansion

REM One-click FRP server installer for Windows.
REM It installs frps, writes C:\frp-server\frps.toml, creates service, and starts it.
REM
REM Usage:
REM   scripts\gen_frps_config.bat [frp_version] [bind_port] [install_dir]
REM
REM Example:
REM   scripts\gen_frps_config.bat 0.67.0 7000 C:\frp-server

set "FRP_VERSION=%~1"
if "%FRP_VERSION%"=="" set "FRP_VERSION=0.67.0"

set "BIND_PORT=%~2"
if "%BIND_PORT%"=="" set "BIND_PORT=7000"

set "INSTALL_DIR=%~3"
if "%INSTALL_DIR%"=="" set "INSTALL_DIR=C:\frp-server"

for /f "delims=0123456789" %%A in ("%BIND_PORT%") do (
  echo bind_port must be numeric
  exit /b 1
)
if %BIND_PORT% LSS 1 (
  echo bind_port must be 1-65535
  exit /b 1
)
if %BIND_PORT% GTR 65535 (
  echo bind_port must be 1-65535
  exit /b 1
)

net session >nul 2>&1
if %errorlevel% neq 0 (
  echo Please run this script as Administrator.
  exit /b 1
)

set "ARCH=amd64"
if /I "%PROCESSOR_ARCHITECTURE%"=="ARM64" set "ARCH=arm64"

set "PKG=frp_%FRP_VERSION%_windows_%ARCH%.zip"
set "URL=https://github.com/fatedier/frp/releases/download/v%FRP_VERSION%/%PKG%"

echo Downloading: %URL%
powershell -NoProfile -Command "Invoke-WebRequest -Uri '%URL%' -OutFile '%TEMP%\%PKG%'"
if errorlevel 1 (
  echo Download failed.
  exit /b 1
)

if exist "%TEMP%\frp_extract" rd /s /q "%TEMP%\frp_extract"
mkdir "%TEMP%\frp_extract"
powershell -NoProfile -Command "Expand-Archive -Path '%TEMP%\%PKG%' -DestinationPath '%TEMP%\frp_extract' -Force"
if errorlevel 1 (
  echo Extract failed.
  exit /b 1
)

set "EXTRACT_DIR=%TEMP%\frp_extract\frp_%FRP_VERSION%_windows_%ARCH%"
if not exist "%EXTRACT_DIR%\frps.exe" (
  echo frps.exe not found in package.
  exit /b 1
)

if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
copy /y "%EXTRACT_DIR%\frps.exe" "%INSTALL_DIR%\frps.exe" >nul
if errorlevel 1 (
  echo Failed to copy frps.exe
  exit /b 1
)

for /f %%i in ('powershell -NoProfile -Command "$c='ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'; -join (1..8 ^| ForEach-Object { $c[(Get-Random -Maximum $c.Length)] })"') do set "TOKEN=%%i"

(
  echo bindAddr = "0.0.0.0"
  echo bindPort = %BIND_PORT%
  echo.
  echo [auth]
  echo method = "token"
  echo token = "%TOKEN%"
)> "%INSTALL_DIR%\frps.toml"

sc query frps >nul 2>&1
if %errorlevel% equ 0 (
  sc stop frps >nul 2>&1
  sc delete frps >nul 2>&1
  timeout /t 2 >nul
)

sc create frps binPath= "\"%INSTALL_DIR%\frps.exe\" -c \"%INSTALL_DIR%\frps.toml\"" start= auto DisplayName= "FRP Server"
if errorlevel 1 (
  echo Failed to create Windows service.
  exit /b 1
)

sc start frps >nul 2>&1

echo ========================================
echo FRP server installed and started.
echo frps binary : %INSTALL_DIR%\frps.exe
echo config file : %INSTALL_DIR%\frps.toml
echo service     : frps
echo token       : %TOKEN%
echo port        : %BIND_PORT%
echo status      : sc query frps
echo ========================================

exit /b 0
