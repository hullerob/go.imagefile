// See LICENSE file for copyright and license details.

package imagefile

import (
	"encoding/binary"
	"image"
	"image/color"
	"io"
)

const (
	imagefileHeader string = "imagefile????????"
)

func init() {
	image.RegisterFormat("imagefile", imagefileHeader, Decode, DecodeConfig)
}

func Decode(r io.Reader) (image.Image, error) {
	cfg, err := DecodeConfig(r)
	if err != nil {
		return nil, err
	}
	size := 4 * cfg.Width * cfg.Height
	buff := make([]uint8, size)
	_, err = io.ReadFull(r, buff)
	img := &image.NRGBA{
		Pix:    buff,
		Stride: cfg.Width * 4,
		Rect: image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{cfg.Width, cfg.Height},
		},
	}
	return img, err
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	var cfg image.Config
	buff := make([]uint8, 17)
	_, err := io.ReadFull(r, buff)
	cfg.ColorModel = color.RGBAModel
	be := binary.BigEndian
	cfg.Width = int(be.Uint32(buff[9:13]))
	cfg.Height = int(be.Uint32(buff[13:17]))
	return cfg, err
}

func Encode(w io.Writer, m image.Image) error {
	header := []uint8(imagefileHeader)
	be := binary.BigEndian
	width := m.Bounds().Max.X - m.Bounds().Min.X
	height := m.Bounds().Max.Y - m.Bounds().Min.Y
	be.PutUint32(header[9:13], uint32(width))
	be.PutUint32(header[13:17], uint32(height))
	_, err := w.Write(header)
	if err != nil {
		return err
	}
	switch img := m.(type) {
	case *image.NRGBA:
		_, err = w.Write(img.Pix)
	default:
		pix := toRGBA(img)
		_, err = w.Write(pix)
	}
	return err
}

func toRGBA(m image.Image) []uint8 {
	width := m.Bounds().Max.X - m.Bounds().Min.X
	height := m.Bounds().Max.Y - m.Bounds().Min.Y
	size := 4 * width * height
	pix := make([]uint8, size)
	pos := 0
	for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
		for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			pix[pos] = uint8(r >> 8)
			pix[pos+1] = uint8(g >> 8)
			pix[pos+2] = uint8(b >> 8)
			pix[pos+3] = uint8(a >> 8)
			pos += 4
		}
	}
	return pix
}
