package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := []byte("admin888")
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("生成密码哈希失败: %v\n", err)
		return
	}
	fmt.Printf("密码哈希值: %s\n", string(hashedPassword))

	// 验证生成的哈希值是否正确
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		fmt.Printf("验证失败: %v\n", err)
	} else {
		fmt.Println("验证成功！")
	}
}
