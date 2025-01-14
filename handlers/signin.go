package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	//"html/template"
	"io/ioutil"
	"net/http"
)

const registeredUsers = "users"

//User ...
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin/" {
		fmt.Fprintf(w, "<script>alert(\"page %s not found\")</script>", r.URL.Path)
		return
	}
	if r.Method == http.MethodPost {
		inputEmail := r.FormValue("inputEmail")
		inputPassword := r.FormValue("inputPassword")
		fmt.Println("Login is:", inputEmail)
		fmt.Println("Password is:", inputPassword)
		//TODO: Validation

		file, err := os.Open(registeredUsers)
		defer file.Close()
		if err != nil {
			fmt.Println("file open error:", err)
			return
		}
		reader := bufio.NewReader(file)
		user := User{}
		emailFound := false
		for {
			line, err := reader.ReadBytes(byte('\n')) //file should has '\n' as last symbol
			if err != nil {
				break
			}
			err = json.Unmarshal(line, &user)
			if err != nil {
				fmt.Println("json unmarshal error:", err)
				return
			}
			if user.Login == inputEmail {
				emailFound = true
				if user.Password == inputPassword {
					http.Redirect(w, r, "/blog/", http.StatusPermanentRedirect) // redirect make full path itself from r.Host
					return
				}
				fmt.Fprintln(w, "<script>alert(\"wrong password\")</script>")
				break
			}
		}
		if emailFound == false {
			fmt.Fprintln(w, "<script>alert(\"wrong login\")</script>")
		}
	}
	fileContent, err := ioutil.ReadFile("static/pages/signin.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(fileContent)
}
