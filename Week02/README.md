# Week 02 <!-- omit in toc -->

## Table of contents <!-- omit in toc -->

- [Homework](#homework)
  - [Question](#question)
  - [Answer](#answer)
  - [Run](#run)
- [Note](#note)
  - [Error](#error)
  - [Error vs Exception](#error-vs-exception)
  - [Error handling strategies](#error-handling-strategies)
    - [Sentinel Error](#sentinel-error)
    - [Error Type](#error-type)
    - [Opaque Error](#opaque-error)

## Homework

### Question

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

### Answer

一般情况来说，与标准库和第三方库互动或是自有基础库，大多数情况下是可以直接 Wrap error 后抛给上层的。但某些情况下不会选择这样做，比如

- 需要隐藏底层的错误，避免上层直接依赖特定实现的错误，导致底层更换实作时出现问题。
  - 比如 DAO 层直接返回 sql.ErrNoRows 或是直接 Wrap，上层就有可能是直接判断是否为 sql.ErrNoRows。此时可以选择返回自定义类型的错误来封装，使上层相依在 DAO 层定义的错误，这样 DAO 层更换存储时就不会出现问题。举例来说，对于一个使用者 ID 查询使用者资讯的 API，就可以这样实现。
- 业务类型允许返回降级资料。
  - 比如列出所有使用者的 API，是可能出现没有使用者的情况，此时可处理错误，对上层返回空值。

### Run

```bash
$ go run main.go
[GIN-debug] GET    /user/:id                 --> github.com/crosserclaws/Go-000/Week02/api.GetUserInfoByID (3 handlers)
[GIN-debug] GET    /users                    --> github.com/crosserclaws/Go-000/Week02/api.GetUsers (3 handlers)
[GIN-debug] GET    /ping                     --> main.main.func1 (3 handlers)
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080
...
```

## Note

学习笔记

### Error

在 Go 中，error 是一個普通的 interface，普通的值。

```go
// Ref: https://golang.org/pkg/builtin/#error
type error interface {
    Error() string
}
```

一般會使用 errors.New() 來創建一個 error。

- 注意是回傳 pointer，所以比較時即使值相同也會比較是否為相同的 struct。
  - 可以作為 Sentinel Error 使用。

```go
// Ref: https://golang.org/src/errors/errors.go

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
    return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
```

Std library 中常見一些自定義的 error。

```go
// Ref: https://golang.org/pkg/bufio/
var (
    ErrInvalidUnreadByte = errors.New("bufio: invalid use of UnreadByte")
    ErrInvalidUnreadRune = errors.New("bufio: invalid use of UnreadRune")
    ErrBufferFull        = errors.New("bufio: buffer full")
    ErrNegativeCount     = errors.New("bufio: negative count")
)
```

### Error vs Exception

Go 的並不引入 exception 進行錯誤處理，而是支持多參數回傳值，同時返回結果與錯誤值。

- 如果需要使用結果，必須先檢查錯誤值；反過來，如果不在乎結果則不需要處理。
- 對於那些不可恢復的程式錯誤則使用 panic。
  - Ex: Stack overflow, Out of index, 無法恢復的環境問題。
- Summary
  - 簡單。
  - Plan for failure, not success.
  - 沒有隱藏的控制流。
  - 完全交由使用者控制 error。
  - `Error are values`。

### Error handling strategies

#### Sentinel Error

Sentinel Error 是一種預定義的錯誤。

- 源於 programming 中使用一個特定值來表示無法進行進一步處理的作法。
- 在 Go 中有時使用此概念來表示特定錯誤，如 io.EOF。
- 建議: 盡量避免使用 Sentinel Error。
  - 因為這是最不靈活的錯誤處理策略。
    - 調用方必須使用 == 與預先聲明的值進行比較，也就無法對一個錯誤添加更多的資訊。
    - Sentinel Error 會成為 API 的公共部分並產生依賴關係。

#### Error Type

Error Type 是實現了 error interface 的自定義類型。

- 相較於單純的錯誤值，可以封裝更多並提供更多錯誤資訊。
  - 使用者可以藉由 Type Assertion 將 error 轉換成對應類型來獲取更多資訊。
  - Ex: os.PathError。
- 建議: 盡量避免使用 Error Type。
  - 雖然比 Sentinel Error 好，能提供更多的錯誤資訊，但有許多相同的問題。
  - 調用方要使用 Type Assertion，就需要 public Error Type，導致強耦合。

#### Opaque Error

Assert error for behavior, not type.