
# Order Management System - Consul Configuration Setup
# Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
$ErrorActionPreference = "Stop"

$ConsulAddress = "http://127.0.0.1:8500"
$ConfigKey = "$ConsulAddress/v1/kv/order-management/config"

Write-Host ""
Write-Host "========================================"
Write-Host " Order Management - Consul Configuration"   
Write-Host "========================================"

# Step 1: Check Consul
Write-Host "`n[1/3] Checking Consul..."

try {
    $leader = Invoke-RestMethod "$ConsulAddress/v1/status/leader"

    if ([string]::IsNullOrWhiteSpace([string]$leader)) {
        throw "Consul has no elected leader."
    }

    Write-Host "Consul is running. Leader: $leader"
}
catch {
    Write-Host "ERROR: Consul is not available."
    exit 1
}

# Step 2: Store configuration
Write-Host "`n[2/3] Storing configuration..."

$config = @{
    order_service_url = "http://localhost:8081"
    product_service_url = "http://localhost:8082"
    auth_service_url = "http://localhost:8083"
    http_client_timeout_seconds = 5
    http_client_max_retries = 3
    circuit_breaker_failure_threshold = 5
    circuit_breaker_reset_timeout_seconds = 30
} | ConvertTo-Json

try {
    $result = Invoke-RestMethod `
        -Method Put `
        -Uri $ConfigKey `
        -ContentType "application/json" `
        -Body $config

    if (-not $result) {
        throw "Consul did not confirm the write."
    }

    Write-Host "Configuration stored successfully."
}
catch {
    Write-Host "ERROR: Could not store configuration."
    Write-Host $_.Exception.Message
    exit 1
}

# Step 3: Read and verify configuration
Write-Host "`n[3/3] Verifying saved configuration..."

try {
    $response = Invoke-WebRequest `
    -Uri "$($ConfigKey)?raw" `
    -Method Get `
    -UseBasicParsing

    $savedConfig = $response.Content | ConvertFrom-Json

    $savedConfig | Format-List

    Write-Host "Configuration verified successfully."
}
catch {
    Write-Host "ERROR: Could not verify the saved configuration."
    Write-Host $_.Exception.Message
    exit 1
}

Write-Host "`nDone. Configuration is saved in Consul."
