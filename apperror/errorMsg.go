package apperror

import "strings"

var PARSE_ERROR = "cannot parse request body"
var INVALID_INPUT = "invalid request body"
var AUTHENTICATE_FAILED = "username or password is incorrect"
var MARSHAL_FAILED = "failed to marshal object to JSON"
var SEND_MSG_FAILED = "failed to send message to queue"
var MISSING_USER_ID = "user id not found in context"

func ServiceError(serviceName string) string {
	return strings.ToLower(serviceName) + " service failed"
}
