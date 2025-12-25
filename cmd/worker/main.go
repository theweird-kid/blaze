package worker

import (
	"log"
	"net/http"

	"github.com/theweird-kid/blaze/internal/worker"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/execute", worker.HandleExecute)

	log.Println("worker listeneing on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
