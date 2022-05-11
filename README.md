# Zero Hash

Zero Hash - Code Challenge


## How to Run

To run this application, you can run locally:
- `make run` or `source .env.local && go run cmd/subscriber/main.go`

or using kubernetes:

- `kubectl apply -f kube/base/deployment.yml` -> you may need to have `minikube` or `k3s` for this to work

### Tests

- You can run the tests with `make test`

## Design Decisions

- The chosen architecture was the *Clean Architecture*, separating the responsibilities into different layers, and
facilitating testing and decoupling, to achieve that I separated the code in these layers:
  - entity: This package represents objects/structs that should be used by the service layer, decoupling it from external formats.
  - service: This package defines the use cases for the application, where all the business logic resides.
  - transport: This package is responsible for the outside communication, dealing with the websocket connection, and converting it to a service-readable
  format (entitiy), and passing it through the service.
![clean architecture layers](https://www.virality.de/wp-content/uploads/2018/05/Clean-Architecture.png "clean architecture")

- The only external library that I used is the gorilla websocket, as described in the code challenge description, apart from that I used the standard library for testing and everything else basically.
- This solution enables easy swapping on the transport layer, so we can receive data from many different sources, even events if needed.
- It was added an extensible solution to print the output of the VWAP in different formats, like a File, Stdout or a buffer (I use this on my tests for example)
- The code was designed focusing on testing, decoupling and efficiency.
