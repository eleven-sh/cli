#!/usr/bin/env pwsh
# Copyright 2018-2022 the Deno authors. All rights reserved. MIT license.
# TODO(everyone): Keep this script simple and easily auditable.

$ErrorActionPreference = 'Stop'

if ($v) {
  $Version = "v${v}"
}
if ($Args.Length -eq 1) {
  $Version = $Args.Get(0)
}

$Releases = "https://api.github.com/repos/eleven-sh/cli/releases"
$LatestVersion = (Invoke-WebRequest $Releases | ConvertFrom-Json)[0].tag_name

$Arch = $env:PROCESSOR_ARCHITECTURE

$ElevenInstall = $env:ELEVEN_INSTALL
$BinDir = if ($ElevenInstall) {
  "${ElevenInstall}\bin"
} else {
  "${Home}\.eleven\bin"
}

$ElevenTar = "${BinDir}\eleven.tar.gz"
$ElevenExe = "${BinDir}\eleven.exe"

$DownloadUrl = if (!$Version) {
  "https://github.com/eleven-sh/cli/releases/latest/download/cli_$($LatestVersion.substring(1))_windows_$($Arch.ToLower()).tar.gz"
} else {
  "https://github.com/eleven-sh/cli/releases/download/${Version}/cli_$($Version.substring(1))_windows_$($Arch.ToLower()).tar.gz"
}

if (!(Test-Path $BinDir)) {
  New-Item $BinDir -ItemType Directory | Out-Null
}

curl.exe -Lo $ElevenTar $DownloadUrl

tar.exe xf $ElevenTar -C $BinDir

Remove-Item $ElevenTar

$User = [System.EnvironmentVariableTarget]::User
$Path = [System.Environment]::GetEnvironmentVariable('Path', $User)
if (!(";${Path};".ToLower() -like "*;${BinDir};*".ToLower())) {
  [System.Environment]::SetEnvironmentVariable('Path', "${Path};${BinDir}", $User)
  $Env:Path += ";${BinDir}"
}

Write-Output "Eleven was installed successfully to ${ElevenExe}"
Write-Output "Run 'eleven --help' to get started"
Write-Output "Stuck? Open a new issue at https://github.com/eleven-sh/cli/issues/new"