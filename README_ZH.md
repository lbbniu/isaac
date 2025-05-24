# ISAAC

ISAAC 是一个由 Robert J. Jenkins Jr. 在 1996 年设计的密码学安全的伪随机数生成器（CSPRNG）和流密码。这个 Go 实现提供了 32 位和 64 位版本的 ISAAC，以及一个支持两种类型的泛型实现。

## 特性

- 纯 Go 实现
- 支持 `uint32` 和 `uint64` 类型的泛型实现
- 密码学安全
- 快速高效
- 线程安全
- 无外部依赖
- 可自定义初始值

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
    // 创建一个新的 ISAAC 实例（使用 uint32）
    rng := isaac.New[uint32]()
    
    // 生成随机数
    for i := 0; i < 5; i++ {
        fmt.Println(rng.Rand())
    }
}
```

### 使用 uint64

```go
// 创建一个新的 ISAAC 实例（使用 uint64）
rng := isaac.New[uint64]()
```

### 设置种子

```go
// 创建并设置种子的 ISAAC 实例
rng := isaac.New[uint32]()
rng.Seed(12345) // 使用 uint32 值作为种子

// 或者使用自定义初始值
rng.Seed(12345, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9,
         0x9e3779b9, 0x9e3779b9, 0x9e3779b9, 0x9e3779b9) // 用于 uint32
// 对于 uint64，使用 0x9e3779b97f4a7c13 替代
```

## 实现细节

该实现包括：

- `isaac.go` 中的泛型实现
- `isaac32.go` 中的 32 位特定实现
- `isaac64.go` 中的 64 位特定实现
- 全面的测试覆盖

## 安全性

ISAAC 被设计为密码学安全的。但是请注意：

1. 始终使用密码学安全的种子
2. 不要将相同的种子用于不同的用途
3. 对于新应用，考虑使用更现代的 CSPRNG

## 许可证

本项目采用 MIT 许可证 - 详见 LICENSE 文件。

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 参考资料

- [ISAAC: 一个快速的密码学随机数生成器](http://burtleburtle.net/bob/rand/isaac.html)
- [ISAAC 和 RC4](http://burtleburtle.net/bob/rand/isaacafa.html)
- [GNU Coreutils ISAAC 测试](https://github.com/coreutils/coreutils/blob/master/gl/tests/test-rand-isaac.c)
- [GNU Coreutils ISAAC 实现](https://github.com/coreutils/coreutils/blob/master/gl/lib/rand-isaac.c) 