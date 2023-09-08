package _msg

const (
	LogErrStartPrepareContext  = "failed to start prepared statement"
	LogErrClosePrepareContext  = "failed to close prepared context"
	LogErrExecContext          = "failed to query row context prepared statement"
	LogErrQueryRowContextScan  = "cannot scan query row context"
	LogErrHttpNewRequest       = "failed create http request"
	LogErrSetRedisClient       = "failed set data in redis"
	LogErrExistsRedisClient    = "failed cek data in redis"
	LogErrExpireRedisClient    = "failed set expire data in redis"
	LogErrGetRedisClient       = "failed get data in redis"
	LogErrDelRedisClient       = "failed delete data in redis"
	LogErrHttpClientDo         = "failed get response from http request post\nfailed get response from http request post"
	LogErrResponseBodyClose    = "failed close response body\nfailed close response body"
	LogErrJsonNewDecoderDecode = "failed decode response to struct\nfailed decode response to struct"
)
