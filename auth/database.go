package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	DataBase *pgxpool.Pool
}

func(d *DB) OnDatabase()(*pgxpool.Pool,error){
	url:= os.Getenv("DB_URL")
	if url==""{
		return nil, fmt.Errorf("env DB url kosong")
	}
	c,cancel:= context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	config, err:= pgxpool.ParseConfig(url)
	if err!=nil{
		return nil, fmt.Errorf("gagal buat config db")
	}

	pool,err:= pgxpool.NewWithConfig(c,config)
	if err!=nil{
		return nil, fmt.Errorf("gagal buat pool db")
	}

	err= pool.Ping(c)
	if err!=nil{
		return nil, fmt.Errorf("gagal buat ping db")
	}
	
	d.DataBase = pool
	return pool, nil
}