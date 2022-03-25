package main

import "antares-me/monitoring-system/internal/app"

const configPath = "configs/main"

func main() {
	app.Run(configPath)
}
