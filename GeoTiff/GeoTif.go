package GeoTiff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
)

type FileOffset = int64

// 这是二进制文件中的结构
// this is the struct of the file
type geoAttribute struct {
	Tag               AttributeTag
	Type              DataType
	Len               uint32
	SourceValue       []byte
	Offset            uint32
	GeoAttributeValue geoAttributeValue
}

func (gAttribute geoAttribute) Bytes() uint32 {
	return gAttribute.Len * gAttribute.Type.Bytes()
}

var geoFileAttributeSize = int64(12)

type geoAttributeValue struct {
	rValue interface{}
	BYTE   []uint8
	ASCII  string
	SHORT  []uint16
	LONG   []uint32
	FLOAT  []float32
	DOUBLE []float64
	uint   []uint
}

type GeoAttributes []geoAttribute
type geoTifHeader struct {
	offset    int64
	Attribute GeoAttributes
}

type Meta struct {
	Columns           uint
	Rows              uint
	BitsPerSample     []uint
	samplesPerPixel   uint
	SampleFormat      uint
	PhotometricInterp uint
	mode              ImageMode
	palette           []uint32
	NodataValue       string
	RasterPixelIsArea bool
	EPSGCode          uint
}
type GeoTif struct {
	// 文件路径
	FilePath     string
	tFile        io.ReaderAt
	byteOrder    binary.ByteOrder
	GeoTifHeader geoTifHeader
	GeoKeys      GeoAttributes
	Meta         Meta
	Data         GeoData
	Transform    transform
}

func (g GeoTif) String() string {
	return fmt.Sprintf(`GeoTif{
	FilePath: %s
	byteOrder: %s
	GeoTifHeader: {
		offset: %d
		Attribute[%d]: [
%s
		]
	},
	GeoKeys: [
%s
	]
}`, g.FilePath, g.byteOrder.String(), g.GeoTifHeader.offset, len(g.GeoTifHeader.Attribute), g.GeoTifHeader.Attribute.toString(3), g.GeoKeys.toString(2))
}

var indents = "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"

func (gAttributes GeoAttributes) toString(indent int) string {
	if indent <= 0 {
		indent = 0
	}
	sis := indents[:indent]
	is := indents[:indent+1]
	strs := []string{}
	for _, gAttribute := range gAttributes {
		strs = append(strs, fmt.Sprintf(`%s{
%sTag: %d,
%sType: %d
%sLen: %d
%sSourceValue: %v
%sOffset: %d
%sGeoAttributeValue: %v
%s}`,
			sis,
			is, gAttribute.Tag,
			is, gAttribute.Type,
			is, gAttribute.Len,
			is, gAttribute.SourceValue,
			is, gAttribute.Offset,
			is, gAttribute.GeoAttributeValue.rValue, sis))
	}
	return strings.Join(strs, "\r\n")
}

var gEC = NewGeoErrorCreator("")

func (geoTif GeoTif) readFile(offset FileOffset, dataLen int) ([]byte, error) {
	data := make([]byte, dataLen, dataLen)
	n, err := geoTif.tFile.ReadAt(data, offset)
	if n != dataLen {
		return data, gEC(WithFunction("readFile"), WithErrorText(fmt.Sprintf("read [%d] bytes, but only get [%d] bytes", dataLen, n)))
	}
	if err != nil {
		return data, gEC(WithFunction("readFile"), WithError(err))
	}
	return data, nil
}
func (attributes GeoAttributes) getAttributeByTag(tag AttributeTag) (geoAttribute, error) {
	for i := 0; i < len(attributes); i++ {
		if attributes[i].Tag == tag {
			return attributes[i], nil
		}
	}
	return geoAttribute{}, gEC(WithFunction("getAttributeByTag"), WithMsg(fmt.Sprintf("can not found attribute [%v]", tag)))
}

func newGeoAttribute(data []byte, order binary.ByteOrder) (*geoAttribute, error) {
	if len(data) != 12 {
		return nil, gEC(WithFunction("readFile"), WithErrorText(fmt.Sprintf("require data len is 12, but data len is %d", len(data))))
	}
	gAttribute := geoAttribute{}
	gAttribute.Tag = AttributeTag(order.Uint16(data[0:2]))
	gAttribute.Type = DataType(order.Uint16(data[2:4]))
	gAttribute.Len = order.Uint32(data[4:8])
	gAttribute.SourceValue = data[8:12]
	gAttribute.Offset = 0
	return &gAttribute, nil
}

