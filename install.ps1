param($Version = "v0.1.0", $InstallPath = (Get-Location))
$ErrorActionPreference = "Stop"

$zipFileName = "proofable-image_windows_amd64.zip"
$downloadLink = "https://github.com/SouthbankSoftware/proofable-image/releases/download/$Version/$zipFileName"
$zipFilePath = Join-Path -Path $InstallPath -ChildPath $zipFileName

Write-Host "Installing from ``$zipFileName`` to ``$InstallPath``..."
(New-Object Net.WebClient).DownloadFile($downloadLink, $zipFilePath)
Expand-Archive -Path $zipFilePath -DestinationPath $InstallPath -Force
Remove-Item -Path $zipFilePath
