package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Credentials struct {
	id     string
	secret string
}

func readSecret(file string, cr *Credentials) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var f map[string]interface{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		return err
	}

	cred := f["installed"].(map[string]interface{})

	cr.id = cred["client_id"].(string)
	cr.secret = cred["client_secret"].(string)
	return nil
}

func main() {
	cred := Credentials{
		id:     "",
		secret: "",
	}
	err := readSecret(".client_secret.json", &cred)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("ID is %s\n", cred.id)
		fmt.Printf("Secret is %s\n", cred.secret)
	}
}
