@echo off
REM Batch file to extract Canvus database schema
REM Usage: extract-schema.bat

echo 🔍 Extracting Canvus Database Schema
echo =====================================

REM Check if PowerShell is available
powershell -Command "Get-Command powershell" >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ PowerShell not found. Please install PowerShell or run the .ps1 file directly.
    pause
    exit /b 1
)

REM Run the PowerShell script
echo 📄 Running PowerShell schema extraction script...
powershell -ExecutionPolicy Bypass -File "%~dp0extract-schema.ps1"

if %errorlevel% equ 0 (
    echo.
    echo ✅ Schema extraction completed successfully!
    echo 📄 Check the output file: canvus-schema.sql
) else (
    echo.
    echo ❌ Schema extraction failed. Check the error messages above.
)

echo.
pause
