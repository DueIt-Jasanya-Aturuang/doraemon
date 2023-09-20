package helper

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func EmailFormat(email string) string {
	emailArr := strings.Split(email, "@")
	if len(emailArr) != 2 {
		log.Warn().Msgf("email user tidak valid | total : %d", len(emailArr))
		return email
	}
	emailString := fmt.Sprintf("%c••••%c@%s", emailArr[0][0], emailArr[0][len(emailArr[0])-1], emailArr[1])
	return emailString
}
