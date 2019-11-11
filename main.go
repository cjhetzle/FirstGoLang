package main

import (
	"fmt"
	"io/ioutil"
	"log"
	//"io"
	"encoding/json"
	"net/http"
	"strings"
)

var sessionid string
var vmware_sid_key string = "vmware-api-session-id"

var globalSessions *session.Manager

// Then, initialize the session manager
func init() {
	globalSessions = NewManager("memory", "gosessionid", 3600)
}

func main() {

	var status = login()
	//fmt.Println(status)

	var result map[string]interface{}
	json.Unmarshal([]byte(status), &result)
	fmt.Println(result)
	fmt.Println(result["value"])
	sessionid := result["value"]
	//status = logout()
	//fmt.Println(status)
}

func login() string {
	var username string = "ur mom"
	var passwd string = "haha_dumb hackers"
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, "https://deehvcr013ccpv1.ssm.sdc.gts.ibm.com/rest/com/vmware/cis/session", nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	fmt.Println("Returning: " + s)
	return s
}

func logout() string {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "https://deehvcr013ccpv1.ssm.sdc.gts.ibm.com/rest/com/vmware/cis/session", nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var sb strings.Builder
	bodyText, err := ioutil.ReadAll(resp.Body)
	sb.WriteString("SUCCESS: LOG OUT -> ")
	sb.WriteString(string(bodyText))
	return sb.String()
}

type Session interface {
	Set(key, value interface{}) error //set session value
	Get(key interface{}) interface{}  //get session value
	Delete(key interface{}) error     //delete session value
	SessionID() string                //back current sessionID
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

type Manager struct {
	cookieName  string     //private cookiename
	lock        sync.Mutex // protects session
	provider    Provider
	maxlifetime int64
}

func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}
