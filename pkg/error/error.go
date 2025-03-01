package errorcode

type ErrorCode = int32

const (
	// Common
	Ok               ErrorCode = 0
	InvalidParam     ErrorCode = 10
	RpcRequestFailed ErrorCode = 11
	// Login Service
	UserNotExist       ErrorCode = 100
	LoginPasswordError ErrorCode = 101
	UserAlreadyExist   ErrorCode = 102
	// Product Service
	ProductAlreadyExist ErrorCode = 200

	// Auth Service
	AuthDeliveryTokenError ErrorCode = 400
	AuthVerifyTokenError   ErrorCode = 400

	// Stock Service
	FlashNoStock ErrorCode = 701

	UnknowError ErrorCode = 999
)