func (gAttribute *geoAttribute) parseValue(order binary.ByteOrder) error {
	gAttribute.GeoAttributeValue = geoAttributeValue{
		BYTE:   nil,
		ASCII:  "",
		SHORT:  nil,
		LONG:   nil,
		FLOAT:  nil,
		DOUBLE: nil,
		uint:   []uint{0},
	}
	switch gAttribute.Type {
	case BYTE:
		gAttribute.GeoAttributeValue.BYTE = gAttribute.SourceValue
		gAttribute.GeoAttributeValue.rValue = gAttribute.SourceValue
		gAttribute.GeoAttributeValue.uint = make([]uint, gAttribute.Len)
		for i := 0; i < int(gAttribute.Len); i++ {
			gAttribute.GeoAttributeValue.uint[i] = uint(gAttribute.SourceValue[i])
		}
	case ASCII:
		gAttribute.GeoAttributeValue.ASCII = string(gAttribute.SourceValue)
		gAttribute.GeoAttributeValue.rValue = string(gAttribute.SourceValue)
		gAttribute.GeoAttributeValue.uint = make([]uint, gAttribute.Len)
		for i := 0; i < int(gAttribute.Len); i++ {
			gAttribute.GeoAttributeValue.uint[i] = uint(gAttribute.SourceValue[i])
		}
	case SHORT:
		gAttribute.GeoAttributeValue.SHORT = make([]uint16, gAttribute.Len)
		gAttribute.GeoAttributeValue.uint = make([]uint, gAttribute.Len)
		for i := 0; i < int(gAttribute.Len); i++ {
			v := order.Uint16(gAttribute.SourceValue[i*2 : i*2+2])
			gAttribute.GeoAttributeValue.SHORT[i] = v
			gAttribute.GeoAttributeValue.uint[i] = uint(v)
		}
		gAttribute.GeoAttributeValue.rValue = gAttribute.GeoAttributeValue.SHORT
	case LONG:
		gAttribute.GeoAttributeValue.LONG = make([]uint32, gAttribute.Len)
		gAttribute.GeoAttributeValue.uint = make([]uint, gAttribute.Len)
		for i := 0; i < int(gAttribute.Len); i++ {
			v := order.Uint32(gAttribute.SourceValue[i*4 : i*4+4])
			gAttribute.GeoAttributeValue.LONG[i] = v
			gAttribute.GeoAttributeValue.uint[i] = uint(v)
		}
		gAttribute.GeoAttributeValue.rValue = gAttribute.GeoAttributeValue.LONG
	case FLOAT:
		gAttribute.GeoAttributeValue.FLOAT = make([]float32, gAttribute.Len)
		for i := 0; i < int(gAttribute.Len); i++ {
			gAttribute.GeoAttributeValue.FLOAT[i] = math.Float32frombits(order.Uint32(gAttribute.SourceValue[i*4 : i*4+4]))
		}
		gAttribute.GeoAttributeValue.rValue = gAttribute.GeoAttributeValue.FLOAT
	case DOUBLE:
		gAttribute.GeoAttributeValue.DOUBLE = make([]float64, gAttribute.Len)
		for i := 0; i < int(gAttribute.Len); i++ {
			gAttribute.GeoAttributeValue.DOUBLE[i] = math.Float64frombits(order.Uint64(gAttribute.SourceValue[i*8 : i*8+8]))
		}
		gAttribute.GeoAttributeValue.rValue = gAttribute.GeoAttributeValue.DOUBLE
	default:
		return gEC(WithFunction("parseValue"), WithError(errors.New(fmt.Sprintf("unknow datatype [%v]", gAttribute.Type))), WithMsg("to default"))
	}
	return nil
}
func (gAttribute geoAttribute) toFloat64() []float64 {
	if gAttribute.Type == DOUBLE {
		return gAttribute.GeoAttributeValue.DOUBLE
	} else if gAttribute.Type == FLOAT {
		var ret = make([]float64, gAttribute.Len)
		for i := 0; i < int(gAttribute.Len); i++ {
			ret[i] = float64(gAttribute.GeoAttributeValue.FLOAT[i])
		}
		return ret
	} else {
		var ret = make([]float64, gAttribute.Len)
		for i := 0; i < int(gAttribute.Len); i++ {
			ret = append(ret, float64(gAttribute.GeoAttributeValue.uint[i]))
		}
		return ret
	}
}

// todo: for save to file
func (gAttribute geoAttribute) toBytes() {

}

func (gAttribute geoAttribute) getValue() interface{} {
	return gAttribute.GeoAttributeValue.rValue
}
