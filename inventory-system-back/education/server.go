package education

import (
	"context"
	"fmt"
	"github.com/jochem11/inventory-system-back/education/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
)

type grpcServer struct {
	pb.UnimplementedEducationServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterEducationServiceServer(serv, &grpcServer{
		service: s,
	})
	return serv.Serve(lis)
}

// --- Course Methods ---

func (s *grpcServer) PostCourse(ctx context.Context, req *pb.PostCourseRequest) (*pb.PostCourseResponse, error) {
	c, err := s.service.PostCourse(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	updatedAt := timestamppb.New(c.UpdatedAt)
	createdAt := timestamppb.New(c.CreatedAt)

	return &pb.PostCourseResponse{Course: &pb.Course{
		Id:        c.ID,
		Name:      c.Name,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}}, nil
}

func (s *grpcServer) GetCourse(ctx context.Context, req *pb.GetCourseRequest) (*pb.GetCourseResponse, error) {
	c, err := s.service.GetCourse(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	updatedAt := timestamppb.New(c.UpdatedAt)
	createdAt := timestamppb.New(c.CreatedAt)

	return &pb.GetCourseResponse{Course: &pb.Course{
		Id:        c.ID,
		Name:      c.Name,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}}, nil
}

func (s *grpcServer) GetCourses(ctx context.Context, req *pb.GetCoursesRequest) (*pb.GetCoursesResponse, error) {
	res, err := s.service.GetCourses(ctx, &req.Skip, &req.Take)
	if err != nil {
		return nil, err
	}

	courses := []*pb.Course{}

	for _, c := range res {
		updatedAt := timestamppb.New(c.UpdatedAt)
		createdAt := timestamppb.New(c.CreatedAt)
		courses = append(courses, &pb.Course{
			Id:        c.ID,
			Name:      c.Name,
			UpdatedAt: updatedAt,
			CreatedAt: createdAt,
		})
	}

	return &pb.GetCoursesResponse{Courses: courses}, nil
}

func (s *grpcServer) UpdateCourse(ctx context.Context, req *pb.UpdateCourseRequest) (*pb.UpdateCourseResponse, error) {
	c, err := s.service.UpdateCourse(ctx, req.Id, req.Name)
	if err != nil {
		return nil, err
	}

	updatedAt := timestamppb.New(c.UpdatedAt)
	createdAt := timestamppb.New(c.CreatedAt)

	return &pb.UpdateCourseResponse{Course: &pb.Course{
		Id:        c.ID,
		Name:      c.Name,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}}, nil
}

func (s *grpcServer) DeleteCourse(ctx context.Context, req *pb.DeleteCourseRequest) (*pb.DeleteCourseResponse, error) {
	err := s.service.DeleteCourseByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCourseResponse{}, nil
}

// --- Class Methods ---

func (s *grpcServer) PostClass(ctx context.Context, req *pb.PostClassRequest) (*pb.PostClassResponse, error) {
	c, err := s.service.PostClass(ctx, req.Name, req.CourseId)
	if err != nil {
		return nil, err
	}

	updatedAt := timestamppb.New(c.UpdatedAt)
	createdAt := timestamppb.New(c.CreatedAt)

	courseUpdatedAt := timestamppb.New(c.Course.UpdatedAt)
	courseCreatedAt := timestamppb.New(c.Course.CreatedAt)

	return &pb.PostClassResponse{Class: &pb.Class{
		Id:       c.ID,
		Name:     c.Name,
		CourseId: c.CourseID,
		Course: &pb.Course{
			Id:        c.Course.ID,
			Name:      c.Course.Name,
			CreatedAt: courseCreatedAt,
			UpdatedAt: courseUpdatedAt,
		},
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}}, nil
}

func (s *grpcServer) GetClass(ctx context.Context, req *pb.GetClassRequest) (*pb.GetClassResponse, error) {
	c, err := s.service.GetClass(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	updatedAt := timestamppb.New(c.UpdatedAt)
	createdAt := timestamppb.New(c.CreatedAt)

	courseUpdatedAt := timestamppb.New(c.Course.UpdatedAt)
	courseCreatedAt := timestamppb.New(c.Course.CreatedAt)

	return &pb.GetClassResponse{Class: &pb.Class{
		Id:       c.ID,
		Name:     c.Name,
		CourseId: c.CourseID,
		Course: &pb.Course{
			Id:        c.Course.ID,
			Name:      c.Course.Name,
			UpdatedAt: courseUpdatedAt,
			CreatedAt: courseCreatedAt,
		},
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}}, nil
}

func (s *grpcServer) GetClasses(ctx context.Context, req *pb.GetClassesRequest) (*pb.GetClassesResponse, error) {
	res, err := s.service.GetClasses(ctx, &req.Skip, &req.Take)
	if err != nil {
		return nil, err
	}

	classes := []*pb.Class{}

	for _, c := range res {
		updatedAt := timestamppb.New(c.UpdatedAt)
		createdAt := timestamppb.New(c.CreatedAt)
		courseUpdatedAt := timestamppb.New(c.Course.UpdatedAt)
		courseCreatedAt := timestamppb.New(c.Course.CreatedAt)
		classes = append(classes, &pb.Class{
			Id:       c.ID,
			Name:     c.Name,
			CourseId: c.CourseID,
			Course: &pb.Course{
				Id:        c.Course.ID,
				Name:      c.Course.Name,
				CreatedAt: courseCreatedAt,
				UpdatedAt: courseUpdatedAt,
			},
			UpdatedAt: updatedAt,
			CreatedAt: createdAt,
		})
	}

	return &pb.GetClassesResponse{Classes: classes}, nil
}

func (s *grpcServer) UpdateClass(ctx context.Context, req *pb.UpdateClassRequest) (*pb.UpdateClassResponse, error) {
	c, err := s.service.UpdateClass(
		ctx,
		req.Id,
		req.Name,
		req.CourseId)
	if err != nil {
		return nil, err
	}

	updatedAt := timestamppb.New(c.UpdatedAt)
	createdAt := timestamppb.New(c.CreatedAt)

	courseUpdatedAt := timestamppb.New(c.Course.UpdatedAt)
	courseCreatedAt := timestamppb.New(c.Course.CreatedAt)

	return &pb.UpdateClassResponse{Class: &pb.Class{
		Id:       c.ID,
		Name:     c.Name,
		CourseId: c.CourseID,
		Course: &pb.Course{
			Id:        c.Course.ID,
			Name:      c.Course.Name,
			UpdatedAt: courseUpdatedAt,
			CreatedAt: courseCreatedAt,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}}, nil
}

func (s *grpcServer) DeleteClass(ctx context.Context, req *pb.DeleteClassRequest) (*pb.DeleteClassResponse, error) {
	err := s.service.DeleteClassByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteClassResponse{}, nil
}
