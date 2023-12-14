package GeoTiff

import (
	"fmt"
	"github.com/SunIBAS/gotool/GeoTiff"
	"testing"
)

var tifFile = "C:\\Users\\11340\\Documents\\paper\\中期\\zq\\附件2：培养环节考核要求\\附件2-1.研究生开题、中期登记表和报告（模板）\\中文版\\图\\数据\\ndvi.2020.20.tif"
var golangTifFile = "C:\\Users\\11340\\go\\pkg\\mod\\golang.org\\x\\image@v0.14.0\\testdata\\bw-deflate.tiff"

func TestReadGeoTif(t *testing.T) {
	geo, err := GeoTiff.OpenGeoTif(tifFile)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(geo)
}
