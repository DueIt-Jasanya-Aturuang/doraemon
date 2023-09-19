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
	LogErrOpenFile  = "failed open file | file : %v | err : %v"
	LogErrCloseFile = "failed close file | file : %v | err : %v"
)

const (
	LogErrPutObjectMinio   = "failed put object minio | err : %v"
	LogErrDelObjectMinio   = "failed delete object minio | err : %v"
	LogInfoFileUploadMinio = "file upload info | msg : %v"
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
