package banner

import (
	"embed"
	"io"
)

var Banner = mustLoadBanner()

//go:embed *
var assets embed.FS

func mustLoadBanner() string{
	f, err := assets.Open("banner.txt")
	if err != nil {
		panic(err)
	}
	
	defer func(){
		derr := f.Close()
		if derr != nil{
			panic(derr)
		}
	}()
	
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	
	return string(data)
}