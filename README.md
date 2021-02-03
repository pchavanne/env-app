# env-app
An app to test golang env and flag packages

The app.go file is a very basic usage of [caarlos0/env](https://github.com/caarlos0/env) and the flag package.

```golang
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
```
we build it
```bash
$ go build -o app
```
and run it
```bash
$ ./app
FOO from caarlos0/env: BAR
FOO from os: 
```
As expected Foo has a default value "BAR", but there is no Environement variable. Let's set FOO
```bash
$ export FOO=BAZ
$ ./app
FOO from caarlos0/env: BAZ
FOO from os: BAZ
```
now BAZ is the default value and os.getenv returns a value.

We can also pass foo value as an argument
```bash
$ ./app -foo TOTO
FOO from caarlos0/env: TOTO
FOO from os: BAZ
```
So this is working really great!! Let's try to dockerize it.

Here is the Dockerfile

```dockerfile
FROM golang:1.15-alpine

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./app .

CMD ["./app"]
```
and the docker-compose.yaml
```yaml
version: "3.7"
services:
  app:
    build:
      context: .
```
Let's run it
```bash
$ docker-compose up
...
app_1  | FOO from caarlos0/env: BAR
app_1  | FOO from os: 
...
```
Let's add an env var to the container
```yaml
version: "3.7"
services:
  app:
    build:
      context: .
    environment:
    - FOO=BAZ
```
```bash
$ docker-compose up
...
app_1  | FOO from caarlos0/env: BAZ
app_1  | FOO from os: BAZ
...
```
now let's try a command line argument
```yaml
version: "3.7"
services:
  app:
    build:
      context: .
    command: "./app -foo TOTO"
    environment:
    - FOO=BAZ
```
```bash
$ docker-compose up
...
app_1  | FOO from caarlos0/env: TOTO
app_1  | FOO from os: BAZ
...
```
So again this is working really great.

I am gona build a Docker image, push it to Dockerhub and use it in Kubernetes.
```bash
$ docker build -t pchavanne/env-app .
$ docker push pchavanne/env-app
``` 
here is a kubenetes job manifest of this image job.yaml
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: app
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: env-test
        image: pchavanne/env-app
```
let's run it on the cluster
```bash
kubectl apply -f job.yaml 
```
and check the logs
```bash
kubectl logs app-bvxsp
FOO from caarlos0/env: BAR
FOO from os:
```
Let's add an env var
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: app
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: env-test
        image: pchavanne/env-app
        env:
        - name: FOO
          value: BAZ

```
let's run it on the cluster
```bash
kubectl apply -f job.yaml 
```
and check the logs
```bash
kubectl logs app-bvxsp
FOO from caarlos0/env: BAZ
FOO from os: BAZ
```

It is working great!!!