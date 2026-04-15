package app

import (
	"log"
	"net"

	"appointment-service/internal/client"
	"appointment-service/internal/repository/memory"
	transportgrpc "appointment-service/internal/transport/grpc"
	"appointment-service/internal/usecase"
	appointmentpb "appointment-service/proto"

	doctorpb "doctor-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run() {
	repo := memory.NewAppointmentRepo()
	doctorConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to doctor-service: %v", err)
	}
	defer doctorConn.Close()

	doctorClient := client.NewDoctorGRPCClient(doctorpb.NewDoctorServiceClient(doctorConn))
	uc := usecase.NewAppointmentUsecase(repo, doctorClient)
	handler := transportgrpc.NewAppointmentHandler(uc)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen on :50052: %v", err)
	}

	server := grpc.NewServer()
	appointmentpb.RegisterAppointmentServiceServer(server, handler)

	log.Println("appointment-service gRPC server listening on :50052")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("appointment-service gRPC server failed: %v", err)
	}
}
