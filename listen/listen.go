package listen

import (
	"io/ioutil"
	"log"

	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

// Listen creates a sound file in /var/tmp/loquacious/audio/recorded/ and writes input from the microphone to it
// then uploads it to Google Speech and logs the transcription
func Listen(timeLimit int) {
	log.Print("Recording audio from microphone")

	// TODO: this currently returns filepath of audio -- is it possible to return bytes instead?
	outputSoundFile := Record(timeLimit)

	log.Print("Loading sound clip")
	soundFile := readSoundFile(outputSoundFile)

	log.Print("Uploading sound clip to Speech client and fetching transcription")
	uploadSoundClip(soundFile)
}

//func createSpeechClient() {
//}

func readSoundFile(filename string) []byte {
	soundFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Something went wrong attempting to read audiofile %v", filename)
	}

	return soundFile
}

func uploadSoundClip(soundFile []byte) {
	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal("Something went wrong creating new Speech client. Make sure that you have a valid Google Cloud Platform service key configured and added to your environment as GOOGLE_APPLICATION_CREDENTIALS")
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: soundFile},
		},
	})

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			log.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
}

// Send the contents of the audio file with the encoding and
// and sample rate information to be transcripted.
