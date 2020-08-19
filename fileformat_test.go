package fileformat_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"pipelined.dev/audio/fileformat"
)

func TestFilePump(t *testing.T) {
	var tests = []struct {
		fileName string
		negative bool
	}{
		{
			fileName: "test.wav",
		},
		{
			fileName: "test.mp3",
		},
		{
			fileName: "test.flac",
		},
		{
			fileName: "",
			negative: true,
		},
	}

	for _, test := range tests {
		format, err := fileformat.FormatByPath(test.fileName)
		if test.negative {
			assert.NotNil(t, err)
		} else {
			assert.NotNil(t, format)
			source := format.Source(nil)
			assert.NotNil(t, source)
		}
	}
}

func TestExtensions(t *testing.T) {
	var tests = []struct {
		format   fileformat.Format
		expected int
	}{
		{
			fileformat.WAV,
			2,
		},
		{
			fileformat.MP3,
			1,
		},
		{
			fileformat.FLAC,
			1,
		},
	}

	for _, test := range tests {
		exts := test.format.Extensions()
		assert.Equal(t, test.expected, len(exts))
	}
}

func TestWalk(t *testing.T) {
	testPositive := func(path string, recursive bool, expected int) func(*testing.T) {
		return func(t *testing.T) {
			processed := 0
			fn := func(f fileformat.Format, path string, fi os.FileInfo) error {
				processed++
				return nil
			}
			walkFn := fileformat.Walk(fn, recursive)
			err := filepath.Walk(path, walkFn)
			assert.Nil(t, err)
			assert.Equal(t, expected, processed)
		}
	}
	testFailedWalk := func() func(*testing.T) {
		return func(t *testing.T) {
			err := filepath.Walk("nonexistentfileformat.wav", fileformat.Walk(nil, false))
			assert.Error(t, err)
		}
	}
	testFailedPipe := func(path string) func(*testing.T) {
		return func(t *testing.T) {
			err := filepath.Walk(path,
				fileformat.Walk(func(fileformat.Format, string, os.FileInfo) error {
					return fmt.Errorf("pipe error")
				}, false))
			assert.Error(t, err)
		}
	}
	t.Run("recursive", testPositive("_testdata", true, 2))
	t.Run("nonrecursive", testPositive("_testdata", false, 0))
	t.Run("nonexistent ext", testPositive("_testdata/test.nonexistentext", false, 0))
	t.Run("nonexistent file", testFailedWalk())
	t.Run("failed pipe", testFailedPipe("_testdata/test.wav"))
}
