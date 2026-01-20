@echo off
REM Helper script to run Go commands with refreshed PATH

REM Refresh PATH from system and user environment variables
for /f "tokens=2*" %%a in ('reg query "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v Path 2^>nul') do set "SystemPath=%%b"
for /f "tokens=2*" %%a in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "UserPath=%%b"
set "PATH=%SystemPath%;%UserPath%"

REM Run the go command with all arguments
go %*
