// What it does:
//
// This example opens a video capture device, then streams MJPEG from it.
// Once running point your browser to the hostname/port you passed in the
// command line (for example http://localhost:8080) and you should see
// the live video stream.
//
// How to run:
//
// mjpeg-streamer [camera ID] [host:port]
//
//		go get -u github.com/hybridgroup/mjpeg
// 		go run ./cmd/mjpeg-streamer/main.go 1 0.0.0.0:8080
//

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	_ "net/http/pprof"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

var (
	deviceID int
	err      error
	webcam   *gocv.VideoCapture
	stream   *mjpeg.Stream
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\tmjpeg-streamer [camera ID] [host:port]")
		return
	}

	// parse args
	deviceID := os.Args[1]
	host := os.Args[2]

	// open webcam
	webcam, err = gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	// create the mjpeg stream
	stream = mjpeg.NewStream()

	// start capturing
	go mjpegCapture()

	fmt.Println("Capturing. Point your browser to " + host)

	// start http server
	http.Handle("/mjpeg", stream)
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("shutdown -h +1")

		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `
<html>
	<head>
	  <title>WeGYB</title>
	</head>
	<body>
	  Failed to shutdwon
	</body>
</html>
			`)
    } else {
			fmt.Fprintf(w, `
<html>
	<head>
		<title>WeGYB</title>
	</head>
	<body>
		Shutdown in 1 minute
	</body>
</html>
			`)
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
<html>
	<head>
	  <title>WeGYB</title>
	</head>
	<style>
	img {
		width: 100vw;
		height: 100vh;
		object-fit: contain;
	}
	</style>

	<body>
	  <img src="/mjpeg">
	</body>
</html>
		`)
	})
	log.Fatal(http.ListenAndServe(host, nil))
}

func mjpegCapture() {
	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf.GetBytes())
		buf.Close()
	}
}
