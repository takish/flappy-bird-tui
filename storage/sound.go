package storage

import (
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
)

const (
	sampleRate = beep.SampleRate(48000)
)

var (
	speakerInitialized bool
	speakerMutex       sync.Mutex
)

// InitSound initializes the audio speaker
// Should be called once at application startup
func InitSound() error {
	speakerMutex.Lock()
	defer speakerMutex.Unlock()

	if speakerInitialized {
		return nil
	}

	// Initialize speaker with sample rate and buffer size
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	if err != nil {
		return err
	}

	speakerInitialized = true
	return nil
}

// PlaySound plays a retro-style beep sound
// If speaker is not initialized, falls back to silent operation
func PlaySound(soundType string) {
	speakerMutex.Lock()
	initialized := speakerInitialized
	speakerMutex.Unlock()

	if !initialized {
		// Fallback: silent (could use fmt.Print("\a") here if desired)
		return
	}

	// Play sound asynchronously to avoid blocking game loop
	go playSound(soundType)
}

// playSound generates and plays the appropriate sound based on type
func playSound(soundType string) {
	switch soundType {
	case "jump":
		// Single high-pitched beep (A4 = 440Hz, 100ms)
		playTone(440, 100*time.Millisecond)

	case "score":
		// Two ascending tones (C5 = 523Hz, E5 = 659Hz)
		playTone(523, 80*time.Millisecond)
		time.Sleep(20 * time.Millisecond)
		playTone(659, 80*time.Millisecond)

	case "gameover":
		// Three descending tones (E5 → C5 → A4)
		playTone(659, 150*time.Millisecond)
		time.Sleep(50 * time.Millisecond)
		playTone(523, 150*time.Millisecond)
		time.Sleep(50 * time.Millisecond)
		playTone(440, 150*time.Millisecond)
	}
}

// playTone generates and plays a tone at the specified frequency and duration
func playTone(frequency int, duration time.Duration) {
	// Generate sine wave tone (retro game sound)
	tone, err := generators.SinTone(sampleRate, frequency)
	if err != nil {
		return // Skip if frequency is invalid
	}

	// Adjust volume to -3dB (slightly quieter than default)
	volume := &effects.Volume{
		Streamer: tone,
		Base:     2,
		Volume:   -3,
		Silent:   false,
	}

	// Limit duration
	limited := beep.Take(sampleRate.N(duration), volume)

	// Play and wait for completion
	done := make(chan bool)
	speaker.Play(beep.Seq(limited, beep.Callback(func() {
		done <- true
	})))

	<-done
}
