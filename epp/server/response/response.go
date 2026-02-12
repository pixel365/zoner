package response

import (
	"encoding/xml"
)

type Marshaller interface {
	Marshal() ([]byte, error)
}

type EPPResponse[RD, E any] struct {
	XMLName  xml.Name        `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Response Response[RD, E] `xml:"response"`
}

type Response[RD, E any] struct {
	ResData       *ResData[RD]   `xml:"resData,omitempty"`
	MsgQ          *MsgQ          `xml:"msgQ,omitempty"`
	TransactionID *TransactionID `xml:"trID,omitempty"`
	Extensions    *Extension[E]  `xml:"extension,omitempty"`
	Results       []Result       `xml:"result"`
}

func NewResponse[R, E any](code int, msg string) *EPPResponse[R, E] {
	result := Result{Code: code, Message: msg}
	result.CheckMessage()

	return &EPPResponse[R, E]{
		Response: Response[R, E]{
			Results: []Result{result},
		},
	}
}

func (r *EPPResponse[RD, E]) Marshal() ([]byte, error) {
	body, err := xml.Marshal(r)
	if err != nil {
		return nil, err
	}

	out := make([]byte, 0, len(XmlHeader)+len(body))
	out = append(out, XmlHeader...)
	out = append(out, body...)

	return out, nil
}

func (r *EPPResponse[RD, E]) WithMsgQ(id string, count int) *EPPResponse[RD, E] {
	r.Response.MsgQ = &MsgQ{ID: id, Count: count}
	return r
}

func (r *EPPResponse[RD, E]) WithTransactionID(clientId, serverId string) *EPPResponse[RD, E] {
	r.Response.TransactionID = &TransactionID{clientId, serverId}
	return r
}

func (r *EPPResponse[RD, E]) WithResData(data RD) *EPPResponse[RD, E] {
	r.Response.ResData = &ResData[RD]{&data}
	return r
}

func (r *EPPResponse[RD, E]) SetExtensions(ext ...E) *EPPResponse[RD, E] {
	r.Response.Extensions = &Extension[E]{Items: append([]E(nil), ext...)}
	return r
}

func (r *EPPResponse[RD, E]) SetResults(results ...Result) *EPPResponse[RD, E] {
	r.Response.Results = make([]Result, 0, len(results))

	for _, result := range results {
		result.CheckMessage()
		r.Response.Results = append(r.Response.Results, result)
	}

	return r
}

func (r *EPPResponse[RD, E]) AppendResults(result ...Result) *EPPResponse[RD, E] {
	if r.Response.Results == nil {
		r.Response.Results = make([]Result, 0, len(result))
	}

	for _, res := range result {
		res.CheckMessage()
		r.Response.Results = append(r.Response.Results, res)
	}

	return r
}

func (r *EPPResponse[RD, E]) AppendExtensions(ext ...E) *EPPResponse[RD, E] {
	if r.Response.Extensions == nil {
		r.Response.Extensions = &Extension[E]{Items: make([]E, 0, len(ext))}
	}

	r.Response.Extensions.Items = append(r.Response.Extensions.Items, ext...)

	return r
}
