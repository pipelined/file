// Package fileformat provides functionality to process files with pipelined
// framework.
package fileformat

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"pipelined.dev/audio/flac"
	"pipelined.dev/audio/mp3"
	"pipelined.dev/audio/wav"
	"pipelined.dev/pipe"
)

// Format of the file that contains audio signal.
type Format struct {
	defaultExtension string
	extensions       []string
	source           func(io.ReadSeeker) pipe.SourceAllocatorFunc
}

var (
	wavFormat = Format{
		defaultExtension: ".wav",
		extensions: []string{
			".wav",
			".wave",
		},
		source: func(rs io.ReadSeeker) pipe.SourceAllocatorFunc {
			return wav.Source(rs)
		},
	}

	mp3Format = Format{
		defaultExtension: ".mp3",
		extensions: []string{
			".mp3",
		},
		source: func(rs io.ReadSeeker) pipe.SourceAllocatorFunc {
			return mp3.Source(rs)
		},
	}

	flacFormat = Format{
		defaultExtension: ".flac",
		extensions: []string{
			".flac",
		},
		source: func(rs io.ReadSeeker) pipe.SourceAllocatorFunc {
			return flac.Source(rs)
		},
	}

	// formatByExtension = mapFormatByExtension(WAV, MP3, FLAC)
	formatByExtension = func(formats ...*Format) map[string]*Format {
		m := make(map[string]*Format)
		for _, format := range formats {
			for _, ext := range format.Extensions() {
				if _, ok := m[ext]; ok {
					panic(fmt.Sprintf("multiple formats have same extension: %s", ext))
				}
				m[ext] = format
			}
		}
		return m
	}(&wavFormat, &mp3Format, &flacFormat)
)

// WAV returns Waveform Audio file format.
func WAV() *Format {
	return &wavFormat
}

// MP3 returns MPEG-1 or MPEG-2 Audio Layer III file format.
func MP3() *Format {
	return &mp3Format
}

// FLAC returns Free Lossless Audio Codec file format.
func FLAC() *Format {
	return &flacFormat
}

// FormatByPath determines file format by file extension
// extracted from path. If extension belongs to unsupported
// format, nil is returned.
func FormatByPath(path string) *Format {
	return formatByExtension[strings.ToLower(filepath.Ext(path))]
}

// MatchExtension checks if ext matches to one of the format's
// extensions. Case is ignored.
func (f *Format) MatchExtension(ext string) bool {
	format, ok := formatByExtension[strings.ToLower(ext)]
	if !ok {
		return false
	}
	return f == format
}

// Source returns pipe.Source for corresponding format
// with injected ReadSeeker.
func (f *Format) Source(rs io.ReadSeeker) pipe.SourceAllocatorFunc {
	return f.source(rs)
}

// DefaultExtension of the format.
func (f *Format) DefaultExtension() string {
	return f.defaultExtension
}

// Extensions returns a slice of format's extensions.
func (f *Format) Extensions() []string {
	return append(f.extensions[:0:0], f.extensions...)
}

// PipeFunc is user-defined function that allows to process files during
// filewalk.
type PipeFunc func(*Format, string, os.FileInfo) error

// Walk takes user-defined pipe function and return filepath.WalkFunc.
// It allows to use it with filepath.Walk function and execute pipe func
// with every file in a path. This function will try to parse file format
// from it's extension.
func Walk(fn PipeFunc, recursive bool) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error during walk: %w", err)
		}
		if fi.IsDir() {
			if recursive {
				return nil
			}
			// skip processing subdirs
			return filepath.SkipDir
		}

		format := FormatByPath(path)
		if format == nil {
			return nil
		}

		if err = fn(format, path, fi); err != nil {
			return fmt.Errorf("error execution pipe func: %w", err)
		}
		return nil
	}
}
