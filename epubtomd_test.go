package epubtomd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestConvert(t *testing.T) {
	convey.Convey("Convert", t, func() {
		epubPath, err := filepath.Abs("./testdata/history.epub")
		convey.So(err, convey.ShouldBeNil)
		dstDir := filepath.Join(filepath.Dir(epubPath), "history")

		convey.Convey("Convert_OK", func() {
			err := Convert(epubPath, dstDir)
			convey.So(err, convey.ShouldBeNil)
			_ = os.RemoveAll(dstDir)
		})
	})
}
