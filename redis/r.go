package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var ctx = context.Background()

func ConnectRedis() (*redis.Client, error) {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		return nil, fmt.Errorf("env redis url kosong")
	}

	// Gunakan ParseURL untuk membaca format URL Upstash (rediss://...)
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("gagal parse redis url: %v", err)
	}

	Redis = redis.NewClient(opt)

	// Tes koneksi
	_, err = Redis.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("gagal connect ke Redis: %v", err)
	}

	return Redis, nil
}

func SetC(n string, isi interface{}, t time.Duration)error{
	
	err := Redis.Set(ctx,n,isi,t).Err()
	if err!=nil{
		fmt.Print("gagal set cache redis")
		return err
	}
	return nil
}

func Getc(n string)(string,error){
	
	ada, err:= Redis.Get(ctx,n).Result()
	if err!=nil{
		fmt.Print("gagal baca redis/kosong")
		return "", err
	}

	return ada, nil
	
}

func Delc(n string)error{
	_= Redis.Del(ctx,n).Err()
	
	return  nil
}