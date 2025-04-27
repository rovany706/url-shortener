package service

import (
	"testing"

	"github.com/rovany706/url-shortener/internal/models"
	"github.com/stretchr/testify/assert"
)

func BenchmarkDeleteBuffer(b *testing.B) {
	count := 10_000
	b.Run("Add", func(b *testing.B) {
		dBuf := NewDeleteBuffer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < count; i++ {
				dBuf.Add(make(chan models.UserDeleteRequest))
			}
		}
	})
	b.Run("AddAndFlush", func(b *testing.B) {
		dBuf := NewDeleteBuffer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < count; i++ {
				dBuf.Add(make(chan models.UserDeleteRequest))
			}

			dBuf.Flush()
		}
	})
}

func TestAddAndFlush(t *testing.T) {
	wantCount := 10

	want := make([]chan models.UserDeleteRequest, wantCount)
	for i := 0; i < wantCount; i++ {
		want[i] = make(chan models.UserDeleteRequest)
	}

	dBuf := NewDeleteBuffer()
	for _, ch := range want {
		dBuf.Add(ch)
	}

	actual := dBuf.Flush()
	assert.Equal(t, wantCount, len(actual))
	assert.Equal(t, want, actual)
}
