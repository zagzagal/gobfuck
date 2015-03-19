package main

import (
	"fmt"
	"io"
)

const CELLSIZE int = 300000

type brainfuck struct {
	inst    []byte
	instPtr int
	cells   []byte
	ptr     int
	depth   int
	verbose bool
	debug   bool
	out     string
}

func NewBrainFuck() brainfuck {
	var a brainfuck
	a.cells = make([]byte, 1, CELLSIZE)
	a.ptr = 0
	a.depth = 0
	a.instPtr = 0
	a.verbose = false
	a.out = ""
	return a
}

func (b *brainfuck) incPtr() {
	if b.ptr != CELLSIZE-1 {
		if n := len(b.cells); b.ptr+1 >= n {
			b.log(fmt.Sprintf("*expanding %d -> %d [%d]", n-1, n, cap(b.cells)))
			b.cells = append(b.cells, 0)
		}
		b.ptr++
		b.log(fmt.Sprintf("incPtr: [btr: %d] [val: %v]",
			b.ptr,
			b.cells[b.ptr]))
	}
}

func (b *brainfuck) decPtr() {
	if b.ptr > 0 {
		b.ptr--
		b.log(fmt.Sprintf("decPtr: [ptr: %d] [val: %v]",
			b.ptr,
			b.cells[b.ptr]))
	} else {
		b.log(fmt.Sprintf("*Attempting to go past tape\n\t[inst: %d]",
			b.inst))
	}
}

func (b *brainfuck) incCell() {
	b.cells[b.ptr]++
	b.log(fmt.Sprintf("incCell: [cell: %d] [val: %v -> %v]",
		b.ptr,
		b.cells[b.ptr]-1,
		b.cells[b.ptr]))
}

func (b *brainfuck) decCell() {
	b.cells[b.ptr]--
	b.log(fmt.Sprintf("decCell: [cell: %d] [val: %v -> %v]",
		b.ptr,
		b.cells[b.ptr]-1,
		b.cells[b.ptr]))
}

func (b *brainfuck) output() byte {
	if b.verbose {
		b.out += string(b.cells[b.ptr])
	}
	b.log(fmt.Sprintf("output: [cell: %d] [val: %v] [ascii: %s]",
		b.ptr,
		b.cells[b.ptr],
		string(b.cells[b.ptr]),
	))
	return b.cells[b.ptr]
}

func (b *brainfuck) input(bt byte) {
	b.cells[b.ptr] = bt
}

func (b *brainfuck) startBrace() {
	if b.cells[b.ptr] != 0 {
		b.log(fmt.Sprintf("StartBrace [start: %d] [val: %d] [depth: %d]",
			b.instPtr,
			b.cells[b.ptr],
			b.depth))
		b.depth++
		return
	}
	startDepth := b.depth
	startLoc := b.instPtr
	for b.cells[b.ptr] != ']' && b.depth != startDepth {
		switch b.inst[b.instPtr] {
		case '[':
			b.depth++
		case ']':
			b.depth--
		}
		b.instPtr++

	}
	b.instPtr++
	b.log(fmt.Sprintf("startBrace: [start: %d] [finish: %d] [val: %d ] [depth: %d]",
		startLoc,
		b.instPtr,
		b.cells[b.ptr],
		b.depth))
}

func (b *brainfuck) endBrace() {
	if b.cells[b.ptr] == 0 {
		b.log(fmt.Sprintf("Endbrace [end: %d] [depth: %d]",
			b.instPtr,
			b.depth))
		return
	}
	startDepth := b.depth
	startLoc := b.instPtr
	b.depth--
	for b.depth != startDepth {
		b.instPtr--
		switch b.inst[b.instPtr] {
		case '[':
			b.depth++
		case ']':
			b.depth--
		}
	}
	b.log(fmt.Sprintf("endBrace: [start: %d] [finish: %d] [val: %d] [depth: %d]",
		startLoc,
		b.instPtr,
		b.cells[b.ptr],
		b.depth))
}

func (b *brainfuck) instructions(i []byte) {
	b.inst = i
}

func (b *brainfuck) exec(in io.ByteReader, out io.Writer, v bool, bug bool) {
	if b.inst == nil {
		return
	}

	b.verbose = v
	b.debug = bug

	for !(b.instPtr >= len(b.inst)) {
		gotCmd := true
		switch b.inst[b.instPtr] {
		case '>':
			b.incPtr()
		case '<':
			b.decPtr()
		case '+':
			b.incCell()
		case '-':
			b.decCell()
		case '.':
			char := b.output()
			if !b.verbose {
				out.Write([]byte{char})
			}
		case ',':
			bt, err := in.ReadByte()
			if err != nil {
				break
			}
			b.input(bt)
		case '[':
			b.startBrace()
		case ']':
			b.endBrace()
		default:
			gotCmd = false
		}
		b.instPtr++
		if gotCmd {
			b.Debug()
		}
	}

	if b.verbose {
		fmt.Printf("\n%v\n", b.out)
	}
}

func (b *brainfuck) log(s string) {
	if b.verbose {
		fmt.Println(s)
	}
}

func (b *brainfuck) Debug() {
	if !b.debug {
		return
	}
	fmt.Scanln()
}
