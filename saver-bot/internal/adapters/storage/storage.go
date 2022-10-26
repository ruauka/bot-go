package storage

import (
	"context"

	"saver-bot/internal/domain/entities"
)

type Storage interface {
	SaveManic(ctx context.Context, m *entities.Manic) error
	SaveMassage(ctx context.Context, m *entities.Massage) error
	GetAllEvents(ctx context.Context) ([]entities.Manic, []entities.Massage, error)
	//PickRandom(ctx context.Context, userName string) (*Page, error)
	//Remove(ctx context.Context, p *Page) error
	//IsExists(ctx context.Context, p *Page) (bool, error)
}
