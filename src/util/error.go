package util

const (
	StatusOK          = iota
	CreateUserByLogin = 1

	InvalidForm = 1001
	InvalidURI  = 1002
	InvalidBody = 1003

	PasswordIncorrect = 2000
	CookieNotSet      = 2001
	AuthorizedError   = 2002
	ContextInfoNotSet = 2003

	WebServerUnavailable  = 3001
	DataBaseUnavailable   = 3002
	RedisUnavailable      = 3003
	FileSystemUnavailable = 3004

	DatabaseInsertFailed = 4001

	UnauthorizedOperation = 5001

	RegisterCodeIncorrect = 6001

	WriteFileError = 7001
	FileDamaged    = 7002
)
