package util

const (
	LogErrDBConn              = "Failed Open Database | err : %v"
	LogErrDBConnClose         = "Failed close Database | err : %v"
	LogErrBeginTx             = "Failed Start Transaction | err : %v"
	LogErrRollback            = "failed rollback data | err : %v"
	LogInfoRollback           = "rollback data | err : %v"
	LogErrCommit              = "failed commit data | err %v"
	LogInfoCommit             = "commit data"
	LogErrPrepareContextClose = "failed close prepared context | err : %v"
	LogErrPrepareContext      = "failed open prepared context | err : %v"
	LogErrExecContext         = "failed exec context | err : %v"
	LogErrQueryRowContextScan = "failed scan data | err : %v"
	LogErrQueryRows           = "failed Rows data | err : %v"
	LogErrQueryRowsClose      = "failed Close Rows data | err : %v"
	LogErrQueryRowsScan       = "failed scan data | err : %v"
)

const (
	LogErrDecode    = "failed decode | data : %v | err : %v"
	LogErrEncode    = "failed encode | data : %v | err : %v"
	LogErrUnmarshal = "failed unmarshal | data : %v | err : %v "
	LogErrMarshal   = "failed marshal | data : %v | err : %v"
)

const (
	LogErrHttpNewRequest = "failed new request | err : %v"
	LogErrClientDo       = "failed client do http | err : %v"
	LogErrClientDoClose  = "failed close client do http | err : %v"
)

const (
	LogErrRedisClientExists = "failed check data | err : %v"
	LogErrRedisClientGet    = "failed get data | err : %v"
	LogErrRedisClientExpire = "failed set expire data | err : %v"
	LogErrRedisClientSet    = "failed set expire data | err : %v"
	LogErrRedisClientDel    = "failed delete data | err : %v"
)

const (
	LogErrKafkaWriterClose  = "failed close kafka writer | err : %v"
	LogErrKafkaWriteMessage = "failed publish message kafka | err : %v"
)

const (
	LogErrFailedClaimJwt = "failed claims jwt token | err : %v | token : %s"
)
