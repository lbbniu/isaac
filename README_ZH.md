# ISAAC

ISAAC 是由 Robert J. Jenkins Jr. 在 1996 年设计的密码学安全的伪随机数生成器（CSPRNG）和流密码。这个 Go 语言实现提供了 32 位和 64 位两个版本的 ISAAC，并通过泛型实现同时支持这两种类型。

## 特性

- 纯 Go 语言实现
- 泛型实现，同时支持 `uint32` 和 `uint64` 类型
- 密码学安全
- 快速高效
- 线程安全
- 无外部依赖

## 安装

```bash
go get github.com/lbbniu/isaac
```

## 使用方法

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/lbbniu/isaac"
)

func main() {
    // 创建一个新的 uint32 类型的 ISAAC 实例
    rng := isaac.New[uint32]()
    
    // 生成随机数
    for i := 0; i < 5; i++ {
        fmt.Println(rng.Rand())
    }
}
```

### 使用 uint64

```go
// 创建一个新的 uint64 类型的 ISAAC 实例
rng := isaac.New[uint64]()
```

### 设置种子

```go
// 创建并设置种子的 ISAAC 实例
rng := isaac.New[uint32]()
rng.Seed(12345) // 使用 uint32 类型的值作为种子
```

## 实现细节

实现包括：

- 泛型实现在 `isaac.go` 中
- 32 位特定实现在 `isaac32.go` 中
- 64 位特定实现在 `isaac64.go` 中
- 完整的测试覆盖

## 安全性

ISAAC 被设计为密码学安全的。但请注意：

1. 始终使用密码学安全的种子
2. 不要在不同用途中重用相同的种子
3. 对于新应用，考虑使用更现代的 CSPRNG

## 许可证

本项目采用 MIT 许可证 - 详见 LICENSE 文件。

## 贡献

欢迎贡献代码！请随时提交 Pull Request。

## 参考资料

- [ISAAC: 一个快速的密码学随机数生成器](http://burtleburtle.net/bob/rand/isaac.html)
- [ISAAC 和 RC4](http://burtleburtle.net/bob/rand/isaacafa.html)
- [GNU Coreutils ISAAC 测试](https://github.com/coreutils/coreutils/blob/master/gl/tests/test-rand-isaac.c)
- [GNU Coreutils ISAAC 实现](https://github.com/coreutils/coreutils/blob/master/gl/lib/rand-isaac.c) 