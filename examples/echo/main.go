package main

func main() {
	c := NewClient("张三", "localhost:2001")
	c.Start()
}
