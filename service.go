package images

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"github.com/pedrobarbosak/go-errors"
)

type processor struct {
	config Config
}

func (p processor) Config() Config {
	return p.config
}

func (p processor) Optimize(_ context.Context, r io.Reader) ([]byte, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, errors.New("failed to decode image:", err)
	}

	if p.config.Resize.Enabled {
		img = p.resize(img)
	}

	if p.config.Conversion.Enabled {
		format = p.config.Conversion.Format
	}

	return p.convert(img, format)
}

func (p processor) Resize(_ context.Context, r io.Reader) ([]byte, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	img = p.resize(img)

	return p.convert(img, format)
}

func (p processor) resize(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	if width <= p.config.Resize.MaxWidth && height <= p.config.Resize.MaxHeight {
		return img
	}

	// Calculate new dimensions while preserving aspect ratio
	var newWidth, newHeight uint
	if width*p.config.Resize.MaxHeight > height*p.config.Resize.MaxWidth {
		newWidth = uint(p.config.Resize.MaxWidth)
		newHeight = uint(height * p.config.Resize.MaxWidth / width)
	} else {
		newHeight = uint(p.config.Resize.MaxHeight)
		newWidth = uint(width * p.config.Resize.MaxHeight / height)
	}

	// Resize with Lanczos resampling (good balance between quality and speed)
	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}

func (p processor) Convert(_ context.Context, r io.Reader, toFormat string) ([]byte, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, errors.New("failed to decode image:", err)
	}

	return p.convert(img, toFormat)
}

func (p processor) convert(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: p.config.JPEG.Quality})
		if err != nil {
			return nil, errors.New("failed to encode jpeg:", err)
		}
	case "png":
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, errors.New("failed to encode png:", err)
		}
	case "webp":
		err := webp.Encode(&buf, img, &webp.Options{Lossless: p.config.WebP.Lossless, Quality: p.config.WebP.Quality})
		if err != nil {
			return nil, errors.New("failed to encode webp:", err)
		}
	case "gif":
		err := gif.Encode(&buf, img, &gif.Options{})
		if err != nil {
			return nil, errors.New("failed to encode gif:", err)
		}

	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return buf.Bytes(), nil
}

func New() ImageProcessor {
	return &processor{
		config: DefaultConfig(),
	}
}

func NewWithConfig(cfg Config) ImageProcessor {
	return &processor{
		config: cfg,
	}
}
