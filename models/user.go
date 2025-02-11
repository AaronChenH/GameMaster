package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"-"`    // 密码不返回给前端
	Role      string             `bson:"role" json:"role"`     // admin或service
	Status    int                `bson:"status" json:"status"` // 1:启用 0:禁用
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Permission struct {
	Role        string   `bson:"role"`
	Permissions []string `bson:"permissions"`
}
