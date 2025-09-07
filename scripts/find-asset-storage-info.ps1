# PowerShell script to find information about Canvus asset storage structure
# Usage: .\find-asset-storage-info.ps1

param(
    [string]$CanvusPath = "C:\ProgramData\MultiTaction\canvus"
)

Write-Host "🔍 Finding Canvus Asset Storage Information" -ForegroundColor Green
Write-Host "===========================================" -ForegroundColor Green

# Check if Canvus folder exists
if (-not (Test-Path $CanvusPath)) {
    Write-Error "❌ Canvus folder not found: $CanvusPath"
    Write-Host "💡 Please check the path or run as administrator if needed." -ForegroundColor Yellow
    exit 1
}

Write-Host "📁 Canvus folder found: $CanvusPath" -ForegroundColor Green

# Look for configuration files
Write-Host "`n🔍 Searching for configuration files..." -ForegroundColor Cyan

$configFiles = @()
$configFiles += Get-ChildItem -Path $CanvusPath -Recurse -Include "*.ini", "*.conf", "*.config", "*.json", "*.yaml", "*.yml" -ErrorAction SilentlyContinue

if ($configFiles.Count -gt 0) {
    Write-Host "   ✅ Found $($configFiles.Count) configuration files:" -ForegroundColor Green
    foreach ($file in $configFiles) {
        Write-Host "      📄 $($file.FullName)" -ForegroundColor White
    }
} else {
    Write-Host "   ⚠️  No configuration files found" -ForegroundColor Yellow
}

# Look for documentation files
Write-Host "`n🔍 Searching for documentation files..." -ForegroundColor Cyan

$docFiles = @()
$docFiles += Get-ChildItem -Path $CanvusPath -Recurse -Include "*.md", "*.txt", "*.doc", "*.docx" -ErrorAction SilentlyContinue

if ($docFiles.Count -gt 0) {
    Write-Host "   ✅ Found $($docFiles.Count) documentation files:" -ForegroundColor Green
    foreach ($file in $docFiles) {
        Write-Host "      📄 $($file.FullName)" -ForegroundColor White
    }
} else {
    Write-Host "   ⚠️  No documentation files found" -ForegroundColor Yellow
}

# Look for log files that might contain asset path information
Write-Host "`n🔍 Searching for log files..." -ForegroundColor Cyan

$logFiles = @()
$logFiles += Get-ChildItem -Path $CanvusPath -Recurse -Include "*.log", "*.log.*" -ErrorAction SilentlyContinue

if ($logFiles.Count -gt 0) {
    Write-Host "   ✅ Found $($logFiles.Count) log files:" -ForegroundColor Green
    foreach ($file in $logFiles) {
        Write-Host "      📄 $($file.FullName)" -ForegroundColor White
    }
} else {
    Write-Host "   ⚠️  No log files found" -ForegroundColor Yellow
}

# Search for asset-related content in files
Write-Host "`n🔍 Searching for asset-related content..." -ForegroundColor Cyan

$searchTerms = @("assets", "hash", "file", "storage", "path", "folder")
$foundContent = @()

foreach ($file in $configFiles + $docFiles + $logFiles) {
    try {
        $content = Get-Content -Path $file.FullName -ErrorAction SilentlyContinue | Out-String
        foreach ($term in $searchTerms) {
            if ($content -match $term) {
                $foundContent += @{
                    File = $file.FullName
                    Term = $term
                    Content = $content
                }
                break
            }
        }
    } catch {
        # Skip files that can't be read
    }
}

if ($foundContent.Count -gt 0) {
    Write-Host "   ✅ Found asset-related content in $($foundContent.Count) files:" -ForegroundColor Green
    foreach ($item in $foundContent) {
        Write-Host "      📄 $($item.File) (contains: $($item.Term))" -ForegroundColor White
    }
} else {
    Write-Host "   ⚠️  No asset-related content found in files" -ForegroundColor Yellow
}

# Check for subfolders in the Canvus directory
Write-Host "`n🔍 Checking Canvus directory structure..." -ForegroundColor Cyan

$subfolders = Get-ChildItem -Path $CanvusPath -Directory -ErrorAction SilentlyContinue
if ($subfolders.Count -gt 0) {
    Write-Host "   ✅ Found $($subfolders.Count) subfolders:" -ForegroundColor Green
    foreach ($folder in $subfolders) {
        $fileCount = (Get-ChildItem -Path $folder.FullName -File -ErrorAction SilentlyContinue).Count
        $subfolderCount = (Get-ChildItem -Path $folder.FullName -Directory -ErrorAction SilentlyContinue).Count
        Write-Host "      📁 $($folder.Name): $fileCount files, $subfolderCount subfolders" -ForegroundColor White
    }
} else {
    Write-Host "   ⚠️  No subfolders found" -ForegroundColor Yellow
}

# Look for any files that might contain asset path information
Write-Host "`n🔍 Searching for files containing asset paths..." -ForegroundColor Cyan

$assetPathPatterns = @(
    "assets",
    "C:\\ProgramData\\MultiTaction\\canvus\\assets",
    "hash",
    "file.*path",
    "storage.*path"
)

$foundPaths = @()
foreach ($file in $configFiles + $docFiles + $logFiles) {
    try {
        $content = Get-Content -Path $file.FullName -ErrorAction SilentlyContinue | Out-String
        foreach ($pattern in $assetPathPatterns) {
            if ($content -match $pattern) {
                $foundPaths += @{
                    File = $file.FullName
                    Pattern = $pattern
                    Content = $content
                }
                break
            }
        }
    } catch {
        # Skip files that can't be read
    }
}

if ($foundPaths.Count -gt 0) {
    Write-Host "   ✅ Found files containing asset path information:" -ForegroundColor Green
    foreach ($item in $foundPaths) {
        Write-Host "      📄 $($item.File) (pattern: $($item.Pattern))" -ForegroundColor White
    }
} else {
    Write-Host "   ⚠️  No files containing asset path information found" -ForegroundColor Yellow
}

Write-Host "`n✅ Asset storage information search complete!" -ForegroundColor Green

Write-Host "`n💡 Next Steps:" -ForegroundColor Cyan
Write-Host "   1. Review the found configuration and documentation files" -ForegroundColor White
Write-Host "   2. Check log files for asset path information" -ForegroundColor White
Write-Host "   3. Look for any asset storage configuration" -ForegroundColor White
Write-Host "   4. Run examine-asset-structure.ps1 to see actual folder structure" -ForegroundColor White
Write-Host "   5. Update KPMG DB Solver based on findings" -ForegroundColor White
