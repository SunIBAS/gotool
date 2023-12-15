package GeoTiff

const (
	littleEndian uint32 = 0x49492A00
	bigEndian    uint32 = 0x4D4D002A
)

type AttributeTag uint16

const (
	NewSubfileType            AttributeTag = 254
	ImageWidth                AttributeTag = 256
	ImageLength               AttributeTag = 257
	BitsPerSample             AttributeTag = 258
	Compression               AttributeTag = 259
	PhotometricInterpretation AttributeTag = 262
	FillOrder                 AttributeTag = 266
	DocumentName              AttributeTag = 269
	PlanarConfiguration       AttributeTag = 284

	StripOffsets    AttributeTag = 273
	Orientation     AttributeTag = 274
	SamplesPerPixel AttributeTag = 277
	RowsPerStrip    AttributeTag = 278
	StripByteCounts AttributeTag = 279

	TileWidth      AttributeTag = 322
	TileLength     AttributeTag = 323
	TileOffsets    AttributeTag = 324
	TileByteCounts AttributeTag = 325

	XResolution    AttributeTag = 282
	YResolution    AttributeTag = 283
	ResolutionUnit AttributeTag = 296

	Software     AttributeTag = 305
	Predictor    AttributeTag = 317
	ColorMap     AttributeTag = 320
	ExtraSamples AttributeTag = 338
	SampleFormat AttributeTag = 339

	GDAL_METADATA AttributeTag = 42112
	GDAL_NODATA   AttributeTag = 42113

	ModelPixelScaleTag     AttributeTag = 33550
	ModelTransformationTag AttributeTag = 34264
	ModelTiepointTag       AttributeTag = 33922
	GeoKeyDirectoryTag     AttributeTag = 34735
	GeoDoubleParamsTag     AttributeTag = 34736
	GeoAsciiParamsTag      AttributeTag = 34737
	IntergraphMatrixTag    AttributeTag = 33920

	GTModelTypeGeoKey              AttributeTag = 1024
	GTRasterTypeGeoKey             AttributeTag = 1025
	GTCitationGeoKey               AttributeTag = 1026
	GeographicTypeGeoKey           AttributeTag = 2048
	GeogCitationGeoKey             AttributeTag = 2049
	GeogGeodeticDatumGeoKey        AttributeTag = 2050
	GeogPrimeMeridianGeoKey        AttributeTag = 2051
	GeogLinearUnitsGeoKey          AttributeTag = 2052
	GeogLinearUnitSizeGeoKey       AttributeTag = 2053
	GeogAngularUnitsGeoKey         AttributeTag = 2054
	GeogAngularUnitSizeGeoKey      AttributeTag = 2055
	GeogEllipsoidGeoKey            AttributeTag = 2056
	GeogSemiMajorAxisGeoKey        AttributeTag = 2057
	GeogSemiMinorAxisGeoKey        AttributeTag = 2058
	GeogInvFlatteningGeoKey        AttributeTag = 2059
	GeogAzimuthUnitsGeoKey         AttributeTag = 2060
	GeogPrimeMeridianLongGeoKey    AttributeTag = 2061
	ProjectedCSTypeGeoKey          AttributeTag = 3072
	PCSCitationGeoKey              AttributeTag = 3073
	ProjectionGeoKey               AttributeTag = 3074
	ProjCoordTransGeoKey           AttributeTag = 3075
	ProjLinearUnitsGeoKey          AttributeTag = 3076
	ProjLinearUnitSizeGeoKey       AttributeTag = 3077
	ProjStdParallel1GeoKey         AttributeTag = 3078
	ProjStdParallel2GeoKey         AttributeTag = 3079
	ProjNatOriginLongGeoKey        AttributeTag = 3080
	ProjNatOriginLatGeoKey         AttributeTag = 3081
	ProjFalseEastingGeoKey         AttributeTag = 3082
	ProjFalseNorthingGeoKey        AttributeTag = 3083
	ProjFalseOriginLongGeoKey      AttributeTag = 3084
	ProjFalseOriginLatGeoKey       AttributeTag = 3085
	ProjFalseOriginEastingGeoKey   AttributeTag = 3086
	ProjFalseOriginNorthingGeoKey  AttributeTag = 3087
	ProjCenterLongGeoKey           AttributeTag = 3088
	ProjCenterLatGeoKey            AttributeTag = 3089
	ProjCenterEastingGeoKey        AttributeTag = 3090
	ProjCenterNorthingGeoKey       AttributeTag = 3091
	ProjScaleAtNatOriginGeoKey     AttributeTag = 3092
	ProjScaleAtCenterGeoKey        AttributeTag = 3093
	ProjAzimuthAngleGeoKey         AttributeTag = 3094
	ProjStraightVertPoleLongGeoKey AttributeTag = 3095
	VerticalCSTypeGeoKey           AttributeTag = 4096
	VerticalCitationGeoKey         AttributeTag = 4097
	VerticalDatumGeoKey            AttributeTag = 4098
	VerticalUnitsGeoKey            AttributeTag = 4099

	Photoshop AttributeTag = 34377
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

type ImageMode int

const (
	mBilevel ImageMode = iota
	mPaletted
	mGray
	mGrayInvert
	mRGB
	mRGBA
	mNRGBA
)

type PhotoInterpretation = uint

// Photometric interpretation values (see p. 37 of the spec).
const (
	PI_WhiteIsZero PhotoInterpretation = 0
	PI_BlackIsZero PhotoInterpretation = 1
	PI_RGB         PhotoInterpretation = 2
	PI_Paletted    PhotoInterpretation = 3
	PI_TransMask   PhotoInterpretation = 4 // transparency mask
	PI_CMYK        PhotoInterpretation = 5
	PI_YCbCr       PhotoInterpretation = 6
	PI_CIELab      PhotoInterpretation = 8
)

type CompressionType = uint

// Compression types (defined in various places in the spec and supplements).
const (
	cNone       CompressionType = 1
	cCCITT      CompressionType = 2
	cG3         CompressionType = 3 // Group 3 Fax.
	cG4         CompressionType = 4 // Group 4 Fax.
	cLZW        CompressionType = 5
	cJPEGOld    CompressionType = 6 // Superseded by cJPEG.
	cJPEG       CompressionType = 7
	cDeflate    CompressionType = 8 // zlib compression.
	cPackBits   CompressionType = 32773
	cDeflateOld CompressionType = 32946 // Superseded by cDeflate.
)

// Values for the tPredictor tag (page 64-65 of the spec).
const (
	prNone       = 1
	prHorizontal = 2
)
