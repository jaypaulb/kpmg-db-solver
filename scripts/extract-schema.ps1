# PowerShell script to extract PostgreSQL schema from mt-canvus-server.ini
# Usage: .\extract-schema.ps1

param(
    [string]$IniPath = "C:\ProgramData\MultiTaction\canvus\mt-canvus-server.ini",
    [string]$OutputPath = "canvus-schema.sql"
)

Write-Host "🔍 Extracting Canvus Database Schema" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green

# Check if ini file exists
if (-not (Test-Path $IniPath)) {
    Write-Error "❌ INI file not found: $IniPath"
    Write-Host "Please ensure the Canvus server is installed and the INI file exists." -ForegroundColor Yellow
    exit 1
}

Write-Host "📄 Reading configuration from: $IniPath" -ForegroundColor Cyan

# Parse INI file
$iniContent = Get-Content $IniPath
$config = @{}

$currentSection = ""
foreach ($line in $iniContent) {
    $line = $line.Trim()

    # Skip empty lines and comments
    if ($line -eq "" -or $line.StartsWith(";")) {
        continue
    }

    # Check for section headers
    if ($line.StartsWith("[") -and $line.EndsWith("]")) {
        $currentSection = $line.Substring(1, $line.Length - 2)
        continue
    }

    # Parse key=value pairs
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
Write-Host "   Password: [HIDDEN]" -ForegroundColor White

# Check if psql is available
$psqlPath = Get-Command psql -ErrorAction SilentlyContinue
if (-not $psqlPath) {
    # Try to find psql in common PostgreSQL installation locations
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
        Write-Host "✅ Found psql at: $foundPath" -ForegroundColor Green
        $psqlPath = $foundPath
    } else {
        Write-Error "❌ psql command not found. Please ensure PostgreSQL client tools are installed."
        Write-Host "Common installation locations checked:" -ForegroundColor Yellow
        Write-Host "  - C:\Program Files\PostgreSQL\*\bin\psql.exe" -ForegroundColor Yellow
        Write-Host "  - C:\Program Files (x86)\PostgreSQL\*\bin\psql.exe" -ForegroundColor Yellow
        Write-Host "  - C:\PostgreSQL\*\bin\psql.exe" -ForegroundColor Yellow
        Write-Host "You can download PostgreSQL client tools from: https://www.postgresql.org/download/" -ForegroundColor Yellow
        exit 1
    }
}

if ($psqlPath -is [string]) {
    Write-Host "✅ Found psql at: $psqlPath" -ForegroundColor Green
} else {
    Write-Host "✅ Found psql at: $($psqlPath.Source)" -ForegroundColor Green
}

# Set PGPASSWORD environment variable for password authentication
$env:PGPASSWORD = $password

Write-Host "🔗 Connecting to database..." -ForegroundColor Cyan

# Create the schema extraction command
$schemaCommands = @"
-- Canvus Database Schema Export
-- Generated on: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
-- Database: $databaseName

-- Set output formatting
\set ON_ERROR_STOP on
\timing on

-- Export database schema (structure only)
\echo '=== DATABASE SCHEMA ==='
\dn+  -- List all schemas with details

\echo '=== TABLES ==='
\dt+  -- List all tables with details

\echo '=== TABLE STRUCTURES ==='
"@

# Get list of tables first
$tablesQuery = "SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY tablename;"
$tables = & $psqlPath -h localhost -p $port -U $username -d $databaseName -t -c $tablesQuery

if ($LASTEXITCODE -ne 0) {
    Write-Error "❌ Failed to connect to database. Please check your credentials and connection."
    exit 1
}

Write-Host "✅ Successfully connected to database" -ForegroundColor Green

# Add table structure commands for each table
foreach ($table in $tables) {
    $table = $table.Trim()
    if ($table -ne "") {
        $schemaCommands += "`n\echo '=== TABLE: $table ==='`n\d+ $table`n"
    }
}

$schemaCommands += @"

\echo '=== INDEXES ==='
\di+  -- List all indexes

\echo '=== SEQUENCES ==='
\ds+  -- List all sequences

\echo '=== FUNCTIONS ==='
\df+  -- List all functions

\echo '=== VIEWS ==='
\dv+  -- List all views

\echo '=== FOREIGN KEYS ==='
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

\echo '=== SCHEMA EXPORT COMPLETE ==='
"@

Write-Host "📝 Extracting schema to: $OutputPath" -ForegroundColor Cyan

# Execute schema extraction
$schemaCommands | & $psqlPath -h localhost -p $port -U $username -d $databaseName -o $OutputPath

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Schema successfully extracted to: $OutputPath" -ForegroundColor Green

    # Show file size and line count
    $fileInfo = Get-Item $OutputPath
    $lineCount = (Get-Content $OutputPath | Measure-Object -Line).Lines

    Write-Host "📊 Export Statistics:" -ForegroundColor Cyan
    Write-Host "   File Size: $([math]::Round($fileInfo.Length / 1KB, 2)) KB" -ForegroundColor White
    Write-Host "   Lines: $lineCount" -ForegroundColor White

    Write-Host "`n🎉 Schema extraction completed successfully!" -ForegroundColor Green
    Write-Host "You can now share the schema file: $OutputPath" -ForegroundColor Yellow
} else {
    Write-Error "❌ Schema extraction failed. Check the error messages above."
    exit 1
}

# Clean up environment variable
Remove-Item Env:PGPASSWORD

Write-Host "`n💡 Next steps:" -ForegroundColor Cyan
Write-Host "   1. Review the schema file: $OutputPath" -ForegroundColor White
Write-Host "   2. Share the schema file for analysis" -ForegroundColor White
Write-Host "   3. Use the schema to understand the database structure" -ForegroundColor White
