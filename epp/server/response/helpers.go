package response

const (
	XmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>`

	CommandCompletedSuccessfully                  = "Command completed successfully"
	CommandCompletedSuccessfullyWithActionPending = "Command completed successfully; action pending"
	CommandCompleteSuccessfullyWithNoMessages     = "Command completed successfully; no messages"
	CommandCompleteSuccessfullyAckToDequeue       = "Command completed successfully; ack to dequeue"
	CommandCompleteSuccessfullyEndingSession      = "Command completed successfully; ending session"
	UnknownCommand                                = "Unknown command"
	CommandSyntaxError                            = "Command syntax error"
	CommandUseError                               = "Command use error"
	RequiredParameterMissing                      = "Required parameter missing"
	ParameterValueRangeError                      = "Parameter value range error"
	ParameterValueSyntaxError                     = "Parameter value syntax error"
	UnimplementedProtocolVersion                  = "Unimplemented protocol version"
	UnimplementedCommand                          = "Unimplemented command"
	UnimplementedOption                           = "Unimplemented option"
	UnimplementedExtension                        = "Unimplemented extension"
	BillingFailure                                = "Billing failure"
	ObjectIsNotEligibleForRenewal                 = "Object is not eligible for renewal"
	ObjectIsNotEligibleForTransfer                = "Object is not eligible for transfer"
	AuthenticationError                           = "Authentication error"
	AuthorizationError                            = "Authorization error"
	InvalidAuthorizationInformation               = "Invalid authorization information"
	ObjectPendingTransfer                         = "Object pending transfer"
	ObjectNotPendingTransfer                      = "Object not pending transfer"
	ObjectExists                                  = "Object exists"
	ObjectDoesNotExist                            = "Object does not exist"
	ObjectStatusProhibitsOperation                = "Object status prohibits operation"
	ObjectAssociationProhibitsOperation           = "Object association prohibits operation"
	ParameterValuePolicyError                     = "Parameter value policy error"
	UnimplementedObjectService                    = "Unimplemented object service"
	DataManagementPolicyViolation                 = "Data management policy violation"
	CommandFailed                                 = "Command failed"
	CommandFailedServerClosingConnection          = "Command failed; server closing connection"
	AuthenticationErrorServerClosingConnection    = "Authentication error; server closing connection"
	SessionLimitExceededServerClosingConnection   = "Session limit exceeded; server closing connection"
)

func AnyError(code int, msg string) *EPPResponse[struct{}, struct{}] {
	return NewResponse[struct{}, struct{}](code, msg)
}

//nolint:gocyclo,cyclop
func defaultMessage(code int) string {
	// see https://datatracker.ietf.org/doc/html/rfc5730#section-3
	switch code {
	case 1000:
		return CommandCompletedSuccessfully
	case 1001:
		return CommandCompletedSuccessfullyWithActionPending
	case 1300:
		return CommandCompleteSuccessfullyWithNoMessages
	case 1301:
		return CommandCompleteSuccessfullyAckToDequeue
	case 1500:
		return CommandCompleteSuccessfullyEndingSession
	case 2000:
		return UnknownCommand
	case 2001:
		return CommandSyntaxError
	case 2002:
		return CommandUseError
	case 2003:
		return RequiredParameterMissing
	case 2004:
		return ParameterValueRangeError
	case 2005:
		return ParameterValueSyntaxError
	case 2100:
		return UnimplementedProtocolVersion
	case 2101:
		return UnimplementedCommand
	case 2102:
		return UnimplementedOption
	case 2103:
		return UnimplementedExtension
	case 2104:
		return BillingFailure
	case 2105:
		return ObjectIsNotEligibleForRenewal
	case 2106:
		return ObjectIsNotEligibleForTransfer
	case 2200:
		return AuthenticationError
	case 2201:
		return AuthorizationError
	case 2202:
		return InvalidAuthorizationInformation
	case 2300:
		return ObjectPendingTransfer
	case 2301:
		return ObjectNotPendingTransfer
	case 2302:
		return ObjectExists
	case 2303:
		return ObjectDoesNotExist
	case 2304:
		return ObjectStatusProhibitsOperation
	case 2305:
		return ObjectAssociationProhibitsOperation
	case 2306:
		return ParameterValuePolicyError
	case 2307:
		return UnimplementedObjectService
	case 2308:
		return DataManagementPolicyViolation
	case 2400:
		return CommandFailed
	case 2500:
		return CommandFailedServerClosingConnection
	case 2501:
		return AuthenticationErrorServerClosingConnection
	case 2502:
		return SessionLimitExceededServerClosingConnection
	default:
		return "Unknown code"
	}
}
