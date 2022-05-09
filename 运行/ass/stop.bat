
@echo off 

ass_svr.exe stop

if "%errorlevel%" neq "0" (
	echo.
	echo Service stop failed.
	echo.
	pause
) else (
	echo.
	echo Service stop success.
	echo.
)