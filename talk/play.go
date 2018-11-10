package talk

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// Play accepts an audio file in  WAV or AIFF format, and plays it through default audio device.
func Play(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Something went wrong opening the provided file")
	}
	defer file.Close()

	// decode a WAV or AIFF and open a
	streamer, format, err := wav.Decode(file)
	if err != nil {
		log.Fatal("Something went wrong decoding the provided audio file. Make sure it's either a WAV or AIFF")
	}

	// create a channel  which signals end of playback
	done := make(chan struct{})

	// initialize speaker with file's sample rate and buffer size of 0.1 seconds
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	<-done
}
