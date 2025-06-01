# Go Audio Loopback with PortAudio

A minimal Go application that captures audio from the default microphone and plays it back through the default speakers with a user-specified latency (10 ms–300 ms) using PortAudio.

## Overview

When you run this program, it:

1. Opens your system’s default input device (microphone).
2. Opens your system’s default output device (speakers).
3. Buffers incoming audio samples in a ring buffer sized to the chosen latency.
4. Continuously reads from the ring buffer and writes to the output, creating a fixed delay.

## Requirements

- Go 1.18+ installed on your machine.
- PortAudio development headers and libraries (see instructions below for macOS, Linux, or Windows).

## Installation

1. **Install PortAudio**  
   - **macOS (Homebrew)**  
     ```bash
     brew install portaudio
     ```  
   - **Ubuntu/Debian**  
     ```bash
     sudo apt-get update
     sudo apt-get install portaudio19-dev
     ```  
   - **Windows**  
     1. Download a precompiled PortAudio binary or build from source:  
        https://www.portaudio.com/download.html  
     2. Make sure the `portaudio.dll` (or `.lib`) and headers are accessible to your Go compiler (e.g., via `CGO_CFLAGS`/`CGO_LDFLAGS` or by placing them in a standard include/lib path).

2. **Fetch the Go binding**  
   ```bash
   go get github.com/gordonklaus/portaudio
    ````

3. **Clone (or copy) the source file**
   Save the code (from the original example) into a file named `main.go` in your project folder.

4. **Build the binary**

   ```bash
   go build -o audioloop main.go
   ```

   This will produce an executable named `audioloop` (or `audioloop.exe` on Windows) in the current directory.

---

## How to Use

```bash
./audioloop -latency <milliseconds>
```

* `-latency`
  Desired loopback latency in milliseconds. Must be between **10** and **300**.

    * Example: `-latency 100` sets a 100 ms delay.

### Example Runs

* **Start with 50 ms latency**

  ```bash
  ./audioloop -latency 50
  ```

  You should see:

  ```
  Starting loopback with 50ms latency…
  Running. Press Ctrl+C to stop.
  ```

* **Start with 200 ms latency**

  ```bash
  ./audioloop -latency 200
  ```

  You should see:

  ```
  Starting loopback with 200ms latency…
  Running. Press Ctrl+C to stop.
  ```

To stop the loopback, press **Ctrl+C** in the terminal. The program will close both audio streams and exit cleanly.

---

## Troubleshooting

* If the build fails due to missing PortAudio headers or libraries:

    * Verify that PortAudio is installed and that your environment variables (e.g., `CGO_CFLAGS`, `CGO_LDFLAGS`) point to the correct include/lib directories.
    * On Linux, confirm that `portaudio19-dev` is installed.
    * On macOS, make sure `brew install portaudio` succeeded and that `/usr/local/include` (or `/opt/homebrew/include` on Apple Silicon) is in your include path.

* If you get audio distortion or dropouts:

    * Try reducing `framesPerBuffer` or testing a different `sampleRate` (e.g., 48000).
    * Ensure no other application holds exclusive access to your audio device.

---

## License

This project is released under the MIT License. Feel free to copy and modify as needed.

