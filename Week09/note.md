# Ch09 网络编程 <!-- omit in toc -->

## Table of contents <!-- omit in toc -->

- [Demo](#demo)
  - [top](#top)
  - [nmon](#nmon)
  - [nload](#nload)
  - [tcpflow](#tcpflow)
  - [ifconfig](#ifconfig)
  - [netstat](#netstat)
  - [ss](#ss)
  - [vmstat](#vmstat)
  - [iostat](#iostat)
  - [iotop](#iotop)
  - [pid](#pid)
  - [perf](#perf)
  - [ethtool](#ethtool)
  - [Q&A](#qa)
- [网络通信协议](#网络通信协议)
  - [HTTP 超文本传输协议-演进](#http-超文本传输协议-演进)
- [通过 Go 实现网络编程](#通过-go-实现网络编程)
  - [I/O 模型](#io-模型)
- [Goim 长连接 TCP 编程](#goim-长连接-tcp-编程)
  - [Overview](#overview)
  - [负载均衡](#负载均衡)
  - [心跳保活机制](#心跳保活机制)
  - [唯一 ID 設計](#唯一-id-設計)
  - [IM 私信系統](#im-私信系統)

## Demo

- [Linux Performance](http://www.brendangregg.com/linuxperf.html)
- [Linux Extended BPF (eBPF) Tracing Tools](http://www.brendangregg.com/ebpf.html)

### top

Hardware Interrupt (HI)

Software Interrupt (SI)

理想上盡量讓 SI 均衡

- [MYSQL数据库网卡软中断不平衡问题及解决方案](http://blog.yufeng.info/archives/2037)
- ByPass Kernel
- DPDK

### nmon

注意一下 Kernel 的 Context Switch。

### nload

### tcpflow

### ifconfig

可以看網卡有沒有錯誤，丟包。

### netstat

連線數多的時候容易造成負載過大，盡量用 `netstat -s`。

### ss

- ss -s

### vmstat

### iostat

### iotop

### pid

- lsof - p $pid
- strace -p $pid

### perf

- perf top

### ethtool

### Q&A

- conntrack-table
- LinkedList Head/Tail hlist

## 网络通信协议

- Socket
  - 接口抽象層。
- TCP / UDP
  - TCP 面向連接 (可靠)
  - UDP 無連接 (不可靠)

### HTTP 超文本传输协议-演进

HTTP2

- 头部压缩，通过 HPACK 压缩格式。

## 通过 Go 实现网络编程

### I/O 模型

Linux 下主要的 IO 模型

- Blocking IO
- Nonblocking IO
- IO multiplexing
  - select 有限制，新的程式比較多使用 epoll。
- Signal-driven IO
- Asynchronous IO
  - Redis 有用。

## Goim 长连接 TCP 编程

[GOIM](https://goim.io/)

### Overview

Caveat

- TCP port 是 16-bit unsigned，上限 65535。
  - 新增虛擬網卡可解決。
- ulimit 是 fd 數量上限，不同於 port 上限的問題。
- Connection tracking (“conntrack”)。

Terminology

- 5-Tuple ([RFC 6146](https://www.ietf.org/rfc/rfc6146.txt))
  - The tuple (source IP address, source port, destination IP address, destination port, transport protocol).
  - A 5-tuple uniquely identifies a UDP/TCP session.

### 负载均衡

### 心跳保活机制

高效维持长连接方案

- 进程保活（防止进程被杀死）
- 心跳保活（阻止 NAT 超时）
- 断线重连（断网以后重新连接网络）

自适应心跳时间

- 心跳可选区间。
  - min=60s, max=300s.
- 心跳增加步长。
  - step=30s.
- 心跳周期探测。
  - success=current + step.
  - fail=current - step.

### 唯一 ID 設計

Snowflake

[Go Snowflake](https://github.com/Terry-Mao/gosnowflake)

### IM 私信系統

[IM消息ID技术专题(一)：微信的海量IM聊天消息序列号生成实践（算法原理篇）](http://www.52im.net/thread-1998-1-1.html)
