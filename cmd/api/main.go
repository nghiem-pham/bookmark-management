package main

import "github.com/nghiem-pham/bookmark-management/internal/api"

func main() {
	// create app config
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
