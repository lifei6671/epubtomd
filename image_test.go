package epubtomd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLocalImageHandler_ProcessImage(t *testing.T) {
	convey.Convey("TestLocalImageHandler_ProcessImage", t, func() {
		reader := NewZipEpubReader()
		convey.So(reader, convey.ShouldNotBeNil)
		defer SaleClose(reader)

		r, err := reader.Extract("./testdata/one-hundred-years-of-solitude.epub")
		convey.So(err, convey.ShouldBeNil)
		data, err := reader.ParseMetadata(r)
		convey.So(err, convey.ShouldBeNil)
		convey.So(data, convey.ShouldNotBeNil)
		fullTestDataPath, err := filepath.Abs("./testdata")
		convey.So(err, convey.ShouldBeNil)

		convey.Convey("TestLocalImageHandler_ParseMetadata_OK", func() {
			handler := NewLocalImageHandler(r, fullTestDataPath)
			imagePath, err := handler.CopyImage(data.ImageFiles[0], "cover.jpeg")
			defer func() {
				_ = os.Remove(filepath.Join(fullTestDataPath, imagePath))
			}()
			convey.So(err, convey.ShouldBeNil)
			convey.So(imagePath, convey.ShouldNotBeNil)
			convey.So(imagePath, convey.ShouldEqual, "cover.jpeg")

		})

		convey.Convey("TestLocalImageHandler_ProcessImage_Open_Err", func() {
			handler := NewLocalImageHandler(r, fullTestDataPath)
			imagePath, err := handler.CopyImage("image/covey.jpg", "image/test.jpg")
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(imagePath, convey.ShouldBeEmpty)
		})
	})
}
