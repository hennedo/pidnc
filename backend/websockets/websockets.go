package websockets

type WebsocketMessageType struct {
	Type string `json:"type"`
}

var WsHub = NewHub()
