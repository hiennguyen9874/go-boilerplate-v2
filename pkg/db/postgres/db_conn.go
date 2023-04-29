package postgres

import (
	"fmt"

	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/ent"

	_ "github.com/lib/pq"
)

func NewPsqlClient(cfg *config.Config) (*ent.Client, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Dbname, cfg.Postgres.Port,
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func Close(client *ent.Client) error {
	return client.Close()
}
