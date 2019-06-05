package wabin

type Sessioner interface {
	Init(*Conn)
	Close()
}

type Server interface {
	NewSession() Sessioner
}
