package app

import (
	"log"
	"net"

	"doctor-service/internal/repository/memory"
	transportgrpc "doctor-service/internal/transport/grpc"
	"doctor-service/internal/usecase"
	doctorpb "doctor-service/proto"

	"google.golang.org/grpc"
)

func Run() {
	repo := memory.NewDoctorRepo()
	uc := usecase.NewDoctorUseCase(repo)
	handler := transportgrpc.NewDoctorHandler(uc)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on :50051: %v", err)
	}

	server := grpc.NewServer()
	doctorpb.RegisterDoctorServiceServer(server, handler)

	log.Println("doctor-service gRPC server listening on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("doctor-service gRPC server failed: %v", err)
	}
}
