package isaac

import "math"

type UINT64_C = uint64

// Isaac64 对应 struct isaac_state
type Isaac64 struct {
	m []uint64 // 状态表
	r []uint64 // 结果表
	a uint64
	b uint64
	c uint64
}

// 常量定义，对齐 C 版本
const (
	ISAAC_BITS      = 64
	ISAAC_WORDS     = 1 << 8
	ISAAC_WORDS_LOG = 8
)

func just64(a uint64) uint64 {
	// return a & ((1 << 1 << (ISAAC_BITS - 1)) - 1)
	return a & math.MaxUint64
}

// ind64 原始C里的宏：ind64(mm, x) = *(ub8*)((ub1*)(mm) + ((x) & ((RANDSIZ-1)<<3)))
// 解释：对 mm 做"按字节"的偏移，然后再取 64 位整型。
// 等价于在 Go 中： mm[( (x) & ((RANDSIZ-1)<<3)) >> 3]。
func ind64(m []uint64, x uint64) uint64 {
	return m[(x&((ISAAC_WORDS-1)*8))>>3]
}

// mix64 对应原始C里的宏 mix64(a,b,c,d,e,f,g,h)
func mix64(a, b, c, d, e, f, g, h uint64) (na, nb, nc, nd, ne, nf, ng, nh uint64) {
	a -= e
	f ^= (just64(h) >> 9)
	h += a
	b -= f
	g ^= (a << 9)
	a += b
	c -= g
	h ^= (just64(b) >> 23)
	b += c
	d -= h
	a ^= (c << 15)
	c += d
	e -= a
	b ^= (just64(d) >> 14)
	d += e
	f -= b
	c ^= (e << 20)
	e += f
	g -= c
	d ^= (just64(f) >> 17)
	f += g
	h -= d
	e ^= (g << 14)
	g += h
	return a, b, c, d, e, f, g, h
}

// isaac_refill 对应 C 版本的 isaac_refill 函数
func (s *Isaac64) isaac_refill(r []uint64) {
	a := s.a
	b := s.b + (s.c + 1)
	s.c++

	m := s.m
	HALF := ISAAC_WORDS / 2

	// isaac_step 对应 C 语言中的 ISAAC_STEP 宏
	step := func(i int, off int, mix uint64) {
		a = (0 ^ mix) + m[off+i]
		x := m[i]
		y := ind64(s.m, x) + a + b
		m[i] = y
		b = just64(ind64(s.m, y>>ISAAC_WORDS_LOG) + x)
		r[i] = b
	}

	// 前半段
	for i := 0; i < HALF; i += 4 {
		// step1: a = ^ (a ^ (a << 21))
		step(i, HALF, ^(a ^ (a << 21)))
		// step2: a = a ^ (a >> 5)
		step(i+1, HALF, a^(just64(a)>>5))
		// step3: a = a ^ (a << 12)
		step(i+2, HALF, a^(a<<12))
		// step4: a = a ^ (a >> 33)
		step(i+3, HALF, a^(just64(a)>>33))
	}

	// 后半段
	for i := HALF; i < ISAAC_WORDS; i += 4 {
		// step1: a = ^ (a ^ (a << 21))
		step(i, -HALF, ^(a ^ (a << 21)))
		// step2: a = a ^ (just (a) >>  5)
		step(i+1, -HALF, a^(just64(a)>>5))
		// step3: a = a^(a<<12)
		step(i+2, -HALF, a^(a<<12))
		// step4: a = a^(just64(a)>>33)
		step(i+3, -HALF, a^(just64(a)>>33))
	}

	s.a = a
	s.b = b
}

// isaac_seed 对应 C 版本的 isaac_seed 函数
func (s *Isaac64) isaac_seed() {
	// 使用与 C 版本相同的初始值
	a := uint64(0x647c4677a2884b7c)
	b := uint64(0xb9f8b322c73ac862)
	c := uint64(0x8c0ea5053d4712a0)
	d := uint64(0xb29b2e824a595524)
	e := uint64(0x82f053db8355e0ce)
	f := uint64(0x48fe4a0fa5a09315)
	g := uint64(0xae985bf2cbfc89ed)
	h := uint64(0x98f5704f6c44c0ab)

	// Mix S->m so that every part of the seed affects every part of the state
	// 二遍混合
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
			a, b, c, d, e, f, g, h = mix64(a, b, c, d, e, f, g, h)
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

func NewIsaac64() *Isaac64 {
	return &Isaac64{
		m: make([]uint64, ISAAC_WORDS),
	}
}

// Seed initializes ISAAC64
func (isaac *Isaac64) Seed(seed uint64, initValues ...uint64) {
	// Use the same initial values as the C version
	var a, b, c, d, e, f, g, h uint64
	if len(initValues) >= 8 {
		a = initValues[0]
		b = initValues[1]
		c = initValues[2]
		d = initValues[3]
		e = initValues[4]
		f = initValues[5]
		g = initValues[6]
		h = initValues[7]
	} else {
		a = 0x647c4677a2884b7c
		b = 0xb9f8b322c73ac862
		c = 0x8c0ea5053d4712a0
		d = 0xb29b2e824a595524
		e = 0x82f053db8355e0ce
		f = 0x48fe4a0fa5a09315
		g = 0xae985bf2cbfc89ed
		h = 0x98f5704f6c44c0ab
	}

	// Initialize m array
	for i := 0; i < ISAAC_WORDS; i++ {
		isaac.m[i] = 0
	}

	// Initialize m array with seed
	isaac.m[0] = seed

	// Mix S->m so that every part of the seed affects every part of the state
	// Two rounds of mixing
	for range [2]struct{}{} {
		for i := 0; i < ISAAC_WORDS; i += 8 {
			a += isaac.m[i]
			b += isaac.m[i+1]
			c += isaac.m[i+2]
			d += isaac.m[i+3]
			e += isaac.m[i+4]
			f += isaac.m[i+5]
			g += isaac.m[i+6]
			h += isaac.m[i+7]
			a, b, c, d, e, f, g, h = mix64(a, b, c, d, e, f, g, h)
			isaac.m[i] = a
			isaac.m[i+1] = b
			isaac.m[i+2] = c
			isaac.m[i+3] = d
			isaac.m[i+4] = e
			isaac.m[i+5] = f
			isaac.m[i+6] = g
			isaac.m[i+7] = h
		}
	}

	isaac.a = 0
	isaac.b = 0
	isaac.c = 0
}

func (s *Isaac64) Refill(r []uint64) {
	s.isaac_refill(r)
}

func (s *Isaac64) Uint64() uint64 {
	if r := s.r; len(r) == 0 {
		r = make([]uint64, ISAAC_WORDS)
		s.Refill(r)
		s.r = r
	}
	r := s.r[0]
	s.r = s.r[1:]
	return r
}
