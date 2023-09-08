package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
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
					if time.Since(client.lastSeen) > 1*time.Hour {
						delete(clients, key)
						log.Info().Msg("deleted ip")
						log.Info().Msgf("%v", clients[key])
					}
				} else {
					if time.Since(client.lastSeen) > 6*time.Hour {
						delete(clients, key)
						log.Info().Msg("deleted ip")
						log.Info().Msgf("%v", clients[key])
					}
				}

			}
			mu.Unlock()
		}
	}()
}

func DeletedClientHelper(key string) {
	delete(clients, key)
}

func RateLimiterOTP(req *dto.OTPGenerateReq) error {
	key := req.Email + ":" + req.Type
	_, found := clients[key]

	if !found {
		log.Info().Msg("set limiter")
		clients[key] = &client{
			limit: 1,
		}
	} else {
		clients[key].limit = clients[key].limit + 1
	}

	clients[key].lastSeen = time.Now()
	if clients[key].limit >= 3 {
		log.Warn().Msg("rate limited")
		return _error.ErrStringDefault(http.StatusTooManyRequests)
	}
	return nil
}
