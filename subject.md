# Backend Exercise

Build a backend REST API for managing a fleet of cars and the garages that house them. No frontend is expected

- **Tokens** are used to authenticate requests — endpoints require a valid token in the request header.
- **Cars** are identified by their numberplate and hold basic information (model, color, serial number).
- **Garages** have a capacity limit and reference cars by their numberplate. A car referenced by a garage cannot be deleted until it is removed from that garage.

A full database implementation is not required. The persistence strategy is up to you.


Once finished, you will present your work and answer questions about your code, the choices you made, and what could be improved.

## Technical constraints

The language must be one of: **Java**, **Scala**, **Rust**, **Go**, or **JavaScript**.
The choice of framework is free.

## Documentation

Include a file explaining how to build and run your application.

## Specification

The full API contract is defined in [`openapi.yaml`](./openapi.yaml). It contains all the details on endpoints, request/response schemas, authentication.