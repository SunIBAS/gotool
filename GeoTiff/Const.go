package GeoTiff

const (
	littleEndian uint32 = 0x49492A00
	bigEndian    uint32 = 0x4D4D002A
)

type AttributeTag uint16

const (
	ImageWidth                AttributeTag = 256 // ImageWidth
	ImageLength               AttributeTag = 257 // ImageLength
	BitsPerSample             AttributeTag = 258 // BitsPerSample
	Compression               AttributeTag = 259 // Compression
	PhotometricInterpretation AttributeTag = 262 // PhotometricInterpretation
	FillOrder                 AttributeTag = 266 // FillOrder
	StripOffsets              AttributeTag = 273 // StripOffsets
	SamplesPerPixel           AttributeTag = 277 // SamplesPerPixel
	RowsPerStrip              AttributeTag = 278 // RowsPerStrip
	StripByteCounts           AttributeTag = 279 // StripByteCounts
	PlanarConfiguration       AttributeTag = 284 // PlanarConfiguration
	T4Options                 AttributeTag = 292 // CCITT Group 3 options, a set of 32 flag bits.
	T6Options                 AttributeTag = 293 // CCITT Group 4 options, a set of 32 flag bits.
	TileWidth                 AttributeTag = 322 // TileWidth
	TileLength                AttributeTag = 323 // TileLength
	TileOffsets               AttributeTag = 324 // TileOffsets
	TileByteCounts            AttributeTag = 325 // TileByteCounts
	XResolution               AttributeTag = 282 // XResolution
	YResolution               AttributeTag = 283 // YResolution
	ResolutionUnit            AttributeTag = 296 // ResolutionUnit
	Predictor                 AttributeTag = 317 // Predictor
	ColorMap                  AttributeTag = 320 // ColorMap
	ExtraSamples              AttributeTag = 338 // ExtraSamples
	SampleFormat              AttributeTag = 339 // SampleFormat

	// GeoTIFF Specific Tags
	GeoKeyDirectory     AttributeTag = 34735
	GeoDoubleParams     AttributeTag = 34736
	GeoASCIIParams      AttributeTag = 34737
	ModelPixelScale     AttributeTag = 33550
	ModelTiepoint       AttributeTag = 33922
	ModelTransformation AttributeTag = 34264
)

// From the Tiff 6.0 Specification (p.16)
type DataType int16

const (
	NONE      DataType = 0  // None      =  Default for no field type
	BYTE      DataType = 1  // BYTE      =  8-bit unsigned integer.
	ASCII     DataType = 2  // ASCII     =  8-bit byte that contains a 7-bit ASCII code; the last byte must be NULL (binary zero).
	SHORT     DataType = 3  // SHORT     =  16-bit (2-byte) unsigned integer.
	LONG      DataType = 4  // LONG      =  32-bit (4-byte) unsigned integer.
	RATIONAL  DataType = 5  // RATIONAL  =  Two LONGS: the first represents the numerator of a
	SBYTE     DataType = 6  // SBYTE     =  An 8-bit signed (twos-complement) integer.
	UNDEFINED DataType = 7  // UNDEFINED =  An 8-bit byte that may contain anything, depending on the definition of the field.
	SSHORT    DataType = 8  // SSHORT    =  A 16-bit (2-byte) signed (twos-complement) integer.
	SLONG     DataType = 9  // SLONG     =  A 32-bit (4-byte) signed (twos-complement) integer.
	SRATIONAL DataType = 10 // SRATIONAL =  Two SLONG’s: the first represents the numerator of a fraction, the second the denominator.
	FLOAT     DataType = 11 // FLOAT     =  Single precision (4-byte) IEEE format.
	DOUBLE    DataType = 12 // DOUBLE    =  Double precision (8-byte) IEEE format
)
const (
	zeroByte  = 0
	oneByte   = 1
	twoByte   = 2
	fourByte  = 4
	eightByte = 8
)

var DataTypeLen = [...]uint32{
	zeroByte, oneByte, oneByte, twoByte,
	fourByte, eightByte, oneByte, oneByte,
	twoByte, fourByte, eightByte, fourByte, eightByte,
}

func (dt DataType) Bytes() uint32 {
	if dt == 0 || int(dt) > len(DataTypeLen) {
		return DataTypeLen[0]
	}
	return DataTypeLen[int(dt)]
}
