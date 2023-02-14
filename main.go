package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"
)

func main() {
	// Define the FFmpeg command to extract the audio stream
	cmd := exec.Command("ffmpeg", "-i", "input.mkv", "-vn", "-acodec", "pcm_s16le", "-f", "wav", "-")

	// Run the FFmpeg command and capture its output
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// Create a reader to read the binary data produced by FFmpeg
	reader := bytes.NewReader(stdout.Bytes())

	// Parse the binary data and compute the RMS of the audio for each frame
	var frame [1024]int16
	for {
		if err := binary.Read(reader, binary.LittleEndian, &frame); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("error: %v\n", err)
			return
		}

		var sum int64
		for _, sample := range frame {
			sum += int64(sample) * int64(sample)
		}
		rms := float64(sum) / float64(len(frame))
		fmt.Printf("RMS: %f\n", rms)
	}
}
