# GoSSE (/ɡo͞os/)

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
$ docker run -p 8080:8080 gosse
```

Then visit your browser at

```
http://<DOCKER_IP>:8080/
```

### Go

```
$ go run gosse.go
```

Then visit your browser at

```
http://localhost:8080/
```

## EventSource

To access the EventSource API directly simply visit `<HOST>/events`

