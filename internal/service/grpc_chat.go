package service

type grpcConnection struct {
	hub IHub

	// The grpc connection
	// todo

	// Buffered channel for outbound messages
	send chan []byte
}
