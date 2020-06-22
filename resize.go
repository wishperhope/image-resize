package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/discordapp/lilliput"
)

var EncodeOptions = map[string]map[int]int{
	".jpeg": {lilliput.JpegQuality: 85},
	".png":  {lilliput.PngCompression: 85},
	".webp": {lilliput.WebpQuality: 85},
}

func resize(inputBuf []byte, outputHeight int, outputWidth int) ([]byte, error) {

	decoder, err := lilliput.NewDecoder(inputBuf)
	// this error reflects very basic checks,
	// mostly just for the magic bytes of the file to match known image formats
	if err != nil {
		log.Printf("error decoding image, %s\n", err)
		return nil, err
	}
	defer decoder.Close()

	header, err := decoder.Header()
	// this error is much more comprehensive and reflects
	// format errors
	if err != nil {
		log.Printf("error reading image header, %s\n", err)
		return nil, err
	}

	// print some basic info about the image
	log.Printf("file type: %s\n", decoder.Description())
	log.Printf("%dpx x %dpx\n", header.Width(), header.Height())

	if decoder.Duration() != 0 {
		fmt.Printf("duration: %.2f s\n", float64(decoder.Duration())/float64(time.Second))
	}

	// get ready to resize image,
	// using 8192x8192 maximum resize buffer size
	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	// create a buffer to store the output image, 50MB in this case
	outputImg := make([]byte, 50*1024*1024)

	outputType := "." + strings.ToLower(decoder.Description())

	if outputWidth == 0 {
		outputWidth = header.Width()
	}

	if outputHeight == 0 {
		outputHeight = header.Height()
	}

	resizeMethod := lilliput.ImageOpsResize

	if outputWidth == header.Width() && outputHeight == header.Height() {
		resizeMethod = lilliput.ImageOpsNoResize
	}

	opts := &lilliput.ImageOptions{
		FileType:             outputType,
		Width:                outputWidth,
		Height:               outputHeight,
		ResizeMethod:         resizeMethod,
		NormalizeOrientation: true,
		EncodeOptions:        EncodeOptions[outputType],
	}

	// resize and transcode image
	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		log.Printf("error transforming image, %s\n", err)
		return nil, err
	}

	return outputImg, nil
}
