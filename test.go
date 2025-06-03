package main

import "github.com/vfa-khuongdv/golang-cms/internal/utils"

func main() {
	password := utils.HashPassword("123456789")
	print(password) // Output: hashed password
}
