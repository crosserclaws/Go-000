# 微服务可用性设计 <!-- omit in toc -->

## Table of contents <!-- omit in toc -->

- [Resource](#resource)
- [隔离](#隔离)
  - [服务隔离](#服务隔离)
  - [轻重隔离](#轻重隔离)
  - [物理隔离](#物理隔离)
  - [Case study](#case-study)
- [超时控制](#超时控制)
- [过载保护](#过载保护)
  - [令牌桶算法](#令牌桶算法)
  - [漏桶算法](#漏桶算法)
  - [比較](#比較)
  - [过载保护](#过载保护-1)
- [限流](#限流)
- [降级](#降级)
- [重试](#重试)
- [负载均衡](#负载均衡)
- [最佳实践](#最佳实践)

## Resource

- [Machine Learning](https://www.bilibili.com/video/BV164411S78V)
- logging
  - https://dave.cheney.net/2015/11/05/lets-talk-about-logging
  - https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html
  - https://www.ardanlabs.com/blog/2017/05/design-philosophy-on-logging.html
  - https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern
- [类型转换-Struct转换](https://itician.org/pages/viewpage.action?pageId=1114345)
- [Copier](https://pkg.go.dev/github.com/jinzhu/copier)
- [seta/seta](https://github.com/seata/seata)

## 隔离

隔离，本质上是对系统或资源进行分割，从而实现当系统发生故障时能限定传播范围和影响范围，即发生故障后只有出问题的服务不可用，保证其他服务仍然可用。

- 服务隔离
  - 动静分离、读写分离
- 轻重隔离
  - 核心、快慢、热点
- 物理隔离
  - 线程、进程、集群、机房

[DDD 中的那些模式 — CQRS](https://zhuanlan.zhihu.com/p/115685384)

### 服务隔离

### 轻重隔离

### 物理隔离

线程隔离

主要通过线程池进行隔离，也是实现服务隔离的基础。把业务进行分类并交给不同的线程池进行处理，当某个线程池处理一种业务请求发生问题时，不会讲故障扩散和影响到其他线程池，保证服务可用。

对于 Go 来说，所有 IO 都是 Nonblocking，且托管给了 Runtime，只会阻塞Goroutine，不阻塞 M，我们只需要考虑 Goroutine 总量的控制，不需要线程模型语言的线程隔离。

- Go 不阻塞 thread。

// TODO

进程隔离

// TODO

### Case study

## 超时控制

超时控制，我们的组件能够快速失效(fail fast)，因为我们不希望等到断开的实例直到超时。没有什么比挂起的请求和无响应的界面更令人失望。这不仅浪费资源，而且还会让用户体验变得更差。我们的服务是互相调用的，所以在这些延迟叠加前，应该特别注意防止那些超时的操作。

- 网路传递具有不确定性。
- 客户端和服务端不一致的超时策略导致资源浪费。
- “默认值”策略。
- 高延迟服务导致 client 浪费资源等待，使用超时传递: 进程间传递 + 跨进程传递。

超时控制是微服务可用性的第一道关，良好的超时策略，可以尽可能让服务不堆积请求，尽快清空高延迟的请求，释放 Goroutine。

[Go 语言性能优化](https://cch123.github.io/perf_opt/)

## 过载保护

### 令牌桶算法

是一个存放固定容量令牌的桶，按照固定速率往桶里添加令牌。

token-bucket rate limit algorithm ([x/time/rate](https://pkg.go.dev/golang.org/x/time/rate))

### 漏桶算法

作为计量工具(The Leaky Bucket Algorithm as a Meter)时，可以用于流量整形(Traffic Shaping)和流量控制(TrafficPolicing)。

leaky-bucket rate limit algorithm ([/go.uber.org/ratelimit](https://pkg.go.dev/go.uber.org/ratelimit))

### 比較

漏斗桶/令牌桶确实能够保护系统不被拖垮, 但不管漏斗桶还是令牌桶, 其防护思路都是设定一个指标, 当超过该指标后就阻止或减少流量的继续进入，当系统负载降低到某一水平后则恢复流量的进入。但其通常都是被动的，其实际效果取决于限流阈值设置是否合理，但往往设置合理不是一件容易的事情。

- 集群增加机器或者减少机器限流阈值是否要重新设置?
- 设置限流阈值的依据是什么?
- 人力运维成本是否过高?
- 当调用方反馈429时, 这个时候重新设置限流, 其实流量高峰已经过了重新评估限流是否有意义?

这些其实都是采用漏斗桶/令牌桶的缺点, 总体来说就是太被动, 不能快速适应流量变化。
因此我们需要一种自适应的限流算法，即: 过载保护，根据系统当前的负载自动丢弃流量。

Terminology

- 極限壓測: 測試服務在極限狀態下的行為。

### 过载保护

计算系统临近过载时的峰值吞吐作为限流的阈值来进行流量控制，达到系统保护。

- 服务器临近过载时，主动抛弃一定量的负载，目标是自保。
- 在系统稳定的前提下，保持系统的吞吐量。
- 常见做法：利特尔法则
- CPU、内存作为信号量进行节流。
- 队列管理: 队列长度、LIFO。
- 可控延迟算法: [CoDel](https://blog.csdn.net/dog250/article/details/72849893)。
  - CoDel + BBR。

如何计算接近峰值时的系统吞吐？

- CPU: 使用一个独立的线程采样，每隔 250ms 触发一次。在计算均值时，使用了简单滑动平均去除峰值的影响。
  - 簡單滑動均值: Vt = Beta \* Vt-1 + (1-Beta) \* Theta t
  - 指數型滑動均值 (EMV)。
- Inflight: 当前服务中正在进行的请求的数量。
- Pass&RT: 最近5s，pass 为每100ms采样窗口内成功请求的数量，rt 为单个采样窗口中平均响应时间。

[Exponential Backoff And Jitter](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)

## 限流

## 降级

## 重试

## 负载均衡

## 最佳实践
