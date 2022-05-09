
@echo off 

ass_svr.exe start

if "%errorlevel%" neq "0" (
	echo.
	echo Service start failed.
	echo.
	pause
) else (
	echo.
	echo Service start success.
	echo.
)