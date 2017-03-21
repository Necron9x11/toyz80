package z80

// mode is the instruction's addressing mode.
type mode int

const (
	invalid   mode = iota // invalid addressing mode
	immediate             // source only
	immediateExtended
	modifiedPageZero
	relative
	extended
	indexed
	register
	implied
	registerIndirect
	bit
	condition
	displacement
	indirect
)

// opcode describes an instruction.
type opcode struct {
	noCycles  uint64   // cycles instruction costs
	mnemonic  []string // human readable mnemonic z80 = 0, 8080 = 1
	noBytes   uint16   // no of bytes per instruction
	src       mode     // source addressing mode
	srcR      []string // source disassembly cheat
	dst       mode     // destination addressing mode
	dstR      []string // destination disassembly cheat
	multiByte bool     // use alternative array instead for lookup
}

// opcodes are all possible instructions 16 bit instructions.  We just throw
// mmemory at the problem.
var (
	// z80 only 0xed opcodes
	opcodesED = []opcode{
		0x44: {
			mnemonic: []string{"neg", ""},
			noBytes:  2,
			noCycles: 8,
		},
	}
	opcodesDD = []opcode{
		0x23: {
			mnemonic: []string{"inc"},
			dst:      register,
			dstR:     []string{"ix"},
			noBytes:  2,
			noCycles: 10,
		},
		0x2b: {
			mnemonic: []string{"dec"},
			dst:      register,
			dstR:     []string{"ix"},
			noBytes:  2,
			noCycles: 10,
		},
	}
	opcodesFD = []opcode{
		0x23: {
			mnemonic: []string{"inc"},
			dst:      register,
			dstR:     []string{"iy"},
			noBytes:  2,
			noCycles: 10,
		},
		0x2b: {
			mnemonic: []string{"dec"},
			dst:      register,
			dstR:     []string{"iy"},
			noBytes:  2,
			noCycles: 10,
		},
	}
)

