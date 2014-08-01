package imagefile

import (
	"bytes"
	"github.com/stretchr/testify/assert"
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
		if !assert.Nil(t, err, "test %s: err not nil: %v", test.descr, err) {
			continue
		}
		if !assert.IsType(t, (*image.NRGBA)(nil), img, "test %s", test.descr) {
			continue
		}
		nrgba := img.(*image.NRGBA)
		assert.Equal(t, test.w, img.Bounds().Dx(), "test %s: different width", test.descr)
		assert.Equal(t, test.h, img.Bounds().Dy(), "test %s: different height", test.descr)
		assert.Equal(t, test.img, nrgba.Pix, "test %s: image data differs", test.descr)
	}
}

func TestDecodeImageBadHeader(t *testing.T) {
	r := bytes.NewBuffer([]byte("imagefile\000\000\000\000\000\000\000"))
	img, err := Decode(r)
	assert.NotNil(t, err, "err is nil")
	assert.Equal(t, io.ErrUnexpectedEOF, err, "err is wrong")
	assert.Nil(t, img, "img is not nil")
}

func TestDecodeImageBad(t *testing.T) {
	r := bytes.NewBuffer([]byte("imagefile\000\000\000\001\000\000\000\001rgb"))
	_, err := Decode(r)
	assert.NotNil(t, err, "err is nil")
	assert.Equal(t, io.ErrUnexpectedEOF, err, "err is wrong")
}
