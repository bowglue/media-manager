package handler

import (
	"context"
	"gateway/http/handlers/baseHandler"
	"gateway/http/registry"
	"gateway/types"
	"log"
	"net/http"
	pb "shared/proto/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserHandler represents a user-related handler
type UserHandler struct {
	*baseHandler.HttpHandler
	userClient pb.UserServiceClient
}

func init() {
	conn, err := grpc.NewClient("user-service:12000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		return
	}

	client := pb.NewUserServiceClient(conn)
	registry.Register(NewUserHandler(client))

}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(client pb.UserServiceClient) *UserHandler {
	return &UserHandler{
		HttpHandler: baseHandler.NewHttpHandler(),
		userClient:  client,
	}
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /user/{id}", h.GetUsers)
}

// GetUsers is an endpoint that retrieves users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	resp, err := h.userClient.GetUser(context.Background(), &pb.GetUserRequest{
		UserId: userID,
	})

	if err != nil {
		response := types.Response{
			Status: http.StatusInternalServerError,
			Error:  "Failed to fetch user: " + err.Error(),
		}
		h.SendResponse(w, response)
		return
	}

	response := types.Response{
		Status: http.StatusOK,
		Data: map[string]string{
			"user_id":  resp.UserId,
			"username": resp.Username,
		},
	}
	h.SendResponse(w, response)

	// users := []string{"Alice", "Bob", "Charlie"}
	// response := types.Response{
	// 	Status: http.StatusOK,
	// 	Data:   users,
	// }
	// h.SendResponse(w, response)
}
