package urlshort

import (
	"context"
	"urlshort/pkg/api"
	memory "urlshort/pkg/memory"
)

type Application struct {
	HashMap map[string]string
}

func (a *Application) ToShortLink(ctx context.Context, req *api.FullURL) (*api.ShortURL, error) {
	fUrl := req.Value
	if value, err := memory.ToHash(a.HashMap, fUrl); err == nil {
		return &api.ShortURL{Value: value}, nil
	} else {
		return nil, err
	}
}

func (a *Application) ToFullLink(ctx context.Context, req *api.ShortURL) (*api.FullURL, error) {
	sUrl := req.Value
	if value, err := memory.FromHash(a.HashMap, sUrl); err == nil {
		return &api.FullURL{Value: value}, nil
	} else {
		return nil, err
	}
}
