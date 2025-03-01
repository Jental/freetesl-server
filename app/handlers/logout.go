package handlers

import "net/http"

func Logout(w http.ResponseWriter, req *http.Request) {
	// for now we do nothing here
	// FE will forget token
	// and as for matches that might be present for logged out user - he may want to log back in soon
	// => we do nothing match-related here
	//    but there'll be a background process, that checks activity of users and after some timeout may disconnect user and abort his match
	w.WriteHeader(http.StatusOK)
}
