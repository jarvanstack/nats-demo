# nats-demo

nats 各种使用方法 demo

## 安装

1.docker 安装

克隆本项目

cd nats-demo

docker-compose up -d


### client nats.go

go mod tidy



修改 main_test.go  的 nats 服务的地址,如果没有修改 docker-compose.yaml 的默认地址就不用修改

跑通其他 *_test.go 文件没有报错就行.


![](https://markdown-1304103443.cos.ap-guangzhou.myqcloud.com/2022-02-0420220523203442.png)


## reference

https://docs.nats.io/

https://github.com/nats-io/nats.go

https://pkg.go.dev/github.com/nats-io/nats.go


## QA

### 1. context deadline exceeded

JetStream 管理Stream的时候失败, 超时

nats 需要打开 js 模式

关掉 docker-compose 的 nats 

docker run --network host -p 4222:4222 nats -js