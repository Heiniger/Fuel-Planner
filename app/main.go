// Pit Wall Fuel Planner desktop launcher.
// Serves the embedded planner on a local port and opens it in an app-style
// Microsoft Edge (or Chrome) window. Exits about 20 s after the window closes.
package main

import (
	_ "embed"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

//go:embed index.html
var page string

const port = "47186" // fixed so saved settings (browser storage) survive restarts

const heartbeat = `<script>setInterval(()=>fetch("/ping").catch(()=>{}),5000);fetch("/ping");</script></body>`

func main() {
	url := "http://127.0.0.1:" + port + "/"
	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil { // already running: just open another window
		openWindow(url)
		return
	}
	html := strings.Replace(page, "</body>", heartbeat, 1)
	var last atomic.Int64
	last.Store(time.Now().Unix() + 60) // grace period while the window opens
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		last.Store(time.Now().Unix())
		w.WriteHeader(204)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Write([]byte(html))
	})
	go http.Serve(ln, mux)
	openWindow(url)
	for {
		time.Sleep(5 * time.Second)
		if time.Now().Unix()-last.Load() > 20 {
			os.Exit(0)
		}
	}
}

func openWindow(url string) {
	pf, pf86, la := os.Getenv("ProgramFiles"), os.Getenv("ProgramFiles(x86)"), os.Getenv("LocalAppData")
	candidates := []string{
		filepath.Join(pf86, `Microsoft\Edge\Application\msedge.exe`),
		filepath.Join(pf, `Microsoft\Edge\Application\msedge.exe`),
		filepath.Join(pf, `Google\Chrome\Application\chrome.exe`),
		filepath.Join(pf86, `Google\Chrome\Application\chrome.exe`),
		filepath.Join(la, `Google\Chrome\Application\chrome.exe`),
	}
	for _, b := range candidates {
		if _, err := os.Stat(b); err == nil {
			if exec.Command(b, "--app="+url, "--window-size=1360,900").Start() == nil {
				return
			}
		}
	}
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start() // default browser
}
