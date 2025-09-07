# PowerShell script to analyze the Canvus asset schema
# Usage: .\analyze-asset-schema.ps1 [schema-file]

param(
    [string]$SchemaFile = "canvus-asset-schema.sql"
)

Write-Host "🔍 Analyzing Canvus Asset Schema" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Green

if (-not (Test-Path $SchemaFile)) {
    Write-Error "❌ Schema file not found: $SchemaFile"
    Write-Host "Please run extract-asset-schema.ps1 first to generate the schema file." -ForegroundColor Yellow
    exit 1
}

Write-Host "📄 Analyzing schema file: $SchemaFile" -ForegroundColor Cyan

$content = Get-Content $SchemaFile

Write-Host "`n📊 SCHEMA ANALYSIS RESULTS" -ForegroundColor Yellow
Write-Host "=========================" -ForegroundColor Yellow

# Find tables
Write-Host "`n🗄️ TABLES FOUND:" -ForegroundColor Cyan
$tables = $content | Where-Object { $_ -match "^\s*(\w+)\s+\|" -and $_ -notmatch "table_name|column_name|Name|Owner" }
foreach ($table in $tables) {
    if ($table.Trim() -ne "" -and $table -notmatch "^\s*[-+]+\s*$") {
        Write-Host "   - $($table.Trim())" -ForegroundColor White
    }
}

# Find hash-related columns
Write-Host "`n🔑 HASH/FILE/PATH RELATED COLUMNS:" -ForegroundColor Cyan
$hashColumns = $content | Where-Object { 
    $_ -match "(hash|file|path|asset|name|filename)" -and 
    $_ -notmatch "table_name|column_name|Name|Owner" -and
    $_ -notmatch "^\s*[-+]+\s*$"
}
foreach ($column in $hashColumns) {
    if ($column.Trim() -ne "") {
        Write-Host "   - $($column.Trim())" -ForegroundColor White
    }
}

# Find sample data
Write-Host "`n📋 SAMPLE DATA FOUND:" -ForegroundColor Cyan
$inSampleData = $false
$sampleData = @()

foreach ($line in $content) {
    if ($line -match "SAMPLE.*ROWS") {
        $inSampleData = $true
        continue
    }
    if ($inSampleData -and $line -match "^\s*[-+]+\s*$") {
        $inSampleData = $false
        continue
    }
    if ($inSampleData -and $line.Trim() -ne "" -and $line -notmatch "table_name|column_name") {
        $sampleData += $line
    }
}

if ($sampleData.Count -gt 0) {
    Write-Host "   Found $($sampleData.Count) sample data rows" -ForegroundColor White
    foreach ($row in $sampleData[0..4]) {  # Show first 5 rows
        Write-Host "   $($row.Trim())" -ForegroundColor White
    }
    if ($sampleData.Count -gt 5) {
        Write-Host "   ... and $($sampleData.Count - 5) more rows" -ForegroundColor Gray
    }
} else {
    Write-Host "   No sample data found" -ForegroundColor Gray
}

# Look for specific patterns
Write-Host "`n🔍 KEY INSIGHTS:" -ForegroundColor Cyan

# Check for hash patterns
$hashPatterns = $content | Where-Object { $_ -match "[a-f0-9]{32,}" }
if ($hashPatterns.Count -gt 0) {
    Write-Host "   ✅ Found potential hash values in data" -ForegroundColor Green
    Write-Host "      Sample: $($hashPatterns[0].Trim())" -ForegroundColor White
} else {
    Write-Host "   ❌ No obvious hash patterns found" -ForegroundColor Red
}

# Check for file paths
$pathPatterns = $content | Where-Object { $_ -match "\\\\|/" }
if ($pathPatterns.Count -gt 0) {
    Write-Host "   ✅ Found potential file paths" -ForegroundColor Green
    Write-Host "      Sample: $($pathPatterns[0].Trim())" -ForegroundColor White
} else {
    Write-Host "   ❌ No obvious file paths found" -ForegroundColor Red
}

# Check for file extensions
$extPatterns = $content | Where-Object { $_ -match "\.(jpg|jpeg|png|gif|pdf|mp4|avi|mov)" }
if ($extPatterns.Count -gt 0) {
    Write-Host "   ✅ Found potential file extensions" -ForegroundColor Green
    Write-Host "      Sample: $($extPatterns[0].Trim())" -ForegroundColor White
} else {
    Write-Host "   ❌ No obvious file extensions found" -ForegroundColor Red
}

Write-Host "`n💡 RECOMMENDATIONS:" -ForegroundColor Yellow
Write-Host "===================" -ForegroundColor Yellow
Write-Host "1. Look for tables with 'hash', 'file', 'path', or 'asset' columns" -ForegroundColor White
Write-Host "2. Check the sample data for actual hash values and file paths" -ForegroundColor White
Write-Host "3. Identify which table stores the mapping between hashes and file locations" -ForegroundColor White
Write-Host "4. Look for foreign key relationships between documents and assets" -ForegroundColor White

Write-Host "`n🎯 Next steps:" -ForegroundColor Cyan
Write-Host "   - Review the full schema file: $SchemaFile" -ForegroundColor White
Write-Host "   - Focus on tables that contain hash or file information" -ForegroundColor White
Write-Host "   - Identify the relationship between database records and disk files" -ForegroundColor White
