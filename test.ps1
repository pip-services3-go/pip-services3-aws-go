#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

# Generate an image name using the data in the "$PSScriptRoot/component.json" file
$component = Get-Content -Path "$PSScriptRoot/component.json" | ConvertFrom-Json
$testImage = "$($component.registry)/$($component.name):$($component.version)-$($component.build)-test"

# Set environment variables
$env:IMAGE = $testImage

# Copy private keys to access git repo
if (-not (Test-Path -Path "$PSScriptRoot/docker/id_rsa")) {
    if (-not [string]::IsNullOrEmpty($env:GIT_PRIVATE_KEY)) {
        Write-Host "Creating docker/id_rsa from environment variable..."
        Set-Content -Path "$PSScriptRoot/docker/id_rsa" -Value $env:GIT_PRIVATE_KEY
    } elseif (Test-Path -Path "~/.ssh/id_rsa") {
        Write-Host "Copying ~/.ssh/id_rsa to docker..."
        Copy-Item -Path "~/.ssh/id_rsa" -Destination "docker"
    } else {
        Write-Host "Missing ~/.ssh/id_rsa file..."
        Set-Content -Path "$PSScriptRoot/docker/id_rsa" -Value ""
    }
}

try {
    # Workaround to remove dangling images
    docker-compose -f "$PSScriptRoot/docker/docker-compose.test.yml" down

    docker-compose -f "$PSScriptRoot/docker/docker-compose.test.yml" up --build --abort-on-container-exit --exit-code-from test

    # Save the result to avoid overwriting it with the "down" command below
    $exitCode = $LastExitCode 
} finally {
    # Workaround to remove dangling images
    docker-compose -f "$PSScriptRoot/docker/docker-compose.test.yml" down
}

# Return the exit code of the "docker-compose.test.yml up" command
exit $exitCode 
