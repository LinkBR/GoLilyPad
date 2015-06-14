package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"net/http"
)

type GameProfile struct {
	Id string `json:"id"`
	Properties []GameProfileProperty `json:"properties"`
}

type GameProfileProperty struct {
	Name string `json:"name"`
	Value string `json:"value"`
	Signature string `json:"signature"`
}

func HasPaid(name string) (b bool) {
	resp, err := http.Get("https://minecraft.net/haspaid.jsp?user=" + name)

	if err == nil {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			fmt.Printf(string(contents))
			check, err := strconv.ParseBool(string(contents))

			if err == nil {
				b = check
				return
			}
		}
	}
	b = true
	return
}

func Authenticate(name string, serverId string, sharedSecret []byte, publicKey []byte) (profile GameProfile, err error) {
	response, err := http.Get(fmt.Sprintf(URL, name, MojangSha1Hex([]byte(serverId), sharedSecret, publicKey)))
	if err != nil {
		return
	}
	jsonDecoder := json.NewDecoder(response.Body)
	profile = GameProfile{}
	err = jsonDecoder.Decode(&profile)
	response.Body.Close()
	if err != nil {
		return
	}
	if len(profile.Id) != 32 {
		err = errors.New(fmt.Sprintf("Id is not 32 characters: %d", len(profile.Id)))
	}
	return
}
