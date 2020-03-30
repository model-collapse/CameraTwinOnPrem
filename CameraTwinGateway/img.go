//+build !darwin

package main

import (
	"fmt"
	cv "gocv.io/x/gocv"
	"log"
)

const explicitThres = 0.0002543
const downsample = 2
const varThres = 1e-5
const entriesThres = 100

func isBrokenJPEG(data []byte) bool {
	log.Printf("data size = %d", len(data))
	img, err := cv.IMDecode(data, cv.IMReadColor)
	if err != nil {
		log.Printf("Invalid format of jpeg file, may be the file is broken, %s", err)
		return true
	}

	return isBrokenImage(img)
}

func histogramScore(img cv.Mat, downsample int) (nvariance float64, entries int32, err error) {
	step := img.Step()
	rows := img.Rows()
	cols := img.Cols()
	channels := img.Channels()
	size := rows * cols * channels
	log.Printf("size = %d, rows = %d, cols = %d, step = %d, chans = %d", size, rows, cols, step, channels)
	thres := int(explicitThres * float64(size))

	if img.Type() != cv.MatTypeCV8UC3 {
		err = fmt.Errorf("image was not of type 8UC3, it is %d", img.Type())
		return
	}

	histograms := make([][]int32, channels)
	for i := 0; i < channels; i++ {
		histograms[i] = make([]int32, 256)
	}

	ptr := img.DataPtrUint8()
	off := 0
	for y := 0; y < rows; y += downsample {
		for x := 0; x < cols; x += downsample {
			for c := 0; c < channels; c++ {
				histograms[c][ptr[off+x*channels+c]]++
			}
		}

		off += downsample * step
	}

	mean := 1 / 256.0
	nvariance = 0.0
	entries = 0
	norm := float64(downsample*downsample) / float64(size)
	for i := 0; i < channels; i++ {
		for j := 0; j < 256; j++ {
			d := float64(histograms[i][j])*norm - mean
			if int(histograms[i][j]) > thres {
				entries++
			}
			nvariance += d * d
		}
	}
	nvariance /= float64(256 * channels)

	return
}

func isBrokenImage(img cv.Mat) bool {
	variance, entries, err := histogramScore(img, downsample)
	if err != nil {
		log.Println(err)
		return true
	}

	log.Printf("var = %f, ent = %d", variance, entries)
	return variance > varThres && entries < entriesThres
}
