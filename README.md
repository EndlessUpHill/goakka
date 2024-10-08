# goAkka

[![Go Version](https://img.shields.io/badge/go-v1.22.5-blue)](https://golang.org) [![License](https://img.shields.io/github/license/EndlessUpHill/goAkka)](LICENSE) ![Build Status](https://github.com/EndlessUpHill/goAkka/actions/workflows/run-tests.yml/badge.svg)


## Overview

goAkka is a Golang system inspired by two powerful actor frameworks: **Elixir's OTP** and **Akka** from the JVM ecosystem. By bringing together the simplicity and fault tolerance of Elixir’s actors and the robust, distributed nature of Akka, goAkka provides a highly resilient, scalable, and concurrent architecture for building reliable distributed systems in Go.

If you’re familiar with **Elixir's Supervisors and Actors** or **Akka's Actors and PubSub**, you'll feel right at home with goAkka. Our goal is to make it easier for Go developers to build systems that can handle failure gracefully and scale efficiently without sacrificing performance.

## Features

*   **Actor Model:** At the core of goAkka is the actor model, where actors are lightweight, isolated entities that communicate via message-passing.
*   **Supervision Trees:** Supervised actors ensure fault tolerance by automatically restarting failed actors, inspired by Elixir's OTP.
*   **PubSub Messaging:** A built-in PubSub system allows for loosely coupled communication between actors, similar to Akka.
*   **Modular Architecture:** Extensible design with plug-and-play modules like Redis and NATS for communication.
*   **Crash Recovery:** Supervisor trees allow for dynamic supervision and recovery of actors in case of failures.
*   **Distributed Actors:** Support for scaling actors across multiple nodes (future feature).

## Tech Stack

*   **Golang:** Leveraging Go’s strong performance and concurrency model, goAkka is built for developers who value speed, simplicity, and reliability.
*   **Redis / NATS:** Support for PubSub communication using Redis and NATS, allowing you to scale your system seamlessly across distributed environments.
*   **Reflection:** Used to dynamically generate actor structures and communication channels, avoiding the need for repetitive code generation.

## Getting Started

### Installation

```
go get github.com/EndlessUpHill/goAkka
```

### Basic Usage

1\. Define your actors and supervisors.  
2\. Set up PubSub channels for communication between actors.  
3\. Run your goAkka system and let the supervision trees handle actor crashes and restarts automatically.

```
package main

import "github.com/EndlessUpHill/goAkka"

func main() {
    // Define actors and supervisors
    // Set up communication channels
    // Start your goAkka application
}
```

A full detailed guide is available in the [Documentation](#) (Coming soon).

### Example

Here’s an example of setting up a simple actor system with a supervisor and PubSub communication:

```
// Coming soon
```

### Running Tests

```
go test ./...
```

## Roadmap

*   \[x\] Core Actor System
*   \[ \] Redis Integration for PubSub
*   \[ \] NATS Integration for Distributed Actors
*   \[ \] WebSocket support for real-time communication
*   \[ \] Detailed documentation and examples

## Documentation

Complete documentation can be found [here](#) (Coming soon).

## How to Contribute

We welcome contributions from the community! To get started:

1.  Fork this repository.
2.  Create a new branch for your feature or bugfix.
3.  Submit a pull request.

For more detailed instructions, please read our [Contributing Guide](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.