# ISAAC

ISAAC is a cryptographically secure pseudorandom number generator (CSPRNG) and stream cipher designed by Robert J. Jenkins Jr. in 1996. This Go implementation provides both 32-bit and 64-bit versions of ISAAC, with a generic implementation that supports both types.

## Features

- Pure Go implementation
- Generic implementation supporting both `uint32` and `uint64` types
- Cryptographically secure
- Fast and efficient
- Thread-safe
- No external dependencies
- Fixed-size array state for better performance

## Installation

```bash
go get github.com/lbbniu/isaac
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/lbbniu/isaac"
)

func main() {
    // Create a new ISAAC instance with uint32
    rng := isaac.New[uint32]()
    
    // Generate random numbers
    for i := 0; i < 5; i++ {
        fmt.Println(rng.Rand())
    }
}
```

### Using uint64

```go
// Create a new ISAAC instance with uint64
rng := isaac.New[uint64]()
```

### Seeding

```go
// Create a new ISAAC instance
rng := isaac.New[uint32]()

// Seed with a fixed-size array
var seed [isaac.ISAAC_WORDS]uint32
rng.Seed(seed)
```

### Refilling

```go
// Create a new ISAAC instance
rng := isaac.New[uint32]()

// Get a batch of random numbers
var result [isaac.ISAAC_WORDS]uint32
rng.Refill(&result)
```

## Implementation Details

The implementation includes:

- Generic implementation in `isaac.go` with fixed-size array state
- 32-bit specific implementation in `isaac32.go`
- 64-bit specific implementation in `isaac64.go`
- Comprehensive test coverage with test vectors from GNU Coreutils

## Security

ISAAC is designed to be cryptographically secure. However, please note:

1. Always use a cryptographically secure seed
2. Do not reuse the same seed for different purposes
3. Consider using a more modern CSPRNG for new applications

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## References

- [ISAAC: a fast cryptographic random number generator](http://burtleburtle.net/bob/rand/isaac.html)
- [ISAAC and RC4](http://burtleburtle.net/bob/rand/isaacafa.html)
- [GNU Coreutils ISAAC Test](https://github.com/coreutils/coreutils/blob/master/gl/tests/test-rand-isaac.c)
- [GNU Coreutils ISAAC Implementation](https://github.com/coreutils/coreutils/blob/master/gl/lib/rand-isaac.c)
