package epubtomd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestS3ImageHandler_CopyImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<PutObjectResult><ETag>"mock-etag"</ETag></PutObjectResult>`))
	}))
	defer server.Close()

	convey.Convey("S3ImageHandler_CopyImage", t, func() {

		fullTestDataPath, err := filepath.Abs("./testdata/")
		convey.So(err, convey.ShouldBeNil)

		r := os.DirFS(fullTestDataPath)

		ins, err := NewS3ImageHandler(r, "test", "test", "test", "test", server.URL, "https://www.iminho.me")

		convey.So(err, convey.ShouldBeNil)
		convey.Convey("S3ImageHandler_CopyImage_OK", func() {

			ret, err := ins.CopyImage("test.jpeg", "images/test.jpeg")
			convey.So(err, convey.ShouldBeNil)
			convey.So(ret, convey.ShouldNotBeNil)
			convey.So(ret, convey.ShouldEqual, "https://www.iminho.me/images/test.jpeg")
		})
		convey.Convey("S3ImageHandler_CopyImage_Fail", func() {
			ret, err := ins.CopyImage("testtest.jpeg", "images/test.jpeg")
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(ret, convey.ShouldBeEmpty)
		})
	})
}
