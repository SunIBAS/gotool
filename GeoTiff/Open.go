package GeoTiff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

func OpenGeoTif(FilePath string) (*GeoTif, error) {
	var gEC = NewGeoErrorCreator("OpenGeoTif")
	geoTif := GeoTif{
		FilePath:     FilePath,
		GeoTifHeader: geoTifHeader{},
	}
	if err := geoTif.open(); err != nil {
		return nil, gEC(WithError(err))
	}
	return &geoTif, nil
}
func (g *GeoTif) open() error {
	//var gEC = NewGeoErrorCreator("open")
	var err error
	g.tFile, err = os.Open(g.FilePath)
	if err != nil {
		return gEC(WithError(err))
	}
	if err := g.checkBigOrLittle(); err != nil {
		return gEC(WithError(err))
	}
	// get header offset

	//if err = binary.Read(g.tFile, g.byteOrder, &g.GeoTifHeader.offset); err != nil {
	//	return gEC(WithError(err))
	//}
	if err = g.readAttribute(); err != nil {
		return gEC(WithFunction("open"), WithError(err))
	}
	if err = g.parseGeoKeys(); err != nil {
		return gEC(WithFunction("open"), WithError(err))
	}

	return nil
}

// https://www.nationalarchives.gov.uk/PRONOM/Format/proFormatSearch.aspx?status=detailReport&id=798&strPageToDisplay=signatures
func (g *GeoTif) checkBigOrLittle() error {
	//var gEC = NewGeoErrorCreator("checkBigOrLittle")
	data, err := g.readFile(0, 4)
	if err != nil {
		return gEC(WithError(err))
	}
	var byteOrder uint32 = binary.BigEndian.Uint32(data)
	switch byteOrder {
	case littleEndian:
		g.byteOrder = binary.LittleEndian
	case bigEndian:
		g.byteOrder = binary.BigEndian
	default:
		return gEC(WithFunction("checkBigOrLittle"), WithError(errors.New(fmt.Sprintf("undefined byte order [% x]", byteOrder))))
	}
	return nil
}

func (g *GeoTif) readAttribute() error {
	//var gEC = NewGeoErrorCreator("readAttribute")
	data, err := g.readFile(4, 4)
	if err != nil {
		return gEC(WithFunction("readAttribute"), WithError(err))
	}
	g.GeoTifHeader.offset = int64(g.byteOrder.Uint32(data))

	if g.GeoTifHeader.offset != 0 {
		// get the number of attribute
		var numAttribute int64
		data, err := g.readFile(g.GeoTifHeader.offset, 2)
		if err != nil {
			return gEC(WithFunction("readAttribute"), WithError(err))
		}
		numAttribute = int64(g.byteOrder.Uint16(data))

		g.GeoTifHeader.Attribute = make([]geoAttribute, numAttribute, numAttribute)
		// read all attribute to []byte
		attributeBytes, err := g.readFile(g.GeoTifHeader.offset+2, int(geoFileAttributeSize*numAttribute))
		if err != nil {
			return gEC(WithFunction("readAttribute"), WithError(err))
		}
		// to parse attribute
		for i := 0; i < int(numAttribute); i++ {
			gAttribute, err := newGeoAttribute(attributeBytes[i*12:i*12+12], g.byteOrder)
			if err != nil {
				return gEC(WithFunction("readAttribute"), WithError(err), WithErrorText(fmt.Sprintf("for[%d]", i)))
			}
			totalBytes := gAttribute.Bytes()
			if totalBytes > 4 {
				gAttribute.Offset = g.byteOrder.Uint32(gAttribute.SourceValue)
				realSourceData, err := g.readFile(FileOffset(gAttribute.Offset), int(totalBytes))
				if err != nil {
					return gEC(WithFunction("readAttribute"), WithError(err), WithErrorText(fmt.Sprintf("for[%d] read realSourceData", i)))
				}
				gAttribute.SourceValue = realSourceData
			}
			if err := gAttribute.parseValue(g.byteOrder); err != nil {
				return gEC(WithFunction("readAttribute"), WithError(err), WithErrorText(fmt.Sprintf("for[%d] parse value", i)))
			}
			g.GeoTifHeader.Attribute[i] = *gAttribute
		}

	}
	return nil
}

func (g *GeoTif) parseGeoKeys() error {
	geoKeyDirectoryAtr, err := g.GeoTifHeader.Attribute.getAttributeByTag(GeoKeyDirectory)
	var geoDoubleDirectoryAtr geoAttribute
	var geoASCIIDirectoryAtr geoAttribute
	if err != nil {
		return gEC(WithFunction("parseGeoKeys"), WithError(err))
	} else {
		geoKeyDirectoryValue := geoKeyDirectoryAtr.GeoAttributeValue.SHORT
		if geoKeyDirectoryValue[3] > 0 {
			geoKeyLen := int(geoKeyDirectoryValue[3])
			g.GeoKeys = make([]geoAttribute, geoKeyLen, geoKeyLen)
			for i := 0; i < geoKeyLen; i++ {
				fromIndex := 4*i + 4
				gAttribute := geoAttribute{
					Len:    uint32(geoKeyDirectoryValue[2+fromIndex]),
					Offset: 0,
				}
				if geoKeyDirectoryValue[fromIndex+1] == 0 {
					b := make([]byte, 2)
					g.byteOrder.PutUint16(b, uint16(geoKeyDirectoryValue[fromIndex+3]))
					gAttribute.SourceValue = b
					gAttribute.Type = SHORT
					gAttribute.parseValue(g.byteOrder)
				} else if geoKeyDirectoryValue[fromIndex+1] == uint16(GeoDoubleParams) {
					if geoDoubleDirectoryAtr.Tag == 0 {
						geoDoubleDirectoryAtr, err = g.GeoTifHeader.Attribute.getAttributeByTag(GeoDoubleParams)
						if err != nil {
							return gEC(WithFunction("parseGeoKeys"), WithMsg("get geoDoubleDirectoryAtr fail"), WithError(err))
						}
					}
					gAttribute.Offset = uint32(geoKeyDirectoryValue[3+fromIndex])
					gAttribute.SourceValue = geoDoubleDirectoryAtr.SourceValue[gAttribute.Offset*8 : gAttribute.Offset*8+gAttribute.Len]
					gAttribute.Type = DOUBLE
					gAttribute.parseValue(g.byteOrder)
				} else if geoKeyDirectoryValue[fromIndex+1] == uint16(GeoASCIIParams) {
					if geoASCIIDirectoryAtr.Tag == 0 {
						geoASCIIDirectoryAtr, err = g.GeoTifHeader.Attribute.getAttributeByTag(GeoASCIIParams)
						if err != nil {
							return gEC(WithFunction("parseGeoKeys"), WithMsg("get geoASCIIDirectoryAtr fail"), WithError(err))
						}
					}
					gAttribute.Offset = uint32(geoKeyDirectoryValue[3+fromIndex])
					gAttribute.SourceValue = geoASCIIDirectoryAtr.SourceValue[gAttribute.Offset : gAttribute.Offset+gAttribute.Len]
					gAttribute.Type = ASCII
					gAttribute.parseValue(g.byteOrder)
				}
				g.GeoKeys[i] = gAttribute
			}
		}
	}
	return nil
}
