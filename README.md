# livemusic-go

swag.exe init -g .\cmd\server\main.go -o .\docs\swagger\

.\docs\protoc.exe -I=docs --go_out=. docs\cmd.proto