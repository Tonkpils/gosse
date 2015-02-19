# GooSSE

Go Server Sent Events Example

## Getting Started

Clone the repository

```
$ git clone git@github.com:Tonkpils/gosse.git
```

### Docker

To run with Docker simply use:

```
$ docker build -t "gosse"
$ docker run -p 3000:3000 gosse
```

Then visit your browser at

```
http://<DOCKER_IP>:3000/
```

### Go

```
$ go run gosse.go
```

Then visit your browser at

```
http://localhost:3000/
```

## EventSource

To access the EventSource API directly simply visit `<HOST>/events`

