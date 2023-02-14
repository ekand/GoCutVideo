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
	var sliceOfRMSValues []RMS
	var frame [1024]int16
	for i := 0; true; i++ {
		if err := binary.Read(reader, binary.LittleEndian, &frame); err != nil {
			if err == io.EOF {
				break
			}
			if err.Error() == "unexpected EOF" {
				fmt.Printf("found EOF\n")
				break
			}
			fmt.Printf("hello 597823 error: %v\n", err)
			return
		}

		var sum int64
		for _, sample := range frame {
			sum += int64(sample) * int64(sample)
		}

		rms := float64(sum) / float64(len(frame))
		rmsType := RMS(rms)
		sliceOfRMSValues = append(sliceOfRMSValues, rmsType)
		fmt.Printf("i: %d, RMS: %f\n", i, rms)
	}
	// fmt.Printf("%v", sliceOfRMSValues)
	foo := DetectSilence(sliceOfRMSValues)
	fmt.Printf("%v\n", foo)
}

// RMS is the root mean square of the audio signal
type RMS float64

// DetectSilence detects the sections of the slice that are silent
func DetectSilence(rms []RMS) [][2]int {
	const Threshold = RMS(10000.1)
	const SampleRate = 44100
	const SilenceDuration = 70 //int64(0.001 * float64(SampleRate))

	var silences [][2]int
	// var start int
	var inSilenceStart int
	// var end int
	var possibleSilenceEnd int
	inSilence := false
	for i, r := range rms {
		if r < Threshold && !inSilence {
			inSilence = true
			inSilenceStart = i
		} else if r < Threshold && inSilence {
			possibleSilenceEnd = i
		}
		if r > Threshold && inSilence {
			inSilence = false // exit silence
			if possibleSilenceEnd-inSilenceStart > SilenceDuration {
				//  and record a start/stop of silence
				silences = append(silences, [2]int{inSilenceStart, possibleSilenceEnd})
			}
		} else if r > Threshold && !inSilence {
			// still exit silence
			inSilence = false
		}
	}
	if inSilence && possibleSilenceEnd-inSilenceStart > SilenceDuration {
		silences = append(silences, [2]int{inSilenceStart, possibleSilenceEnd})
	}
	return silences
}
