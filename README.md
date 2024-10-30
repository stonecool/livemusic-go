# livemusic-go

swag.exe init -g .\cmd\server\main.go -o .\docs\swagger\

.\docs\protoc.exe -I=docs --go_out=. docs\cmd.proto

.\docs\protoc.exe --js_out=import_style=commonjs,binary:. .\docs\cmd.proto

gofmt -w . 

### TODO

- 异常终止，资源释放，windows,linux
- 是否将chrome替换为浏览器

### 设计

- chrome实例
  - 创建本地实例
  - 绑定远程实例
  - 恢复实例，本地、远程
  - 获取当前所有实例
  - 删除某个实例