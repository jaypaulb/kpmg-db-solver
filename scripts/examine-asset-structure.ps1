# PowerShell script to examine the actual Canvus assets folder structure
# Usage: .\examine-asset-structure.ps1

param(
    [string]$AssetsPath = "C:\ProgramData\MultiTaction\canvus\assets",
    [string]$OutputPath = "asset-folder-structure.txt"
)

Write-Host "🔍 Examining Canvus Assets Folder Structure" -ForegroundColor Green
Write-Host "=============================================" -ForegroundColor Green

# Check if assets folder exists
if (-not (Test-Path $AssetsPath)) {
    Write-Error "❌ Assets folder not found: $AssetsPath"
    Write-Host "💡 Please check the path or run as administrator if needed." -ForegroundColor Yellow
    exit 1
}

Write-Host "📁 Assets folder found: $AssetsPath" -ForegroundColor Green

# Get folder structure
Write-Host "`n📊 Analyzing folder structure..." -ForegroundColor Cyan

$structure = @()
$fileCount = 0
$folderCount = 0
$totalSize = 0

# Function to get folder info
function Get-FolderInfo {
    param($Path, $Depth = 0)
    
    $indent = "  " * $Depth
    $folder = Get-Item $Path
    
    $structure += "$indent📁 $($folder.Name) ($($folder.CreationTime.ToString('yyyy-MM-dd HH:mm')))"
    
    # Get subfolders
    $subfolders = Get-ChildItem -Path $Path -Directory -ErrorAction SilentlyContinue | Sort-Object Name
    foreach ($subfolder in $subfolders) {
        $folderCount++
        Get-FolderInfo -Path $subfolder.FullName -Depth ($Depth + 1)
    }
    
    # Get files in this folder
    $files = Get-ChildItem -Path $Path -File -ErrorAction SilentlyContinue | Sort-Object Name
    foreach ($file in $files) {
        $fileCount++
        $totalSize += $file.Length
        $structure += "$indent  📄 $($file.Name) ($([math]::Round($file.Length / 1KB, 2)) KB)"
    }
}

# Start analysis
Get-FolderInfo -Path $AssetsPath

# Get additional statistics
Write-Host "`n📈 Folder Statistics:" -ForegroundColor Cyan
Write-Host "   Total Folders: $folderCount" -ForegroundColor White
Write-Host "   Total Files: $fileCount" -ForegroundColor White
Write-Host "   Total Size: $([math]::Round($totalSize / 1MB, 2)) MB" -ForegroundColor White

# Look for hash patterns in filenames
Write-Host "`n🔍 Analyzing file patterns..." -ForegroundColor Cyan

$hashFiles = 0
$extensionFiles = 0
$otherFiles = 0

Get-ChildItem -Path $AssetsPath -Recurse -File -ErrorAction SilentlyContinue | ForEach-Object {
    $filename = $_.BaseName
    $extension = $_.Extension
    
    # Check if filename looks like a hash (64 hex characters)
    if ($filename -match '^[0-9a-fA-F]{64}$') {
        $hashFiles++
    }
    # Check if it has a common file extension
    elseif ($extension -match '\.(jpg|jpeg|png|gif|bmp|pdf|mp4|avi|mov|webm|mp3|wav)$') {
        $extensionFiles++
    }
    else {
        $otherFiles++
    }
}

Write-Host "   Files with hash-like names (64 hex chars): $hashFiles" -ForegroundColor White
Write-Host "   Files with common extensions: $extensionFiles" -ForegroundColor White
Write-Host "   Other files: $otherFiles" -ForegroundColor White

# Look for specific patterns
Write-Host "`n🎯 Pattern Analysis:" -ForegroundColor Cyan

# Check for subfolder organization
$subfolders = Get-ChildItem -Path $AssetsPath -Directory -ErrorAction SilentlyContinue
if ($subfolders.Count -gt 0) {
    Write-Host "   ✅ Assets are organized in subfolders:" -ForegroundColor Green
    foreach ($subfolder in $subfolders) {
        $fileCountInSubfolder = (Get-ChildItem -Path $subfolder.FullName -File -ErrorAction SilentlyContinue).Count
        Write-Host "      📁 $($subfolder.Name): $fileCountInSubfolder files" -ForegroundColor White
    }
} else {
    Write-Host "   ⚠️  No subfolders found - assets may be in root folder" -ForegroundColor Yellow
}

# Check for hash distribution
if ($hashFiles -gt 0) {
    Write-Host "   ✅ Found $hashFiles files with hash-like names" -ForegroundColor Green
    Write-Host "   💡 These are likely the asset files referenced in the database" -ForegroundColor Cyan
}

# Output to file
Write-Host "`n📝 Writing structure to: $OutputPath" -ForegroundColor Cyan
$structure | Out-File -FilePath $OutputPath -Encoding UTF8

Write-Host "✅ Analysis complete!" -ForegroundColor Green

Write-Host "`n🔍 Key Findings:" -ForegroundColor Cyan
Write-Host "   1. Assets folder structure has been analyzed" -ForegroundColor White
Write-Host "   2. File patterns and organization identified" -ForegroundColor White
Write-Host "   3. Hash file distribution calculated" -ForegroundColor White
Write-Host "   4. Results saved to: $OutputPath" -ForegroundColor White

Write-Host "`n💡 Next Steps:" -ForegroundColor Cyan
Write-Host "   1. Review the folder structure in: $OutputPath" -ForegroundColor White
Write-Host "   2. Check if assets are in subfolders or root folder" -ForegroundColor White
Write-Host "   3. Verify hash file naming patterns" -ForegroundColor White
Write-Host "   4. Update KPMG DB Solver if folder structure is different than expected" -ForegroundColor White
