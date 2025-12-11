package pb

type Task struct {
	Id          int32  `json:"id"`          // task id
	Cron        string `json:"cron"`        // cron expression
	Source      string `json:"source"`      // source image name
	Destination string `json:"destination"` // destination image name
}

type Log struct {
	Id     int32  `json:"id"`     // log id
	TaskId int32  `json:"taskId"` // task id
	Msg    string `json:"msg"`    // log message
	Time   int64  `json:"time"`   // log time
}

type Once struct {
	Source      string `json:"source"`      // source image name
	Destination string `json:"destination"` // destination image name
}
