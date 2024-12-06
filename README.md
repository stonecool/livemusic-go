# livemusic-go

swag.exe init -g .\cmd\server\main.go -o .\docs\swagger\

.\docs\protoc.exe -I=docs --go_out=. docs\cmd.proto

.\docs\protoc.exe --js_out=import_style=commonjs,binary:. .\docs\cmd.proto

gofmt -w . 

go test -v .\internal\account\
go test  .\internal\account\account_test.go  .\internal\account\account.go
go test -v .\internal\chrome\... -v -cover 三个点，所有目录和子目录

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

  internal/
  ├── account/                 # 账号领域
  │   ├── model.go            # 数据模型
  │   ├── repository.go       # 数据访问层
  │   ├── validator.go        # 验证逻辑
  │   ├── factory.go          # 工厂方法
  │   ├── account.go          # 业务对象
  │   └── interface.go        # 接口定义
  │
  ├── task/                   # 任务领域
  │   ├── model.go           
  │   ├── repository.go
  │   ├── validator.go
  │   ├── factory.go
  │   ├── task.go
  │   └── interface.go
  │
  ├── chrome/                 # Chrome 实例管理
  │   ├── instance.go
  │   └── pool.go
  │
  ├── scheduler/              # 任务调度
  │   └── scheduler.go
  │
  └── service/               # 业务服务层
      └── crawler.go         # 协调不同领域的服务