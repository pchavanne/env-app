package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Foo string `env:"FOO" envDefault:"BAR"`
}

var cfg = config{}

func init() {
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	flag.StringVar(&cfg.Foo, "foo", cfg.Foo, "That's the Foo var!!")
}

func main() {

	flag.Parse()

	fmt.Printf("FOO from caarlos0/env: %+v\n", cfg.Foo)
	fmt.Printf("FOO from os: %+v\n", os.Getenv("FOO"))
}
