package epubtomd

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestZipEpubReader_Extract(t *testing.T) {
	convey.Convey("TestZipEpubReader_Extract", t, func() {
		reader := NewZipEpubReader()
		convey.So(reader, convey.ShouldNotBeNil)
		defer reader.Close()
		convey.Convey("TestZipEpubReader_Extract_OK", func() {
			r, err := reader.Extract("./testdata/one-hundred-years-of-solitude.epub")
			convey.So(err, convey.ShouldBeNil)
			convey.So(r, convey.ShouldNotBeNil)
		})
		convey.Convey("TestZipEpubReader_Extract_Fail", func() {
			r, err := reader.Extract("./testdata/")
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(r, convey.ShouldBeNil)
		})
	})
}

func TestZipEpubReader_ParseMetadata(t *testing.T) {
	convey.Convey("TestZipEpubReader_ParseMetadata", t, func() {
		reader := NewZipEpubReader()
		convey.So(reader, convey.ShouldNotBeNil)
		defer SaleClose(reader)
		convey.Convey("TestZipEpubReader_ParseMetadata_OK", func() {
			r, err := reader.Extract("./testdata/one-hundred-years-of-solitude.epub")
			convey.So(err, convey.ShouldBeNil)
			data, err := reader.ParseMetadata(r)
			convey.So(err, convey.ShouldBeNil)
			convey.So(data, convey.ShouldNotBeNil)
			convey.So(data.Version, convey.ShouldEqual, "2.0")
			convey.So(data.Title, convey.ShouldEqual, "百年孤独(根据马尔克斯指定版本翻译,未做任何增删)")
			convey.So(data.Author, convey.ShouldEqual, "加西亚•马尔克斯,范晔")
		})
	})
}
