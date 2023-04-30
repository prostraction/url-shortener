package urlshort

import (
	"context"
	"errors"
	"log"
	"urlshort/pkg/api"
	"urlshort/pkg/database"
	"urlshort/pkg/memory"
)

const (
	memoryMethod = iota
	dbMethod
)

type Application struct {
	HashMap  map[string]string
	Database database.DB
	Method   int
}

func (a *Application) InitDB() (err error) {
	err = a.Database.ConnectDB()
	if err != nil {
		return
	}
	log.Println("PostgreSQL: Connection established")
	err = a.Database.CreateTableDB()
	if err != nil {
		return
	}
	log.Println("PostgreSQL: Table ready")
	return err
}

func (a *Application) ToShortLink(ctx context.Context, req *api.FullURL) (*api.ShortURL, error) {
	fUrl := req.Value
	switch a.Method {
	case memoryMethod:
		if value, err := memory.ToHash(a.HashMap, fUrl); err == nil {
			return &api.ShortURL{Value: value}, nil
		} else {
			return nil, err
		}
	case dbMethod:
		if value, err := a.Database.ToHash(fUrl); err == nil {
			return &api.ShortURL{Value: value}, nil
		} else {
			return nil, err
		}
	default:
		return nil, errors.New("wrong method used")
	}
}

func (a *Application) ToFullLink(ctx context.Context, req *api.ShortURL) (*api.FullURL, error) {
	sUrl := req.Value
	switch a.Method {
	case memoryMethod:
		if value, err := memory.FromHash(a.HashMap, sUrl); err == nil {
			return &api.FullURL{Value: value}, nil
		} else {
			return nil, err
		}
	case dbMethod:
		if value, err := a.Database.FromHash(sUrl); err == nil {
			return &api.FullURL{Value: value}, nil
		} else {
			return nil, err
		}
	default:
		return nil, errors.New("wrong method used")
	}
}
