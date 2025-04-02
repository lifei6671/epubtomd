package epubtomd

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewBasicXHTMLConverter(t *testing.T) {
	htmlContent := "\n<div class=\"chatu\"><img alt=\"\" src=\"Image00020.jpg\"/>\n<p class=\"tushuo\">图20 创作者：杰克逊·波洛克灵感涌现的一刻</p>\n</div>\n\n\n"
	convey.Convey("TestHtmlToMarkdown", t, func() {
		result, err := HtmlToMarkdown(htmlContent)
		convey.So(err, convey.ShouldBeNil)
		fmt.Println(result)
	})
}

func TestBasicXHTMLConverter_Convert(t *testing.T) {
	convey.Convey("TestBasicXHTMLConverter_Convert", t, func() {
		reader := NewZipEpubReader()
		convey.So(reader, convey.ShouldNotBeNil)
		defer SaleClose(reader)

		r, err := reader.Extract("./testdata/one-hundred-years-of-solitude.epub")
		convey.So(err, convey.ShouldBeNil)
		data, err := reader.ParseMetadata(r)
		convey.So(err, convey.ShouldBeNil)
		convey.So(data, convey.ShouldNotBeNil)

		ins := NewBasicXHTMLConverter(r)
		convey.Convey("TestLocalImageHandler_ParseMetadata_OK", func() {
			title, ret, err := ins.Convert(data.TextFiles[4])
			convey.So(err, convey.ShouldBeNil)
			convey.So(title, convey.ShouldNotBeNil)
			convey.So(ret, convey.ShouldNotBeEmpty)
		})
		convey.Convey("TestLocalImageHandler_ParseMetadata_Fail", func() {
			title, ret, err := ins.Convert("test.html")
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(title, convey.ShouldBeEmpty)
			convey.So(ret, convey.ShouldBeEmpty)
		})
	})
}
