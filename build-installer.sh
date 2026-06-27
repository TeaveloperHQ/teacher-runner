#!/usr/bin/env bash
# 설치 프로그램(teaveloper-runner-setup.exe) 빌드.
# 저장소 루트에서 실행. exe 를 먼저 빌드한 뒤 NSIS 로 설치본을 만든다.
# 리눅스(포털/CI)에서 동작: makensis 만 있으면 된다(예: apt-get install nsis).
set -euo pipefail

# 1) 러너 exe (CGO 없이 윈도우 크로스컴파일, 아이콘 .syso 자동 임베드)
./build.sh dist/teaveloper-runner.exe

# 2) 설치 프로그램 (NSIS). File/아이콘 경로는 저장소 루트 기준이라 루트에서 실행.
if ! command -v makensis >/dev/null 2>&1; then
  echo "makensis 가 필요합니다 (예: apt-get install -y nsis)." >&2
  exit 1
fi
makensis installer/installer.nsi

echo "built: dist/teaveloper-runner-setup.exe"
