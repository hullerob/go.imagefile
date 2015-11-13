package imagefile

import (
	"bytes"
	"image"
	"io"
	"testing"
)

var decodeTests = []struct {
	descr string
	input []byte
	img   []byte
	w, h  int
}{
	{
		"empty image",
		[]byte("imagefile\000\000\000\000\000\000\000\000"),
		[]byte{},
		0, 0,
	},
	{
		"image 1x1",
		[]byte("imagefile\000\000\000\001\000\000\000\001rgba"),
		[]byte("rgba"),
		1, 1,
	},
}

func TestDecodeImage(t *testing.T) {
	for _, test := range decodeTests {
		r := bytes.NewBuffer(test.input)
		img, err := Decode(r)
		if err != nil {
			t.Errorf("%s: err not nil: %v", test.descr, err)
			continue
		}
		nrgba, ok := img.(*image.NRGBA)
		if !ok {
			t.Errorf("%s: wrong image type")
			continue
		}
		if dx, dy := img.Bounds().Dx(), img.Bounds().Dy(); dx != test.w || dy != test.h {
			t.Errorf("%s: wrong image size: %d x %d | expected: %d x %d",
				test.descr, dx, dy, test.w, test.h)
		}
		if 0 != bytes.Compare(test.img, nrgba.Pix) {
			t.Errorf("%s: image data differs", test.descr)
		}
	}
}

func TestDecodeImageBadHeader(t *testing.T) {
	r := bytes.NewBuffer([]byte("imagefile\000\000\000\000\000\000\000"))
	img, err := Decode(r)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("err: %#v | expected: %#v", err, io.ErrUnexpectedEOF)
	}
	if img != nil {
		t.Errorf("img is not nil")
	}
}

func TestDecodeImageBad(t *testing.T) {
	r := bytes.NewBuffer([]byte("imagefile\000\000\000\001\000\000\000\001rgb"))
	_, err := Decode(r)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("err: %#v | expected: %#v", err, io.ErrUnexpectedEOF)
	}
}
