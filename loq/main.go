package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/o-kasso/loquacious/cmd"
	"golang.org/x/net/context"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func main() {
	cmd.Execute()
}

func createSavePath() string {
	savePath := "/var/tmp/synthesized_speech"
	if err := os.MkdirAll(savePath, 0747); err != nil {
		log.Fatalf("Something went wrong creating save path in %v\n", savePath)
	}
	return savePath
}

// Writes audio to new output file with format "synthesized_speech_YYYY-MM-DD-HH-MM-SS.wav"
func writeAudioToFile(audioContent []byte) string {
	now := time.Now().Format("2006-01-02-15-04-05")
	savePath := createSavePath()
	filename := savePath + "/synthesized_speech_" + now + ".wav"
	err := ioutil.WriteFile(filename, audioContent, 0644)
	if err != nil {
		log.Fatalf("Something went wrong writing audio to %v\n", filename)
	}

	log.Printf("Audio content written to file: %v\n", filename)
	return filename
}

func readSSML(file string) string {
	bytes, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal("Something went wrong attempting to load SSML file.")
	}

	return string(bytes)
}

// UploadSSML uploads provided ssml to Google Speech API and returns audio stream.
// ssml format: https://www.w3.org/TR/speech-synthesis/
func uploadSSML(ssml string) []byte {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal("Something went wrong creating new TextToSpeech client. Make sure that you have a valid Google Cloud Platform service key configured and added to your environment as GOOGLE_APPLICATION_CREDENTIALS")
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: ssml},
		},

		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Wavenet-D",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal("Something went wrong making request to SynthesizeSpeech client")
	}

	return resp.AudioContent
}

func verifyGoogleCredentials() {
	configFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	config, err := os.Stat(configFile)
	if err != nil || config.Size() <= 0 {
		log.Fatal("Could not find a valid Google Cloud Platform service key set to GOOGLE_APPLICATION_CREDENTIALS environment variable.")
	}
}

func playSynthesizedSpeech(filename string) {
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

// expects path to valid SSML file as argument.
func oldMain() {
	if len(os.Args) <= 1 {
		log.Fatal("Please provide path to a valid SSML file as argument")
	}

	log.Print("Verifying presence of Google Cloud Platform credentials")
	verifyGoogleCredentials()

	log.Print("Uploading SSML to Speech client")
	ssml := readSSML(os.Args[1])
	audioContent := uploadSSML(ssml)

	log.Print("Saving speech audio")
	filename := writeAudioToFile(audioContent)

	log.Print("Playing speech audio")
	playSynthesizedSpeech(filename)
}
