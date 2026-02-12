package response

type Result struct {
	Message string `xml:",chardata"`
	Code    int    `xml:"code,attr"`
}

func (r *Result) CheckMessage() {
	if r.Message == "" {
		r.Message = defaultMessage(r.Code)
	}
}

func (r *Result) IsError() bool {
	return r.Code >= 2000
}

type MsgQ struct {
	ID      string `xml:"id,attr,omitempty"`
	Message string `xml:"msg,omitempty"`
	Count   int    `xml:"count,attr,omitempty"`
}

type TransactionID struct {
	ClientID string `xml:"clTRID,omitempty"`
	ServerID string `xml:"svTRID,omitempty"`
}

type Extension[E any] struct {
	Items []E `xml:",any"`
}

type ResData[R any] struct {
	Value *R `xml:",any"`
}
