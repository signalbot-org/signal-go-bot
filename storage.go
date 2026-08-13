package signalbot

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

// ErrNotFound is returned when a key is not found in the storage
var ErrNotFound = errors.New("key not found")

type Storage interface {
	Exists(ctx context.Context, key string) (bool, error)
	Read(ctx context.Context, key string, v any) error
	Save(ctx context.Context, key string, v any) error
	Delete(ctx context.Context, key string) error
}

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dataSourceName string) (*SQLiteStorage, error) {
	if dataSourceName == "" {
		dataSourceName = ":memory:"
	}

	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	query := `CREATE TABLE IF NOT EXISTS signalbot (key TEXT UNIQUE, value TEXT)`
	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Exists(ctx context.Context, key string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM signalbot WHERE key = ?)`
	err := s.db.QueryRowContext(ctx, query, key).Scan(&exists)
	return exists, err
}

func (s *SQLiteStorage) Read(ctx context.Context, key string, v any) error {
	var val string
	query := `SELECT value FROM signalbot WHERE key = ?`
	err := s.db.QueryRowContext(ctx, query, key).Scan(&val)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return json.Unmarshal([]byte(val), v)
}

func (s *SQLiteStorage) Save(ctx context.Context, key string, v any) error {
	valBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	val := string(valBytes)

	query := `INSERT INTO signalbot (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?`
	_, err = s.db.ExecContext(ctx, query, key, val, val)
	return err
}

func (s *SQLiteStorage) Delete(ctx context.Context, key string) error {
	query := `DELETE FROM signalbot WHERE key = ?`
	_, err := s.db.ExecContext(ctx, query, key)
	return err
}

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(host string, port int, password string) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0,
	})
	return &RedisStorage{client: client}
}

func (s *RedisStorage) Exists(ctx context.Context, key string) (bool, error) {
	n, err := s.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (s *RedisStorage) Read(ctx context.Context, key string, v any) error {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return err
	}
	return json.Unmarshal([]byte(val), v)
}

func (s *RedisStorage) Save(ctx context.Context, key string, v any) error {
	valBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, valBytes, 0).Err()
}

func (s *RedisStorage) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}
