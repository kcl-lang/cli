# Stop on error.
# Note: Before you start tests, please install kcl and kpm
# kcl Installation: https://kcl-lang.io/docs/user_docs/getting-started/install
# kpm Installation: https://kcl-lang.io/docs/user_docs/guides/package-management/installation
$ErrorActionPreference = "Stop"
$pwd = Split-Path -Parent $MyInvocation.MyCommand.Path

$paths = @("configuration", "validation", "abstraction", "definition", "konfig", "mutation", "data-integration", "automation", "package-management", "kubernetes", "codelab", "server", "source")
foreach ($path in $paths) {
    Write-Host "Testing $path ..."
    Set-Location -Path "$pwd\$path"
    try {
        & make test | Out-Host
        Write-Host "Test SUCCESSED - $path" -ForegroundColor Green
    }
    catch {
        Write-Host "Test FAILED - $path" -ForegroundColor Red
        exit 1
    }
    Write-Host ""
}
