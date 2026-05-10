package server

type Server struct {
	port string
}

func NewServer(port string) *Server {
	if port == "" {
		port = "8080"
	}
	return &Server{port: port}

}
