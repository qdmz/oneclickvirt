$baseURL = "http://localhost:8890/api/v1/auth/login"

$payload = @{
    username = "admin"
    password = "admin123"
    loginType = "username"
    userType = "admin"
    captcha = "1234"
    captchaId = "test"
} | ConvertTo-Json -Depth 10

Write-Host "Testing login API..." -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Uri $baseURL -Method Post -Body $payload -ContentType "application/json" -TimeoutSec 10
    Write-Host "Response:" -ForegroundColor Yellow
    Write-Host ($response | ConvertTo-Json -Depth 10) -ForegroundColor White
} catch {
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        Write-Host "Status Code: $($_.Exception.Response.StatusCode)" -ForegroundColor Yellow
        try {
            $errorBody = $_.Exception.Response.GetResponseStream()
            $reader = New-Object System.IO.StreamReader($errorBody)
            $errorText = $reader.ReadToEnd()
            Write-Host "Error Body: $errorText" -ForegroundColor Red
        } {}
    }
}
