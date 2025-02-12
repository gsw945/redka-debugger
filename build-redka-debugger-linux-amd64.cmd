SET GOARCH=amd64
SET GOOS=linux
go build -o "%~dp0redka-debugger_%GOOS%_%GOARCH%" "%~dp0main.go"
