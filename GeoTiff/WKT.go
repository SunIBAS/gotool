package GeoTiff

import (
	"fmt"
	"math"
)

// https://github.com/OSGeo/gdal/blob/66f5be9000b7ec0182aa775f7033aa250513e594/frmts/gtiff/gt_wkt_srs.cpp#L3529C13-L3529C42
type transform struct {
	Data           [6]float64
	PixelIsPoint   bool
	PointGeoIgnore bool
	TilePoints     []tilePoints
	Resolution     [3]float64
}
type tilePoints struct {
	i, j, k float64
	x, y, z float64
}

func getAttributeAndCheck(gAttributes GeoAttributes, atrTag AttributeTag, size int) ([]float64, error) {
	if atr, err := gAttributes.getAttributeByTag(atrTag); err != nil {
		return nil, gEC(WithFunction("getAttributeAndCheck"), WithError(err))
	} else {
		f64 := atr.toFloat64()
		if len(f64) >= size {
			return f64, nil
		} else {
			return nil, gEC(WithFunction("getAttributeAndCheck"), WithErrorText(fmt.Sprintf("require len is %d, but got %d", size, len(f64))))
		}
	}
}

// Init http://geotiff.maptools.org/spec/geotiff2.6.html
func (t *transform) Init(allAttribute ...geoAttribute) error {
	t.Data[0] = 0
	t.Data[1] = 1
	t.Data[2] = 0
	t.Data[3] = 0
	t.Data[4] = 0
	t.Data[5] = 1
	attributeCount := len(allAttribute)
	if attributeCount == 0 {
		return gEC(WithFunction("Transform.Init"), WithErrorText("attribute is empty"))
	}

	t.initResolution(allAttribute...)
	var val []float64
	var err error
	if val, err = getAttributeAndCheck(allAttribute, ModelTransformationTag, 16); err == nil {
		t.Data[0] = val[3]
		t.Data[1] = val[0]
		t.Data[2] = val[1]
		t.Data[3] = val[7]
		t.Data[4] = val[4]
		t.Data[5] = val[5]
	} else if val, err = getAttributeAndCheck(allAttribute, ModelTiepointTag, 6); err == nil {
		//https://github.com/grumets/MiraMonMapBrowser/blob/b997173bc0ee2ebd1d61567a0d4e33d1c44004a4/src/geotiff/geotiffimage.js#L744
		valCount := len(val) / 6
		t.TilePoints = make([]tilePoints, valCount)
		for i := 0; i < len(val); i += 6 {
			t.TilePoints[i%6] = tilePoints{
				i: val[i],
				j: val[i+1],
				k: val[i+2],
				x: val[i+3],
				y: val[i+4],
				z: val[i+5],
			}
		}
		t.Data[0] = t.Resolution[0]
		t.Data[1] = 0
		t.Data[2] = t.TilePoints[0].x
		t.Data[3] = 0
		t.Data[4] = t.Resolution[1]
		t.Data[5] = t.TilePoints[0].y
		//fmt.Println("i don't know how programming")
	} else if val, err = getAttributeAndCheck(allAttribute, ModelPixelScaleTag, 2); err == nil {
		t.Data[1] = val[0]
		t.Data[5] = -math.Abs(val[1])
		if val, err = getAttributeAndCheck(allAttribute, ModelTiepointTag, 6); err == nil {
			t.Data[0] = val[3] - val[0]*t.Data[1]
			t.Data[3] = val[4] - val[1]*t.Data[5]

			if t.PixelIsPoint && !t.PointGeoIgnore {
				t.Data[0] -= t.Data[1]*0.5 + t.Data[2]*0.5
				t.Data[3] -= t.Data[4]*0.5 + t.Data[5]*0.5
			}
		}
	} else {
		return gEC(WithFunction("Transform.Init"), WithErrorText("can not init t.Data"))
	}
	t.Resolution[0] = t.Data[0]
	t.Resolution[1] = t.Data[4]
	t.Resolution[2] = 0
	return nil
}

func (t *transform) initResolution(allAttribute ...geoAttribute) {
	var val []float64
	var err error
	if val, err = getAttributeAndCheck(allAttribute, ModelPixelScaleTag, 3); err == nil {
		t.Resolution = [3]float64{
			val[0],
			-val[1],
			val[2],
		}
	} else if val, err = getAttributeAndCheck(allAttribute, ModelTransformationTag, 16); err == nil {
		t.Resolution = [3]float64{
			val[0],
			val[5],
			val[10],
		}
	}
}
