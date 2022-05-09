
@echo off 

ass_svr.exe install

if "%errorlevel%" neq "0" (
	echo.
	echo Service install failed.
	echo.
	pause
) else (
	echo.
	echo Service install success.
	echo.
)


