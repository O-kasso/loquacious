# **Loquacious**

> lo·qua·cious /lōˈkwāSHəs/ 
> talking or tending to talk much or freely; talkative; chattering; babbling

Loquacious is a small cli program (written in Go) that enables dead-simple **text-to-speech** and **speech-to-text** directly from the command line. Uses Google's [Cloud Text-to-Speech](https://cloud.google.com/text-to-speech/) and [Cloud Speech-to-Text](https://cloud.google.com/speech-to-text/) APIs.



### Setup Credentials

- Setup a new [GCP](https://console.cloud.google.com/projectcreate) project (you may need to enable billing)
- Add a Service account for the Speech API to your project
- Download your Service account key as JSON
- Add the path to your account key as an environment variable named  "GOOGLE_APPLICATION_CREDENTIALS"
  - `export GOOGLE_APPLICATION_CREDENTIALS=$HOME/gcp_service_key.json`
  - (optionally, you can preface every use of `loquacious` with a one-time variable assignment)



## Text-to-Speech

To generate speech from a [SSML](https://en.wikipedia.org/wiki/Speech_Synthesis_Markup_Language) file, provide it as a command line argument:

```
$ ./loquacious ./sample/shakespeare.ssml
```



If you didn't export `GOOGLE_APPLICATION_CREDENTIALS`, run:

```
$ GOOGLE_APPLICATION_CREDENTIALS=$HOME/gcp_service_key.json ./loquacious ./sample/shakespeare.ssml
```



## Speech-to-Text

Coming soon.