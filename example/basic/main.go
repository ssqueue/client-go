package main

import (
	"context"
	"fmt"
	"time"

	client "github.com/ssqueue/client-go"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client.Init(client.Options{Address: "http://127.0.0.1:6560"})

	err := client.SendString(ctx, "foo", "hello world")

	fmt.Printf("err: %v\n", err)

	res, err2 := client.GetString(ctx, "bar", time.Second*5)

	fmt.Printf("res: %v, err2: %v\n", res, err2)
}
