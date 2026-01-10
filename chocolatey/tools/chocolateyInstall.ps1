$ErrorActionPreference = 'Stop'

$packageName = 'rime-wanxiang-updater'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$version = '$VERSION$'

# 根据系统架构选择正确的可执行文件
$architecture = $env:PROCESSOR_ARCHITECTURE
$exeName = switch ($architecture) {
    'AMD64' { "rime-wanxiang-updater-windows-amd64.exe" }
    'ARM64' { "rime-wanxiang-updater-windows-arm64.exe" }
    default { "rime-wanxiang-updater-windows-amd64.exe" }
}

# 根据架构选择下载 URL 和校验和
if ($architecture -eq 'ARM64') {
    $url = "https://github.com/czyt/rime-wanxiang-updater/releases/download/v$version/rime-wanxiang-updater-windows-arm64.exe"
    $checksum = '$CHECKSUMARM64$'
} else {
    # 默认使用 AMD64 版本
    $url = "https://github.com/czyt/rime-wanxiang-updater/releases/download/v$version/rime-wanxiang-updater-windows-amd64.exe"
    $checksum = '$CHECKSUM64$'
}

# 下载并安装
$fileLocation = Join-Path $toolsDir $exeName
Get-ChocolateyWebFile `
    -PackageName $packageName `
    -FileFullPath $fileLocation `
    -Url $url `
    -Checksum $checksum `
    -ChecksumType 'sha256'

# 创建 shim (让可执行文件在 PATH 中可用)
Install-BinFile -Name 'rime-wanxiang-updater' -Path $fileLocation

Write-Host "Rime 万象输入法更新工具已成功安装！" -ForegroundColor Green
Write-Host "使用方法: 在命令行中运行 'rime-wanxiang-updater'" -ForegroundColor Cyan
