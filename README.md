### Load Balancer Go Program

This Go program implements a simple load balancer that distributes incoming HTTP requests among a set of backend servers.

#### Prerequisites
- Go installed on your system.

#### Usage
1. Clone this repository.
2. Navigate to the directory containing the main.go file.
3. Run the program using the `go run main.go` command.
4. Access the load balancer at `localhost:8000` in your web browser.

#### Overview
- `main.go`: Contains the main program logic.
- `Server`: Interface defining methods for server properties and behavior.
- `SimpleServer`: Struct representing a backend server with its address and a reverse proxy.
- `LoadBalancer`: Struct managing a set of backend servers and distributing requests.
- `newSimpleServer`: Creates a new instance of `SimpleServer`.
- `newLoadBalancer`: Creates a new instance of `LoadBalancer`.
- `getNextAvailableServer`: Selects the next available server using round-robin scheduling.
- `serveProxy`: Forwards incoming requests to the selected server.
- `main`: Initializes the load balancer with a set of backend servers and starts serving HTTP requests.

#### Example
```bash
$ go run main.go
serving requests at 'localhost:8000'
```
#### Contributing
Feel free to contribute by reporting issues, suggesting enhancements, or submitting pull requests.