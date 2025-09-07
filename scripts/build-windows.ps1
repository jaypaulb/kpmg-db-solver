# PowerShell script to build KPMG DB Solver for Windows
# Usage: .\scripts\build-windows.ps1

Write-Host "🔨 Building KPMG DB Solver for Windows" -ForegroundColor Green
Write-Host "=======================================" -ForegroundColor Green

# Check if Go is installed
$goVersion = go version 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Error "❌ Go is not installed or not in PATH"
    Write-Host "💡 Please install Go from: https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ Go version: $goVersion" -ForegroundColor Green

# Set build parameters
$buildTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
$version = "1.0.0"
$outputDir = "bin"
$outputFile = "kpmg-db-solver.exe"

# Create output directory if it doesn't exist
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
    Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Cyan
}

Write-Host "`n🔧 Build Configuration:" -ForegroundColor Cyan
Write-Host "   Target OS: Windows (amd64)" -ForegroundColor White
Write-Host "   Output: $outputDir\$outputFile" -ForegroundColor White
Write-Host "   Version: $version" -ForegroundColor White
Write-Host "   Build Time: $buildTime" -ForegroundColor White

# Set environment variables for cross-compilation
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

Write-Host "`n🚀 Building executable..." -ForegroundColor Cyan

# Build the executable
$buildArgs = @(
    "build",
    "-ldflags", "-X main.version=$version -X main.buildTime=$buildTime",
    "-o", "$outputDir\$outputFile",
    "./cmd/kpmg-db-solver"
)

& go @buildArgs

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Build completed successfully!" -ForegroundColor Green
    
    # Get file info
    $fileInfo = Get-Item "$outputDir\$outputFile"
    $fileSize = [math]::Round($fileInfo.Length / 1MB, 2)
    
    Write-Host "`n📊 Build Results:" -ForegroundColor Cyan
    Write-Host "   File: $($fileInfo.FullName)" -ForegroundColor White
    Write-Host "   Size: $fileSize MB" -ForegroundColor White
    Write-Host "   Created: $($fileInfo.CreationTime.ToString('yyyy-MM-dd HH:mm:ss'))" -ForegroundColor White
    
    Write-Host "`n🎉 Windows executable ready!" -ForegroundColor Green
    Write-Host "💡 You can now copy $outputDir\$outputFile to your Windows machine" -ForegroundColor Cyan
    
} else {
    Write-Error "❌ Build failed!"
    exit 1
}

# Clean up environment variables
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
Remove-Item Env:CGO_ENABLED

Write-Host "`n🔍 Next Steps:" -ForegroundColor Cyan
Write-Host "   1. Copy the executable to your Windows machine" -ForegroundColor White
Write-Host "   2. Run: .\kpmg-db-solver.exe discover" -ForegroundColor White
Write-Host "   3. Or run: .\kpmg-db-solver.exe run" -ForegroundColor White
