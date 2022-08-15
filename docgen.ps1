#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

# Generate image and container names using the data in the "component.json" file
$component = Get-Content -Path "$PSScriptRoot/component.json" | ConvertFrom-Json

$docImage = "$($component.registry)/$($component.name):$($component.version)-$($component.build)-docs"
$container = $component.name

# Remove build files
if (Test-Path "$PSScriptRoot/docs") {
    Remove-Item -Recurse -Force -Path "$PSScriptRoot/docs/*"
} else {
    $null = New-Item -ItemType Directory -Force -Path "$PSScriptRoot/docs"
}

# Build docker image
docker build -f "$PSScriptRoot/docker/Dockerfile.docs" -t $docImage "$PSScriptRoot/."

# Run docgen container
docker run -d --name $container $docImage
# Wait it to start
Start-Sleep -Seconds 2
# Generate docs
docker exec -ti $container /bin/bash -c "wget -r -np -N -E -p -k http://localhost:6060/pkg/"
# Copy docs from container
docker cp "$($container):/app/localhost:6060/pkg" "$PSScriptRoot/docs/pkg"
docker cp "$($container):/app/localhost:6060/lib" "$PSScriptRoot/docs/lib"
# Remove docgen container
docker rm $container --force

Write-Output "<head><meta http-equiv='refresh' content='0; URL=./pkg/index.html'></head>" > "$PSScriptRoot/docs/index.html"

# Verify docs 
if (-not (Test-Path "$PSScriptRoot/docs")) {
    Write-Error "docs folder doesn't exist in root dir. Watch logs above."
}