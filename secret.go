// Copyright (C) 2019 Evgeny Kuznetsov (evgeny@kuznetsov.md)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along tihe this program. If not, see <https://www.gnu.org/licenses/>.

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
