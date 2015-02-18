# GooSSE

Go Server Sent Events Example

## Getting Started

Clone the repository

```
$ git clone git@github.com:Tonkpils/goosse.git
```

### Docker

To run with Docker simply use:

```
$ docker build -t "goosse"
$ docker run -p 3000:3000 goosse
```

Then visit your browser at

```
http://<DOCKER_IP>:3000/
```

### Go

```
$ go run goosse.go
```

Then visit your browser at

```
http://localhost:3000/
```

## EventSource

To access the EventSource API directly simply visit `<HOST>/events`

