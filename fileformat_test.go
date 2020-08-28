package fileformat_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

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
		format := fileformat.FormatByPath(test.fileName)
		if test.negative {
			var nilFormat *fileformat.Format
			assertEqual(t, "format", format, nilFormat)
		} else {
			assertNotNil(t, "format", format)
			source := format.Source(nil)
			assertNotNil(t, "source", source)
		}
	}
}

func TestExtensions(t *testing.T) {
	var tests = []struct {
		format   *fileformat.Format
		expected int
	}{
		{
			fileformat.WAV(),
			2,
		},
		{
			fileformat.MP3(),
			1,
		},
		{
			fileformat.FLAC(),
			1,
		},
	}

	for _, test := range tests {
		exts := test.format.Extensions()
		assertEqual(t, "extensions", len(exts), test.expected)
	}
}

func TestWalk(t *testing.T) {
	testPositive := func(path string, recursive bool, expected int) func(*testing.T) {
		return func(t *testing.T) {
			processed := 0
			fn := func(f *fileformat.Format, path string, fi os.FileInfo) error {
				processed++
				return nil
			}
			walkFn := fileformat.Walk(fn, recursive)
			err := filepath.Walk(path, walkFn)
			assertNil(t, "error", err)
			assertEqual(t, "processed", processed, expected)
		}
	}
	testFailedWalk := func() func(*testing.T) {
		return func(t *testing.T) {
			err := filepath.Walk("nonexistentfileformat.wav", fileformat.Walk(nil, false))
			assertNotNil(t, "error", err)
		}
	}
	testFailedPipe := func(path string) func(*testing.T) {
		return func(t *testing.T) {
			err := filepath.Walk(path,
				fileformat.Walk(func(*fileformat.Format, string, os.FileInfo) error {
					return fmt.Errorf("pipe error")
				}, false))
			assertNotNil(t, "error", err)
		}
	}
	t.Run("recursive", testPositive("_testdata", true, 2))
	t.Run("nonrecursive", testPositive("_testdata", false, 0))
	t.Run("nonexistent ext", testPositive("_testdata/test.nonexistentext", false, 0))
	t.Run("nonexistent file", testFailedWalk())
	t.Run("failed pipe", testFailedPipe("_testdata/test.wav"))
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, expected, expected)
	}
}

func assertNil(t *testing.T, name string, result interface{}) {
	t.Helper()
	assertEqual(t, name, result, nil)
}

func assertNotNil(t *testing.T, name string, result interface{}) {
	t.Helper()
	if reflect.DeepEqual(nil, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, nil, nil)
	}
}
