$ErrorActionPreference = 'Stop'

$packageName = 'rime-wanxiang-updater'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

# 移除 shim
Uninstall-BinFile -Name 'rime-wanxiang-updater'

Write-Host "Rime 万象输入法更新工具已成功卸载！" -ForegroundColor Green
