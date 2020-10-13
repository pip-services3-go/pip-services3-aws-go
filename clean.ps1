#!/usr/bin/env pwsh

$component = Get-Content -Path "component.json" | ConvertFrom-Json
$buildImage="$($component.registry)/$($component.name):$($component.version)-build"
$testImage="$($component.registry)/$($component.name):$($component.version)-test"

# Clean up build directories
if (Test-Path "dist") {
    Remove-Item -Recurse -Force -Path "dist"
}

# Remove docker images
docker rmi $buildImage --force
docker rmi $testImage --force
docker image prune --force

# Remove existed containers
docker ps -a | Select-String -Pattern "Exit" | foreach($_) { docker rm $_.ToString().Split(" ")[0] }
