SET GOARCH=amd64
SET GOOS=windows
go build -o "%~dp0redka-debugger_%GOARCH%.exe" "%~dp0main.go"
