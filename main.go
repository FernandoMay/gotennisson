package main

func main() {
	router := newRouter()
	router.Run(":8080")
}
