# Week 04 <!-- omit in toc -->

## Table of contents <!-- omit in toc -->

- [Homework](#homework)
  - [Question](#question)
  - [Answer](#answer)
  - [Run](#run)
- [Build](#build)
- [Note](#note)

## Homework

### Question

按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

### Answer

假定有一個視頻串流平台，我們實現一個取得視頻詳細訊息的服務。包含視頻的名稱、觀看數以及是否被舉報。

Project layout

- api: proto files & generated code.
- bin: binaries. (After build)
- cmd: application code. (client & server)
- internal: code not for share.
  - biz: Business layer. (Business logics. Declare the repository interface.)
  - data: Repository layer. (Implement repository interfaces. Return data.)
  - domain: DO(domain object).
  - service Service layer. (Implement API interfaces and call businesses. DTO <-> DO.)

### Run

Server

```shell
$ ./bin/server
2020/12/18 22:25:06 [Program] Start.
2020/12/18 22:25:06 [Program] Wait routines.
2020/12/18 22:25:06 [Server] Start.
2020/12/18 22:25:06 [Notifier] Start.
2020/12/18 22:25:06 [Catcher] Start.
2020/12/18 22:25:06 Register signals: [interrupt terminated quit]
2020/12/18 22:25:11 [Video info svc] Get video info: ID= -1
2020/12/18 22:25:14 [Video info svc] Get video info: ID= 0
2020/12/18 22:25:17 [Video info svc] Get video info: ID= 1
2020/12/18 22:25:22 [Video info svc] Get video info: ID= 2
^C2020/12/18 22:25:27 Ask server to shutdown when capturing a registered signal: interrupt
2020/12/18 22:25:27 [Catcher] End.
2020/12/18 22:25:27 [Notifier] The context is done.
2020/12/18 22:25:27 [Notifier] End.
2020/12/18 22:25:27 [Server] End.
2020/12/18 22:25:27 Failed in an error group: Capture a registered signal: interrupt
2020/12/18 22:25:27 [Program] End.
```

Client

```shell
$ ./bin/client -1
2020/12/18 22:25:11 Failed to get video info: rpc error: code = Unknown desc = Invalid video id: -1
$ ./bin/client 0
2020/12/18 22:25:14 Failed to get video info: rpc error: code = Unknown desc = The video id=0 is reported. Cannot retrieve any information from it
$ ./bin/client 1
2020/12/18 22:25:17 Video info: name=Video-1, count=1
$ ./bin/client 2
2020/12/18 22:25:22 Video info: name=Video-2, count=2
```

## Build

Multiple targets

```shell
# Build client + server + proto + wire
$ make all

# Build client + server
$ make build
```

Single target

```shell
# Build one of client, server, proto or wire.
# Usage: make [client | server | proto | wire]
# Take proto as an example.
$ make proto
```

## Note

学习笔记 is [HERE](note.md).
