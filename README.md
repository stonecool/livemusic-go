# livemusic-go

swag.exe init -g .\cmd\server\main.go -o .\docs\swagger\

.\docs\protoc.exe -I=docs --go_out=. docs\cmd.proto

.\docs\protoc.exe --js_out=import_style=commonjs,binary:. .\docs\cmd.proto

gofmt -w . 

### TODO

- 异常终止，资源释放，windows,linux