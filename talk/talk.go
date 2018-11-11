package talk

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

// Talk uploads valid SSML to Google Speech, saves audio to /var/tmp/loquacious/audio/
// and streams through default sound device.
func Talk(ssmlPath string) {
	log.Print("Parsing SSML file")
	ssml := readSSML(ssmlPath)

	log.Print("Uploading SSML to Speech client")
	audioContent := uploadSSML(ssml)

	log.Print("Saving speech audio")
	filename := writeAudioToFile(audioContent)

	log.Print("Playing speech audio")
	Play(filename)
}

// Demo synthesizes speech like Talk but with included sample SSML for testing purposes
func Demo() {
	ssml := sampleSSML()
	log.Print("Uploading sample SSML to Speech client")
	audioContent := uploadSSML(ssml)

	log.Print("Saving speech audio")
	filename := writeAudioToFile(audioContent)

	log.Print("Playing speech audio")
	Play(filename)
}

func createSavePath() string {
	savePath := "/var/tmp/loquacious/audio/synthesized"
	if err := os.MkdirAll(savePath, 0747); err != nil {
		log.Fatalf("Something went wrong creating save path in %v\n", savePath)
	}
	return savePath
}

// Writes audio to new output file with format "synthesized_speech_YYYY-MM-DD-HH-MM-SS.wav"
func writeAudioToFile(audioContent []byte) string {
	now := time.Now().Format("2006-01-02-15-04-05")
	savePath := createSavePath()
	filename := savePath + "/speech_" + now + ".wav"
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

	log.Println("Building Speech API request with provided SSML")
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

	log.Println("Requesting speech audio")
	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal("Something went wrong making request to SynthesizeSpeech client")
	}

	return resp.AudioContent
}

func sampleSSML() string {
	return `
<speak>
  Tomorrow, <break time="400ms"/> and tomorrow, and tomorrow,
  <break time="300ms"/>
  Creeps in this petty pace <break time="250ms"/> from day to day,
  To the last syllable of recorded time;
  <break time="350ms"/>
  And all our yesterdays have lighted fools
  <break time="350ms"/>
  The way to dusty death.
  Out, out, brief candle!
  <break time="350ms"/>
  Life's but a walking shadow, <break time="300ms"/> a poor player,
  That struts and frets his hour upon the stage,
  <break time="350ms"/>
  And then is heard no more. It is a tale
  <break time="250ms"/>
  Told by an idiot, <break time="250ms"/> full of sound and fury,
  <break time="250ms"/>
  Signifying <break time="300ms"/> <emphasis level="moderate">nothing.</emphasis>
</speak>
`
}
