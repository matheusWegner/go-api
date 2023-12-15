package main

func main() {
	server := newAPIServer(":8080")
	server.run()
}
