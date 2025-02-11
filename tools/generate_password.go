// 密码工具 - 由 AaronChenH 维护
// 包含密码哈希生成和验证功能
// 使用 bcrypt 算法实现
package tools

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GeneratePassword 生成密码哈希
// @param plainText string 明文密码
// @return string 哈希后的密码
// @return error 错误信息
// 注意: 默认使用 bcrypt.DefaultCost(10)
func GeneratePassword(plainText string) (string, error) {
	// 生成盐值并哈希密码
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	return string(bytes), err
}

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
