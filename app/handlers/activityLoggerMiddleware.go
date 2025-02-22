package handlers

import (
	"log"
	"net/http"

	"github.com/jental/freetesl-server/services"
)

func ActivityLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var playerID int = -1
		contextVal := req.Context().Value("userID")
		if contextVal == nil {
			log.Println("player id is not found in a context")
		} else {
			var ok bool = false
			playerID, ok = contextVal.(int)
			if !ok {
				log.Println("player id from a context has invalid type")
				return
			}
		}
		if playerID >= 0 {
			log.Printf("ActivityLoggerMiddleware: Player id: %d", playerID)
			services.UpdatePlayerLastActivityTime(playerID)
		}

		next.ServeHTTP(w, req)
	})
}
