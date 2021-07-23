package sidx

import (
	"bytes"
	"encoding/binary"
)

// ReferenceItem for one range item
type ReferenceItem struct {
	ReferenceType      int
	ReferenceSize      uint32
	SubsegmentDuration uint32
	DurationSec        float32
	StartTimeSec       float32
	StartRange         uint32
	EndRange           uint32
}

// ParsedInfo for sidx result
type ParsedInfo struct {
	Version     uint32
	Flags       uint32
	ReferenceID uint32
	TimeScale   uint32
	References  []*ReferenceItem
}

// Parser for sidx
type Parser struct {
	data  []byte
	index int
}

// NewParser sidx parser
func NewParser(data []byte) *Parser {
	return &Parser{
		data,
		0,
	}
}

func (p *Parser) read(n int) []byte {
	var r = p.data[p.index : p.index+n]
	p.index += n
	return r
}

func (p *Parser) readUInt8() uint8 {
	return readUInt8(p.read(1))
}

func (p *Parser) readUInt16() uint16 {
	return readUInt16(p.read(2))
}

func (p *Parser) readUInt32() uint32 {
	return readUInt32(p.read(4))
}
func (p *Parser) readUInt64() uint64 {
	return readUInt64(p.read(8))
}

// Parse return raw parsed info
func (p *Parser) Parse(indexEndoffset uint32) *ParsedInfo {
	// 前8字节是固定的 box header,略过
	p.read(8)
	var (
		versionAndFlags          = p.readUInt32()
		version                  = versionAndFlags >> 24
		flags                    = versionAndFlags & 0xFFFFFF
		referenceID              = p.readUInt32()
		timeScale                = p.readUInt32()
		earliestPresentationTime uint64
		firstOffset              uint64
	)
	if version == 0 {
		earliestPresentationTime = uint64(p.readUInt32())
		firstOffset = uint64(p.readUInt32())
	} else {
		earliestPresentationTime = p.readUInt64()
		firstOffset = p.readUInt64()
	}
	var (
		_              = p.read(2) // skip reserved
		referenceCount = p.readUInt16()
		i              uint16
		references     = []*ReferenceItem{}
		offset         = uint32(firstOffset) + indexEndoffset + 1
		time           = earliestPresentationTime
	)

	for i = 0; i < referenceCount; i++ {
		// 由于reference_type的值都是0,最终结果无影响,所以referenced_size作为32位处理
		referenceType := 0
		referenceSize := p.readUInt32()
		subsegmentDuration := p.readUInt32()

		// 下面是 starts_with_SAP(1), SAP_type(3), SAP_delta_time(28) 没用到,这里忽略掉
		p.read(4)

		startRange := offset
		endRange := offset + referenceSize

		item := &ReferenceItem{
			ReferenceType:      referenceType,
			ReferenceSize:      referenceSize,
			SubsegmentDuration: subsegmentDuration,
			DurationSec:        float32(subsegmentDuration) / float32(timeScale),
			StartTimeSec:       float32(time) / float32(timeScale),
			StartRange:         startRange,
			EndRange:           endRange,
		}

		references = append(references, item)
		offset += referenceSize
		time += uint64(subsegmentDuration)
	}
	var info = &ParsedInfo{
		Version:     version,
		Flags:       flags,
		ReferenceID: referenceID,
		TimeScale:   timeScale,
		References:  references,
	}
	return info
}

func readUInt8(data []byte) (ret uint8) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}

func readUInt16(data []byte) (ret uint16) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}

func readUInt32(data []byte) (ret uint32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}

func readUInt64(data []byte) (ret uint64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}
