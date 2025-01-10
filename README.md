# Podio - Audio Compilation Library
Podio is an audio processing library that allows you to build and manipulate audio in various ways, including looping, padding, background music, volume adjustments, fading, and saving durations. This library allows you to compile and process audio using WebSocket communication with a Podio server.

Table of Contents
Getting Started
Installation
Usage
Available Operations
Examples
Error Handling
License
Getting Started
To use Podio, you will need to sign up at poxate.com, create a project, and generate an API key. The API key is used to authenticate your requests to the Podio server.

Steps:
Sign up: Go to poxate.com and create an account.
Create a Project: Once signed in, create a new project from the dashboard.
Generate API Key: After creating your project, generate an API key that will be used to interact with the Podio server.
Installation
To install and use Podio in your Go project, follow these steps:

Install Go: If you don’t have Go installed, you can get it from golang.org.

Clone the repository:

bash
Copy code
git clone https://github.com/poxate/podio.git
Install dependencies: In your Go project directory, run:

bash
Copy code
go mod tidy
Usage
To start using the Podio client, first create a new Client with your API key. Then, you can call methods to build and compile audio.

Example:
go
Copy code
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/poxate/podio"
)

func main() {
	// Replace with your actual API key
	client := podio.NewClient("your-api-key")

	// Create an audio builder with remote audio
	audioBuilder := podio.Remote("https://example.com/audio.mp3")

	// Compile the audio into an MP3 format and write to a buffer
	var outputBuffer bytes.Buffer
	ctx := context.Background()
	if err := client.Compile(ctx, podio.MP3, audioBuilder, &outputBuffer); err != nil {
		log.Fatalf("Error compiling audio: %v", err)
	}

	// Now you can use the compiled audio (stored in outputBuffer)
	fmt.Println("Audio compiled successfully")
}
Parameters:
API Key: A valid API key obtained from the Podio dashboard.
AudioBuilder: The audio you want to build, which can be customized with various operations (e.g., looping, padding, background music, etc.).
Format: The format to compile the audio into (e.g., MP3, WAV).
Destination (io.Writer): The destination where the compiled audio will be written (e.g., *bytes.Buffer).
Available Operations
Podio provides several audio manipulation operations that you can apply using AudioBuilder. These operations can be chained to customize the audio output.

Remote(url string): Load audio from a remote URL.
Loop(count int): Loop the audio count times (the audio will be played count + 1 times in total).
PadLeft(duration time.Duration): Add silence to the beginning of the audio.
PadRight(duration time.Duration): Add silence to the end of the audio.
WithBackground(backgroundAudio AudioBuilder): Add background audio that loops until the main audio finishes.
Volume(volume float64): Adjust the volume of the audio (e.g., 0.5 for half the volume).
FadeIn(duration time.Duration): Apply a fade-in effect over the specified duration.
FadeOut(duration time.Duration): Apply a fade-out effect over the specified duration.
SaveDuration(duration *time.Duration): Save the duration of the audio for later use.
Examples
Here are a few examples to show you how to use the various operations.

Example 1: Looping an audio
go
Copy code
audio := podio.Remote("https://example.com/audio.mp3").Loop(3)
This will create an audio that loops 3 times (so the audio will play 4 times in total).

Example 2: Adding silence to the beginning and end
go
Copy code
audio := podio.Remote("https://example.com/audio.mp3").
    PadLeft(5 * time.Second).
    PadRight(2 * time.Second)
This will add 5 seconds of silence at the beginning and 2 seconds at the end of the audio.

Example 3: Adjusting the volume
go
Copy code
audio := podio.Remote("https://example.com/audio.mp3").Volume(1.5) // Increase volume to 150%
Example 4: Adding background music
go
Copy code
background := podio.Remote("https://example.com/background.mp3")
audio := podio.Remote("https://example.com/main.mp3").WithBackground(background)
This will add the background audio that loops until the main audio finishes.

Example 5: Fading the audio in and out
go
Copy code
audio := podio.Remote("https://example.com/audio.mp3").
    FadeIn(2 * time.Second).
    FadeOut(3 * time.Second)
Error Handling
In case of errors, such as network issues, invalid audio files, or unexpected server behavior, the library will return errors that can be handled appropriately.

Common Errors:
Connection Errors: Issues connecting to the Podio server will result in a DialError or ConnectionError.
Invalid Audio Format: If the requested audio format is not supported, an error will be returned.
Audio Builder Errors: If an operation is chained incorrectly (e.g., missing required data), an error will occur.
License
This project is licensed under the MIT License - see the LICENSE file for details.

This README gives a straightforward guide for getting started with Podio, including setup, basic usage, examples, and error handling. If you have more features or specific instructions to include, feel free to adjust or expand this document!