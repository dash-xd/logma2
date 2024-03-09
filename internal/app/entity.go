package entity

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type BaseEntity struct {
	Env        string
	EntityName string
	EntityID   int
}

func (e *BaseEntity) Register(client *redis.Client) (string, error) {
	// Implement registration logic using Redis
	registrationString := fmt.Sprintf("%s:%d", e.EntityName, e.EntityID)
	result, err := client.Do(context.Background(), "FCALL", "REGISTER", registrationString).Result()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), nil
}

func (e *BaseEntity) Unregister(client *redis.Client) (string, error) {
	// Implement unregistration logic using Redis
	unregistrationString := fmt.Sprintf("%s:%d", e.EntityName, e.EntityID)
	result, err := client.Do(context.Background(), "FCALL", "UNREGISTER", unregistrationString).Result()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), nil
}

type Entity struct {
	BaseEntity
	ParentNamespace string
	ChildNamespace  string
}

type ActiveSubscription struct {
	Entity
}

type SavedSubscription struct {
	Entity
}

func DefaultEntity(env string, entityName string, entityID int) Entity {
	return Entity{
		BaseEntity: BaseEntity{
			Env:        env,
			EntityName: entityName,
			EntityID:   entityID,
		},
		ParentNamespace: "myOrg",
		ChildNamespace:  "global",
	}
}
