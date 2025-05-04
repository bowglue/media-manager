package services

func (s *UserService) GetUserById(id string) map[string]string {
	// Simulate fetching user data from a database or service
	userData := map[string]string{
		"id":   id,
		"name": "John Doe",
	}

	return userData

}
