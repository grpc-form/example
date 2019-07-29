package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"

	"example/backend/api"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	database := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	connStr := fmt.Sprintf("host=database port=5432 user=%s dbname=%s password=%s sslmode=disable", user, database, password)
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	api.RegisterDatabaseServer(s, &server{db: db})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

type server struct {
	db *sql.DB
}

func (s *server) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.User, error) {
	var username string
	var age int
	var car string
	row := s.db.QueryRow("SELECT username, age, car FROM usercars WHERE username=$1", req.GetName())
	err := row.Scan(&username, &age, &car)
	if err != nil {
		return &api.User{}, nil
	}
	return &api.User{Name: username, Age: int64(age), Car: car}, nil
}

func (s *server) InsertUser(ctx context.Context, user *api.User) (*api.User, error) {
	u, _ := s.GetUser(ctx, &api.GetUserRequest{Name: user.GetName()})
	if u.GetName() == user.GetName() {
		return &api.User{}, nil
	}
	_, err := s.db.Exec("INSERT INTO usercars (username, age, car) VALUES ($1, $2, $3)",
		user.GetName(), int(user.GetAge()), user.GetCar())
	if err != nil {
		return &api.User{}, nil
	}
	return user, nil
}
