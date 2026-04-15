# grpcurl Commands (Assignment 2 Testing Artifact)

## Prerequisites

- Doctor Service running on `localhost:50051`
- Appointment Service running on `localhost:50052`
- `grpcurl` installed
- Run these commands from the repository root, or adjust the paths to match your current directory.

## 1) Doctor Service RPCs

List Doctor service methods:

```bash
grpcurl -plaintext -import-path . -proto doctor-service/proto/doctor.proto localhost:50051 list
grpcurl -plaintext -import-path . -proto doctor-service/proto/doctor.proto localhost:50051 describe doctor.DoctorService
```

Create doctor (success):

```bash
grpcurl -plaintext -import-path doctor-service/proto -proto doctor.proto \
  -d '{"full_name":"Dr. Sarah Connor","specialization":"Cardiology","email":"sarah.connor@clinic.com"}' \
  localhost:50051 doctor.DoctorService/CreateDoctor
```

Create doctor (duplicate email -> `AlreadyExists`):

```bash
grpcurl -plaintext -import-path doctor-service/proto -proto doctor-service/proto/doctor.proto \
  -d '{"full_name":"Dr. Other","specialization":"Neurology","email":"sarah.connor@clinic.com"}' \
  localhost:50051 doctor.DoctorService/CreateDoctor
```

Get doctor (replace `<doctor_id>`):

```bash
grpcurl -plaintext -import-path doctor-service/proto -proto doctor-service/proto/doctor.proto \
  -d '{"id":"<doctor_id>"}' \
  localhost:50051 doctor.DoctorService/GetDoctor
```

Get doctor (not found -> `NotFound`):

```bash
grpcurl -plaintext -import-path doctor-service/proto -proto doctor-service/proto/doctor.proto \
  -d '{"id":"missing-id"}' \
  localhost:50051 doctor.DoctorService/GetDoctor
```

List doctors:

```bash
grpcurl -plaintext -import-path doctor-service/proto -proto doctor-service/proto/doctor.proto \
  -d '{}' \
  localhost:50051 doctor.DoctorService/ListDoctors
```

## 2) Appointment Service RPCs

List Appointment service methods:

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment.proto localhost:50052 list
grpcurl -plaintext -import-path appointment-service/proto -proto appointment.proto localhost:50052 describe appointment.AppointmentService
```

Create appointment (success; use existing `<doctor_id>`):

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment.proto \
  -d '{"title":"Initial consultation","description":"Discuss treatment plan","doctor_id":"<doctor_id>"}' \
  localhost:50052 appointment.AppointmentService/CreateAppointment
```

Create appointment (missing title -> `InvalidArgument`):

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment-service/proto/appointment.proto \
  -d '{"title":"","description":"Missing title","doctor_id":"<doctor_id>"}' \
  localhost:50052 appointment.AppointmentService/CreateAppointment
```

Create appointment (doctor missing remotely -> `FailedPrecondition`):

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment-service/proto/appointment.proto \
  -d '{"title":"Checkup","description":"Remote doctor check","doctor_id":"missing-doctor-id"}' \
  localhost:50052 appointment.AppointmentService/CreateAppointment
```

Get appointment (replace `<appointment_id>`):

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment-service/proto/appointment.proto \
  -d '{"id":"<appointment_id>"}' \
  localhost:50052 appointment.AppointmentService/GetAppointment
```

Get appointment (not found -> `NotFound`):

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment-service/proto/appointment.proto \
  -d '{"id":"missing-appointment-id"}' \
  localhost:50052 appointment.AppointmentService/GetAppointment
```

List appointments:

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment-service/proto/appointment.proto \
  -d '{}' \
  localhost:50052 appointment.AppointmentService/ListAppointments
```

Update status to in_progress:

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment-service/proto/appointment.proto \
  -d '{"id":"<appointment_id>","status":"in_progress"}' \
  localhost:50052 appointment.AppointmentService/UpdateAppointmentStatus
```

Update status to done:

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment-service/proto/appointment.proto \
  -d '{"id":"<appointment_id>","status":"done"}' \
  localhost:50052 appointment.AppointmentService/UpdateAppointmentStatus
```

Forbidden transition done -> new (`InvalidArgument`):

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment-service/proto/appointment.proto \
  -d '{"id":"<appointment_id>","status":"new"}' \
  localhost:50052 appointment.AppointmentService/UpdateAppointmentStatus
```

## 3) Failure Scenario: Doctor Service Unavailable

1. Stop Doctor Service (`Ctrl+C` in its terminal).
2. Run:

```bash
grpcurl -plaintext -import-path appointment-service/proto -proto appointment-service/proto/appointment.proto \
  -d '{"title":"Dependency down case","description":"doctor service unavailable","doctor_id":"any-id"}' \
  localhost:50052 appointment.AppointmentService/CreateAppointment
```

Expected result: `code = Unavailable` with descriptive dependency error message.
