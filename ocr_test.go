package lookup

import (
	"embed"
	_ "embed"
	. "github.com/smartystreets/goconvey/convey"
	"image"
	_ "image/png"
	"testing"
)

//go:embed testdata
var reconfiguration embed.FS

func TestOCR(t *testing.T) {
	Convey("Given an OCR object", t, func() {
		ocr := NewOCR(0.8)

		Convey("When I try to load an invalid font directory", func() {
			err := ocr.LoadFont("testdata/NON_EXISTENT")

			Convey("It returns an error", func() {
				So(err.Error(), ShouldNotBeNil)
			})
		})

		Convey("When I load a valid font on it", func() {
			err := ocr.LoadFont("testdata/font_1")

			Convey("It loads the fonts successfully", func() {
				So(err, ShouldBeNil)
			})

			Convey("It stores the font family", func() {
				So(ocr.fontFamilies, ShouldContainKey, "font_1")
				So(ocr.fontFamilies, ShouldHaveLength, 1)
				So(ocr.fontFamilies["font_1"], ShouldHaveLength, 13)
			})

			Convey("It updates the totalSymbols", func() {
				So(ocr.allSymbols, ShouldHaveLength, 13)
			})

			Convey("And when I pass an image to be recognized", func() {
				img := loadImageColor("testdata/test3.png")
				text, _ := ocr.Recognize(img)

				Convey("It recognizes the text in the image", func() {
					So(text, ShouldEqual, "3662\n32€/€")
				})
			})

			Convey("And when I pass an subimage to be recognized", func() {
				img := loadImageColor("testdata/full.png")
				text, _ := ocr.Recognize(img.(*image.NRGBA).SubImage(image.Rect(1280, 646, 1280+61, 646+31)))

				Convey("It only recognizes the text inside the subimage", func() {
					So(text, ShouldEqual, "4339")
				})
			})
		})

	})
}

func BenchmarkOCR(b *testing.B) {
	b.StopTimer()
	ocr := NewOCR(0.7)
	if err := ocr.LoadFont("testdata/font_1"); err != nil {
		panic(err)
	}
	img := loadImageGray("testdata/test3.png")
	if _, err := ocr.Recognize(img); err != nil {
		panic(err)
	}
	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ocr.Recognize(img)
	}
}

func BenchmarkOCRParallel(b *testing.B) {
	b.StopTimer()
	ocr := NewOCR(0.7, 5)
	if err := ocr.LoadFont("testdata/font_1"); err != nil {
		panic(err)
	}
	img := loadImageGray("testdata/test3.png")
	if _, err := ocr.Recognize(img); err != nil {
		panic(err)
	}
	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ocr.Recognize(img)
	}
}

func TestOCR_LoadFontFromFs(t *testing.T) {
	Convey("Given an OCR object", t, func() {
		ocr := NewOCR(0.8)

		Convey("When I try to load an invalid font directory", func() {
			err := ocr.LoadFontFromFs(reconfiguration, "/")

			Convey("It returns an error", func() {
				So(err.Error(), ShouldNotBeNil)
			})
		})

		Convey("When I load a valid font on it", func() {
			err := ocr.LoadFont("testdata/font_1")

			Convey("It loads the fonts successfully", func() {
				So(err, ShouldBeNil)
			})

			Convey("It stores the font family", func() {
				So(ocr.fontFamilies, ShouldContainKey, "font_1")
				So(ocr.fontFamilies, ShouldHaveLength, 1)
				So(ocr.fontFamilies["font_1"], ShouldHaveLength, 13)
			})

			Convey("It updates the totalSymbols", func() {
				So(ocr.allSymbols, ShouldHaveLength, 13)
			})

			Convey("And when I pass an image to be recognized", func() {
				img := loadImageColor("testdata/test3.png")
				text, _ := ocr.Recognize(img)

				Convey("It recognizes the text in the image", func() {
					So(text, ShouldEqual, "3662\n32€/€")
				})
			})

			Convey("And when I pass an subimage to be recognized", func() {
				img := loadImageColor("testdata/full.png")
				text, _ := ocr.Recognize(img.(*image.NRGBA).SubImage(image.Rect(1280, 646, 1280+61, 646+31)))

				Convey("It only recognizes the text inside the subimage", func() {
					So(text, ShouldEqual, "4339")
				})
			})
		})

	})
}
