package response

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResponse_MarshalUnmarshal(t *testing.T) {
	resp := NewResponse[struct{}, struct{}](1000, CommandCompletedSuccessfully)

	data, err := resp.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var parsed EPPResponse[struct{}, struct{}]
	require.NoError(t, xml.Unmarshal(data, &parsed))

	require.Len(t, parsed.Response.Results, 1)
	assert.Equal(t, 1000, parsed.Response.Results[0].Code)
	assert.Equal(t, CommandCompletedSuccessfully, parsed.Response.Results[0].Message)
	assert.False(t, parsed.Response.Results[0].IsError())

	assert.Nil(t, parsed.Response.MsgQ)
	assert.Nil(t, parsed.Response.TransactionID)
	assert.Nil(t, parsed.Response.ResData)
	assert.Nil(t, parsed.Response.Extensions)
}

func TestNewResponse_EmptyMsg_UsesDefaultMessage(t *testing.T) {
	resp := NewResponse[struct{}, struct{}](2002, "")

	data, err := resp.Marshal()
	require.NoError(t, err)

	var parsed EPPResponse[struct{}, struct{}]
	require.NoError(t, xml.Unmarshal(data, &parsed))

	require.Len(t, parsed.Response.Results, 1)
	assert.Equal(t, 2002, parsed.Response.Results[0].Code)
	assert.NotEmpty(t, parsed.Response.Results[0].Message)
	assert.Equal(t, defaultMessage(2002), parsed.Response.Results[0].Message)

	for _, result := range parsed.Response.Results {
		assert.True(t, result.IsError())
	}
}

func TestWithMsgQ(t *testing.T) {
	resp := AnyError(2001, CommandSyntaxError).
		WithMsgQ("Q-1", 3)

	data, err := resp.Marshal()
	require.NoError(t, err)

	var parsed EPPResponse[struct{}, struct{}]
	require.NoError(t, xml.Unmarshal(data, &parsed))

	require.NotNil(t, parsed.Response.MsgQ)
	assert.Equal(t, "Q-1", parsed.Response.MsgQ.ID)
	assert.Equal(t, 3, parsed.Response.MsgQ.Count)
}

func TestWithTransactionID(t *testing.T) {
	resp := AnyError(1000, CommandCompletedSuccessfully).
		WithTransactionID("cl-123", "sv-456")

	data, err := resp.Marshal()
	require.NoError(t, err)

	var parsed EPPResponse[struct{}, struct{}]
	require.NoError(t, xml.Unmarshal(data, &parsed))

	require.NotNil(t, parsed.Response.TransactionID)
	assert.Equal(t, "cl-123", parsed.Response.TransactionID.ClientID)
	assert.Equal(t, "sv-456", parsed.Response.TransactionID.ServerID)
}

func TestWithResData(t *testing.T) {
	type LoginResData struct {
		XMLName xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 loginData"`
		Result  string   `xml:"result"`
	}

	resp := NewResponse[LoginResData, struct{}](1000, CommandCompletedSuccessfully).
		WithResData(LoginResData{Result: "ok"})

	data, err := resp.Marshal()
	require.NoError(t, err)

	var parsed EPPResponse[LoginResData, struct{}]
	require.NoError(t, xml.Unmarshal(data, &parsed))

	require.NotNil(t, parsed.Response.ResData)
	assert.Equal(t, "ok", parsed.Response.ResData.Value.Result)
}

func TestSetResults_CheckMessageApplied(t *testing.T) {
	resp := AnyError(1000, "ignored").
		SetResults(
			Result{Code: 2002, Message: ""},
			Result{Code: 2400, Message: "Custom error"},
		)

	data, err := resp.Marshal()
	require.NoError(t, err)

	var parsed EPPResponse[struct{}, struct{}]
	require.NoError(t, xml.Unmarshal(data, &parsed))

	require.Len(t, parsed.Response.Results, 2)
	assert.Equal(t, 2002, parsed.Response.Results[0].Code)
	assert.Equal(t, defaultMessage(2002), parsed.Response.Results[0].Message)

	assert.Equal(t, 2400, parsed.Response.Results[1].Code)
	assert.Equal(t, "Custom error", parsed.Response.Results[1].Message)

	for _, result := range parsed.Response.Results {
		assert.True(t, result.IsError())
	}
}

func TestAppendResults_CheckMessageApplied(t *testing.T) {
	resp := AnyError(2001, CommandSyntaxError).
		AppendResults(Result{Code: 2002, Message: ""})

	data, err := resp.Marshal()
	require.NoError(t, err)

	var parsed EPPResponse[struct{}, struct{}]
	require.NoError(t, xml.Unmarshal(data, &parsed))

	require.Len(t, parsed.Response.Results, 2)
	assert.Equal(t, 2001, parsed.Response.Results[0].Code)
	assert.Equal(t, CommandSyntaxError, parsed.Response.Results[0].Message)

	assert.Equal(t, 2002, parsed.Response.Results[1].Code)
	assert.Equal(t, defaultMessage(2002), parsed.Response.Results[1].Message)

	for _, result := range parsed.Response.Results {
		assert.True(t, result.IsError())
	}
}

func TestSetExtensions_ContainerShape(t *testing.T) {
	type FeeChkData struct {
		XMLName xml.Name `xml:"urn:ietf:params:xml:ns:fee-0.7 chkData"`
		Name    string   `xml:"name"`
	}

	resp := NewResponse[struct{}, FeeChkData](1000, CommandCompletedSuccessfully).
		SetExtensions(FeeChkData{Name: "example.com"})

	data, err := resp.Marshal()
	require.NoError(t, err)

	s := string(data)
	assert.Contains(t, s, "<extension>")
	assert.Contains(t, s, "</extension>")
	assert.NotContains(t, s, "</extension><extension>")

	var parsed EPPResponse[struct{}, FeeChkData]
	require.NoError(t, xml.Unmarshal(data, &parsed))

	require.NotNil(t, parsed.Response.Extensions)
	require.Len(t, parsed.Response.Extensions.Items, 1)
	assert.Equal(t, "example.com", parsed.Response.Extensions.Items[0].Name)
}

func TestAppendExtensions_Appends(t *testing.T) {
	type Ext1 struct {
		XMLName xml.Name `xml:"urn:example:ext ext1"`
		V       string   `xml:"v"`
	}
	type Ext2 struct {
		XMLName xml.Name `xml:"urn:example:ext ext2"`
		V       string   `xml:"v"`
	}

	type AnyExt struct {
		Inner   any `xml:",any"`
		XMLName xml.Name
	}

	resp := NewResponse[struct{}, AnyExt](1000, CommandCompletedSuccessfully).
		AppendExtensions(
			AnyExt{Inner: Ext1{V: "a"}},
			AnyExt{Inner: Ext2{V: "b"}},
		)

	data, err := resp.Marshal()
	require.NoError(t, err)

	var parsed EPPResponse[struct{}, AnyExt]
	require.NoError(t, xml.Unmarshal(data, &parsed))

	require.NotNil(t, parsed.Response.Extensions)
	require.Len(t, parsed.Response.Extensions.Items, 2)
}

func TestMarshal_IncludesXMLHeader(t *testing.T) {
	resp := AnyError(2002, CommandUseError)

	data, err := resp.Marshal()
	require.NoError(t, err)

	assert.Contains(t, string(data), `<?xml version="1.0" encoding="UTF-8" standalone="no"?>`)
}
