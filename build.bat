
@ECHO OFF

SET compile=%GOROOT%\pkg\tool\%GOOS%_%GOARCH%\compile.exe
SET link=%GOROOT%\pkg\tool\%GOOS%_%GOARCH%\link.exe

mklink /d src vendor
SET save=%GOPATH%
SET GOPATH=%cd%
go install github.com/cihub/seelog
go install github.com/fatih/set
go install golang.org/x/net/context

SET incdir=pkg\%GOOS%_%GOARCH%

@rem ----- cryptology -----
cd cryptology
%compile% -pack -o cryptology.a interface.go ca.go key.go cert.go
move /y cryptology.a ..
cd ..

@rem ----- router -----
cd router
%compile% -I ..\%incdir% -pack -o router.a router.go
move /y router.a ..
cd ..

@rem ----- meeting -----
cd meeting
%compile% -pack -o meeting.a proposal.go decide.go
move /y meeting.a ..
cd ..

@rem ----- core -----
cd core
%compile% -I ..\%incdir% -I .. -pack -o core.a types.go user.go signature.go broadcast.go interface.go proposal.go decide.go dig.go
move /y core.a ..
cd ..

%compile% -I . -I %cd%\%incdir% -pack -o main.a m.go
%link% -L . -L %cd%\%incdir% -o m.exe main.a

del /f /q main.a cryptology.a router.a meeting.a core.a

SET GOPATH=%save%
