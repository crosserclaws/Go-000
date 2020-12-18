# Go 工程化实践 <!-- omit in toc -->

## Table of contents <!-- omit in toc -->

- [Resource](#resource)
- [工程项目结构](#工程项目结构)
  - [Standard Go Project Layout](#standard-go-project-layout)
  - [Kit Project Layout](#kit-project-layout)
  - [Service Application Project Layout](#service-application-project-layout)
  - [Service Application Project](#service-application-project)
- [API设计](#api设计)
  - [gRPC](#grpc)
  - [API project](#api-project)
  - [API layout](#api-layout)
  - [API compatibility](#api-compatibility)
  - [API Naming Conventions](#api-naming-conventions)
  - [API Primitive Fields](#api-primitive-fields)
  - [API Errors](#api-errors)
  - [API Design](#api-design)
- [配置管理](#配置管理)
  - [Configuration](#configuration)
  - [Functional options](#functional-options)
  - [Hybrid APIs](#hybrid-apis)
  - [Configuration & APIs](#configuration--apis)
  - [Configuration best practice](#configuration-best-practice)
- [包管理](#包管理)
  - [go mod](#go-mod)
- [测试](#测试)

## Resource

- [Standard Go Project Layout]
- [Google API Design Guide CN]

[Standard Go Project Layout]: https://github.com/golang-standards/project-layout
[Google API Design Guide CN]: https://www.bookstack.cn/read/API-design-guide/API-design-guide-README.md

## 工程项目结构

Code 的組織與資源如何分層及注入。

- [Standard Go Project Layout]

### Standard Go Project Layout

- /cmd
  - Main applications for this project.
  - Example
    - /cmd/prometheus/main.go
    - /cmd/kubernetes/kube-proxy/proxy.go
- /internal
  - Private application and library code.
  - This is the code you don't want others importing in their applications or libraries.
  - Note that this layout pattern is enforced by the Go compiler itself.
- /pkg
  - Library code that's ok to use by external applications.

### Kit Project Layout

[Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html)

> To this end, the Kit project is not allowed to have a vendor folder. If any of packages are dependent on 3rd party packages, they must always build against the latest version of those dependencies.

kit 项目必须具备的特点

- 统一。
- 标准库方式布局。
- 高度抽象。
  - 定義好抽象，有簡單的實現或是交由插件實現。
- 支持插件。

### Service Application Project Layout

- 目錄
  - /api
    - API 協議定義目錄，如 protobuf 文件以及生成的 go 文件。
    - 通常直接在 proto 中撰寫 doc。
  - /configs
  - /test
    - Go 會忽略以 `.` 或 `_` 開頭的目錄或文件。
- 注意事項
  - 不應該包含 /src。

### Service Application Project

Monorepo vs. Polyrepo ([Github](https://github.com/joelparkerhenderson/monorepo_vs_polyrepo))

- Monorepo
  - 一個 project 中放多個微服務的 app。
  - App 目录内的每个微服务按照自己的全局唯一名称，来建立目录。
    - Ex: account.service.vip。
      - 業務.服務.子服務。
- Polyrepo
  - 一個 group 裡多個 projects，每個 project 對應一個 app。

微服务中的 app 服务类型分为4类: interface, service, job, admin。

- interface
  - 對外的 BFF 服務，接受來自用戶的請求。
    - Ex: 暴露的 HTTP/gRPC interface。
- service
  - 純粹對內的微服務，只接受內部其他服務或 gateway 的請求。
- admin
  - 不同於 service，面向運營管理的服務。
  - 通常權限更高，隔離帶來更好的代碼級別安全。
  - 破壞共享模式，即 2 個服務共享一個 DB/Cache。
- job
  - 流式任務處理的服務，上游一般依賴 message broker。
- task
  - 定時任務，類似 cronjob，部署到 task 託管平台之中。

Layout 演進

- v1
  - /service
    - api
      - 定義了 API proto file 和生成的 stub code。
      - 所生成的 interface 會在 service 中實現。
    - cmd
    - configs
    - internal
      - model
        - 對應 Java POJO 的概念。
        - 存放對應存儲層的結構體，是對存儲的一一映射。
          - Ex: 映射 MySQL 中 table 的結構體。
      - dao
        - Data access object.
        - 依賴 model，面向 table 的概念。
        - 數據讀寫層，DB/Cache 在這層統一處理。
          - 包含 cache miss 的處理。
      - service
        - 依賴 DAO 層的抽象。
        - 組合各種數據訪問來構建業務邏輯。
      - server
        - 依賴 service。
        - 依賴 proto 定義的服務作為入參，提供快捷的啟動服務全局方法。
        - 此層可以消除。
          - 在 main 中使用 gRPC 將 service 註冊給 transport server。
  - 會有人為求快速，DTO 直接在業務邏輯層/甚至 DAO 層中使用。
    - 因為 service 的方法簽名隱式實現了 API 的接口定義。
    - 缺乏 DTO -> DO 的對象轉換。
  - 參考 [kratos-demo](https://github.com/bilibili/kratos-demo)
  - Terminology
    - DTO (Data Transfer Object)
      - An object that carries data between processes.
      - 數據傳輸對象。概念源於 J2EE 的設計模式。
      - 在這裡泛指展示層/API 層與服務層(業務邏輯層)之間的數據傳輸對象。
        - api <-> DTO <-> server <-deep copy->  service
        - [What is Data Transfer Object?](https://stackoverflow.com/questions/1051182/)
    - DO (Domain Object)
      - 領域對象，就是從現實世界中抽象出來的有形或無形的業務實體。
    - [初识失血，贫血，充血，胀血四种模型](https://www.jianshu.com/p/d0d1c8c0de07)
      - 上述定義似乎有待考證，詳見 [设计概念的统一语言](http://zhangyi.xyz/ubiquitous-language-of-design-concept/)。
- v2
  - /internal
    - 為了避免同業務下有人跨目錄引用了內部的 biz, data, service 等內部 struct。
    - biz
      - 業務邏輯的組裝層。
        - 類似 DDD 的 domain 層。
        - repo 接口在這裡定義，使用依賴倒置的原則。
    - data
      - 業務數據訪問，包含 cache, db 等封裝，實現了 biz 的 repo 接口。
        - data 類似 DDD 的 repo。
        - 我們可能會把 data 與 dao 的概念搞混；data 偏重業務的含義，它所要做的是將領域對象重新拿出來，這裡我們去掉了 DDD 的 infra 層。
    - service
      - 實現 API 定義的服務層。
        - 類似 DD 的 application 層。
      - 處理 DTO 到領域實體的轉換 (DTO -> DO)，同時協同各類 biz 交互但是不應處理複雜邏輯。
        - 只有組合邏輯，沒有業務邏輯。
  - Terminology
    - PO (Persistent Object)
      - 持久化對象。
      - 跟持久層的數據結構形成一一對應的映射關係。
        - 如果持久層是 RDB，那麼數據表中的每個(或 N 個) 字段就對應 PO 的一個(或 N 個) 屬性。
      - [facebook/ent](https://github.com/facebook/ent)
  - 參考 [service-layout](https://github.com/go-kratos/service-layout)

Dependency inversion principle

- 上層不應該依賴於下層模塊，它們共同依賴於一個抽象。
- 抽象不能依賴於具象，具象依賴於抽象。

Lifecycle

需要考慮服務應用的對象初始化以及生命週期的管理。

- 所以 HTTP/gRPC 以阿里的前置資源初始化。之後再啟動監聽服務。
  - 包括 data, biz, service。
- 可以使用 wire 來管理所有資源的依賴注入。
- Reference
  - [google/wire](https://github.com/google/wire)
  - [Kratos-app.go](https://github.com/go-kratos/kratos/blob/v2/app.go)
    - 參考 [uber-go/fx](https://github.com/uber-go/fx)，並加以改進。
  - [Kratos-service.go-RegisterService](https://github.com/go-kratos/kratos/blob/v2/transport/http/service.go)

Wire

[Compile-time Dependency Injection With Go Cloud's Wire](https://blog.golang.org/wire)

- 手刻資源的初始化和關閉是繁瑣且易出錯的。
- DI 的思想搭配 wire 工具，可以生成靜態的代碼，方便診斷與查看，而非 runtime 時利用 reflection 實現。

## API设计

### gRPC

// TODO

### API project

統一的 API 倉庫

- 方便統一檢索和規範 API。
- API，方便跨部門協作。
- 基於 git 做版本管理。
- Lint, design review, diff。
- 權限管理。
  - 目錄 owners。

### API layout

// TODO

### API compatibility

向後兼容(非破壞性)的修改

- 給 API 服務定義添加 API interface。
- 給 Request 添加字段。
- 給 Response 添加字段。

非向後兼容(破壞性)的修改

- 刪除或重命名服務、字段、方法或枚舉值。
- 修改字段的類型。
- 修改現有請求的可見行為。
- 給資源消息添加 讀取/寫入 字段。

### API Naming Conventions

包名为应用的标识(APP_ID)，用于生成 gRPC 请求路径，或者 proto 之间进行引用 Message。

```txt
# RequestURL: /<package_name>.<version>.<service_name>/{method}
package <package_name>.<version>;

# Example
package google.example.library.v1;
```

### API Primitive Fields

gRPC 默認使用 Protobuf v3 格式，默認都是 optional 字段。

- 如果沒有賦值的字段，默認為類型字段的默認值。
- Protobuf v3 中可使用 [Wrapper](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/wrappers.proto)。
  - Wrapper 類型的字段，即包裝一個 message，使用時變為指針。

### API Errors

使用一小組標準錯誤配合大量資源。

- [Error details](https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto)

錯誤傳播。

- 如果改 API service 依賴於其他 service，則不應盲目地將錯誤傳播到 client side，需要翻譯錯誤。

- 隱藏實現詳細信息和機密信息。
- 调整负责该错误的一方。例如，从另一个服务接收 INVALID_ARGUMENT 错误的服务器应该将 INTERNAL 传播给它自己的调用者。

全局錯誤碼。

- 全局错误码，是松散、易被破坏契约的，基于我们上述讨论的，在每个服务传播错误的时候，做一次翻译，这样保证每个服务 + 错误枚举，应该是唯一的，而且在 proto 定义中是可以写出来文档的。
- [Kratos v2 errors](https://github.com/go-kratos/kratos/tree/v2/errors)
- [Kratos examples](https://github.com/go-kratos/kratos/tree/v2/examples/kratos-demo/api/kratos/demo/errors)
- Service error --> gRPC error --> Service error
  - Error 放在 metadata 中。
  - // TODO
  - https://github.com/go-kratos/kratos/blob/v2/examples/kratos-demo/api/kratos/demo/v1/greeter_grpc.pb.go

### API Design

- // TODO
- FieldMask
- [Google API Design Guide CN]

## 配置管理

### Configuration

- 環境變量。
  - 從平台注入，不建議放在 config file 中。
  - Example: Region, Zone, Cluster, Environment, Color, Discovery, AppID, Host etc.
- 靜態配置。
  - 不建議 on-the-fly 的修改配置，建議還是重啟服務。
- 動態配置。
  - 用基礎類型做一些在線的開關。
  - [expvar](https://pkg.go.dev/expvar)
- 全局配置。

### Functional options

- [Self-referential functions and the design of options](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html)
- [Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

### Hybrid APIs

// TODO

### Configuration & APIs

// TODO
Tool: Margin note

### Configuration best practice

// TODO

## 包管理

### go mod

// TODO

## 测试

Unit test 基本要求

- 快速。
- 環境一致。
- 任意順序。
- 並行。

Book

- Google 測試之道
- 微軟測試之道
