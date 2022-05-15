package others

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/liulihaocai/YetAnotherControlPanel/util"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Session  string `json:"session"` // session token should be changed every time user login/logout/change password
	MD5Hash  string `json:"-"`
}

var SESSION_TOKEN_LENGTH = 128

func (a *Account) UpdateHash() {
	(*a).MD5Hash = util.GetMD5Hash(a.Username + a.Password)
}

func (a *Account) UpdateSession() {
	(*a).Session = util.RandStringRunes(SESSION_TOKEN_LENGTH)
	SaveAccounts() // TODO: save slowly to reduce disk IO
}

var accounts []Account

func InitAccounts() error {
	accountFile := configDir + "account.json"

	// check file exists
	if _, err := os.Stat(accountFile); os.IsNotExist(err) {
		accounts = []Account{
			{
				Username: util.RandStringRunes(8),
				Password: util.RandStringRunes(8),
			},
		}
		log.Println("Your initial account is:\n\tusername=" + accounts[0].Username + "\n\tpassword=" + accounts[0].Password)
		err := SaveAccounts()
		if err != nil {
			return err
		}
	} else {
		// read file
		file, err := ioutil.ReadFile(accountFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(file, &accounts)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(accounts); i++ {
		accounts[i].UpdateHash()
	}

	return nil
}

func FindAccountHash(hash string) *Account {
	for i := 0; i < len(accounts); i++ {
		if accounts[i].MD5Hash == hash {
			return &accounts[i]
		}
	}

	return nil
}

func CheckSession(session string) bool {
	for i := 0; i < len(accounts); i++ {
		if accounts[i].Session == session {
			return true
		}
	}

	return false
}

func SaveAccounts() error {
	accountFile := configDir + "account.json"

	file, err := json.MarshalIndent(accounts, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(accountFile, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
