param(
    [switch]$SkipFrontend,
    [switch]$SkipAdmin,
    [switch]$SkipClient,
    [switch]$SkipBackend,
    [switch]$SkipTests
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

function Ensure-Command {
    param(
        [Parameter(Mandatory = $true)][string]$Name
    )
    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "Required command '$Name' not found in PATH."
    }
}

function Resolve-RepoRoot {
    return (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
}

function Detect-NodeTool {
    if (Get-Command pnpm -ErrorAction SilentlyContinue) { return "pnpm" }
    if (Get-Command npm -ErrorAction SilentlyContinue) { return "npm" }
    throw "Neither 'pnpm' nor 'npm' found. Please install one of them."
}

function Build-Frontend {
    param(
        [Parameter(Mandatory = $true)][string]$Dir,
        [Parameter(Mandatory = $true)][string]$Name,
        [Parameter(Mandatory = $true)][ValidateSet("pnpm", "npm")][string]$Tool
    )

    Write-Host "Building frontend ($Name)" -ForegroundColor Cyan
    Push-Location $Dir
    try {
        if ($Tool -eq "pnpm") {
            Write-Host "Using pnpm" -ForegroundColor Green
            if (Test-Path (Join-Path $Dir "pnpm-lock.yaml")) {
                pnpm install --frozen-lockfile
            } else {
                pnpm install
            }
            pnpm run build
        } else {
            Write-Host "Using npm" -ForegroundColor Yellow
            if (Test-Path (Join-Path $Dir "package-lock.json")) {
                npm ci
            } else {
                npm install
            }
            npm run build
        }
    }
    finally {
        Pop-Location
    }
}

$RepoRoot = Resolve-RepoRoot

$AdminDir = Join-Path $RepoRoot "frontend\admin"
$ClientDir = Join-Path $RepoRoot "frontend\client"

$BinDir = Join-Path $RepoRoot "bin"
$BackendOut = Join-Path $BinDir "mou1ght.exe"

Write-Host "Repo: $RepoRoot" -ForegroundColor DarkGray

if (-not $SkipFrontend) {
    $nodeTool = Detect-NodeTool

    if (-not $SkipAdmin) {
        if (-not (Test-Path $AdminDir)) { throw "Admin frontend directory not found: $AdminDir" }
        Build-Frontend -Dir $AdminDir -Name "frontend/admin" -Tool $nodeTool
    }

    if (-not $SkipClient) {
        if (Test-Path $ClientDir) {
            Build-Frontend -Dir $ClientDir -Name "frontend/client" -Tool $nodeTool
        } else {
            Write-Host "Skipping frontend/client (directory not found)" -ForegroundColor DarkYellow
        }
    }
} else {
    Write-Host "Skipping all frontend builds" -ForegroundColor DarkYellow
}

if (-not $SkipBackend) {
    Write-Host "Building backend (Go)" -ForegroundColor Cyan
    Ensure-Command go
    if (-not (Test-Path $BinDir)) { New-Item -ItemType Directory -Path $BinDir | Out-Null }
    Push-Location $RepoRoot
    try {
        if (-not $SkipTests) {
            Write-Host "Running go test ./..." -ForegroundColor Cyan
            go test ./...
        }
        go build -o $BackendOut ./cmd
    }
    finally {
        Pop-Location
    }
    Write-Host "Backend output: $BackendOut" -ForegroundColor Green
} else {
    Write-Host "Skipping backend build" -ForegroundColor DarkYellow
}

Write-Host "Build complete." -ForegroundColor Green
