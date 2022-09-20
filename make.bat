@echo off
REM This is a Windows 10 Batch file for building dataset command
REM from the command prompt.
REM
REM It requires: go version 1.12.4 or better and the cli for git installed
REM
go version
echo Getting ready to build the dataset.exe

echo Using jq to extract version string from codemeta.json
jq .version codemeta.json > version.txt
SET /P DS_VERSION= < version.txt
DEL version.txt
echo Building version: %DS_VERSION%
echo package eprinttools > version.go
echo.  >> version.go
echo // Version of package >> version.go
echo const Version = %DS_VERSION% >> version.go
go build -o bin\doi2eprintxml.exe "cmd\doi2eprintxml\doi2eprintxml.go"
go build -o bin\ep3api.exe "cmd\ep3api\ep3api.go"
go build -o bin\epfmt.exe "cmd\epfmt\epfmt.go"
go build -o bin\eputil.exe "cmd\eputil\eputil.go"


echo Checking compile should see version number of dataset
.\bin\epfmt.exe -version

echo If OK, you can now copy the doi2eprintxml.exe to %USERPROFILE%\go\bin
echo.
echo       copy bin\doi2eprintxml.exe %USERPROFILE%\go\bin
echo.
