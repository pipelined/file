module pipelined.dev/file

go 1.13

require (
	github.com/mewkiz/pkg v0.0.0-20200702171441-dd47075182ea // indirect
	github.com/stretchr/testify v1.4.0
	pipelined.dev/audio/flac v0.0.0-00010101000000-000000000000
	pipelined.dev/audio/mp3 v0.0.0-00010101000000-000000000000
	pipelined.dev/audio/wav v0.0.0-00010101000000-000000000000
	pipelined.dev/pipe v0.8.2
	pipelined.dev/signal v0.7.3 // indirect
)

replace (
	pipelined.dev/audio/flac => ../flac
	pipelined.dev/audio/mp3 => ../mp3
	pipelined.dev/audio/wav => ../wav
)
