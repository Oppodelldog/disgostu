package api

import "net/http"

func RunApi() {
	http.HandleFunc("/test/capture/start", captureStartHandler)
	http.ListenAndServe("0.0.0.0:7676", nil)
}

func captureStartHandler(rw http.ResponseWriter, r *http.Request) {

}
