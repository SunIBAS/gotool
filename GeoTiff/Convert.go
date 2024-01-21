package GeoTiff

import "encoding/binary"

// https://github.com/grumets/MiraMonMapBrowser/blob/b997173bc0ee2ebd1d61567a0d4e33d1c44004a4/src/geotiff/compression/packbits.js#L4

type bufConverter struct {
	order binary.ByteOrder
}

func (bfc bufConverter) ConverInt64ToFloat64(bytes []byte) float64 {
	return float64(int64(bfc.order.Uint64(bytes)))
}

func (bfc bufConverter) ConverInt32ToFloat64(bytes []byte) float64 {
	return float64(int32(bfc.order.Uint32(bytes)))
}

func (bfc bufConverter) ConverInt16ToFloat64(bytes []byte) float64 {
	return float64(int16(bfc.order.Uint16(bytes)))
}

func (bfc bufConverter) ConverInt8ToFloat64(bytes byte) float64 {
	return float64(int8(bytes))
}

func (g *GeoData) ToConver(ymin, ymax, xmin, xmax, width int, convert func([]byte) float64) {
	for y := ymin; y < ymax; y++ {
		for x := xmin; x < xmax; x++ {
			i := y*width + x
			g.Data[i] = convert(g.buf[g.off : g.off+2])
			g.off += 2
		}
	}
}
