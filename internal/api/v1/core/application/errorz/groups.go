package errorz

var InvalidGroupIdError = Error_{
	StatusCode: 400,
	Message:    "Invalid Group ID. It must be int64.",
}
