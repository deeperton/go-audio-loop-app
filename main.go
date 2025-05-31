package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gordonklaus/portaudio"
)

func main() {
	var latencyMs int
	flag.IntVar(&latencyMs, "latency", 50, "desired latency in milliseconds (10–300)")
	flag.Parse()
	if latencyMs < 10 || latencyMs > 300 {
		log.Fatalf("invalid latency %dms: must be between 10 and 300", latencyMs)
	}
	fmt.Printf("Starting loopback with %dms latency…\n", latencyMs)

	// Initialize PortAudio
	if err := portaudio.Initialize(); err != nil {
		log.Fatalf("failed to initialize PortAudio: %v", err)
	}
	defer portaudio.Terminate()

	const sampleRate = 44100   // 44.1 kHz
	const channels = 1         // mono
	const framesPerBuffer = 64 // small I/O chunk
	bufferFrames := latencyMs * sampleRate / 1000

	// Create our ring buffer
	ring := make([]float32, bufferFrames)
	writeIdx, readIdx := 0, 0

	inBuf := make([]float32, framesPerBuffer)
	outBuf := make([]float32, framesPerBuffer)

	// Open input (microphone) stream
	inStream, err := portaudio.OpenDefaultStream(channels, 0, sampleRate, len(inBuf), &inBuf)
	if err != nil {
		log.Fatalf("opening input stream error: %v", err)
	}
	defer inStream.Close()

	// Open output (speaker) stream
	outStream, err := portaudio.OpenDefaultStream(0, channels, sampleRate, len(outBuf), &outBuf)
	if err != nil {
		log.Fatalf("opening output stream error: %v", err)
	}
	defer outStream.Close()

	// Start both streams
	if err := inStream.Start(); err != nil {
		log.Fatalf("starting input stream: %v", err)
	}
	defer inStream.Stop()
	if err := outStream.Start(); err != nil {
		log.Fatalf("starting output stream: %v", err)
	}
	defer outStream.Stop()

	// Handle Ctrl+C to exit cleanly
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	fmt.Println("Running. Press Ctrl+C to stop.")

	// Main loop: read from mic, write into ring buffer; read from ring, write to speakers
loop:
	for {
		select {
		case <-sig:
			fmt.Println("\nInterrupted; shutting down.")
			break loop
		default:
			// 1. Read one chunk from microphone
			if err := inStream.Read(); err != nil {
				log.Fatalf("input read error: %v", err)
			}
			// 2. Write into ring buffer
			for i := 0; i < len(inBuf); i++ {
				ring[writeIdx] = inBuf[i]
				writeIdx = (writeIdx + 1) % len(ring)
			}
			// 3. Read delayed samples from ring buffer
			for i := 0; i < len(outBuf); i++ {
				outBuf[i] = ring[readIdx]
				readIdx = (readIdx + 1) % len(ring)
			}
			// 4. Send to speakers
			if err := outStream.Write(); err != nil {
				log.Fatalf("output write error: %v", err)
			}
		}
	}
}
