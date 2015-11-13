package imagefile

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func TestEncodeEmptyImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	w := new(bytes.Buffer)
	err := Encode(w, img)
	if err != nil {
		t.Errorf("err is not nil: %v", err)
	}
	if 0 != bytes.Compare([]byte("imagefile\000\000\000\000\000\000\000\000"), w.Bytes()) {
		t.Errorf("encoded image differs")
	}
}

func TestEncodeSmallImage(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Pix = []byte("abcd")
	w := new(bytes.Buffer)
	err := Encode(w, img)
	if err != nil {
		t.Errorf("err is not nil: %v", err)
	}
	if 0 != bytes.Compare([]byte("imagefile\000\000\000\001\000\000\000\001abcd"), w.Bytes()) {
		t.Errorf("encoded image differs")
	}
}

func TestEncodeGrayImage(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White)
	w := new(bytes.Buffer)
	err := Encode(w, img)
	if err != nil {
		t.Errorf("err is not nil: %v", err)
	}
	if 0 != bytes.Compare([]byte("imagefile\000\000\000\001\000\000\000\001\xff\xff\xff\xff"), w.Bytes()) {
		t.Errorf("encoded image differs")
	}
}

func TestEncodeSubImage(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			c := uint8(y*4*4 + x*4)
			img.SetNRGBA(x, y, color.NRGBA{c, c + 1, c + 2, c + 3})
		}
	}
	w := new(bytes.Buffer)
	err := Encode(w, img.SubImage(image.Rect(1, 1, 3, 3)))
	if err != nil {
		t.Errorf("err is not nil: %v", err)
	}
	if 0 != bytes.Compare([]byte("imagefile\000\000\000\002\000\000\000\002"+
		"\x14\x15\x16\x17\x18\x19\x1a\x1b"+
		"\x24\x25\x26\x27\x28\x29\x2a\x2b"),
		w.Bytes()) {
		t.Errorf("encoded image differs")
	}
}
