package isaac

import (
	"math"
)

// ISAAC struct using generic type
type ISAAC[T uint32 | uint64] struct {
	m []T
	r []T
	a T
	b T
	c T
}

// New creates a new ISAAC instance
func New[T uint32 | uint64]() *ISAAC[T] {
	var isaac ISAAC[T]
	isaac.m = make([]T, ISAAC_WORDS)
	return &isaac
}

// Seed initializes ISAAC
func (s *ISAAC[T]) Seed(seed T) {
	// Use the same initial values as the C version
	var a, b, c, d, e, f, g, h T
	switch any(a).(type) {
	case uint32:
		a32 := uint32(0x1367df5a)
		b32 := uint32(0x95d90059)
		c32 := uint32(0xc3163e4b)
		d32 := uint32(0x0f421ad8)
		e32 := uint32(0xd92a4a78)
		f32 := uint32(0xa51a3c49)
		g32 := uint32(0xc4efea1b)
		h32 := uint32(0x30609119)
		a, b, c, d, e, f, g, h = T(a32), T(b32), T(c32), T(d32), T(e32), T(f32), T(g32), T(h32)
	case uint64:
		a64 := uint64(0x647c4677a2884b7c)
		b64 := uint64(0xb9f8b322c73ac862)
		c64 := uint64(0x8c0ea5053d4712a0)
		d64 := uint64(0xb29b2e824a595524)
		e64 := uint64(0x82f053db8355e0ce)
		f64 := uint64(0x48fe4a0fa5a09315)
		g64 := uint64(0xae985bf2cbfc89ed)
		h64 := uint64(0x98f5704f6c44c0ab)
		a, b, c, d, e, f, g, h = T(a64), T(b64), T(c64), T(d64), T(e64), T(f64), T(g64), T(h64)
	}

	// Initialize m array
	for i := 0; i < ISAAC_WORDS; i++ {
		s.m[i] = 0
	}

	// Initialize m array with seed
	s.m[0] = seed

	// Mix S->m so that every part of the seed affects every part of the state
	// Two rounds of mixing
	for range [2]struct{}{} {
		for i := 0; i < ISAAC_WORDS; i += 8 {
			a += s.m[i]
			b += s.m[i+1]
			c += s.m[i+2]
			d += s.m[i+3]
			e += s.m[i+4]
			f += s.m[i+5]
			g += s.m[i+6]
			h += s.m[i+7]
			a, b, c, d, e, f, g, h = mix(a, b, c, d, e, f, g, h)
			s.m[i] = a
			s.m[i+1] = b
			s.m[i+2] = c
			s.m[i+3] = d
			s.m[i+4] = e
			s.m[i+5] = f
			s.m[i+6] = g
			s.m[i+7] = h
		}
	}

	s.a = 0
	s.b = 0
	s.c = 0
}

// Refill replenishes the random number array
func (isaac *ISAAC[T]) Refill(r []T) {
	a := isaac.a
	b := isaac.b + (isaac.c + 1)
	isaac.c++

	m := isaac.m
	HALF := ISAAC_WORDS / 2

	// isaac_step corresponds to the ISAAC_STEP macro in C
	step := func(i int, off int, mix T) {
		switch any(a).(type) {
		case uint32:
			a = (a ^ mix) + m[off+i]
		case uint64:
			a = (0 ^ mix) + m[off+i]
		}
		x := m[i]
		y := ind(m, x) + a + b
		m[i] = y
		b = just(ind(m, y>>ISAAC_WORDS_LOG) + x)
		r[i] = b
	}

	// First half
	for i := 0; i < HALF; i += 4 {
		switch any(a).(type) {
		case uint32:
			// step1: a = (a << 13)
			step(i, HALF, a<<13)
			// step2: a = (a >> 6)
			step(i+1, HALF, a>>6)
			// step3: a = (a << 2)
			step(i+2, HALF, a<<2)
			// step4: a = (a >> 16)
			step(i+3, HALF, a>>16)
		case uint64:
			// step1: a = ^ (a ^ (a << 21))
			step(i, HALF, ^(a ^ (a << 21)))
			// step2: a = a ^ (a >> 5)
			step(i+1, HALF, a^(just(a)>>5))
			// step3: a = a ^ (a << 12)
			step(i+2, HALF, a^(a<<12))
			// step4: a = a ^ (a >> 33)
			//nolint:staticcheck // >>33 is only executed in uint64 branch, uint32 won't trigger
			step(i+3, HALF, a^(just(a)>>33))
		}
	}

	// Second half
	for i := HALF; i < ISAAC_WORDS; i += 4 {
		switch any(a).(type) {
		case uint32:
			// step1: a = (a << 13)
			step(i, -HALF, a<<13)
			// step2: a = (a >> 6)
			step(i+1, -HALF, a>>6)
			// step3: a = (a << 2)
			step(i+2, -HALF, a<<2)
			// step4: a = (a >> 16)
			step(i+3, -HALF, a>>16)
		case uint64:
			// step1: a = ^ (a ^ (a << 21))
			step(i, -HALF, ^(a ^ (a << 21)))
			// step2: a = a ^ (a >> 5)
			step(i+1, -HALF, a^(just(a)>>5))
			// step3: a = a ^ (a << 12)
			step(i+2, -HALF, a^(a<<12))
			// step4: a = a ^ (a >> 33)
			//nolint:staticcheck // >>33 is only executed in uint64 branch, uint32 won't trigger
			step(i+3, -HALF, a^(just(a)>>33))
		}
	}

	isaac.a = a
	isaac.b = b
}

// Rand returns the next random number
func (isaac *ISAAC[T]) Rand() T {
	if len(isaac.r) == 0 {
		r := make([]T, ISAAC_WORDS)
		isaac.Refill(r)
		isaac.r = r
	}
	result := isaac.r[0]
	isaac.r = isaac.r[1:]
	return result
}

// Generic implementation of ind function
func ind[T uint32 | uint64](m []T, x T) T {
	var shift T
	switch any(x).(type) {
	case uint32:
		shift = 2
	case uint64:
		shift = 3
	}
	return m[(x&((ISAAC_WORDS-1)<<shift))>>shift]
}

// Generic implementation of just function
func just[T uint32 | uint64](a T) T {
	switch v := any(a).(type) {
	case uint32:
		return T(v & math.MaxUint32)
	case uint64:
		return T(v & math.MaxUint64)
	default:
		return a
	}
}

// Generic implementation of mix function
func mix[T uint32 | uint64](a, b, c, d, e, f, g, h T) (T, T, T, T, T, T, T, T) {
	switch any(a).(type) {
	case uint32:
		a ^= b << 11
		d += a
		b += c
		b ^= just(c) >> 2
		e += b
		c += d
		c ^= d << 8
		f += c
		d += e
		d ^= just(e) >> 16
		g += d
		e += f
		e ^= f << 10
		h += e
		f += g
		f ^= just(g) >> 4
		a += f
		g += h
		g ^= h << 8
		b += g
		h += a
		h ^= just(a) >> 9
		c += h
		a += b
	case uint64:
		a -= e
		f ^= just(h) >> 9
		h += a
		b -= f
		g ^= a << 9
		a += b
		c -= g
		h ^= just(b) >> 23
		b += c
		d -= h
		a ^= c << 15
		c += d
		e -= a
		b ^= just(d) >> 14
		d += e
		f -= b
		c ^= e << 20
		e += f
		g -= c
		d ^= just(f) >> 17
		f += g
		h -= d
		e ^= g << 14
		g += h
	}
	return a, b, c, d, e, f, g, h
}
