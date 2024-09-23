package main

func main() {
	app := App{}
	app.Initialize(DBUser, DBPassword, DBName)
	app.Run("localhost:10000")
}