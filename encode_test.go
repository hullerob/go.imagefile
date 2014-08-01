package imagefile

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"testing"
)

func TestEncodeEmptyImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	w := new(bytes.Buffer)
	err := Encode(w, img)
	assert.Nil(t, err, "err is not nil: %v", err)
	assert.Equal(t, []byte("imagefile\000\000\000\000\000\000\000\000"),
		w.Bytes(), "encoded image differs")
}

func TestEncodeSmallImage(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Pix = []byte("abcd")
	w := new(bytes.Buffer)
	err := Encode(w, img)
	assert.Nil(t, err, "err is not nil: %v", err)
	assert.Equal(t, []byte("imagefile\000\000\000\001\000\000\000\001abcd"),
		w.Bytes(), "encoded image differs")
}

func TestEncodeGrayImage(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White)
	w := new(bytes.Buffer)
	err := Encode(w, img)
	assert.Nil(t, err, "err is not nil: %v", err)
	assert.Equal(t, []byte("imagefile\000\000\000\001\000\000\000\001\xff\xff\xff\xff"),
		w.Bytes(), "encoded image differs")
}

func TestEncodePremultiplied(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.NRGBA{0x33, 0x44, 0x55, 0x66})
	w := new(bytes.Buffer)
	err := Encode(w, img)
	assert.Nil(t, err, "err is not nil: %v", err)
	assert.Equal(t, []byte("imagefile\000\000\000\001\000\000\000\001"),
		w.Bytes()[:17], "encoded image differs")
	assert.True(t, testAlmostEqual(2, 0x33, int(w.Bytes()[17])), "encoded image differs too much")
	assert.True(t, testAlmostEqual(2, 0x44, int(w.Bytes()[18])), "encoded image differs too much")
	assert.True(t, testAlmostEqual(2, 0x55, int(w.Bytes()[19])), "encoded image differs too much")
	assert.True(t, testAlmostEqual(2, 0x66, int(w.Bytes()[20])), "encoded image differs too much")
}

func testAlmostEqual(diff, expected, got int) bool {
	adiff := expected - got
	if adiff < 0 {
		adiff = -adiff
	}
	return diff >= adiff
}
