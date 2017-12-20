package main

func main() {
	var app = App{}
	app.Initialize("root", "sunday", "todo")
	app.Run(":3000")
}
