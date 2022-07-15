package client

// Header creates and adds new headers
type Header interface {
	AddNewHeader()
	AddDefaultHeaders()
}

// HeaderStruct holds key value pairs of headers
type HeaderStruct struct {
	Header []KeyValue
}

// KeyValue pairs struct
type KeyValue struct {
	Key   string
	Value string
}

// NewHeader is a constructor for creading a new Header
func NewHeader() *HeaderStruct {
	return &HeaderStruct{}
}

func newKeyValue() *KeyValue {
	return &KeyValue{}
}

// AddNewHeader adds a new key value pair to the array of headers
func (h *HeaderStruct) AddNewHeader(key, value string) *HeaderStruct {
	newKeyValue := newKeyValue().addNewHeader(key, value)
	h.Header = append(h.Header, *newKeyValue)
	return h
}

func (k *KeyValue) addNewHeader(key, value string) *KeyValue {
	newKeyValue := newKeyValue()
	newKeyValue.Key = key
	newKeyValue.Value = value
	return newKeyValue
}

// AddDefaultHeaders adds default headers
func (h *HeaderStruct) AddDefaultHeaders() {
	h.AddNewHeader("Content-Type", "application/json")
	h.AddNewHeader("Access-Control-Allow-Origin", "*")
	h.AddNewHeader("User-Agent", "cs-code-review-bot")
}