// Opcodes are all possible instructions 8 bit instructions.
var (
	opcodes = []opcode{
		// 0x00 nop
		opcode{
			mnemonic: []string{"nop", "nop"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x01 ld bc,nn
		opcode{
			mnemonic: []string{"ld", "lxi"},
			dst:      register,
			dstR:     []string{"bc", "b"},
			src:      immediateExtended,
			srcR:     []string{"", ""},
			noBytes:  3,
			noCycles: 10,
		},
		// 0x02 ld (bc),a
		opcode{
			mnemonic: []string{"ld", "stax"},
			dst:      registerIndirect,
			dstR:     []string{"bc", "b"},
			src:      register,
			srcR:     []string{"a", ""},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x03 inc bc
		opcode{
			mnemonic: []string{"inc", "inx"},
			dst:      register,
			dstR:     []string{"bc", "b"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x04 inc b
		opcode{
			mnemonic: []string{"inc", "inr"},
			dst:      register,
			dstR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x05
		opcode{},
		// 0x06 ld b,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x07
		opcode{},
		// 0x08
		opcode{},
		// 0x09
		opcode{},
		// 0x0a ld a,(bc)
		opcode{
			mnemonic: []string{"ld", "ldax"},
			dst:      register,
			dstR:     []string{"a", ""},
			src:      registerIndirect,
			srcR:     []string{"bc", "b"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x0b dec bc
		opcode{
			mnemonic: []string{"dec", "dcx"},
			dst:      register,
			dstR:     []string{"bc", "d"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x0c
		opcode{},
		// 0x0d
		opcode{},
		// 0x0e ld c,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x0f
		opcode{},

		// 0x10
		opcode{},
		// 0x11
		opcode{},
		// 0x12 ld (de),a
		opcode{
			mnemonic: []string{"ld", "stax"},
			dst:      registerIndirect,
			dstR:     []string{"de", "d"},
			src:      register,
			srcR:     []string{"a", ""},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x13 inc de
		opcode{
			mnemonic: []string{"inc", "inx"},
			dst:      register,
			dstR:     []string{"de", "d"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x14
		opcode{},
		// 0x15
		opcode{},
		// 0x16 ld d,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x17
		opcode{},
		// 0x18
		opcode{
			mnemonic: []string{"jr", ""},
			dst:      displacement,
			noBytes:  2,
			noCycles: 12,
		},
		// 0x19
		opcode{},
		// 0x1a ld a,(de)
		opcode{
			mnemonic: []string{"ld", "ldax"},
			dst:      register,
			dstR:     []string{"a", ""},
			src:      registerIndirect,
			srcR:     []string{"de", "d"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x1b dec de
		opcode{
			mnemonic: []string{"dec", "dcx"},
			dst:      register,
			dstR:     []string{"de", "d"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x1c
		opcode{},
		// 0x1d
		opcode{},
		// 0x1e ld e,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x1f
		opcode{},

		// 0x20
		opcode{},
		// 0x21
		opcode{},
		// 0x22
		opcode{},
		// 0x23 inc de
		opcode{
			mnemonic: []string{"inc", "inx"},
			dst:      register,
			dstR:     []string{"hl", "h"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x24
		opcode{},
		// 0x25
		opcode{},
		// 0x26 ld h,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x27
		opcode{},
		// 0x28 // jr z,d
		opcode{
			mnemonic: []string{"jr", ""},
			dst:      condition,
			dstR:     []string{"z", ""},
			src:      displacement,
			noBytes:  2,
			noCycles: 12, // XXX or 7
		},
		// 0x29
		opcode{},
		// 0x2a
		opcode{},
		// 0x2b dec hl
		opcode{
			mnemonic: []string{"dec", "dcx"},
			dst:      register,
			dstR:     []string{"hl", "h"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x2c
		opcode{},
		// 0x2d
		opcode{},
		// 0x2e ld l,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x2f
		opcode{
			mnemonic: []string{"cpl", "cma"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x30
		opcode{},
		// 0x31 ld sp,nn
		opcode{
			mnemonic: []string{"ld", "lxi"},
			dst:      register,
			dstR:     []string{"sp", "sp"},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0x32
		opcode{},
		// 0x33 inc de
		opcode{
			mnemonic: []string{"inc", "inx"},
			dst:      register,
			dstR:     []string{"sp", "h"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x34
		opcode{},
		// 0x35
		opcode{},
		// 0x36 ld (hl),n
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "m"},
			src:      immediate,
			noBytes:  2,
			noCycles: 10,
		},
		// 0x37
		opcode{
			mnemonic: []string{"scf", "stc"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x38
		opcode{},
		// 0x39
		opcode{},
		// 0x3a ld a,(nn)
		opcode{
			mnemonic: []string{"ld", "lda"},
			dst:      register,
			dstR:     []string{"a", ""},
			src:      extended,
			noBytes:  3,
			noCycles: 7,
		},
		// 0x3b dec sp
		opcode{
			mnemonic: []string{"dec", "dcx"},
			dst:      register,
			dstR:     []string{"sp", "sp"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x3c
		opcode{},
		// 0x3d
		opcode{},
		// 0x3e ld a,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x3f
		opcode{
			mnemonic: []string{"ccf", "cmc"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x40 ld b,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x41 ld b,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x42 ld b,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x43 ld b,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x44 ld b,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x45 ld b,l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x46 ld b,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x47 ld b,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x48 ld c,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x49 ld c,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x4a ld c,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x4b ld c,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x4c ld c,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x4d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x4e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x4f ld c,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x50 ld d,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x51 ld d,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x52 ld d,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x53 ld d,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x54 ld d,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x55 ld d,l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x56 ld d,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x57 ld d,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x58 ld e,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x59 ld e,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x5a ld e,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x5b ld e,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x5c ld e,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x5d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x5e ld e,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x5f ld e,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x60 ld h,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x61 ld h,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x62 ld h,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x63 ld h,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x64 ld h,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x65 ld h,l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x66 ld h,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x67 ld h,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x68 ld l,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x69 ld l,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x6a ld l,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x6b ld l,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x6c ld l,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x6d ld l,l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x6e ld l,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x6f ld l,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x70 ld (hl),b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x71 ld (hl),c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x72 ld (hl),d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x73 ld (hl),e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x74 ld (hl),h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x75 ld (hl),l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x76 halt
		opcode{
			mnemonic: []string{"halt", "hlt"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x77 ld (hl),a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x78 ld a,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x79 ld a,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x7a ld a,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x7b ld a,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x7c ld a,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x7d ld a,l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x7e ld a,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x7f ld a,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x80
		opcode{},
		// 0x81
		opcode{},
		// 0x82
		opcode{},
		// 0x83
		opcode{},
		// 0x84
		opcode{},
		// 0x85
		opcode{},
		// 0x86
		opcode{},
		// 0x87
		opcode{},
		// 0x88
		opcode{},
		// 0x89
		opcode{},
		// 0x8a
		opcode{},
		// 0x8b
		opcode{},
		// 0x8c
		opcode{},
		// 0x8d
		opcode{},
		// 0x8e
		opcode{},
		// 0x8f
		opcode{},

		// 0x90
		opcode{},
		// 0x91
		opcode{},
		// 0x92
		opcode{},
		// 0x93
		opcode{},
		// 0x94
		opcode{},
		// 0x95
		opcode{},
		// 0x96
		opcode{},
		// 0x97
		opcode{},
		// 0x98
		opcode{},
		// 0x99
		opcode{},
		// 0x9a
		opcode{},
		// 0x9b
		opcode{},
		// 0x9c
		opcode{},
		// 0x9d
		opcode{},
		// 0x9e
		opcode{},
		// 0x9f
		opcode{},

		// 0xa0
		opcode{},
		// 0xa1
		opcode{},
		// 0xa2
		opcode{},
		// 0xa3
		opcode{},
		// 0xa4
		opcode{},
		// 0xa5
		opcode{},
		// 0xa6
		opcode{},
		// 0xa7
		opcode{},
		// 0xa8
		opcode{},
		// 0xa9
		opcode{},
		// 0xaa
		opcode{},
		// 0xab
		opcode{},
		// 0xac
		opcode{},
		// 0xad
		opcode{},
		// 0xae
		opcode{},
		// 0xaf
		opcode{},

		// 0xb0
		opcode{},
		// 0xb1
		opcode{},
		// 0xb2
		opcode{},
		// 0xb3
		opcode{},
		// 0xb4
		opcode{},
		// 0xb5
		opcode{},
		// 0xb6
		opcode{},
		// 0xb7
		opcode{},
		// 0xb8
		opcode{},
		// 0xb9
		opcode{},
		// 0xba
		opcode{},
		// 0xbb
		opcode{},
		// 0xbc
		opcode{},
		// 0xbd
		opcode{},
		// 0xbe
		opcode{},
		// 0xbf
		opcode{
			mnemonic: []string{"cp", "cmp"},
			dst:      register,
			dstR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0xc0
		opcode{},
		// 0xc1
		opcode{},
		// 0xc2
		opcode{
			mnemonic: []string{"jp", "jnz"},
			dst:      condition,
			dstR:     []string{"nz", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xc3 jp, nn
		opcode{
			mnemonic: []string{"jp", "jmp"},
			dst:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xc4
		opcode{},
		// 0xc5
		opcode{},
		// 0xc6
		opcode{},
		// 0xc7
		opcode{},
		// 0xc8
		opcode{},
		// 0xc9
		opcode{},
		// 0xca
		opcode{
			mnemonic: []string{"jp", "jz"},
			dst:      condition,
			dstR:     []string{"z", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xcb
		opcode{},
		// 0xcc
		opcode{},
		// 0xcd
		opcode{},
		// 0xce
		opcode{},
		// 0xcf
		opcode{},

		// 0xd0
		opcode{},
		// 0xd1
		opcode{},
		// 0xd2
		opcode{},
		// 0xd3
		opcode{
			mnemonic: []string{"out", "out"},
			dst:      indirect,
			src:      register,
			srcR:     []string{"a", ""},
			noBytes:  2,
			noCycles: 11,
		},
		// 0xd4
		opcode{},
		// 0xd5
		opcode{},
		// 0xd6
		opcode{},
		// 0xd7
		opcode{},
		// 0xd8
		opcode{},
		// 0xd9
		opcode{},
		// 0xda
		opcode{},
		// 0xdb
		opcode{},
		// 0xdc
		opcode{},
		// 0xdd z80 multi byte
		opcode{
			multiByte: true,
		},
		// 0xde
		opcode{},
		// 0xdf
		opcode{},

		// 0xe0
		opcode{},
		// 0xe1
		opcode{},
		// 0xe2
		opcode{},
		// 0xe3
		opcode{},
		// 0xe4
		opcode{},
		// 0xe5
		opcode{},
		// 0xe6
		opcode{},
		// 0xe7
		opcode{},
		// 0xe8
		opcode{},
		// 0xe9
		opcode{},
		// 0xea
		opcode{},
		// 0xeb ex de,hl
		opcode{
			mnemonic: []string{"ex", "xchg"},
			dst:      register,
			dstR:     []string{"de", ""},
			src:      register,
			srcR:     []string{"hl", ""},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xec
		opcode{},
		// 0xed z80 multi byte
		opcode{
			multiByte: true,
		},
		// 0xee
		opcode{},
		// 0xef
		opcode{},

		// 0xf0
		opcode{},
		// 0xf1
		opcode{},
		// 0xf2
		opcode{},
		// 0xf3
		opcode{},
		// 0xf4
		opcode{},
		// 0xf5
		opcode{},
		// 0xf6
		opcode{},
		// 0xf7
		opcode{},
		// 0xf8
		opcode{},
		// 0xf9
		opcode{},
		// 0xfa
		opcode{},
		// 0xfb
		opcode{},
		// 0xfc
		opcode{},
		// 0xfd z80 multi byte
		opcode{
			multiByte: true,
		},
		// 0xfe
		opcode{},
		// 0xff
		opcode{},
	}
)
