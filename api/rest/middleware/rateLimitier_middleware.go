package middleware

import (
	"sync"
	"time"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
)

type client struct {
	lastSeen time.Time
	limit    int
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func DeletedClient(typeReq []string) {
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			for key, client := range clients {
				if typeReq[0] == "activasi-account" {
					if time.Since(client.lastSeen) > 10*time.Minute {
						delete(clients, key)
						log.Info().Msgf("deleted key %s", key)
					}
				} else {
					if time.Since(client.lastSeen) > 10*time.Minute {
						delete(clients, key)
						log.Info().Msgf("deleted key %s", key)
					}
				}

			}
			mu.Unlock()
		}
	}()
}

func DeletedClientHelper(key string) {
	log.Info().Msgf("Deleted otp %s", key)
	delete(clients, key)
}

func RateLimiterOTP(req *domain.RequestGenerateOTP) error {
	key := req.Email + ":" + req.Type
	_, found := clients[key]

	if !found {
		log.Info().Msgf("set limiter %s", key)
		clients[key] = &client{
			limit: 1,
		}
	} else {
		clients[key].limit = clients[key].limit + 1
	}

	clients[key].lastSeen = time.Now()
	if clients[key].limit >= 5 {
		log.Warn().Msg("rate limited")
		return _error.HttpErrString("anda sudah mencapai batas limit", response.CM12)
	}
	return nil
}
