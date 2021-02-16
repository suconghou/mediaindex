package ebml

import (
	"encoding/hex"
	"encoding/json"
	"math"
	"strconv"
)

// Parser for ebml
type Parser struct {
	parser *ParserElement
	result []*ParserElement
}

type vIntItem struct {
	id   string
	data []byte
}

// ParserBuffer for buffer
type ParserBuffer struct {
	data  []byte
	index int
	len   int
}

// ParserElement for ele
type ParserElement struct {
	ID       string
	Name     string
	Type     string
	Value    []byte
	buffer   *ParserBuffer
	Children []*ParserElement
}

// NewParser for ebml
func NewParser(data []byte) *Parser {
	return &Parser{
		NewParserElement(data),
		nil,
	}
}

// NewParserBuffer create buffer
func NewParserBuffer(data []byte) *ParserBuffer {
	return &ParserBuffer{
		data,
		0,
		len(data),
	}
}

// NewParserElement create new ele
func NewParserElement(data []byte) *ParserElement {
	return &ParserElement{
		buffer: NewParserBuffer(data),
	}
}

// Parse for ebml
func (p *Parser) Parse() []*ParserElement {
	if p.result == nil {
		p.result = p.parser.parseElements()
	}
	return p.result
}

// JSON turn into json
func (p *Parser) JSON() ([]byte, error) {
	if p.result == nil {
		p.result = p.parser.parseElements()
	}
	return json.Marshal(p.result)
}

// Parse one buffer
func (p *ParserBuffer) parse() *vIntItem {
	if p.index >= p.len {
		return nil
	}
	return p.readVint()
}

func (p *ParserBuffer) readVint() *vIntItem {
	var s = p.read(1)
	var w = vIntWidth(s[0]) + 1
	var id = p.rewind(1).read(int(w))
	s = p.read(1)
	w = vIntWidth(s[0]) + 1
	var len = p.rewind(1).read(int(w))
	var lenNum = vIntNum(len)
	var data = p.read(int(lenNum))
	return &vIntItem{
		data: data,
		id:   hex.EncodeToString(id),
	}
}

func (p *ParserBuffer) read(n int) []byte {
	r := p.data[p.index : p.index+n]
	p.index += n
	return r
}

func (p *ParserBuffer) rewind(n int) *ParserBuffer {
	p.index -= n
	return p
}

func (p *ParserElement) parseElements() []*ParserElement {
	for {
		item := p.buffer.parse()
		if item == nil {
			break
		}
		p.Children = append(p.Children, p.parseElement(item.id, item.data))
	}
	return p.Children
}

func (p *ParserElement) parseElement(id string, data []byte) *ParserElement {
	var ele = NewParserElement(data)
	ele.ID = id
	var meta = getOne(id)
	if meta != nil {
		ele.Name = meta.name
		ele.Type = meta.vtype
		if meta.vtype == "m" {
			ele.parseElements()
		} else if meta.vtype == "u" {
			ele.Value = data
		}
	}
	return ele
}

// 传入一个字节8位, 判断前多少个bit是0, 返回值可能为 0 - 7
func vIntWidth(b uint8) uint8 {
	var width uint8 = 0
	var num = float64(b)
	for {
		if num >= math.Pow(2, float64(7-width)) {
			break
		}
		width++
	}
	return width
}

// 传入若干个字节 最高位1置为0后转为十进制数
func vIntNum(data []byte) uint64 {
	n := VarNum(data)
	return n - hibit(n)
}

// VarNum  支持0,1,2,3,4,5字节, Big endian https://stackoverflow.com/questions/45000982/convert-3-bytes-to-int-in-go
func VarNum(b []byte) uint64 {
	var l = len(b)
	if l == 0 {
		return 0
	} else if l == 1 {
		return uint64(b[0])
	} else if l == 2 {
		return uint64(int(b[1]) | int(b[0])<<8)
	} else if l == 3 {
		return uint64(int(b[2]) | int(b[1])<<8 | int(b[0])<<16)
	} else if l == 4 {
		return uint64(int(b[3]) | int(b[2])<<8 | int(b[1])<<16 | int(b[0])<<24)
	} else if l == 5 {
		return uint64(int(b[4]) | int(b[3])<<8 | int(b[2])<<16 | int(b[1])<<24 | int(b[0])<<32)
	} else {
		panic("unsupport vint len " + strconv.Itoa(l))
	}
}

// https://stackoverflow.com/questions/53161/find-the-highest-order-bit-in-c
func hibit(n uint64) uint64 {
	n |= (n >> 1)
	n |= (n >> 2)
	n |= (n >> 4)
	n |= (n >> 8)
	n |= (n >> 16)
	n |= (n >> 32)
	return n - (n >> 1)
}
