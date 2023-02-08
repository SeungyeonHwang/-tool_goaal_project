package main

import (
	"html/template"
	"os"
)

type User struct {
	Name  string
	Email string
	Age   int
}

func (u User) IsOld() bool {
	return u.Age > 30
}

func main() {
	user := User{Name: "seungyeon", Email: "syhwang.web@gmail.com", Age: 31}
	user2 := User{Name: "aaa", Email: "aaa.web@gmail.com", Age: 18}
	users := []User{user, user2}
	tmpl, err := template.New("Tmpl1").ParseFiles("templateApp/templates/tmpl1.tmpl", "templateApp/templates/tmpl2.tmpl")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", users)
}
