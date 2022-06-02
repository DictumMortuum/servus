package deta

import (
	"github.com/DictumMortuum/servus/pkg/config"
	api "github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
)

func New(name string) (*base.Base, error) {
	d, err := api.New(api.WithProjectKey(config.App.Deta.ProjectKey))
	if err != nil {
		return nil, err
	}

	db, err := base.New(d, name)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Populate(db *base.Base, item interface{}) error {
	_, err := db.Put(&item)
	if err != nil {
		return err
	}
	return nil
}
