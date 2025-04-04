package main

import (
	"context"
	"image/png"
	"io"
)

type ImageProcessor interface {
	Config() Config
	Optimize(ctx context.Context, r io.Reader) ([]byte, error)
	Resize(ctx context.Context, r io.Reader) ([]byte, error)
	Convert(ctx context.Context, r io.Reader, toFormat string) ([]byte, error)
}

type Format = string

const (
	JPEG Format = "jpeg"
	JPG  Format = "jpg"
	GIF  Format = "gif"
	PNG  Format = "png"
	WebP Format = "webp"
)

type Config struct {
	Conversion Conversion
	JPEG       JPEGOptions
	PNG        PNGOptions
	WebP       WebPOptions
	Resize     Resize
}

type Conversion struct {
	Enabled bool
	Format  Format
}

type JPEGOptions struct {
	Quality int
}

type PNGOptions struct {
	Compression png.CompressionLevel
	Lossless    bool
	Optimize    bool
}

type WebPOptions struct {
	Quality  float32 // 0-100 for lossy, 0-9 for lossless
	Lossless bool
}

type Resize struct {
	Enabled   bool
	MaxWidth  int
	MaxHeight int
}

func DefaultConfig() Config {
	return Config{
		Conversion: Conversion{
			Enabled: true,
			Format:  WebP,
		},
		JPEG: JPEGOptions{
			Quality: 100,
		},
		PNG: PNGOptions{
			Compression: png.BestCompression,
			Lossless:    true,
			Optimize:    true,
		},
		WebP: WebPOptions{
			Lossless: true,
			Quality:  100,
		},
		Resize: Resize{
			Enabled: false,
		},
	}
}

func ValidExtension(ext string) bool {
	return ext == JPEG || ext == JPG || ext == GIF || ext == PNG || ext == WebP
}
