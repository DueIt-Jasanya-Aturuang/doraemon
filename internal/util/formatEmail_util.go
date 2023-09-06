package util

import (
	"fmt"
	"strings"
)

func EmailFormat(email string) string {
	emailArr := strings.Split(email, "@")
	emailString := fmt.Sprintf("%c••••%c@%s", emailArr[0][0], emailArr[0][len(emailArr[0])-1], emailArr[1])
	return emailString
}
