package main

import (
	"context"
	"log"
	"neon_products/internal/http"
	"neon_products/internal/store/inmemory"
)

func main() {
	store := inmemory.Init()

	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()

}
