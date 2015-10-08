@echo off

barc -log-level=DEBUG up -endpoint=http://spin.local:3000/v1 !bar*.cmd !bar-spec*.json

for /f %%i in ('barc -log-level=DEBUG spec-export -endpoint=http://spin.local:3000/v1 -upload !bar*.cmd !bar-spec*.json') do set VAR=%%i
REM barc -log-level=DEBUG spec-export -endpoint=http://spin.local:3000/v1 -upload !bar*.cmd !bar-spec*.json
start http://spin.local:3000/v1/spec/%VAR%
echo press any key...
pause >nul