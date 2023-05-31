module github.com/kbinani/screenshot

go 1.19

require (
	github.com/gen2brain/shm v0.0.0-20221026125803-c33c9e32b1c8
	github.com/jezek/xgb v1.1.0
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e
)

replace github.com/kbinani/screenshot => ./

require golang.org/x/sys v0.8.0 // indirect
