// See LICENSE file for copyright and license details.

// Package imagefile implements a imagefile decoder and encoder.
//
// The imagefile specification can be found at http://git.2f30.org/imagefile/
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

// Decode reads an imagefile from r and returns it as image.NRGBA.
func Decode(r io.Reader) (image.Image, error) {
	cfg, err := DecodeConfig(r)
	if err != nil {
		return nil, err
	}
	img := image.NewNRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))
	_, err = io.ReadFull(r, img.Pix)
	return img, err
}

// DecodeConfig returns dimensions of imagefile image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	var cfg image.Config
	buff := make([]uint8, 17)
	_, err := io.ReadFull(r, buff)
	cfg.ColorModel = color.NRGBAModel
	be := binary.BigEndian
	cfg.Width = int(be.Uint32(buff[9:13]))
	cfg.Height = int(be.Uint32(buff[13:17]))
	return cfg, err
}

// Encode writes m to w in imagefile format. If m is not image.NRGBA,
// it will be converted lossily.
func Encode(w io.Writer, m image.Image) error {
	header := []uint8(imagefileHeader)
	be := binary.BigEndian
	width := m.Bounds().Dx()
	height := m.Bounds().Dy()
	be.PutUint32(header[9:13], uint32(width))
	be.PutUint32(header[13:17], uint32(height))
	_, err := w.Write(header)
	if err != nil {
		return err
	}
	switch img := m.(type) {
	case *image.NRGBA:
		pix := img.Pix
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			_, err = w.Write(pix[:width*4])
			if err != nil {
				return err
			}
			pix = pix[img.Stride:]
		}
	default:
		pix := make([]uint8, width*4)
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			encodeLine(pix, img, y)
			_, err = w.Write(pix)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func encodeLine(pix []uint8, m image.Image, y int) {
	pos := 0
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		c := color.NRGBAModel.Convert(m.At(x, y)).(color.NRGBA)
		pix[pos] = c.R
		pix[pos+1] = c.G
		pix[pos+2] = c.B
		pix[pos+3] = c.A
		pos += 4
	}
}
