package usecase

// DoctorClient defines the methods AppointmentService expects from DoctorService
type DoctorClient interface {
	DoctorExists(doctorID string) (bool, error)
}
