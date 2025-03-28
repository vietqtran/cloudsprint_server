package constants

type HttpStatusCode int

const (
	StatusContinue                      HttpStatusCode = 100
	StatusSwitchingProtocols            HttpStatusCode = 101
	StatusProcessing                    HttpStatusCode = 102
	StatusEarlyHints                    HttpStatusCode = 103
	StatusOK                            HttpStatusCode = 200
	StatusCreated                       HttpStatusCode = 201
	StatusAccepted                      HttpStatusCode = 202
	StatusNonAuthoritativeInformation   HttpStatusCode = 203
	StatusNoContent                     HttpStatusCode = 204
	StatusResetContent                  HttpStatusCode = 205
	StatusPartialContent                HttpStatusCode = 206
	StatusMultiStatus                   HttpStatusCode = 207
	StatusAlreadyReported               HttpStatusCode = 208
	StatusIMUsed                        HttpStatusCode = 226
	StatusMultipleChoices               HttpStatusCode = 300
	StatusMovedPermanently              HttpStatusCode = 301
	StatusFound                         HttpStatusCode = 302
	StatusSeeOther                      HttpStatusCode = 303
	StatusNotModified                   HttpStatusCode = 304
	StatusUseProxy                      HttpStatusCode = 305
	StatusTemporaryRedirect             HttpStatusCode = 307
	StatusPermanentRedirect             HttpStatusCode = 308
	StatusBadRequest                    HttpStatusCode = 400
	StatusUnauthorized                  HttpStatusCode = 401
	StatusPaymentRequired               HttpStatusCode = 402
	StatusForbidden                     HttpStatusCode = 403
	StatusNotFound                      HttpStatusCode = 404
	StatusMethodNotAllowed              HttpStatusCode = 405
	StatusNotAcceptable                 HttpStatusCode = 406
	StatusProxyAuthRequired             HttpStatusCode = 407
	StatusRequestTimeout                HttpStatusCode = 408
	StatusConflict                      HttpStatusCode = 409
	StatusGone                          HttpStatusCode = 410
	StatusLengthRequired                HttpStatusCode = 411
	StatusPreconditionFailed            HttpStatusCode = 412
	StatusRequestEntityTooLarge         HttpStatusCode = 413
	StatusRequestURITooLong             HttpStatusCode = 414
	StatusUnsupportedMediaType          HttpStatusCode = 415
	StatusRequestedRangeNotSatisfiable  HttpStatusCode = 416
	StatusExpectationFailed             HttpStatusCode = 417
	StatusTeapot                        HttpStatusCode = 418
	StatusMisdirectedRequest            HttpStatusCode = 421
	StatusUnprocessableEntity           HttpStatusCode = 422
	StatusLocked                        HttpStatusCode = 423
	StatusFailedDependency              HttpStatusCode = 424
	StatusTooEarly                      HttpStatusCode = 425
	StatusUpgradeRequired               HttpStatusCode = 426
	StatusPreconditionRequired          HttpStatusCode = 428
	StatusTooManyRequests               HttpStatusCode = 429
	StatusRequestHeaderFieldsTooLarge   HttpStatusCode = 431
	StatusUnavailableForLegalReasons    HttpStatusCode = 451
	StatusInternalServerError           HttpStatusCode = 500
	StatusNotImplemented                HttpStatusCode = 501
	StatusBadGateway                    HttpStatusCode = 502
	StatusServiceUnavailable            HttpStatusCode = 503
	StatusGatewayTimeout                HttpStatusCode = 504
	StatusHTTPVersionNotSupported       HttpStatusCode = 505
	StatusVariantAlsoNegotiates         HttpStatusCode = 506
	StatusInsufficientStorage           HttpStatusCode = 507
	StatusLoopDetected                  HttpStatusCode = 508
	StatusNotExtended                   HttpStatusCode = 510
	StatusNetworkAuthenticationRequired HttpStatusCode = 511
)
