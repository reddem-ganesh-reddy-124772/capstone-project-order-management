# Order Management System - Start All Services
#Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass

$ErrorActionPreference = "Stop"

$ProjectRoot = $PSScriptRoot
$ConsulAddress = "127.0.0.1:8500"
$ConsulExe = "C:\Users\reddem.reddy\Downloads\consul_2.0.4_windows_amd64\consul.exe"

$Services = @(
    @{ Name = "auth-service";    ID = "auth-service-1";    Path = "cmd/auth-service";    Port = 8083 },
    @{ Name = "product-service"; ID = "product-service-1"; Path = "cmd/product-service"; Port = 8082 },
    @{ Name = "order-service";   ID = "order-service-1";   Path = "cmd/order-service";   Port = 8081 },
    @{ Name = "api-gateway";     ID = "api-gateway-1";     Path = "cmd/api-gateway";     Port = 8080 }
)

function Test-TcpPort {
    param([string]$HostName, [int]$Port)

    $Client = New-Object System.Net.Sockets.TcpClient

    try {
        $Result = $Client.BeginConnect($HostName, $Port, $null, $null)

        if (-not $Result.AsyncWaitHandle.WaitOne(1000, $false)) {
            return $false
        }

        $Client.EndConnect($Result)
        return $true
    }
    catch {
        return $false
    }
    finally {
        $Client.Close()
    }
}

function Wait-ForPort {
    param([int]$Port, [int]$TimeoutSeconds = 60)

    $Deadline = (Get-Date).AddSeconds($TimeoutSeconds)

    while ((Get-Date) -lt $Deadline) {
        if (Test-TcpPort -HostName "127.0.0.1" -Port $Port) {
            return $true
        }

        Start-Sleep -Seconds 2
    }

    return $false
}

function Register-ConsulService {
    param([string]$Name, [string]$ID, [int]$Port)

    $Registration = @{
        ID      = $ID
        Name    = $Name
        Address = "127.0.0.1"
        Port    = $Port
        Check   = @{
            TCP      = "127.0.0.1:$Port"
            Interval = "40s"
            Timeout  = "2s"
        }
    }

    $Body = $Registration | ConvertTo-Json -Depth 5

    Invoke-RestMethod `
        -Method Put `
        -Uri "http://$ConsulAddress/v1/agent/service/register" `
        -ContentType "application/json" `
        -Body $Body | Out-Null

    Write-Host "Registered $Name with Consul on port $Port" -ForegroundColor Green
}

Write-Host "`n=== ORDER MANAGEMENT STARTUP ===" -ForegroundColor Cyan

# Confirm the root .env file exists.
$EnvFile = Join-Path $ProjectRoot ".env"

if (-not (Test-Path $EnvFile)) {
    throw ".env file not found at: $EnvFile"
}

# Confirm each service has a main.go file.
foreach ($Service in $Services) {
    $MainFile = Join-Path $ProjectRoot (Join-Path $Service.Path "main.go")

    if (-not (Test-Path $MainFile)) {
        throw "Service main.go not found: $MainFile"
    }
}

# Find Consul.
$ConsulCommand = Get-Command consul -ErrorAction SilentlyContinue

if ($ConsulCommand) {
    $ConsulExe = $ConsulCommand.Source
}

# Start Consul if it is not already listening.
if (-not (Test-TcpPort -HostName "127.0.0.1" -Port 8500)) {
    if (-not (Test-Path $ConsulExe)) {
        throw "Consul executable not found: $ConsulExe"
    }

    Write-Host "Starting Consul..." -ForegroundColor Yellow

    Start-Process `
        -FilePath $ConsulExe `
        -ArgumentList @("agent", "-dev", "-client=127.0.0.1") `
        -WorkingDirectory (Split-Path $ConsulExe)

    if (-not (Wait-ForPort -Port 8500 -TimeoutSeconds 30)) {
        throw "Consul did not start on port 8500."
    }
}

Invoke-RestMethod -Uri "http://$ConsulAddress/v1/status/leader" | Out-Null
Write-Host "Consul is ready." -ForegroundColor Green

# Build all Go packages before starting the services.
Write-Host "Building Go project..." -ForegroundColor Yellow

Push-Location $ProjectRoot
try {
    go build ./...

    if ($LASTEXITCODE -ne 0) {
        throw "Go build failed."
    }
}
finally {
    Pop-Location
}

Write-Host "Build successful." -ForegroundColor Green

# Start each service from the PROJECT ROOT.
# This allows godotenv.Load() to find the root .env file.
foreach ($Service in $Services) {
    $Port = [int]$Service.Port

    Write-Host "`nStarting $($Service.Name)..." -ForegroundColor Cyan

    if (Test-TcpPort -HostName "127.0.0.1" -Port $Port) {
        Write-Host "Port $Port is already listening; checking/registering it." -ForegroundColor Yellow
    }
    else {
        $RelativePath = $Service.Path

        $Command = "Set-Location -LiteralPath '$ProjectRoot'; go run ./$RelativePath"

        Start-Process `
            -FilePath "powershell.exe" `
            -WorkingDirectory $ProjectRoot `
            -ArgumentList @("-NoExit", "-Command", $Command)
    }

    Write-Host "Waiting for $($Service.Name) on port $Port..." -ForegroundColor Yellow

    if (-not (Wait-ForPort -Port $Port -TimeoutSeconds 60)) {
        Write-Warning "$($Service.Name) failed to open port $Port. Check its service window."
        continue
    }

    try {
        Register-ConsulService `
            -Name $Service.Name `
            -ID $Service.ID `
            -Port $Port
    }
    catch {
        Write-Warning "Consul registration failed for $($Service.Name): $($_.Exception.Message)"
    }
}

Write-Host "`n=== STARTUP SCRIPT FINISHED ===" -ForegroundColor Cyan
Write-Host "Consul UI:   http://localhost:8500/ui/"
Write-Host "API Gateway: http://localhost:8080"
Write-Host "Check the individual service windows for errors." -ForegroundColor Yellow