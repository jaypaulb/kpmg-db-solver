# PowerShell script to query asset hashes directly from the database
# Usage: .\query-asset-hashes.ps1

param(
    [string]$IniPath = "C:\ProgramData\MultiTaction\canvus\mt-canvus-server.ini",
    [string]$OutputPath = "database-asset-hashes.sql"
)

Write-Host "🔍 Querying Asset Hashes from Database" -ForegroundColor Green
Write-Host "=======================================" -ForegroundColor Green

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

# Create asset hash query
$hashQuery = @"
-- Asset Hash Query from Canvus Database
-- Generated on: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
-- Database: $databaseName

\set ON_ERROR_STOP on
\timing on

\echo '=== ASSET HASHES FROM DATABASE ==='

\echo '=== ALL CANVAS DOCUMENTS WITH PREVIEW HASHES ==='
SELECT 
    cd.id as widget_id,
    cd.name as widget_name,
    cd.type as widget_type,
    cd.canvas_id,
    d.name as canvas_name,
    cd.preview as asset_hash,
    CASE 
        WHEN cd.preview IS NOT NULL AND LENGTH(cd.preview) = 64 THEN 'Valid Hash'
        WHEN cd.preview IS NOT NULL THEN 'Invalid Hash Length'
        ELSE 'No Hash'
    END as hash_status
FROM canvas_documents cd
JOIN documents d ON cd.document_id = d.id
WHERE cd.preview IS NOT NULL 
AND cd.preview != ''
ORDER BY d.name, cd.name;

\echo '=== HASH STATISTICS ==='
SELECT 
    COUNT(*) as total_widgets,
    COUNT(CASE WHEN preview IS NOT NULL AND preview != '' THEN 1 END) as widgets_with_hashes,
    COUNT(CASE WHEN preview IS NOT NULL AND LENGTH(preview) = 64 THEN 1 END) as valid_hashes,
    COUNT(CASE WHEN preview IS NOT NULL AND LENGTH(preview) != 64 THEN 1 END) as invalid_hashes
FROM canvas_documents;

\echo '=== UNIQUE HASHES ==='
SELECT 
    preview as asset_hash,
    COUNT(*) as usage_count,
    STRING_AGG(DISTINCT d.name, ', ') as used_in_canvases
FROM canvas_documents cd
JOIN documents d ON cd.document_id = d.id
WHERE cd.preview IS NOT NULL 
AND cd.preview != ''
AND LENGTH(cd.preview) = 64
GROUP BY preview
ORDER BY usage_count DESC, preview;

\echo '=== CANVAS SUMMARY ==='
SELECT 
    d.id as canvas_id,
    d.name as canvas_name,
    COUNT(cd.id) as total_widgets,
    COUNT(CASE WHEN cd.preview IS NOT NULL AND cd.preview != '' THEN 1 END) as widgets_with_hashes
FROM documents d
LEFT JOIN canvas_documents cd ON d.id = cd.document_id
WHERE d.type = 'mt-canvus'
GROUP BY d.id, d.name
ORDER BY d.name;

\echo '=== ASSET HASH QUERY COMPLETE ==='
"@

Write-Host "📝 Querying asset hashes to: $OutputPath" -ForegroundColor Cyan

# Execute query
$hashQuery | & $psqlPath -h localhost -p $port -U $username -d $databaseName -o $OutputPath

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Asset hash query completed successfully!" -ForegroundColor Green
    
    $fileInfo = Get-Item $OutputPath
    $lineCount = (Get-Content $OutputPath | Measure-Object -Line).Lines
    
    Write-Host "📊 Query Results:" -ForegroundColor Cyan
    Write-Host "   File Size: $([math]::Round($fileInfo.Length / 1KB, 2)) KB" -ForegroundColor White
    Write-Host "   Lines: $lineCount" -ForegroundColor White
    
    Write-Host "`n🎉 Asset hash query completed!" -ForegroundColor Green
    Write-Host "`n💡 This will show you:" -ForegroundColor Cyan
    Write-Host "   - All widgets with asset hashes" -ForegroundColor White
    Write-Host "   - Hash validation (64-character hex strings)" -ForegroundColor White
    Write-Host "   - Usage statistics and duplicate detection" -ForegroundColor White
    Write-Host "   - Canvas-by-canvas breakdown" -ForegroundColor White
} else {
    Write-Error "❌ Asset hash query failed."
    exit 1
}

# Clean up
Remove-Item Env:PGPASSWORD

Write-Host "`n🔍 Next steps:" -ForegroundColor Cyan
Write-Host "   1. Review the query results: $OutputPath" -ForegroundColor White
Write-Host "   2. Compare with API results from the KPMG DB Solver" -ForegroundColor White
Write-Host "   3. Verify that database hashes match API hashes" -ForegroundColor White
Write-Host "   4. Use this to validate the asset discovery process" -ForegroundColor White
