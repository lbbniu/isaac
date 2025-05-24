package isaac

import "math"

type UINT32_C = uint32

// Isaac32 corresponds to struct isaac_state
type Isaac32 struct {
	m [Words]uint32 // state table
	r []uint32      // result table
	a uint32
	b uint32
	c uint32
}

func just32(a uint32) uint32 {
	// return a & ((1 << 1 << (32 - 1)) - 1)
	return a & math.MaxUint32
}

// ind32 corresponds to the C macro: ind(mm, x) = *(ub4*)((ub1*)(mm) + ((x) & ((RANDSIZ-1)<<2)))
// Explanation: Perform byte-level offset on mm, then take 32-bit integer.
// Equivalent in Go: mm[( (x) & ((RANDSIZ-1)<<2)) >> 2].
func ind32(m [Words]uint32, x uint32) uint32 {
	return m[(x&((Words-1)*4))>>2]
}

// mix32 corresponds to the C macro mix(a,b,c,d,e,f,g,h)
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

// isaac_refill corresponds to the C version of isaac_refill function
func (s *Isaac32) isaac_refill(r *[Words]uint32) {
	a := s.a
	b := s.b + (s.c + 1)
	s.c++

	HALF := Words / 2

	// isaac_step corresponds to the C ISAAC_STEP macro
	step := func(i int, off int, mix uint32) {
		a = (a ^ mix) + s.m[off+i]
		x := s.m[i]
		y := ind32(s.m, x) + a + b
		s.m[i] = y
		b = just32(ind32(s.m, y>>WordsLog) + x)
		r[i] = b
	}

	// First half
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

	// Second half
	for i := HALF; i < Words; i += 4 {
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

func NewIsaac32() *Isaac32 {
	s := &Isaac32{}
	s.Seed([Words]uint32{})
	return s
}

// Seed initializes ISAAC32
// Corresponds to the C isaac_seed function
func (isaac *Isaac32) Seed(seed [Words]uint32, initValues ...uint32) {
	if len(initValues) > 0 && len(initValues) != 8 {
		panic("isaac: need exactly 8 initial values for uint32")
	}

	// Use the same initial values as the C version
	var a, b, c, d, e, f, g, h uint32
	if len(initValues) == 8 {
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
	for i := 0; i < Words; i++ {
		isaac.m[i] = seed[i]
	}

	// Mix S->m so that every part of the seed affects every part of the state
	// Two rounds of mixing
	for range [2]struct{}{} {
		for i := 0; i < Words; i += 8 {
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

func (s *Isaac32) Refill(r *[Words]uint32) {
	s.isaac_refill(r)
}

func (s *Isaac32) Uint32() uint32 {
	if r := s.r; len(r) == 0 {
		var r [Words]uint32
		s.Refill(&r)
		s.r = r[:]
	}
	r := s.r[0]
	s.r = s.r[1:]
	return r
}
