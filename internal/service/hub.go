package service

// Hub maintains the set of active connections and broadcasts messages to the
// connections
type Hub struct {
	// Registered connections
	chats map[string]map[Chat]bool

	// Inbound messages from the connections.
	broadcast chan Message

	// Register requests
	register chan Chat

	// Unregister requests
	unregister chan Chat
}

type IHub interface {
	Run() error
	Unregister(c Chat)
	Register(c Chat)
	Broadcast(m Message)
}

func NewHub() IHub {
	return Hub{
		broadcast:  make(chan Message),
		register:   make(chan Chat),
		unregister: make(chan Chat),
		chats:      make(map[string]map[Chat]bool),
	}
}

func (h Hub) Run() error {
	for {
		select {
		case chat := <-h.register:
			connections := h.chats[chat.GetID()]
			if connections == nil {
				connections = make(map[Chat]bool)
				h.chats[chat.GetID()] = connections
			}
			h.chats[chat.GetID()][chat] = true
		case chat := <-h.unregister:
			connections := h.chats[chat.GetID()]
			if connections != nil {
				if _, ok := connections[chat]; ok {
					delete(connections, chat)
					close(chat.GetSendChan())
					if len(connections) == 0 {
						delete(h.chats, chat.GetID())
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.chats[m.ChatID]
			for c := range connections {
				select {
				case c.GetSendChan() <- m.Data:
				default:
					close(c.GetSendChan())
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.chats, m.ChatID)
					}
				}
			}
		}
	}
}

func (h Hub) Unregister(c Chat) {
	h.unregister <- c
}

func (h Hub) Register(c Chat) {
	h.register <- c
}

func (h Hub) Broadcast(m Message) {
	h.broadcast <- m
}
