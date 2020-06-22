package main

// register all route here
func (s *Server) route() {
	s.router.HandleFunc("/new", s.authMiddleware(s.postHandler())).Methods("POST")
	s.router.HandleFunc("/image/uploads/{fileName}", s.getHandler()).Methods("Get")
}
