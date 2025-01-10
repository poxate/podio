package podio

import (
	"time"

	"github.com/google/uuid"
	"github.com/sanity-io/litter"
)

type remoteBody string
type loopBody struct {
	body  AudioBuilder
	count int
}
type concatBody []AudioBuilder
type padLeftBody struct {
	body     AudioBuilder
	duration time.Duration
}
type padRightBody struct {
	body     AudioBuilder
	duration time.Duration
}
type backgroundBody struct {
	body       AudioBuilder
	background AudioBuilder
}
type volumeBody struct {
	body   AudioBuilder
	volume float64
}
type fadeInBody struct {
	body     AudioBuilder
	duration time.Duration
}
type fadeOutBody struct {
	body     AudioBuilder
	duration time.Duration
}
type saveDurationBody struct {
	body     AudioBuilder
	duration *time.Duration
}

type audioConstruct interface {
	construct()
}

func (remoteBody) construct()       {}
func (loopBody) construct()         {}
func (concatBody) construct()       {}
func (padLeftBody) construct()      {}
func (padRightBody) construct()     {}
func (backgroundBody) construct()   {}
func (volumeBody) construct()       {}
func (fadeInBody) construct()       {}
func (fadeOutBody) construct()      {}
func (saveDurationBody) construct() {}

type AudioBuilder struct {
	audioConstruct
}

// Generators
func Remote(url string) AudioBuilder {
	return AudioBuilder{remoteBody(url)}
}

func Concat(audios ...AudioBuilder) AudioBuilder {
	return AudioBuilder{concatBody(audios)}
}

// Utils
// How many times will the audio be replayed (not played), so if provided count is 3, the audio will be played a total of 4 times
func (ab AudioBuilder) Loop(count int) AudioBuilder {
	return AudioBuilder{loopBody{body: ab, count: count}}
}

// Add silence to the beginning of the audio
func (ab AudioBuilder) PadLeft(duration time.Duration) AudioBuilder {
	return AudioBuilder{padLeftBody{ab, duration}}
}

// Add silence to the end of the audio
func (ab AudioBuilder) PadRight(duration time.Duration) AudioBuilder {
	return AudioBuilder{padRightBody{ab, duration}}
}

// Background music will be looped until the audio ends
func (ab AudioBuilder) WithBackground(backgroundAudio AudioBuilder) AudioBuilder {
	return AudioBuilder{backgroundBody{
		body:       ab,
		background: backgroundAudio,
	}}
}

// If we want our volume to be half of the input volume:
//
// podio.Remote(url).Volume(0.5)
//
// 150% of current volume:
//
// podio.Remote(url).Volume(1.5)
func (ab AudioBuilder) Volume(volume float64) AudioBuilder {
	return AudioBuilder{volumeBody{
		body:   ab,
		volume: volume,
	}}
}

func (ab AudioBuilder) FadeIn(duration time.Duration) AudioBuilder {
	return AudioBuilder{fadeInBody{
		body:     ab,
		duration: duration,
	}}
}

func (ab AudioBuilder) FadeOut(duration time.Duration) AudioBuilder {
	return AudioBuilder{fadeOutBody{
		body:     ab,
		duration: duration,
	}}
}

// Saves to a reference of time.Duration
// Useful if you want to save the duration of the audio
//
// Example:
//
//	var duration time.Duration
//	podio.Remote(url).SaveDuration(&duration).Fetch(podio.MP3, file)
func (ab AudioBuilder) SaveDuration(duration *time.Duration) AudioBuilder {
	return AudioBuilder{saveDurationBody{
		body:     ab,
		duration: duration,
	}}
}

func (ab AudioBuilder) toJson(ctx *fetchContext) map[string]any {
	type obj = map[string]any

	switch op := ab.audioConstruct.(type) {
	case remoteBody:
		return obj{"type": "remote", "url": string(op)}
	case loopBody:
		return obj{"type": "loop", "count": op.count, "body": op.body.toJson(ctx)}
	case concatBody:
		arr := []obj{}
		for _, v := range op {
			arr = append(arr, v.toJson(ctx))
		}
		return obj{"type": "concat", "list": arr}
	case padLeftBody:
		return obj{"type": "padLeft", "duration": op.duration.Nanoseconds(), "body": op.body.toJson(ctx)}
	case padRightBody:
		return obj{"type": "padRight", "duration": op.duration.Nanoseconds(), "body": op.body.toJson(ctx)}
	case backgroundBody:
		return obj{"type": "background", "body": op.body.toJson(ctx), "background": op.background.toJson(ctx)}
	case volumeBody:
		return obj{"type": "volume", "body": op.body.toJson(ctx), "volume": op.volume}
	case fadeInBody:
		return obj{"type": "fadeIn", "body": op.body.toJson(ctx), "duration": op.duration.Nanoseconds()}
	case fadeOutBody:
		return obj{"type": "fadeOut", "body": op.body.toJson(ctx), "duration": op.duration.Nanoseconds()}
	case saveDurationBody:
		ref := uuid.NewString()
		ctx.durations[ref] = op.duration
		return obj{"type": "saveDuration", "body": op.body.toJson(ctx), "duration_ref": ref}
	default:
		return obj{"type": "unknown", "body": litter.Sdump(ab)}
	}
}
