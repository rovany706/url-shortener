package repository

type MockRepository struct {
	mockMap map[string]string
}

func (r *MockRepository) GetFullURL(shortID string) (fullURL string, ok bool) {
	v, ok := r.mockMap[shortID]
	return v, ok
}

func (r *MockRepository) SaveEntry(shortID string, fullURL string) error {
	r.mockMap[shortID] = fullURL

	return nil
}

func NewMockRepository(m map[string]string) *MockRepository {
	return &MockRepository{
		mockMap: m,
	}
}
