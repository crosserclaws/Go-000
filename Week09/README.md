# Week 09 <!-- omit in toc -->

## Table of contents <!-- omit in toc -->

- [Homework](#homework)
  - [Question](#question)
  - [Answer](#answer)
  - [Run](#run)
- [Note](#note)

## Homework

### Question

实现一个 tcp server ，用两个 goroutine 读写 conn，两个 goroutine 通过 chan 可以传递 message，能够正确退出。

### Answer

实现一个 echo server 作为范例。

- 每个 connection 会由 2 个 go routines (reader/writer) 来处理。
  - 两个 routines 之间经由 channel 来传递 message。
  - 其中一个 routine 出现错误时会通知另一方结束。
- Server 注册监听的 signals。
  - 收到 signals 会通知 routines 结束。

### Run

Server

```shell
$ go run main.go
2021/01/28 21:33:40 Server starts.
2021/01/28 21:33:40 Register signals: [interrupt terminated quit]
2021/01/28 21:33:43 Writer[1] starts.
2021/01/28 21:33:43 Reader[1] starts.
2021/01/28 21:33:47 Reader[1] receives a msg: Hi, I am Alice.
2021/01/28 21:33:47 Writer[1] echos a response: Hi, I am Alice.
2021/01/28 21:33:49 Writer[2] starts.
2021/01/28 21:33:49 Reader[2] starts.
2021/01/28 21:33:56 Reader[2] receives a msg: Hello, I am Bob.
2021/01/28 21:33:56 Writer[2] echos a response: Hello, I am Bob.
^C2021/01/28 21:34:08 Ask server to shutdown when capturing a registered signal: interrupt
2021/01/28 21:34:08 Server waits for closing connections.
2021/01/28 21:34:08 Writer[1] gets the program is going to shutdown.
2021/01/28 21:34:08 Writer[1] ends.
2021/01/28 21:34:08 Writer[2] gets the program is going to shutdown.
2021/01/28 21:34:08 Reader[1] fails to read: read tcp 127.0.0.1:8080->127.0.0.1:50735: use of closed network connection
2021/01/28 21:34:08 Reader[1] ends.
2021/01/28 21:34:08 Writer[2] ends.
2021/01/28 21:34:08 Reader[2] fails to read: read tcp 127.0.0.1:8080->127.0.0.1:50737: use of closed network connection
2021/01/28 21:34:08 Reader[2] ends.
2021/01/28 21:34:08 Server ends.
```

Client 1

```shell
$ nc -4 localhost 8080
Hi, I am Alice.
Hi, I am Alice.[Notification] Server is going to shutdown.
```

Client 2

```shell
$ nc -4 localhost 8080
Hello, I am Bob.
Hello, I am Bob.[Notification] Server is going to shutdown.
```

## Note

学习笔记 is [HERE](note.md).
