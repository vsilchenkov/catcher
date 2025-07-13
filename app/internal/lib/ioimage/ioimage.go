package ioimage

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
	"github.com/cockroachdb/errors"
)

func Compress(input []byte, percent uint) ([]byte, error) {

	const op = "ioimage.Compress"

	if percent == 0 || percent > 100 {
		return nil, errors.Errorf("%s - percent должен быть в диапазоне 1-100, передано %v", op, percent)
	}

	img, format, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, errors.WithMessagef(err, "%s - ошибка декодирования изображения", op)
	}

	origBounds := img.Bounds()
	origWidth := uint(origBounds.Dx())
	origHeight := uint(origBounds.Dy())

	newWidth := origWidth * percent / 100
	newHeight := origHeight * percent / 100

	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 80})
	case "png":
		err = png.Encode(&buf, resizedImg)
	default:
		return nil, errors.Errorf("%s - неподдерживаемыq формат image: %s", op, format)
	}
	if err != nil {
		return nil, errors.WithMessagef(err, "%s - ошибка кодирования изображения", op)
	}

	return buf.Bytes(), nil
}
