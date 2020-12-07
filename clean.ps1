#!/usr/bin/env pwsh

# Recreate image names using the data in the "component.json" file
$component = Get-Content -Path "component.json" | ConvertFrom-Json
$buildImage="$($component.registry)/$($component.name):$($component.version)-$($component.build)-build"
$docsImage="$($component.registry)/$($component.name):$($component.version)-$($component.build)-docs"
$testImage="$($component.registry)/$($component.name):$($component.version)-$($component.build)-test"

# Clean up build directories
Get-ChildItem -Path "." -Include "exe" -Recurse | foreach($_) { Remove-Item -Force -Recurse $_.FullName }

# Remove docker images
docker rmi $buildImage --force
docker rmi $docsImage --force
docker rmi $testImage --force
docker image prune --force
docker rmi -f $(docker images -f "dangling=true" -q) # remove build container if build fails

# Remove existed containers
$exitedContainers = docker ps -a | Select-String -Pattern "Exit"
foreach($c in $exitedContainers) { docker rm $c.ToString().Split(" ")[0] }

# Remove unused volumes
docker volume rm -f $(docker volume ls -f "dangling=true")
