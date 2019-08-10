package handlers

import (
	"net/http"
)

/* TODO: implement a CORS middleware handler, as described
in https://drstearns.github.io/tutorials/cors/ that responds
with the following headers to all requests:

  Access-Control-Allow-Origin: *
  Access-Control-Allow-Methods: GET, PUT, POST, PATCH, DELETE
  Access-Control-Allow-Headers: Content-Type, Authorization
  Access-Control-Expose-Headers: Authorization
  Access-Control-Max-Age: 600
*/

//ResponseHeader is a middleware handler that adds a header to the response
type ResponseHeader struct {
  handler http.Handler
}

//NewResponseHeader constructs a new ResponseHeader middleware handler
func NewResponseHeader(handlerToWrap http.Handler) *ResponseHeader {
  return &ResponseHeader{handlerToWrap}
}

//ServeHTTP handles the request by adding the response header
func (rh *ResponseHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  //add the headers
  w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Add("Access-Control-Expose-Headers", "Authorization")
  w.Header().Add("Access-Control-Max-Age", "600")
  
  if r.Method == "OPTIONS" {
    w.WriteHeader(http.StatusOK)
    return
  }
  //call the wrapped handler
  rh.handler.ServeHTTP(w, r)
}