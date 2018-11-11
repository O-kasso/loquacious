package listen

import (
	"encoding/binary"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gordonklaus/portaudio"
)

// Record activates default audio input device (e.g. microphone) and writes audio stream to AIFF file in /var/tmp/
// Returns path to output file
func Record(timeLimit int) string {
	log.Println("Recording.  Press Ctrl-C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	outputFile := getOutputFilePath()
	f, err := os.Create(outputFile)
	//chk(err)

	// form chunk
	_, err = f.WriteString("FORM")
	chk(err)
	chk(binary.Write(f, binary.BigEndian, int32(0))) //total bytes
	_, err = f.WriteString("AIFF")
	chk(err)

	// common chunk
	_, err = f.WriteString("COMM")
	chk(err)
	chk(binary.Write(f, binary.BigEndian, int32(18)))                  //size
	chk(binary.Write(f, binary.BigEndian, int16(1)))                   //channels
	chk(binary.Write(f, binary.BigEndian, int32(0)))                   //number of samples
	chk(binary.Write(f, binary.BigEndian, int16(32)))                  //bits per sample
	_, err = f.Write([]byte{0x40, 0x0e, 0xac, 0x44, 0, 0, 0, 0, 0, 0}) //80-bit sample rate 44100
	chk(err)

	// sound chunk
	_, err = f.WriteString("SSND")
	chk(err)
	chk(binary.Write(f, binary.BigEndian, int32(0))) //size
	chk(binary.Write(f, binary.BigEndian, int32(0))) //offset
	chk(binary.Write(f, binary.BigEndian, int32(0))) //block
	nSamples := 0
	defer func() {
		// fill in missing sizes
		totalBytes := 4 + 8 + 18 + 8 + 8 + 4*nSamples
		_, err = f.Seek(4, 0)
		chk(err)
		chk(binary.Write(f, binary.BigEndian, int32(totalBytes)))
		_, err = f.Seek(22, 0)
		chk(err)
		chk(binary.Write(f, binary.BigEndian, int32(nSamples)))
		_, err = f.Seek(42, 0)
		chk(err)
		chk(binary.Write(f, binary.BigEndian, int32(4*nSamples+8)))
		chk(f.Close())
	}()

	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	chk(err)
	defer stream.Close()

	chk(stream.Start())

	for stay, timeout := true, time.After(time.Second*time.Duration(timeLimit)); stay; {
		chk(stream.Read())
		chk(binary.Write(f, binary.BigEndian, in))
		nSamples += len(in)
		select {
		case <-sig:
			return outputFile
		default:
		}
		select {
		case <-timeout:
			stay = false
		default:
		}
	}

	chk(stream.Stop())
	log.Println("Recording can be found at: " + outputFile)
	return outputFile
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func getOutputFilePath() string {
	saveDir := "/var/tmp/loquacious/audio/recorded/"
	if err := os.MkdirAll(saveDir, 0747); err != nil {
		log.Fatalf("Something went wrong creating save path in %v\n", saveDir)
	}
	now := time.Now().Format("2006-01-02-15-04-05")
	return saveDir + "loq-audio-recording-" + now + ".aiff"
}
