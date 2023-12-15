package GeoTiff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
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
	if err = g.initMeta(); err != nil {
		return gEC(WithFunction("open"), WithError(err))
	}
	if err = g.readData(); err != nil {
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
	geoKeyDirectoryAtr, err := g.GeoTifHeader.Attribute.getAttributeByTag(GeoKeyDirectoryTag)
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
				} else if geoKeyDirectoryValue[fromIndex+1] == uint16(GeoDoubleParamsTag) {
					if geoDoubleDirectoryAtr.Tag == 0 {
						geoDoubleDirectoryAtr, err = g.GeoTifHeader.Attribute.getAttributeByTag(GeoDoubleParamsTag)
						if err != nil {
							return gEC(WithFunction("parseGeoKeys"), WithMsg("get geoDoubleDirectoryAtr fail"), WithError(err))
						}
					}
					gAttribute.Offset = uint32(geoKeyDirectoryValue[3+fromIndex])
					gAttribute.SourceValue = geoDoubleDirectoryAtr.SourceValue[gAttribute.Offset*8 : gAttribute.Offset*8+gAttribute.Len]
					gAttribute.Type = DOUBLE
					gAttribute.parseValue(g.byteOrder)
				} else if geoKeyDirectoryValue[fromIndex+1] == uint16(GeoAsciiParamsTag) {
					if geoASCIIDirectoryAtr.Tag == 0 {
						geoASCIIDirectoryAtr, err = g.GeoTifHeader.Attribute.getAttributeByTag(GeoAsciiParamsTag)
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

// check all the value of the arr is the same value
// 【16 16 16 16】 => (16, true)
func checkAllIsFirst(arr []uint) (uint, bool) {
	f := arr[0]
	for i := 1; i < len(arr); i++ {
		if f != arr[i] {
			return 0, false
		}
	}
	return f, true
}
func (g *GeoTif) initMeta() error {
	var errs []error = []error{}
	var getValue = func(at AttributeTag) []uint {
		atr, err := g.GeoTifHeader.Attribute.getAttributeByTag(at)
		if err != nil {
			errs = append(errs, gEC(WithError(err), WithFunction("initMeta")))
		}
		return atr.GeoAttributeValue.uint
	}

	var err error
	var atr geoAttribute

	g.Meta = Meta{
		Columns:           getValue(ImageWidth)[0],
		Rows:              getValue(ImageLength)[0],
		PhotometricInterp: getValue(PhotometricInterpretation)[0],
		samplesPerPixel:   getValue(SamplesPerPixel)[0],
		SampleFormat:      getValue(SampleFormat)[0],
		TiepointData:      TiepointTransformationParameters{},

		BitsPerSample:     nil,
		RasterPixelIsArea: false,

		EPSGCode:    0,
		NodataValue: "",
		palette:     nil,
		mode:        0,
	}
	if len(errs) > 0 {
		return gEC(WithFunction("initMeta"), WithError(errs[0]))
	}
	if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(BitsPerSample); err == nil {
		g.Meta.BitsPerSample = atr.GeoAttributeValue.uint
	}
	// See if geokeys has GTRasterTypeGeoKey
	if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(GTRasterTypeGeoKey); err == nil {
		v := atr.GeoAttributeValue.uint
		if v[0] == 1 {
			g.Meta.RasterPixelIsArea = true
		} else {
			g.Meta.RasterPixelIsArea = false
		}
	}
	// EPSG code
	if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(ProjectionGeoKey); err == nil {
		g.Meta.EPSGCode = atr.GeoAttributeValue.uint[0]
	} else if atr, err := g.GeoTifHeader.Attribute.getAttributeByTag(GeographicTypeGeoKey); err == nil {
		g.Meta.EPSGCode = atr.GeoAttributeValue.uint[0]
	}
	// nodata
	if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(GDAL_NODATA); err == nil {
		g.Meta.NodataValue = atr.GeoAttributeValue.ASCII
	}

	var ok bool
	switch g.Meta.PhotometricInterp {
	case PI_RGB:
		if g.Meta.BitsPerSample[0], ok = checkAllIsFirst(g.Meta.BitsPerSample); !ok {
			return gEC(WithFunction("initMeta"), WithError(err))
		}
		l := len(g.Meta.BitsPerSample)
		if l == 3 {
			g.Meta.mode = mRGB
		} else if l == 4 {
			if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(ExtraSamples); err == nil {
				v := atr.GeoAttributeValue.uint[0]
				if v == 1 {
					g.Meta.mode = mRGBA
				} else if v == 2 {
					g.Meta.mode = mNRGBA
				} else {
					return gEC(WithFunction("initMeta"), WithErrorText(fmt.Sprintf("wrong number of samples for RGB")))
				}
			} else {
				return gEC(WithFunction("initMeta"), WithErrorText(fmt.Sprintf("wrong number of samples for RGB")))
			}
		} else {
			return gEC(WithFunction("initMeta"), WithErrorText(fmt.Sprintf("wrong number of samples for RGB,require 3 or 4,but get %d", l)))
		}
	case PI_Paletted:
		g.Meta.mode = mPaletted
		if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(ColorMap); err == nil {
			numColors := len(atr.GeoAttributeValue.uint) / 3
			vals := atr.GeoAttributeValue.uint
			// the color number should be in 0~256 and the number of all value should be the integer of the 3-time
			if numColors <= 0 || numColors >= 256 || len(vals)%3 != 0 {
				return gEC(WithFunction("initMeta"), WithErrorText(fmt.Sprintf("require 0 < numColors < 256, but is %d and len(colors)%%3 should be 0, but %d", numColors, len(atr.GeoAttributeValue.uint)%3)))
			}
			g.Meta.palette = make([]uint32, numColors)
			for i := 0; i < numColors; i++ {
				red := uint32(float64(vals[i]) / 65535 * 255)
				green := uint32(float64(vals[i+numColors]) / 65535 * 255)
				blue := uint32(float64(vals[i+numColors*2]) / 65535 * 255)
				a := uint32(255)
				val := uint32((a << 24) | (red << 16) | (green << 8) | blue)
				g.Meta.palette[i] = val
			}
		} else {
			return gEC(WithFunction("initMeta"), WithErrorText(fmt.Sprintf("could not found the colormap")))
		}
	case PI_WhiteIsZero:
		g.Meta.mode = mGrayInvert
	case PI_BlackIsZero:
		g.Meta.mode = mGray
	default:
		return gEC(WithFunction("initMeta"), WithErrorText(fmt.Sprintf("unkonw image format:[%d]", g.Meta.PhotometricInterp)))
	}
	return nil
}

func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
func (g *GeoTif) readData() error {
	var err error
	var atr geoAttribute
	var compressionType uint
	atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(Compression)
	if err == nil {
		compressionType = atr.GeoAttributeValue.uint[0]
	} else {
		compressionType = 0
	}
	width := int(g.Meta.Columns)
	height := int(g.Meta.Rows)
	gData := GeoData{
		buf:  []byte{},
		off:  0,
		Data: make([]float64, width*height),
	}
	blockPadding := false
	blockWidth := width
	blockHeight := height
	blocksAcross := 1
	blocksDown := 1
	var blockOffsets, blockCounts []uint

	atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(TileWidth)
	if atr.GeoAttributeValue.uint[0] != 0 {
		blockPadding = true
		blockWidth = int(atr.GeoAttributeValue.uint[0])
		if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(TileLength); err == nil {
			blockHeight = int(atr.GeoAttributeValue.uint[0])
		} else {
			return gEC(WithFunction("readData"), WithErrorText("can not found TileLength"))
		}

		blocksAcross = (width + blockWidth - 1) / blockWidth
		blocksDown = (height + blockHeight - 1) / blockHeight

		if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(TileOffsets); err == nil {
			blockOffsets = atr.GeoAttributeValue.uint
		}
		if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(TileByteCounts); err == nil {
			blockCounts = atr.GeoAttributeValue.uint
		}
	} else {
		if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(RowsPerStrip); err != nil {
			v := int(atr.GeoAttributeValue.uint[0])
			blockHeight = v
		}
		blocksDown = (height + blockHeight - 1) / blockHeight
		if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(TileOffsets); err != nil {
			blockOffsets = atr.GeoAttributeValue.uint
		}
		if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(TileByteCounts); err != nil {
			blockCounts = atr.GeoAttributeValue.uint
		}
	}

	tPredictor := uint(0)
	if atr, err = g.GeoTifHeader.Attribute.getAttributeByTag(Predictor); err == nil {
		tPredictor = atr.GeoAttributeValue.uint[0]
	}

	//bConverter := bufConverter{
	//	order: g.byteOrder,
	//}
	gDataReader := geoDataReader{
		tFile:           g.tFile,
		byteOrder:       g.byteOrder,
		compressionType: compressionType,
	}
	for i := 0; i < blocksAcross; i++ {
		blkW := blockWidth
		if !blockPadding && i == blocksAcross-1 && width%blockWidth != 0 {
			blkW = width % blockWidth
		}
		for j := 0; j < blocksDown; j++ {
			blkH := blockHeight
			if !blockPadding && j == blocksDown-1 && height%blockHeight != 0 {
				blkH = height % blockHeight
			}
			offset := int64(blockOffsets[j*blocksAcross+i])
			n := int64(blockCounts[j*blocksAcross+i])
			//switch compressionType {
			//case cNone:
			//	if b, ok := g.tFile.(*buffer); ok {
			//		gData.buf, err = b.Slice(int(offset), int(n))
			//	} else {
			//		gData.buf = make([]byte, n)
			//		_, err = g.tFile.ReadAt(gData.buf, offset)
			//	}
			//case cLZW:
			//	r := lzw.NewReader(io.NewSectionReader(g.tFile, offset, n), lzw.MSB, 8)
			//	defer r.Close()
			//	gData.buf, err = ioutil.ReadAll(r)
			//	if err != nil {
			//		println(err)
			//		//println("Block X: ", i, "Block Y: ", j, "Offset: ", offset, "n: ", n, "buf len: ", len(g.buf))
			//		//	panic(err)
			//	}
			//case cDeflate, cDeflateOld:
			//	r, err := zlib.NewReader(io.NewSectionReader(g.tFile, offset, n))
			//	if err != nil {
			//		return gEC(WithFunction("readData"), WithError(err))
			//	}
			//	gData.buf, err = ioutil.ReadAll(r)
			//	r.Close()
			//case cPackBits:
			//
			//default:
			//	return gEC(WithErrorText(fmt.Sprintf("Unsupported compression value %d", compressionType)))
			//
			//}
			xmin := i * blockWidth
			ymin := j * blockHeight
			xmax := xmin + blkW
			ymax := ymin + blkH

			xmax = minInt(xmax, width)
			ymax = minInt(ymax, height)
			gData.buf, err = gDataReader.read(offset, n)
			if err != nil {
				return gEC(WithFunction("readData"), WithError(err))
			}

			gData.off = 0

			if tPredictor == prHorizontal {
				// does it make sense to extend this to 32 and 64 bits?
				if g.Meta.BitsPerSample[0] == 16 {
					var off int
					spp := len(g.Meta.BitsPerSample) // samples per pixel
					bpp := spp * 2                   // bytes per pixel
					for y := ymin; y < ymax; y++ {
						off += spp * 2
						for x := 0; x < (xmax-xmin-1)*bpp; x += 2 {
							v0 := g.byteOrder.Uint16(gData.buf[off-bpp : off-bpp+2])
							v1 := g.byteOrder.Uint16(gData.buf[off : off+2])
							g.byteOrder.PutUint16(gData.buf[off:off+2], v1+v0)
							off += 2
						}
					}
				} else if g.Meta.BitsPerSample[0] == 8 {
					var off int
					spp := len(g.Meta.BitsPerSample) // samples per pixel
					for y := ymin; y < ymax; y++ {
						off += spp
						for x := 0; x < (xmax-xmin-1)*spp; x++ {
							gData.buf[off] += gData.buf[off-spp]
							off++
						}
					}
				}
			}

			switch g.Meta.mode {
			case mGray, mGrayInvert:
				switch g.Meta.SampleFormat {
				case 1: // Unsigned integer data
					switch g.Meta.BitsPerSample[0] {
					case 8:
						for y := ymin; y < ymax; y++ {
							for x := xmin; x < xmax; x++ {
								i := y*width + x
								gData.Data[i] = float64(gData.buf[gData.off])
								gData.off++
							}
						}
					case 16:
						for y := ymin; y < ymax; y++ {
							for x := xmin; x < xmax; x++ {
								value := g.byteOrder.Uint16(gData.buf[gData.off : gData.off+2])
								i := y*width + x
								gData.Data[i] = float64(value)
								gData.off += 2
							}
						}
					case 32:
						for y := ymin; y < ymax; y++ {
							for x := xmin; x < xmax; x++ {
								value := g.byteOrder.Uint32(gData.buf[gData.off : gData.off+4])
								i := y*width + x
								gData.Data[i] = float64(value)
								gData.off += 4
							}
						}
					case 64:
						for y := ymin; y < ymax; y++ {
							for x := xmin; x < xmax; x++ {
								value := g.byteOrder.Uint64(gData.buf[gData.off : gData.off+8])
								i := y*width + x
								gData.Data[i] = float64(value)
								gData.off += 8
							}
						}
					default:
						return gEC(WithFunction("readData"), WithErrorText("Unsupported data format"))
					}
				case 2: // Signed integer data
					switch g.Meta.BitsPerSample[0] {
					case 8:
						for y := ymin; y < ymax; y++ {
							for x := xmin; x < xmax; x++ {
								i := y*width + x
								gData.Data[i] = float64(int8(gData.buf[gData.off]))
								gData.off++
							}
						}
					case 16:
						for y := ymin; y < ymax; y++ {
							for x := xmin; x < xmax; x++ {
								value := int16(g.byteOrder.Uint16(gData.buf[gData.off : gData.off+2]))
								i := y*width + x
								gData.Data[i] = float64(value)
								gData.off += 2
							}
						}
					case 32:
						for y := ymin; y < ymax; y++ {
							for x := xmin; x < xmax; x++ {
								value := int32(g.byteOrder.Uint32(gData.buf[gData.off : gData.off+4]))
								i := y*width + x
								gData.Data[i] = float64(value)
								gData.off += 4
							}
						}
					case 64:
						for y := ymin; y < ymax; y++ {
							for x := xmin; x < xmax; x++ {
								value := int64(g.byteOrder.Uint64(gData.buf[gData.off : gData.off+8]))
								i := y*width + x
								gData.Data[i] = float64(value)
								gData.off += 8
							}
						}
					default:
						return gEC(WithFunction("readData"), WithErrorText("Unsupported data format"))
					}
				case 3: // Floating point data
					switch g.Meta.BitsPerSample[0] {
					case 32:
						for y := ymin; y < ymax; y++ {
							for x := xmin; x < xmax; x++ {
								if gData.off <= len(gData.buf) {
									bits := g.byteOrder.Uint32(gData.buf[gData.off : gData.off+4])
									float := math.Float32frombits(bits)
									i := y*width + x
									gData.Data[i] = float64(float)
									gData.off += 4
								}
							}
							if xmax*4 < blockWidth {
								gData.off = gData.off - xmax*4 + blockWidth*4
							}
						}
					case 64:
						for y := ymin; y < ymax; y++ {
							for x := xmin; x < xmax; x++ {
								if gData.off <= len(gData.buf) {
									bits := g.byteOrder.Uint64(gData.buf[gData.off : gData.off+8])
									float := math.Float64frombits(bits)
									i := y*width + x
									gData.Data[i] = float
									gData.off += 8
								}
							}
						}
					default:
						return gEC(WithFunction("readData"), WithErrorText("Unsupported data format"))
					}
				default:
					return gEC(WithFunction("readData"), WithErrorText("Unsupported sample format"))
				}
			case mPaletted:
				for y := ymin; y < ymax; y++ {
					for x := xmin; x < xmax; x++ {
						i := y*width + x
						val := int(gData.buf[gData.off])
						gData.Data[i] = float64(g.Meta.palette[val])
						gData.off++
					}
				}

			case mRGB:
				if g.Meta.BitsPerSample[0] == 8 {
					for y := ymin; y < ymax; y++ {
						for x := xmin; x < xmax; x++ {
							red := uint32(gData.buf[gData.off])
							green := uint32(gData.buf[gData.off+1])
							blue := uint32(gData.buf[gData.off+2])
							a := uint32(255)
							gData.off += 3
							i := y*width + x
							val := uint32((a << 24) | (red << 16) | (green << 8) | blue)
							gData.Data[i] = float64(val)
						}
					}
				} else if g.Meta.BitsPerSample[0] == 16 {
					for y := ymin; y < ymax; y++ {
						for x := xmin; x < xmax; x++ {
							// the spec doesn't talk about 16-bit RGB images so
							// I'm not sure why I bother with this. They specifically
							// say that RGB images are 8-bits per channel. Anyhow,
							// I rescale the 16-bits to an 8-bit channel for simplicity.
							red := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+0:gData.off+2])) / 65535.0 * 255.0)
							green := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+2:gData.off+4])) / 65535.0 * 255.0)
							blue := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+4:gData.off+6])) / 65535.0 * 255.0)
							a := uint32(255)
							gData.off += 6
							i := y*width + x
							val := uint32((a << 24) | (red << 16) | (green << 8) | blue)
							gData.Data[i] = float64(val)
						}
					}
				} else {
					return gEC(WithFunction("readData"), WithErrorText("Unsupported data format"))
				}
			case mNRGBA:
				if g.Meta.BitsPerSample[0] == 8 {
					for y := ymin; y < ymax; y++ {
						for x := xmin; x < xmax; x++ {
							red := uint32(gData.buf[gData.off])
							green := uint32(gData.buf[gData.off+1])
							blue := uint32(gData.buf[gData.off+2])
							a := uint32(gData.buf[gData.off+3])
							gData.off += 4
							i := y*width + x
							val := uint32((a << 24) | (red << 16) | (green << 8) | blue)
							gData.Data[i] = float64(val)
						}
					}
				} else if g.Meta.BitsPerSample[0] == 16 {
					for y := ymin; y < ymax; y++ {
						for x := xmin; x < xmax; x++ {
							red := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+0:gData.off+2])) / 65535.0 * 255.0)
							green := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+2:gData.off+4])) / 65535.0 * 255.0)
							blue := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+4:gData.off+6])) / 65535.0 * 255.0)
							a := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+6:gData.off+8])) / 65535.0 * 255.0)
							gData.off += 8
							i := y*width + x
							val := uint32((a << 24) | (red << 16) | (green << 8) | blue)
							gData.Data[i] = float64(val)
						}
					}
				} else {
					return gEC(WithFunction("readData"), WithErrorText("Unsupported data format"))
				}
			case mRGBA:
				if g.Meta.BitsPerSample[0] == 16 {
					for y := ymin; y < ymax; y++ {
						for x := xmin; x < xmax; x++ {
							red := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+0:gData.off+2])) / 65535.0 * 255.0)
							green := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+2:gData.off+4])) / 65535.0 * 255.0)
							blue := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+4:gData.off+6])) / 65535.0 * 255.0)
							a := uint32(float64(g.byteOrder.Uint16(gData.buf[gData.off+6:gData.off+8])) / 65535.0 * 255.0)
							gData.off += 8
							i := y*width + x
							val := uint32((a << 24) | (red << 16) | (green << 8) | blue)
							gData.Data[i] = float64(val)
						}
					}
				} else {
					for y := ymin; y < ymax; y++ {
						for x := xmin; x < xmax; x++ {
							red := uint32(gData.buf[gData.off])
							green := uint32(gData.buf[gData.off+1])
							blue := uint32(gData.buf[gData.off+2])
							a := uint32(gData.buf[gData.off+3])
							gData.off += 4
							i := y*width + x
							val := uint32((a << 24) | (red << 16) | (green << 8) | blue)
							gData.Data[i] = float64(val)
						}
					}
				}
			}
		}

	}

	g.Data = gData
	return nil
}
