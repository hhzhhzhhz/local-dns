package backends

type RangeResponse struct {
	Kvs []*KeyValue `json:"kvs,omitempty"`
	Count int64 `json:"count,omitempty"`
}

type KeyValue struct {
	// key is the key in bytes. An empty key is not allowed.
	Key []byte `json:"key,omitempty"`
	// value is the value held by the key, in bytes.
	Value []byte `json:"value,omitempty"`
}