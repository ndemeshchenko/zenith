package environment

import "go.mongodb.org/mongo-driver/mongo"

type Environment struct {
	Name string
}

func GetAll(*mongo.Client) ([]Environment, error) {
	// TODO this is a mock
	// TODO implement fetch from DB
	return []Environment{
		{Name: "development"},
		{Name: "testing"},
		{Name: "acceptance"},
		{Name: "production"},
	}, nil
}
