#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

$component = Get-Content -Path "component.json" | ConvertFrom-Json
$version = (Get-Content -Path package.json | ConvertFrom-Json).version

if ($component.version -ne $version) {
    throw "Versions in component.json and package.json do not match"
}

# Automatically login to server
if ($env:NPM_USER -ne $null -and $env:NPM_PASS -ne $null) {
    npm-cli-login
}

# Publish to global repository
Write-Output "Pushing package to npm registry"
npm publish
