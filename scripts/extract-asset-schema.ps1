# PowerShell script to extract specific asset-related schema from Canvus database
# Usage: .\extract-asset-schema.ps1

param(
    [string]$IniPath = "C:\ProgramData\MultiTaction\canvus\mt-canvus-server.ini",
    [string]$OutputPath = "canvus-asset-schema.sql"
)

Write-Host "🔍 Extracting Canvus Asset Schema" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Green

# Check if ini file exists
if (-not (Test-Path $IniPath)) {
    Write-Error "❌ INI file not found: $IniPath"
    exit 1
}

# Parse INI file
$iniContent = Get-Content $IniPath
$config = @{}

$currentSection = ""
foreach ($line in $iniContent) {
    $line = $line.Trim()
    
    if ($line -eq "" -or $line.StartsWith(";")) {
        continue
    }
    
    if ($line.StartsWith("[") -and $line.EndsWith("]")) {
        $currentSection = $line.Substring(1, $line.Length - 2)
        continue
    }
    
    if ($line.Contains("=")) {
        $parts = $line.Split("=", 2)
        $key = $parts[0].Trim()
        $value = $parts[1].Trim()
        
        if ($currentSection -eq "sql") {
            $config[$key] = $value
        }
    }
}

# Extract database connection details
$databaseName = $config["databasename"]
$username = $config["username"]
$password = $config["password"]
$port = if ($config["port"]) { $config["port"] } else { "5432" }

Write-Host "📊 Database Configuration:" -ForegroundColor Cyan
Write-Host "   Database: $databaseName" -ForegroundColor White
Write-Host "   Username: $username" -ForegroundColor White
Write-Host "   Port: $port" -ForegroundColor White

# Find psql
$psqlPath = Get-Command psql -ErrorAction SilentlyContinue
if (-not $psqlPath) {
    $commonPaths = @(
        "C:\Program Files\PostgreSQL\*\bin\psql.exe",
        "C:\Program Files (x86)\PostgreSQL\*\bin\psql.exe",
        "C:\PostgreSQL\*\bin\psql.exe"
    )
    
    $foundPath = $null
    foreach ($path in $commonPaths) {
        $matches = Get-ChildItem -Path $path -ErrorAction SilentlyContinue | Sort-Object Name -Descending
        if ($matches) {
            $foundPath = $matches[0].FullName
            break
        }
    }
    
    if ($foundPath) {
        $psqlPath = $foundPath
    } else {
        Write-Error "❌ psql command not found."
        exit 1
    }
}

Write-Host "✅ Found psql at: $psqlPath" -ForegroundColor Green

# Set PGPASSWORD environment variable
$env:PGPASSWORD = $password

Write-Host "🔗 Connecting to database..." -ForegroundColor Cyan

# Create focused schema extraction commands
$schemaCommands = @"
-- Canvus Asset Schema Export
-- Generated on: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
-- Database: $databaseName

\set ON_ERROR_STOP on
\timing on

\echo '=== ALL TABLES IN DATABASE ==='
SELECT schemaname, tablename, tableowner 
FROM pg_tables 
WHERE schemaname = 'public' 
ORDER BY tablename;

\echo '=== ASSET-RELATED TABLE STRUCTURES ==='

\echo '=== DOCUMENTS TABLE ==='
\d+ documents

\echo '=== CANVAS_DOCUMENTS TABLE ==='
\d+ canvas_documents

\echo '=== SERVER_INDEX_FOLDERS TABLE ==='
\d+ server_index_folders

\echo '=== USERS TABLE ==='
\d+ users

\echo '=== GROUPS TABLE ==='
\d+ groups

\echo '=== ACCESS_TOKENS TABLE ==='
\d+ access_tokens

\echo '=== SAMPLE DATA FROM KEY TABLES ==='

\echo '=== DOCUMENTS TABLE SAMPLE (first 5 rows) ==='
SELECT * FROM documents LIMIT 5;

\echo '=== CANVAS_DOCUMENTS TABLE SAMPLE (first 5 rows) ==='
SELECT * FROM canvas_documents LIMIT 5;

\echo '=== SERVER_INDEX_FOLDERS TABLE SAMPLE (first 5 rows) ==='
SELECT * FROM server_index_folders LIMIT 5;

\echo '=== COLUMNS THAT MIGHT CONTAIN HASHES ==='
SELECT 
    table_name,
    column_name,
    data_type,
    is_nullable
FROM information_schema.columns 
WHERE table_schema = 'public' 
AND (
    column_name ILIKE '%hash%' OR 
    column_name ILIKE '%file%' OR 
    column_name ILIKE '%path%' OR 
    column_name ILIKE '%asset%' OR
    column_name ILIKE '%name%' OR
    column_name ILIKE '%filename%'
)
ORDER BY table_name, column_name;

\echo '=== FOREIGN KEY RELATIONSHIPS ==='
SELECT
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
    AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
ORDER BY tc.table_name, kcu.column_name;

\echo '=== ASSET SCHEMA EXPORT COMPLETE ==='
"@

Write-Host "📝 Extracting asset schema to: $OutputPath" -ForegroundColor Cyan

# Execute schema extraction
$schemaCommands | & $psqlPath -h localhost -p $port -U $username -d $databaseName -o $OutputPath

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Asset schema successfully extracted to: $OutputPath" -ForegroundColor Green
    
    $fileInfo = Get-Item $OutputPath
    $lineCount = (Get-Content $OutputPath | Measure-Object -Line).Lines
    
    Write-Host "📊 Export Statistics:" -ForegroundColor Cyan
    Write-Host "   File Size: $([math]::Round($fileInfo.Length / 1KB, 2)) KB" -ForegroundColor White
    Write-Host "   Lines: $lineCount" -ForegroundColor White
    
    Write-Host "`n🎉 Asset schema extraction completed!" -ForegroundColor Green
} else {
    Write-Error "❌ Schema extraction failed."
    exit 1
}

# Clean up
Remove-Item Env:PGPASSWORD

Write-Host "`n💡 Next steps:" -ForegroundColor Cyan
Write-Host "   1. Review the schema file: $OutputPath" -ForegroundColor White
Write-Host "   2. Look for hash/file/path columns in the output" -ForegroundColor White
Write-Host "   3. Analyze the sample data to understand asset storage" -ForegroundColor White
