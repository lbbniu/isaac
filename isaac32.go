package isaac

import "math"

type UINT32_C = uint32

// Isaac32 对应 struct isaac_state
type Isaac32 struct {
	m []uint32 // 状态表
	r []uint32 // 结果表
	a uint32
	b uint32
	c uint32
}

func just32(a uint32) uint32 {
	// return a & ((1 << 1 << (32 - 1)) - 1)
	return a & math.MaxUint32
}

// ind32 原始C里的宏：ind(mm, x) = *(ub4*)((ub1*)(mm) + ((x) & ((RANDSIZ-1)<<2)))
// 解释：对 mm 做"按字节"的偏移，然后再取 32 位整型。
// 等价于在 Go 中： mm[( (x) & ((RANDSIZ-1)<<2)) >> 2]。
func ind32(m []uint32, x uint32) uint32 {
	return m[(x&((ISAAC_WORDS-1)*4))>>2]
}

// mix32 对应原始C里的宏 mix(a,b,c,d,e,f,g,h)
func mix32(a, b, c, d, e, f, g, h uint32) (na, nb, nc, nd, ne, nf, ng, nh uint32) {
	a ^= b << 11
	d += a
	b += c
	b ^= just32(c) >> 2
	e += b
	c += d
	c ^= d << 8
	f += c
	d += e
	d ^= just32(e) >> 16
	g += d
	e += f
	e ^= f << 10
	h += e
	f += g
	f ^= just32(g) >> 4
	a += f
	g += h
	g ^= h << 8
	b += g
	h += a
	h ^= just32(a) >> 9
	c += h
	a += b
	return a, b, c, d, e, f, g, h
}

// isaac_refill 对应 C 版本的 isaac_refill 函数
func (s *Isaac32) isaac_refill(r []uint32) {
	a := s.a
	b := s.b + (s.c + 1)
	s.c++

	m := s.m
	HALF := ISAAC_WORDS / 2

	// isaac_step 对应 C 语言中的 ISAAC_STEP 宏
	step := func(i int, off int, mix uint32) {
		a = (a ^ mix) + m[off+i]
		x := m[i]
		y := ind32(m, x) + a + b
		m[i] = y
		b = just32(ind32(m, y>>ISAAC_WORDS_LOG) + x)
		r[i] = b
	}

	// 前半段
	for i := 0; i < HALF; i += 4 {
		// step1: a = (a << 13)
		step(i, HALF, a<<13)
		// step2: a = (a >> 6)
		step(i+1, HALF, a>>6)
		// step3: a = (a << 2)
		step(i+2, HALF, a<<2)
		// step4: a = (a >> 16)
		step(i+3, HALF, a>>16)
	}

	// 后半段
	for i := HALF; i < ISAAC_WORDS; i += 4 {
		// step1: a = (a << 13)
		step(i, -HALF, a<<13)
		// step2: a = (a >> 6)
		step(i+1, -HALF, a>>6)
		// step3: a = (a << 2)
		step(i+2, -HALF, a<<2)
		// step4: a = (a >> 16)
		step(i+3, -HALF, a>>16)
	}

	s.a = a
	s.b = b
}

// isaac_seed 对应 C 版本的 isaac_seed 函数
func (s *Isaac32) isaac_seed() {
	// 使用与 C 版本相同的初始值
	a := uint32(0x1367df5a)
	b := uint32(0x95d90059)
	c := uint32(0xc3163e4b)
	d := uint32(0x0f421ad8)
	e := uint32(0xd92a4a78)
	f := uint32(0xa51a3c49)
	g := uint32(0xc4efea1b)
	h := uint32(0x30609119)

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
			a, b, c, d, e, f, g, h = mix32(a, b, c, d, e, f, g, h)
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

func NewIsaac32() *Isaac32 {
	return &Isaac32{
		m: make([]uint32, ISAAC_WORDS),
	}
}

// Seed initializes ISAAC32
func (isaac *Isaac32) Seed(seed uint32, initValues ...uint32) {
	// Use the same initial values as the C version
	var a, b, c, d, e, f, g, h uint32
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
		a = 0x1367df5a
		b = 0x95d90059
		c = 0xc3163e4b
		d = 0x0f421ad8
		e = 0xd92a4a78
		f = 0xa51a3c49
		g = 0xc4efea1b
		h = 0x30609119
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
			a, b, c, d, e, f, g, h = mix32(a, b, c, d, e, f, g, h)
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

func (s *Isaac32) Refill(r []uint32) {
	s.isaac_refill(r)
}

func (s *Isaac32) Uint32() uint32 {
	if r := s.r; len(r) == 0 {
		r = make([]uint32, ISAAC_WORDS)
		s.Refill(r)
		s.r = r
	}
	r := s.r[0]
	s.r = s.r[1:]
	return r
}
