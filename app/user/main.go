package main

import "gcluster/essential/app"

func main() {
	userApp := app.GetGClusterApp()
	userApp.Name = "user"
	userApp.Version = "1.0.0"


}
