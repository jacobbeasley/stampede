@echo off
rem Automatically clean migrations/schema.sql before running buffalo test
if "%1"=="test" (
    go run "%~dp0scripts\clean_schema.go"
)
rem Call the real buffalo CLI. We must avoid calling this batch file recursively.
rem We search for the real buffalo.exe in the PATH.
setlocal enabledelayedexpansion
set "BUFFALO_EXE="
for /f "tokens=*" %%i in ('where buffalo.exe 2^>nul') do (
    set "TEMP_PATH=%%~dpki"
    rem Check if the found executable is in this script's directory
    if not "!TEMP_PATH!"=="%~dp0" (
        if "!BUFFALO_EXE!"=="" (
            set "BUFFALO_EXE=%%i"
        )
    )
)
if not "!BUFFALO_EXE!"=="" (
    "!BUFFALO_EXE!" %*
) else (
    rem Fallback if where didn't find any other buffalo
    buffalo %*
)
endlocal
