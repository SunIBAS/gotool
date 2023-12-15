package GeoTiff

import (
	"compress/lzw"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

type GeoData struct {
	buf  []byte
	off  int // Current offset in buf.
	Data []float64
}

type geoDataReader struct {
	tFile           io.ReaderAt
	byteOrder       binary.ByteOrder
	compressionType CompressionType
}

func (gdr geoDataReader) read(offset, size int64) ([]byte, error) {
	var buf []byte
	var err error
	switch gdr.compressionType {
	case cNone:
		buf, err = readCNone(gdr.tFile, offset, size)
	case cLZW:
		r := lzw.NewReader(io.NewSectionReader(gdr.tFile, offset, size), lzw.MSB, 8)
		defer r.Close()
		buf, err = io.ReadAll(r)
		if err = r.Close(); err != nil {
			return nil, gEC(WithFunction("geoDataReader.read"), WithError(err))
		}
	case cDeflate, cDeflateOld:
		r, err := zlib.NewReader(io.NewSectionReader(gdr.tFile, offset, size))
		if err != nil {
			return nil, gEC(WithFunction("geoDataReader.read"), WithError(err))
		}
		buf, err = io.ReadAll(r)
		if err = r.Close(); err != nil {
			return nil, gEC(WithFunction("geoDataReader.read"), WithError(err))
		}
	case cPackBits:
		return readCPackBits(gdr.tFile, offset, size)
	default:
		return nil, gEC(WithFunction("geoDataReader.read"), WithErrorText(fmt.Sprintf("Unsupported compression value %d", gdr.compressionType)))
	}
	if err != nil {
		return nil, gEC(WithFunction("geoDataReader.read"), WithError(err))
	}
	return buf, nil
}
func readCNone(tFile io.ReaderAt, offset, size int64) ([]byte, error) {
	var buf []byte
	var err error
	if b, ok := tFile.(*buffer); ok {
		buf, err = b.Slice(int(offset), int(size))
	} else {
		buf = make([]byte, size)
		_, err = tFile.ReadAt(buf, offset)
	}
	return buf, err
}

// https://github.com/grumets/MiraMonMapBrowser/blob/b997173bc0ee2ebd1d61567a0d4e33d1c44004a4/src/geotiff/compression/packbits.js#L18
// notice: L18 is (j = 0; j <= header;j++) <=== j <= header, so read (header + 1) data
func readCPackBits(tFile io.ReaderAt, offset, size int64) ([]byte, error) {
	srcBuf, err := readCNone(tFile, offset, size)
	if err == nil {
		buf := make([]byte, 0, len(srcBuf))
		pos := 0
		for pos < len(srcBuf) {
			headerByte := int(srcBuf[pos])
			if headerByte > 127 {
				var copyCount int
				copyCount = 256 - headerByte
				copyByte := srcBuf[pos+1]
				for i := 0; i <= copyCount; i++ {
					buf = append(buf, copyByte)
				}
				pos += 2
			} else {
				headerByte++
				for i := 0; i < headerByte; i++ {
					pos++
					buf = append(buf, srcBuf[pos])
				}
				pos++
			}
		}
		return buf, nil
	}
	return nil, err
}
