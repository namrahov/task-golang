package model

const Exception = "error.task-golang"

const (
	HeaderKeyCustomerID = "DP-Customer-ID"
	HeaderKeyUserID     = "DP-User-ID"
	HeaderKeyUserAgent  = "User-Agent"
	HeaderKeyUserIP     = "X-Forwarded-For"
	HeaderKeyRequestID  = "requestid"
)

const (
	LoggerKeyRequestID  = "REQUEST_ID"
	LoggerKeyOperation  = "OPERATION"
	LoggerKeyCustomerID = "CUSTOMER_ID"
	LoggerKeyUserID     = "USER_ID"
	LoggerKeyUserIP     = "USER_IP"
	LoggerKeyUserAgent  = "USER_AGENT"
	ContextLogger       = "contextLogger"
	ContextHeader       = "contextHeader"
)

type ContextKey string

const (
	ContextUserID     ContextKey = "userID"
	ContextUserRoles  ContextKey = "userRoles"
	ContextAuthHeader ContextKey = "authHeader"
)

type LinkType string

const (
	Registration   LinkType = "REGISTRATION"
	ChangeEmail    LinkType = "CHANGE_EMAIL"
	ForgetPassword LinkType = "FORGET_PASSWORD"
	SetPassword    LinkType = "SET_PASSWORD"
)
