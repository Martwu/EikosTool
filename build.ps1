rsrc -manifest .\resources\thorgui.manifest -ico .\resources\main.ico  -o thorgui.syso
go build -ldflags="-H windowsgui"
