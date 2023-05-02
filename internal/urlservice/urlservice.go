package urlservice

import (
	"context"
	"errors"
	"log"
	"os"
	"urlshort/internal/api"
	"urlshort/pkg/database"
	"urlshort/pkg/memory"
)

const (
	memoryMethod = iota
	dbMethod
)

// Service добавляет в структуру сервиса метод для хранения URL
type Service struct {
	HashMap  map[string]string
	Database database.DB
	Method   int
}

// InitDbWithTable устанавливает соединения с БД и создает таблицу для хранения URL
func (s *Service) InitDbWithTable(table string) (err error) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	err = s.Database.ConnectDB()
	if err != nil {
		return
	}
	log.Println("PostgreSQL: Connection established on port:", os.Getenv("POSTGRES_PORT"))
	err = s.Database.CreateTableDB(table)
	if err != nil {
		return
	}
	log.Println("PostgreSQL: Table ready")
	return err
}

// ToShortLink имплиментированный метод gRPC, принимает полный URL, возвращает короткий URL
func (s *Service) ToShortLink(ctx context.Context, req *api.FullURL) (*api.ShortURL, error) {
	fUrl := req.Value
	switch s.Method {
	case memoryMethod:
		if value, err := memory.ToHash(s.HashMap, fUrl); err == nil {
			return &api.ShortURL{Value: value}, nil
		} else {
			return nil, err
		}
	case dbMethod:
		if value, err := s.Database.ToHash(fUrl); err == nil {
			return &api.ShortURL{Value: value}, nil
		} else {
			return nil, err
		}
	default:
		return nil, errors.New("wrong method used")
	}
}

// ToFullLink имплиментированный метод gRPC, принимает короткий URL, возвращает полный URL
func (s *Service) ToFullLink(ctx context.Context, req *api.ShortURL) (*api.FullURL, error) {
	sUrl := req.Value
	switch s.Method {
	case memoryMethod:
		if value, err := memory.FromHash(s.HashMap, sUrl); err == nil {
			return &api.FullURL{Value: value}, nil
		} else {
			return nil, err
		}
	case dbMethod:
		if value, err := s.Database.FromHash(sUrl); err == nil {
			return &api.FullURL{Value: value}, nil
		} else {
			return nil, err
		}
	default:
		return nil, errors.New("wrong method used")
	}
}
