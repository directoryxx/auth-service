package repository

import (
	"auth/internal/domain"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
)

// UserRepository represent the user's repository contract
type UserRepository interface {
	GetAll(context.Context) ([]*domain.User, error)
	GetOneByID(ctx context.Context, id int) (*domain.User, error)
	GetOneByUsername(ctx context.Context, username string) (*domain.User, error)
	Insert(ctx context.Context, input *domain.User) (*domain.User, error)
	Update(ctx context.Context, id int, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id int) error
	RememberUUID(ctx context.Context, user *domain.User, uuid string) error
	GetUUID(ctx context.Context, uuid string) (string, error)
	Publish(ctx context.Context, data string, topic string) error
}

type UserRepositoryImpl struct {
	DB    *sql.DB
	Redis *redis.Client
	Kafka *kafka.Producer
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewUserRepository(db *sql.DB, Redis *redis.Client, kafkaProducer *kafka.Producer) UserRepository {
	return &UserRepositoryImpl{
		DB:    db,
		Redis: Redis,
		Kafka: kafkaProducer,
	}
}

func (m *UserRepositoryImpl) GetAll(context context.Context) (res []*domain.User, err error) {
	rows, err := m.DB.Query(`SELECT id, name, username, password FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var users []*domain.User

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var user *domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Username,
			&user.Password); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func (m *UserRepositoryImpl) GetOneByID(context context.Context, id int) (res *domain.User, err error) {
	stmt, err := m.DB.PrepareContext(context, "SELECT id, email, name, username, password FROM users WHERE id=$1")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRowContext(context, id)
	var user domain.User

	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserRepositoryImpl) GetOneByUsername(ctx context.Context, username string) (res *domain.User, err error) {
	stmt, err := m.DB.PrepareContext(ctx, "SELECT id, name, email, username, password FROM users WHERE username=$1")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, username)
	var user domain.User

	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserRepositoryImpl) Insert(ctx context.Context, input *domain.User) (user *domain.User, err error) {
	stmt := `insert into users (name, email, username, password)
		values ($1, $2, $3, $4) returning id`

	var newID int

	err = m.DB.QueryRowContext(ctx, stmt,
		input.Name,
		input.Email,
		input.Username,
		input.Password,
	).Scan(&newID)

	if err != nil {
		return nil, err
	}

	user, err = m.GetOneByID(ctx, newID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *UserRepositoryImpl) Update(ctx context.Context, id int, update *domain.User) (user *domain.User, err error) {
	stmt := `update users set
		name = $1,
		username = $2,
		password = $3,
		where id = $4
	`

	_, err = m.DB.ExecContext(ctx, stmt,
		update.Name,
		update.Username,
		update.Password,
		id,
	)

	if err != nil {
		return nil, err
	}

	user, err = m.GetOneByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *UserRepositoryImpl) Delete(ctx context.Context, id int) (err error) {
	stmt := `delete from users where id = $1`

	_, err = m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserRepositoryImpl) RememberUUID(ctx context.Context, user *domain.User, uuid string) error {
	userModel, _ := json.Marshal(user)
	jwtHourExpire := os.Getenv("JWT_EXPIRE_HOUR")
	convJwtHour, _ := strconv.Atoi(jwtHourExpire)
	err := m.Redis.Set(ctx, uuid, userModel, time.Hour*time.Duration(convJwtHour)).Err()
	return err
}

func (m *UserRepositoryImpl) GetUUID(ctx context.Context, uuid string) (res string, err error) {
	res, err = m.Redis.Get(ctx, uuid).Result()
	return res, err
}

func (m *UserRepositoryImpl) Publish(ctx context.Context, data string, topic string) error {
	err := m.Kafka.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(data),
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sending message success: %s", data)

	return err
}
