# Appointment Service

Simple appointment service that stores appointments in memory and validates doctors via the Doctor Service.

## Requirements

- Go 1.25+

## Run Locally

Start Doctor Service first (required for doctor validation):

```bash
cd doctor-service
go run .
```

Then start Appointment Service:

```bash
cd appointment-service
go run .
```

Default ports:

- Doctor Service: `http://localhost:8080`
- Appointment Service: `http://localhost:8081`

## API

### Create Appointment

`POST /appointments`

Request body:

```json
{
  "title": "Checkup",
  "description": "Annual visit",
  "doctor_id": "<doctor-id>"
}
```

Responses:

- `201 Created` with appointment JSON
- `400 Bad Request` validation error or doctor does not exist
- `502 Bad Gateway` doctor service error
- `503 Service Unavailable` doctor service unavailable

### Get Appointment

`GET /appointments/:id`

Responses:

- `200 OK` with appointment JSON
- `404 Not Found` appointment not found

### List Appointments

`GET /appointments`

Responses:

- `200 OK` with list of appointments

### Update Status

`PATCH /appointments/:id/status`

Request body:

```json
{
  "status": "in_progress"
}
```

Allowed values: `new`, `in_progress`, `done`

Responses:

- `200 OK`
- `400 Bad Request` invalid status or invalid transition
- `404 Not Found` appointment not found

## Notes

- Storage is in memory only. Restarting the service clears all data.
- Doctor validation is done via REST call to the Doctor Service with retries on network errors and 429/5xx.
