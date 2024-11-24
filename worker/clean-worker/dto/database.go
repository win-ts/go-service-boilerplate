package dto

// TestEntity represents the data in each row of the test table
type TestEntity struct {
	ID      string `sql:"id"`
	Message string `sql:"message"`
}

// TestModel represents the response model of the test table
type TestModel struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
