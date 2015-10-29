package front

const export_bat_tpl = `
@echo off

@WHERE bar
IF %ERRORLEVEL% NEQ 0 (
	ECHO bar is not found. downloading...
	powershell -command "$clnt = new-object System.Net.WebClient; $clnt.DownloadFile(\"{{.Info.HTTPEndpoint}}/win/bar.exe\", \"bar.exe\")"
)

bar up --log-level=DEBUG --endpoint={{.Info.JoinRPCEndpoints}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}} !bar*.bat !bar-spec*.json !bar.exe !desktop.ini
for /f %%i in ('bar spec export --upload --cc --endpoint={{.Info.JoinRPCEndpoints}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}} !bar*.bat !bar-spec*.json !bar.exe !desktop.ini') do set VAR=%%i

start {{.Info.HTTPEndpoint}}/spec/%VAR%

echo press any key...
pause >nul
`

const import_bat_tpl = `
@echo off

@WHERE bar
IF %ERRORLEVEL% NEQ 0 (
	ECHO bar is not found. downloading...
	powershell -command "$clnt = new-object System.Net.WebClient; $clnt.DownloadFile(\"{{.Info.HTTPEndpoint}}/win/bar.exe\", \"bar.exe\")"
)

for /f %%i in ('bar spec-import --squash --endpoint={{.Info.JoinRPCEndpoints}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}} {{.ID}}') do set VAR=%%i
bar down --log-level=DEBUG --endpoint={{.Info.JoinRPCEndpoints}} --chunk={{.Info.ChunkSize}} --pool={{.Info.PoolSize}} --buffer={{.Info.BufferSize}} %VAR%

echo press any key...
pause >nul
`
