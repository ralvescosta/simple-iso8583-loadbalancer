# Simple ISO8583 Load Balancer

## Overview

The **Simple ISO8583 Load Balancer** is a Golang-based project designed to distribute ISO 8583 financial messages across multiple backends. It is ideal for environments that require high availability and efficient load distribution, such as payment processing systems.

## Features

- **Load Balancing**: Distributes ISO 8583 messages across multiple servers.

## Dependencies

- [moov-io/iso8583](https://github.com/moov-io/iso8583)

## Getting Started

1. **Clone the repository**:
    ```sh
    git clone https://github.com/ralvescosta/simple-iso8583-loadbalancer.git
    cd simple-iso8583-loadbalancer
    ```

2. **Install dependencies**:
    ```sh
    go mod download
    ```

4. **Run the load balancer**:
    ```sh
    go run main.go
    ```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
