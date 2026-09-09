# Fleet Management API

A REST API written in Go for managing vehicles and garages, secured with token-based
authentication.

## Requirements

- Go 1.24 or later
- MariaDB (or MySQL) 10.6+

## Configuration

The application is configured entirely through environment variables.

| Variable        | Description                                                            | Default | Required |
| --------------- | ---------------------------------------------------------------------- | ------- | -------- |
| `PORT`          | Port the HTTP server listens on                                        | `8080`  | No       |
| `TOKEN_DURATION`| Token lifetime, in minutes                                             | `10`    | No       |
| `DB_URL`        | Database address and schema, in Go MySQL driver DSN form: `tcp(host:port)/database` | —       | Yes      |
| `DB_USERNAME`   | User with access to the database                                       | —       | Yes      |
| `DB_PASSWORD`   | Password for that user                                                 | —       | Yes      |

Example:

```sh
export PORT=8080
export TOKEN_DURATION=10
export DB_URL="tcp(127.0.0.1:3306)/fleet"
export DB_USERNAME="fleet"
export DB_PASSWORD="password"
```

## Database setup

Create the schema and a dedicated user:

```sql
CREATE DATABASE fleet CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'fleet'@'localhost' IDENTIFIED BY 'changeme';
GRANT ALL PRIVILEGES ON fleet.* TO 'fleet'@'localhost';
FLUSH PRIVILEGES;
```
```
create table car
(
    id          int auto_increment
        primary key,
    numberplate varchar(20)                           not null,
    model       varchar(20)                           not null,
    color       varchar(20)                           null,
    serial      varchar(255)                          not null,
    created_at  timestamp default current_timestamp() not null,
    constraint car_pk_1
        unique (numberplate)
);

create table garage
(
    id               int auto_increment
        primary key,
    name             varchar(255)                          not null,
    occupation_limit int                                   not null,
    created_at       timestamp default current_timestamp() null,
    constraint garage_pk_2
        unique (name)
);

create table garage_car
(
    garage_id  int                                   not null,
    car_id     int                                   not null,
    created_at timestamp default current_timestamp() null,
    primary key (garage_id, car_id),
    constraint garage_car_car_id_fk
        foreign key (car_id) references car (id)
            on update cascade,
    constraint garage_car_garage_id_fk
        foreign key (garage_id) references garage (id)
            on update cascade
);

create table token
(
    hash       binary(32)                            not null
        primary key,
    expires_at timestamp                             not null,
    created_at timestamp default current_timestamp() not null
);

```

## Build and run

### Build a binary

```sh
go build
./evertrust-backend-exercise
```

`go build` names the executable after the module declared in `go.mod`. On Windows it
produces `evertrust-backend-exercise.exe` instead.

To choose the name yourself:

```sh
go build -o server .
./server
```

The server logs its listening address on startup:

```
Server is running on port 8080
```

## Authentication

Every vehicle and garage route is protected by a middleware that requires a valid,
unexpired token. Requests without one are rejected before reaching the handler.

Tokens are generated with `crypto/rand` and are **only returned once**, at creation
time. Only their hash is persisted, so a leaked database does not expose usable
credentials.

## Endpoints

### Tokens (public)

| Method   | Path              | Description                  | Status |
| -------- |-------------------| ---------------------------- | ------ |
| `POST`   | `/tokens`         | Generate a new token         | Done   |
| `GET`    | `/tokens/{token}` | Show a token's details       | Done   |
| `DELETE` | `/tokens/{token}` | Delete a token               | Done   |

### Vehicles (authenticated)

| Method   | Path                      | Description                    | Status |
| -------- |---------------------------| ------------------------------ | ------ |
| `POST`   | `/vehicles`               | Create a vehicle               | Done   |
| `GET`    | `/vehicles`               | List all available vehicles    | Done   |
| `GET`    | `/vehicles/{numberplate}` | Show a single vehicle          | Done   |
| `PUT`    | `/vehicles/{numberplate}`          | Update a vehicle               | Done   |
| `DELETE` | `/vehicles/{numberplate}`          | Delete a vehicle               | Done   |

### Garages (authenticated)

| Method   | Path              | Description             | Status      |
| -------- |-------------------| ----------------------- | ----------- |
| `POST`   | `/garages`        | Create a garage         | Partial     |
| `GET`    | `/garages`        | List all garages        | Not started |
| `GET`    | `/garages/{name}` | Show a single garage    | Not started |
| `PUT`    | `/garages/{name}` | Update a garage         | Not started |
| `DELETE` | `/garages/{name}` | Delete a garage         | Not started |


## Implementation status

This submission was written under the time limit (3 hours) and is not feature-complete. The
gaps below are known and deliberate, not oversights.

**Garage endpoints.** Creation is partially implemented; read, update and delete are
not started. 

**RFC 7807 error responses.** Errors currently return an appropriate HTTP status code
with no structured body — `application/problem+json` is not implemented. This was the
next item on my list. The intended shape:

```json
{
  "type": "https://api.example.com/problems/unauthorized",
  "title": "Authentication failed",
  "status": 401,
  "detail": "The provided credentials are invalid.",
  "instance": "urn:uuid:..."
}
```