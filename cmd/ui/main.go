package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

//go:embed static
var staticFiles embed.FS

type BombardierRequest struct {
	URL               string            `json:"url"`
	Method            string            `json:"method"`
	NumConnections    int               `json:"numConnections"`
	NumRequests       *int              `json:"numRequests"`
	Duration          *int              `json:"duration"` // in seconds
	Headers           map[string]string `json:"headers"`
	Body              string            `json:"body"`
	Rate              *int              `json:"rate"`
	Timeout           int               `json:"timeout"` // in ms
	Latencies         bool              `json:"latencies"`
	InsecureSkipVerify bool             `json:"insecureSkipVerify"`
	DisableKeepAlives bool              `json:"disableKeepAlives"`
	ClientType        string            `json:"clientType"` // "fasthttp" or "nethttp"
}

type BombardierResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Output       string `json:"output,omitempty"`
}

func main() {
	// Create endpoints for the API
	http.HandleFunc("/api/run", handleRun)
	
	// Serve static files
	http.Handle("/", http.FileServer(http.FS(staticFiles)))
	
	port := 8080
	fmt.Printf("Starting Bombardier UI server on port %d\n", port)
	fmt.Printf("Open your browser and navigate to http://localhost:%d\n", port)
	openBrowser(fmt.Sprintf("http://localhost:%d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request BombardierRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Convert the request into bombardier command arguments
	args := buildBombardierArgs(request)
	
	// Execute the bombardier command
	cmd := exec.Command("bombardier", args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Error running bombardier: %v\n%s", err, output), http.StatusInternalServerError)
		return
	}
	
	// Return the output
	response := BombardierResponse{
		Success: true,
		Output:  string(output),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func buildBombardierArgs(req BombardierRequest) []string {
	var args []string
	
	// Add connections
	args = append(args, "-c", fmt.Sprintf("%d", req.NumConnections))
	
	// Add number of requests or duration
	if req.NumRequests != nil {
		args = append(args, "-n", fmt.Sprintf("%d", *req.NumRequests))
	} else if req.Duration != nil {
		args = append(args, "-d", fmt.Sprintf("%ds", *req.Duration))
	}
	
	// Add method
	args = append(args, "-m", req.Method)
	
	// Add timeout
	args = append(args, "-t", fmt.Sprintf("%dms", req.Timeout))
	
	// Add headers
	for key, value := range req.Headers {
		args = append(args, "-H", fmt.Sprintf("%s: %s", key, value))
	}
	
	// Add body if any
	if req.Body != "" {
		// Write body to a temporary file to avoid command line length limitations
		tmpFile, err := ioutil.TempFile("", "bombardier-body-*.txt")
		if err == nil {
			defer os.Remove(tmpFile.Name())
			tmpFile.WriteString(req.Body)
			tmpFile.Close()
			args = append(args, "-f", tmpFile.Name())
		}
	}
	
	// Add rate if any
	if req.Rate != nil {
		args = append(args, "--rate", fmt.Sprintf("%d", *req.Rate))
	}
	
	// Add latencies if enabled
	if req.Latencies {
		args = append(args, "-l")
	}
	
	// Add insecure if enabled
	if req.InsecureSkipVerify {
		args = append(args, "--insecure")
	}
	
	// Add disable keep-alives if enabled
	if req.DisableKeepAlives {
		args = append(args, "-k")
	}
	
	// Add client type
	if req.ClientType == "nethttp" {
		args = append(args, "--http1")
	}
	
	// Add the URL (always last)
	args = append(args, req.URL)
	
	return args
}

func sendErrorResponse(w http.ResponseWriter, message string, status int) {
	response := BombardierResponse{
		Success:      false,
		ErrorMessage: message,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// openBrowser opens the default browser with the specified URL
func openBrowser(url string) {
	var err error
	
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
} 