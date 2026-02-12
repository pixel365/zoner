package response

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnyError(t *testing.T) {
	resp := AnyError(2002, CommandUseError)

	data, err := resp.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var parsed EPPResponse[struct{}, struct{}]
	require.NoError(t, xml.Unmarshal(data, &parsed))
	require.Len(t, parsed.Response.Results, 1)

	assert.Equal(t, 2002, parsed.Response.Results[0].Code)
	assert.Equal(t, CommandUseError, parsed.Response.Results[0].Message)
	assert.Nil(t, parsed.Response.ResData)
	assert.Nil(t, parsed.Response.Extensions)
}

func TestDefaultMessage(t *testing.T) {
	tests := []struct {
		expected string
		code     int
	}{
		{CommandCompletedSuccessfully, 1000},
		{CommandCompletedSuccessfullyWithActionPending, 1001},
		{CommandCompleteSuccessfullyWithNoMessages, 1300},
		{CommandCompleteSuccessfullyAckToDequeue, 1301},
		{CommandCompleteSuccessfullyEndingSession, 1500},
		{UnknownCommand, 2000},
		{CommandSyntaxError, 2001},
		{CommandUseError, 2002},
		{RequiredParameterMissing, 2003},
		{ParameterValueRangeError, 2004},
		{ParameterValueSyntaxError, 2005},
		{UnimplementedProtocolVersion, 2100},
		{UnimplementedCommand, 2101},
		{UnimplementedOption, 2102},
		{UnimplementedExtension, 2103},
		{BillingFailure, 2104},
		{ObjectIsNotEligibleForRenewal, 2105},
		{ObjectIsNotEligibleForTransfer, 2106},
		{AuthenticationError, 2200},
		{AuthorizationError, 2201},
		{InvalidAuthorizationInformation, 2202},
		{ObjectPendingTransfer, 2300},
		{ObjectNotPendingTransfer, 2301},
		{ObjectExists, 2302},
		{ObjectDoesNotExist, 2303},
		{ObjectStatusProhibitsOperation, 2304},
		{ObjectAssociationProhibitsOperation, 2305},
		{ParameterValuePolicyError, 2306},
		{UnimplementedObjectService, 2307},
		{DataManagementPolicyViolation, 2308},
		{CommandFailed, 2400},
		{CommandFailedServerClosingConnection, 2500},
		{AuthenticationErrorServerClosingConnection, 2501},
		{SessionLimitExceededServerClosingConnection, 2502},
		{"Unknown code", 9999},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("code_%d", tt.code),
			func(t *testing.T) {
				assert.Equal(t, tt.expected, defaultMessage(tt.code))
			},
		)
	}
}
