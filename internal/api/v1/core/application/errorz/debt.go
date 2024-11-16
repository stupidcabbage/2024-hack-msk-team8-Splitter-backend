package errorz

var ErrDebtNotFound = Error_{
	StatusCode: 404,
	Message:    "Debt not found",
}

var ErrDebtAffectingPermissionDenied = Error_{
	StatusCode: 403,
	Message:    "Debt affecting access restricted",
}

var ErrDebtIsClosed = Error_{
	StatusCode: 400,
	Message:    "Cant affect closed debt",
}

var ErrCantDebtYourself = Error_{
	StatusCode: 400,
	Message:    "Cant debt yourself",
}

var ErrCantDebtNotInGroup = Error_{
	StatusCode: 400,
	Message:    "Cant debt not in group",
}