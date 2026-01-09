$ErrorActionPreference = 'Stop'

$packageName = 'rime-wanxiang-updater'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$version = '$VERSION$'

# 根据系统架构选择正确的可执行文件
$architecture = Get-ProcessorArchitecture
$exeName = switch ($architecture) {
    'X64' { "rime-wanxiang-updater-windows-amd64.exe" }
    'ARM64' { "rime-wanxiang-updater-windows-arm64.exe" }
    default { "rime-wanxiang-updater-windows-amd64.exe" }
}

$packageArgs = @{
    packageName    = $packageName
    fileType       = 'EXE'
    url64bit       = "https://github.com/czyt/rime-wanxiang-updater/releases/download/v$version/rime-wanxiang-updater-windows-amd64.exe"
    urlARM64       = "https://github.com/czyt/rime-wanxiang-updater/releases/download/v$version/rime-wanxiang-updater-windows-arm64.exe"
    checksum64     = '$CHECKSUM64$'
    checksumARM64  = '$CHECKSUMARM64$'
    checksumType   = 'sha256'
    silentArgs     = ''
    validExitCodes = @(0)
}

# 下载文件
if ($architecture -eq 'ARM64') {
    $url = $packageArgs.urlARM64
    $checksum = $packageArgs.checksumARM64
} else {
    $url = $packageArgs.url64bit
    $checksum = $packageArgs.checksum64
}

# 下载并安装
$fileLocation = Join-Path $toolsDir $exeName
Get-ChocolateyWebFile `
    -PackageName $packageName `
    -FileFullPath $fileLocation `
    -Url64bit $url `
    -Checksum64 $checksum `
    -ChecksumType 'sha256'

# 创建 shim (让可执行文件在 PATH 中可用)
Install-BinFile -Name 'rime-wanxiang-updater' -Path $fileLocation

Write-Host "Rime 万象输入法更新工具已成功安装！" -ForegroundColor Green
Write-Host "使用方法: 在命令行中运行 'rime-wanxiang-updater'" -ForegroundColor Cyan
