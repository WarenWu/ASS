
@echo off 

ass_svr.exe uninstall

if "%errorlevel%" neq "0" (
	echo.
	echo Service uninstall failed.
	echo.
	pause
) else (
	echo.
	echo Service uninstall success.
	echo.
)