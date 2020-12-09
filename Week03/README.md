# Week 03 <!-- omit in toc -->

## Table of contents <!-- omit in toc -->

- [Homework](#homework)
  - [Question](#question)
  - [Answer](#answer)
    - [Interrupt](#interrupt)
    - [Deadline](#deadline)
    - [Listen](#listen)
- [Note](#note)

## Homework

### Question

基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

### Answer

實現上，由一個 10 seconds timeout 的 context 產生一個 error groups，其包含下列 go routines:

- Server: 啟動 server，失敗則回傳 error。
- Notifier: 當有 group context 取消時，通知 server shutdown。
- Catcher: 取得 signal 後，回傳 error。

#### Interrupt

```bash
2020/12/09 23:38:16 [Program] Start.
2020/12/09 23:38:16 [Program] Wait routines.
2020/12/09 23:38:16 [Notifier] Start.
2020/12/09 23:38:16 [Catcher] Start.
2020/12/09 23:38:16 [Server] Start.
2020/12/09 23:38:16 Register signals: [interrupt terminated quit]
^C2020/12/09 23:38:18 Ask server to shutdown when capturing a registered signal: interrupt
2020/12/09 23:38:18 [Catcher] End.
2020/12/09 23:38:18 [Notifier] The context is done.
2020/12/09 23:38:18 [Notifier] End.
2020/12/09 23:38:18 [Server] End.
2020/12/09 23:38:18 Failed in an error group: Capture a registered signal: interrupt
2020/12/09 23:38:18 [Program] Ensure the shutdown is graceful.
2020/12/09 23:38:18 Do registered shutdown.
2020/12/09 23:38:18 [Program] End.
```

#### Deadline

```bash
$ go run main.go
2020/12/09 23:38:25 [Program] Start.
2020/12/09 23:38:25 [Program] Wait routines.
2020/12/09 23:38:25 [Notifier] Start.
2020/12/09 23:38:25 [Server] Start.
2020/12/09 23:38:25 [Catcher] Start.
2020/12/09 23:38:25 Register signals: [interrupt terminated quit]
2020/12/09 23:38:35 [Notifier] The context is done.
2020/12/09 23:38:35 [Catcher] The context is done.
2020/12/09 23:38:35 [Catcher] End.
2020/12/09 23:38:35 [Notifier] End.
2020/12/09 23:38:35 Do registered shutdown.
2020/12/09 23:38:35 [Server] End.
2020/12/09 23:38:35 Failed in an error group: context deadline exceeded
2020/12/09 23:38:35 [Program] Ensure the shutdown is graceful.
2020/12/09 23:38:35 [Program] End.
```

#### Listen

```bash
2020/12/09 23:36:45 [Program] Start.
2020/12/09 23:36:45 [Program] Wait routines.
2020/12/09 23:36:45 [Notifier] Start.
2020/12/09 23:36:45 [Server] Start.
2020/12/09 23:36:45 [Catcher] Start.
2020/12/09 23:36:45 Register signals: [interrupt terminated quit]
2020/12/09 23:36:45 [Server] End.
2020/12/09 23:36:45 [Catcher] The context is done.
2020/12/09 23:36:45 [Catcher] End.
2020/12/09 23:36:45 [Notifier] The context is done.
2020/12/09 23:36:45 [Notifier] End.
2020/12/09 23:36:45 Failed in an error group: listen tcp :8080: bind: address already in use
2020/12/09 23:36:45 [Program] Ensure the shutdown is graceful.
2020/12/09 23:36:45 Do registered shutdown.
2020/12/09 23:36:45 [Program] End.
```

## Note

学习笔记
