package encryption

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
)

func TestDecryptStringCFB(t *testing.T) {
	key := "rp1LZ7WNMXaq9KhK2LIkMbAA8qbMkFQb"
	text := "6fe9168f11b0ea3cf529e3363d7e74f09f95a84b"
	res, err := DecryptStringCFB(text, key)
	assert.NoError(t, err)
	t.Log(res)
}

func TestEncrypStringCFB(t *testing.T) {
	key := "rp1LZ7WNMXaq9KhK2LIkMbAA8qbMkFQb"
	res, err := EncrypStringCFB("rama", key)
	assert.NoError(t, err)
	t.Log(res)
}
func TestDecryptStringCBC(t *testing.T) {
	key := "rp1LZ7WNMXaq9KhK2LIkMbAA8qbMkFQb"
	iv := "auTHxeaZSrsxfIZI"
	text := "62idN0WO0Ym9cqcKnCPusz717WPOBK+VSxDhoArTVmNiUiSleh/y9v9lulpNMt0WCIkN2Lt9vArYIrsm2Lv3D+m+e980FT13qjL8HMs+L20vzQ1A5ug1kQJGiyfpDGVvZ5V4iKsR9lijfveEToBdDe+rjteXydaUSi1dlUCoJO6NyNtf2nl56KLir/WpESg+rDEEptbWPe+dXWOtnocHLLzk1Bgj5XnFE+E7UUKd1rxSeMgOzOioWyeyeLjOf3MZUrVz81f7idGoGhDqfdM/ZLXaLUcUxirNKFK3jMCBoCrmkFc2d2+Qa/gqyrRPngRfjRVbYv/QEVxDAPGJeAQwCwiqzUC3s6b4p4tOl+mz8s2iO0ZVeufPtRg3nv+IE3SbPu83qDliJD2HoXGa9QYvKWDSQTjR8Ds/7Nvlyf8AOBXgFIF87BcDH+UsJavl0z5D6NPv7VHj8ZO3QF3K86uva2nYIR9erVF7P0jM4aRSzp27nIWNfE1e/cpc7ZoDZQ7SamXODpjNJ/BSfhoVqpmOndZ7+XRqnhTVvV35WRfPAt6vryl6uIUpi7xq6+05vEKwN9pe2082HRR5YSFjEvWLETpcf75Qhh3D1gVSbA4Lasj7ybghw2zddMoxvleRzd6+xeGBK6qEE7dhnUagRlFBLxXTLxBB59tlsAf3gY8AlyHAkp+3FdA/zjS80G3Cht23y5F/ieTdCYvuna4rXG/gjjEj9cRjGCojyT/9Ko5sgeTSNJw3Rh5A8H9pGMTG/Euq4++ZF/N2TAazJnV8FR0pOtZ0czry5a6/CTtHjKEc/cIKEppLHpFOnR93Y5NVpDBdef7NR05E2xfrFIZ/A/64Tn+Pq1f2z3EjDCBwV6VqVITPU7Xhvs1NJxPGNKPX9hhRl4DTTu+3qhsCBaf9VH3rzrkw5x6RqqAcF6CX2aVRzwRr8+qmtwemljTUVkvrT2C43JDvvywcZizNuf+lD/68Y3Jujqdh1PC814Y9mefaTNp67jnInHw1YVQ6LV/eWH2M98ZIQxO1Aza3asqoQqJHXQ58w12R7eWSLimB7PlHplOJYcJRDP12yC4GEl8xR4nEj2HNiIuuKmLCtuTocCZUpllLv8DBWnShvYOS/Xz+s3l1cVz76N0EZk/vspusEdHTLgVIfLOmsZuJdy96NTo2G9k6WyqYUqGnCFL2iJ+z3azn073f+Pz5BpnS0yVrFnoZFK8Bw3h66JVASHqMrFiFujeoTeNp6vC8JBdAbMK0xU0S4WyeiBcIt8xWYm5RjIoMZd+LuSIgm4H8wJhcsHGCe6chz/Jg3qG4hXE35vxLsqMCxb6HjnSW6/bZslFU7/iEtyoVJhQql6JCOz8MKNzr/+JmFqJjxVzfvA+6PSa7yb0A8apv5jVT95kCof5vkUDQZIUDLSvab8rK1R6OZoS4hquVECtdapmjqMH6NzAiwMKzIMM0Q+tMe6Fw7PjqG7G8E6lf7jdqNKkNYdEf4AYRcreHfZ9Y97BCZXcV0Ap9VWyzuFXeCqRozBCttNV+MWEJTGtkcZPt8OlSeUfOfBEiQr2K0P1l3Qgt+VMCzvAjUDAvNGzTZfvPUBpb0nq/CzvCf0nXlVStLr6/sGg0CTTG6dbOkOpfld2j6D2HYpVAPmCEBrg4cmIYGb7mXOQlBbRowLAWuVEj6W9SJY4e0373WZIA3AERo5wnmRM++jKbPUaQKZCgii0ujKyxHXUeeFYMJV+DPEDygzrYXv8yC2VRtWSeV70OKQJ/hL8ulqw8SZCrvFDI3ZAlLeYPQzted90y2JEchXtB3v9xMbBAzf9Ziw=="
	res, err := DecryptStringCBC(text, key, iv)
	var googleOauthToken *domain.Oauth2GoogleToken
	err = json.Unmarshal([]byte(res), &googleOauthToken)
	assert.NoError(t, err)
	t.Log(googleOauthToken.IDToken)
}

func TestEncrypStringCBC(t *testing.T) {
	key := "rp1LZ7WNMXaq9KhK2LIkMbAA8qbMkFQb"
	iv := "mydigit15digit11"

	// req := map[string]string{
	// 	"access_token": "this is access token",
	// 	"id_token":     "this is access token",
	// }
	//
	// reqJson, _ := json.Marshal(req)

	res, err := EncrypStringCBC("rama", key, iv)
	assert.NoError(t, err)
	t.Log(res)
}
