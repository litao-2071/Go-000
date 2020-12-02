# **HomeWork**



- 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

```shell
应该wrap这个error，抛给上层。因为我们要告诉调用者完整的堆栈信息，error只处理一次，不应该在dao层去打log，应该抛给上层
```



# ERROR

## 1.Error vs Exception

### 1.1 Error本质

- Error本质上是一个接口

```go
type error interface{
    Error() string
}
```

- 经常使用errors.New()来返回一个error对象

例如标准库中的error定义, 通过bufio 前缀带上上下文信息

```go
var (
	ErrInvalidUnreadByte = errors.New("bufio: invalid use of UnreadByte")
	ErrInvalidUnreadRune = errors.New("bufio: invalid use of UnreadRune")
	ErrBufferFull        = errors.New("bufio: buffer full")
	ErrNegativeCount     = errors.New("bufio: negative count")
)
```

- errors.New()是返回的error对象的指针

为了防止在error比较时，因为error内部内容定义相同导致两个不同类型的error误判相等

### 1.2 Error和Exception的区别

**各语言演进：**

- C: 一般传入指针，通过返回的int值判断成功还是失败
- C++: 无法知道抛出的什么异常，是否抛出异常（只能通过文档）
- JAVA: 需要抛出异常则方法的所有者必须声明，调用者也必须处理。处理方式、轻重程度都由调用者区分。
- GO: 不引入exception，采用多参数返回，一般最后一个返回值都是error

**error and panic:**

- error是我们预期内的错误
- panic表示不可恢复的程序错误，比如下标越界这种，意味着整个程序退出

**使用error代替exception的好处**：

- 简单
- 考虑失败不是成功
- 没有隐藏的控制流
- error are value

## 2.Error Type

### 2.1 Sentinel Error （不推荐）

```go
ErrInvalidUnreadByte = errors.New("bufio: invalid use of UnreadByte")
```

预定义的错误，缺点：

- 不灵活，调用方使用==去比较错误值是否相等；一旦出现fmt.Errorf这种携带上下文信息的error，会破坏相等性检查
- <u>成为你的公共api;比如io.reader,io.copy这类函数都需要去判断错误类型是否是io.eof,但这并不是一个错误 ？</u>
- 创建了两个包之间的依赖

### 2.2 Error Types

Error type 是实现了 error 接口的自定义类型。例如 MyError 类型记录了文件和行号以展示发生了什么:

os包下的error类型：https://golang.org/src/os/error.go

```go
// PathError records an error and the operation and file path that caused it.
type PathError struct {
	Op   string
	Path string
	Err  error
}

func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }

func (e *PathError) Unwrap() error { return e.Err }
```

- 调用者要使用类型断言和类型 switch，就要让自定义的 error 变为 public。这种模型会导致和调用者产生强耦合，从而导致 API 变得脆弱。
- error types 携带更多的上下文，但也共享 error values 许多相同的问题。
- 建议是避免错误类型，或者至少避免将它们作为公共 API 的一部分。

### 2.3 Opque Error

不透明的错误处理，优势在于：减少代码之间耦合，调用者只需关心成功还是失败，无需关心错误内部



说白了就是把错误类型的判断，封装起来，调用者只用判断err是否为空，或者调用某个函数判断是否是一个err type
