# Podio v0.1 - Easy Audio Processing for Golang

> [!IMPORTANT]  
> We are launching our Audio-as-a-Service in BETA. During this phase, we kindly ask that you report any bugs either under the "Issues" section or directly to me at support@poxate.com. Thank you for your understanding and support.

**Podio** is a simple and powerful audio processing library for Go, designed to make it easy to manipulate and compile audio sequences. It leverages **distributed ffmpeg** processing to scale efficiently, so you don't have to worry about heavy computational tasks on a single server.

## Features

-   **Easy Audio Composition**: Use a fluent API to chain audio operations like looping, padding, volume adjustments, and more.
-   **Distributed Audio Processing**: Offload the heavy lifting of audio transformations to our serverless architecture.
-   **Multiple Format Support**: Work with various audio formats like MP3, WAV, OPUS, and PCM.

> [!CAUTION]
> For the sake of simplicity, we insert API Key as a hard-coded string. In production, you should read it from an Environment Variable.

## Getting Started

### 1. Sign Up at Poxate.com

To start using Podio, you first need to create an account at [Poxate.com](https://poxate.com):

1.  Visit [Poxate.com](https://poxate.com) and sign up.
2.  After logging in, create a new project.
3.  Generate an **API key** for your project, which you will use to authenticate your requests.

### 2. Install the Podio Library

To use Podio in your Go project, you can install the library with the following:

```
go get github.com/poxate/podio
```

### 3. Set Up Your Client

In your Go application, import the Podio package and initialize a new client with your API key.

```go
package main

import (
	"github.com/poxate/podio"
)

func main() {
	// Replace with your own API key from Poxate.com
	client := podio.NewClient("your-api-key-here")

	// Example usage...
}
```

----------

## Audio Composition Examples

Podio uses the builder pattern to easily compose complex audio operations. Below are some common use cases.

### Concatenate Audio Files
To concatenate multiple audio files:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/poxate/podio"
)

func main() {
	tmp, err := os.CreateTemp("", "*.mp3")
	if err != nil {
		log.Fatalln(err)
	}

	// Get an API Key from https://www.poxate.com
	client := podio.NewClient("---api key---")
	voiceSoundtrack := "https://cdn.prayershub.com/scripture/1733531671.mp3"
	musicSoundtrack := "https://cdn.prayershub.com/songs/french/piano/entends-tu-le-chant-joyeux-215.mp3"

	ab := podio.Remote(voiceSoundtrack).WithBackground(
		podio.Remote(musicSoundtrack).FadeIn(2 * time.Second).FadeOut(2 * time.Second),
	)

	err = client.Compile(context.Background(), podio.MP3, ab, tmp)
	if err != nil {
		log.Fatalln("client.Compile:", err)
	}

	fmt.Println("Finished, compiled audio to:", tmp.Name())
}
```

### Concat audios
Audios play one after another, like a playlist.
```go
ab := podio.Concat(
	podio.Remote("https://example.com/first-speech.mp3"),
	podio.Remote("https://example.com/second-speech.mp3"),
	podio.Remote("https://example.com/third-speech.mp3"),
)
```

### Loop Audio

You can loop an audio track a specified number of times. The track will be repeated **n** times (e.g., `count = 3` will play the track a total of 4 times):

```go
ab := podio.Remote("https://example.com/bell.mp3").Loop(3)
```

### Pad Audio (Left/Right)

Add silence to the beginning or end of an audio file:
```go
ab := podio.
	Remote("https://example.com/voice.mp3").
	PadLeft(time.Second * 2).
	PadRight(time.Second * 2)
```

### Adjust Volume

Adjust the volume of an audio file:

 - 0.5 decreases the volume to half
 - 1 has no change
 - 2 doubles the volume

```go
ab := podio.Remote("https://example.com/loud-music.mp3").Volume(0.5)
```

### Add Background Music
You can add background music to an audio file, where the background audio loops until the main audio ends:

```go
ab := podio.Remote("https://example.com/voice.mp3").WithBackground(
	podio.Remote("https://example.com/loud-music.mp3").
	Volume(0.5).
	FadeIn(2 * time.Second).
	FadeOut(2 * time.Second),
)
```

### Save Duration
```go
var voiceDuration time.Duration

ab := podio.
	Remote("https://cdn.prayershub.com/prayer-lists/Lord's-prayer-our-father.mp3").
	SaveDuration(&voiceDuration)
if err := client.Compile(context.Background(), podio.MP3, ab, tmp); err != nil {
	// ... error handling
}
// At this point, voiceDuration is available

fmt.Println("Voice duration:", voiceDuratoin)
// Sample output:"Voice duration: 25.03201814s"
```

#### Complex use case #1: multiple SaveDuration
Let's say you want to combine multiple audio files into one single audio file, but want the duration of each individual audio track.

Traditionally, you'd launch an `ffprobe` process to obtain the duration for each audio track, and then `ffmpeg` to combine each audio.

With Podio, you can create a `[]time.Duration`, and then pass a reference for each into Podio.

It's that simple!
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/poxate/podio"
)

func main() {
	tmpFileWriter, err := os.CreateTemp("", "*.mp3")
	if err != nil {
		log.Fatalln(err)
	}

	// Retrieve API KEY from https://www.poxate.com/auth/register
	client := podio.NewClient("--api-key")
	
	playlist := getAudios()
	durations := make([]time.Duration, len(playlist))
	concatList := []podio.AudioBuilder{}

	for i, audio := range playlist {
		// WARNING: make sure to pass a pointer that avoids copy
		var refDuration *time.Duration = &durations[i]
		concatList = append(concatList, podio.Remote(audio).SaveDuration(refDuration))
	}

	if err = client.Compile(context.TODO(), podio.MP3, podio.Concat(concatList...), tmpFileWriter); err != nil {
		// ... error handling
	}

	for i, audio := range playlist {
		fmt.Printf("Audio (%s) is %s\n", audio, durations[i])
	}

	fmt.Println("Finished, compiled audio to:", tmpFileWriter.Name())
}

func getAudios() []string {
	return []string{
		"https://cdn.prayershub.com/songs/french/piano/tenons-nos-lampes-pretes-120.mp3",
		"https://cdn.prayershub.com/songs/french/piano/semons-des-que-brille-laurore-208.mp3",
		"https://cdn.prayershub.com/songs/french/piano/entends-tu-le-chant-joyeux-215.mp3",
	}
}

```

Sample output:
```
Audio (https://cdn.prayershub.com/songs/french/piano/tenons-nos-lampes-pretes-120.mp3) is 44.747755102s
Audio (https://cdn.prayershub.com/songs/french/piano/semons-des-que-brille-laurore-208.mp3) is 1m12.96s
Audio (https://cdn.prayershub.com/songs/french/piano/entends-tu-le-chant-joyeux-215.mp3) is 29.753469387s
Finished, compiled audio to: /tmp/136159727.mp3
```

### Complex use case #2: get original & final duration of compilation
Not only can you get the duration of remote files, you can get the duration of complex build graphs:
```go
client := podio.NewClient("API_KEY_GOES_HERE")

var originalDuration time.Duration
var finalDuration time.Duration

ab := podio.
	Remote("https://example.com/music.mp3").
	SaveDuration(&originalDuration).
	Loop(3).
	FadeIn(time.Second).
	FadeOut(time.Second).
	SaveDuration(&finalDuration)

if err = client.Compile(context.TODO(), podio.MP3, ab, tmpFile); err != nil {
	// ... error handling
}

fmt.Println("Original duration:", originalDuration)
fmt.Println("Final duration:", finalDuration)
```

## Contributing

If you'd like to contribute to Podio, feel free to fork the repository and submit a pull request. We welcome contributions that improve the functionality, performance, and documentation of Podio.

## License

The Podio SDK is open-source software licensed under the MIT license. See the LICENSE file for more information.

## Conclusion

Podio simplifies audio processing tasks with an intuitive API and distributed backend. It is a perfect solution for developers needing to process audio files efficiently in the cloud or on-demand.

Start by signing up at **Poxate.com**, get your API key, and start building your audio applications today!

----------

Let me know if you need further adjustments or additional details!

In addition, you can contact me at support@poxate.com
