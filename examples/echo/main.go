package main

func main() {
	c := NewClient("张三", "localhost:3001")
	c.Start()
}
