package comms

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	t "github.com/foadmom/yari/types"
)

var _httpServer *http.Server

func InitHTTP(config *t.ConfigType) {
	_logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	_router := http.NewServeMux()
	_router.HandleFunc("/YARI/message", messageHandler)
	_router.HandleFunc("/YARI/shutdown", shutdownHandler)

	_httpServer = &http.Server{Addr: "localhost:8669",
		Handler:      _router,
		ErrorLog:     _logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	_err := _httpServer.ListenAndServe()
	if _err != nil {
		fmt.Println(_err)
	}

	defer _httpServer.Close()
}

func messageHandler(w http.ResponseWriter, r *http.Request) {

	// switch r.Method {
	// case "POST":
	// 	var	_request 	[]byte;
	// 	_request, _err := ioutil.ReadAll(r.Body);
	// 	if (_err == nil) {
	// 		fmt.Println (string(_request));
	// 		_resp := ProcessMessage (nil, _request);
	// 		fmt.Fprintln(w, string(_resp));
	// 	}
	// }
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "shutting down")
	//	time.Sleep(3*time.Second);
	_httpServer.Close()
}
