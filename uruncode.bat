
echo off
cls
FOR /L %%A IN (1,1,2000) DO (
rem go run mainprog.go dataobj.go 
rem go run mygame.go
rem -race 
go run LunaLander.go anim.go
pause
)
