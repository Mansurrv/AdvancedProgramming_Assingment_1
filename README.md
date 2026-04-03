# AP2 Assignment 1 - Medical Scheduling Platform

Two Go microservices built with Clean Architecture and REST:

- Doctor Service: manages doctor profile data.
- Appointment Service: manages appointments and validates doctors via REST.

## Architecture

Clean Architecture is applied inside each service:

- `internal/model` holds domain models (no HTTP/JSON/framework types).
- `internal/usecase` contains business rules and interfaces.
- `internal/repository` contains persistence and outbound clients.
- `internal/transport/http` contains thin HTTP handlers.
- `internal/app` wires dependencies together.

Dependency flow:

`transport/http -> usecase -> repository (interfaces) -> model`

Microservice boundaries:

- Each service owns its own data and repository implementation.
- The Appointment Service calls the Doctor Service over REST only.
- No shared database or cross-service data access.

## Inter-service Communication

The Appointment Service checks doctor existence by calling:

- `GET /doctors/:id`
- `200 OK` means the doctor exists.
- `404 Not Found` means the doctor does not exist.
- Any network error or `5xx` results in a failure response to the client.

## Diagram

```mermaid
flowchart LR
  subgraph Doctor Service
    DHTTP[HTTP API]
    DRepo[(Doctor Data (owned))]
    DHTTP --> DRepo
  end
  subgraph Appointment Service
    AHTTP[HTTP API]
    ARepo[(Appointment Data (owned))]
    AHTTP --> ARepo
  end
  AHTTP -- "REST: GET /doctors/:id" --> DHTTP
```

## How To Run

Requirements:

- Go 1.25+

Start Doctor Service:

```bash
cd doctor-service
go run .
```

Start Appointment Service (in another terminal):

```bash
cd appointment-service
go run .
```

Default ports:

- Doctor Service: `http://localhost:8080`
- Appointment Service: `http://localhost:8081`

## API Examples

### Doctor Service

Create a doctor:

```bash
curl -X POST http://localhost:8080/doctors \
  -H 'Content-Type: application/json' \
  -d '{"full_name":"Dr. Aisha Seitkali","specialization":"Cardiology","email":"a.seitkali@clinic.kz"}'
```

Get a doctor by ID:

```bash
curl http://localhost:8080/doctors/<doctor-id>
```

List doctors:

```bash
curl http://localhost:8080/doctors
```

### Appointment Service

Create an appointment:

```bash
curl -X POST http://localhost:8081/appointments \
  -H 'Content-Type: application/json' \
  -d '{"title":"Initial cardiac consultation","description":"Patient referred for palpitations","doctor_id":"<doctor-id>"}'
```

Get an appointment by ID:

```bash
curl http://localhost:8081/appointments/<appointment-id>
```

List appointments:

```bash
curl http://localhost:8081/appointments
```

Update appointment status:

```bash
curl -X PATCH http://localhost:8081/appointments/<appointment-id>/status \
  -H 'Content-Type: application/json' \
  -d '{"status":"in_progress"}'
```

Status values: `new`, `in_progress`, `done`

## Business Rules

Doctor Service:

- `full_name` is required.
- `email` is required.
- `email` must be unique.

Appointment Service:

- `title` is required.
- `doctor_id` is required.
- Doctor must exist in the Doctor Service.
- `status` must be one of: `new`, `in_progress`, `done`.
- Status cannot transition from `done` back to `new`.

## Failure Handling

- The Appointment Service validates doctor existence on create and status update.
- If the Doctor Service is unavailable, the request fails and an internal log is written by the Doctor client.
- Network failures return `503 Service Unavailable`.
- Doctor Service errors return `502 Bad Gateway`.
- The Doctor client uses a timeout and basic retry strategy; a circuit breaker could be added in the client layer for production workloads.

## Why No Shared Database

Each service owns its data to keep bounded contexts and avoid tight coupling. The Appointment Service never reads Doctor data directly; it only uses the Doctor Service API.

## Trade-offs (Microservices vs Monolith)

- Pros: clear ownership, independent evolution, explicit boundaries, easier to scale parts independently.
- Cons: network latency, failure modes across services, more operational complexity.

