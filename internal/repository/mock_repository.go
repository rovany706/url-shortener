package repository

import "context"

type MockRepository struct {
	mockMap map[string]string
}

func (r *MockRepository) GetFullURL(ctx context.Context, shortID string) (fullURL string, ok bool) {
	v, ok := r.mockMap[shortID]

	return v, ok
}

func (r *MockRepository) SaveEntry(ctx context.Context, shortID string, fullURL string) error {
	r.mockMap[shortID] = fullURL

	return nil
}

func (r *MockRepository) Close() error {
	return nil
}

func (r *MockRepository) Ping(ctx context.Context) error {
	return ErrPingNotSupported
}

func (r *MockRepository) SaveEntries(ctx context.Context, shortIDMap map[string]string) error {
	for shortID, fullURL := range shortIDMap {
		err := r.SaveEntry(ctx, shortID, fullURL)

		if err != nil {
			return err
		}
	}

	return nil
}

func NewMockRepository(m map[string]string) *MockRepository {
	return &MockRepository{
		mockMap: m,
	}
}
