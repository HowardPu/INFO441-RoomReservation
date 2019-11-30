package handlers

import (
	U "INFO441-RoomReservation/servers/gateway/models/users"
	S "INFO441-RoomReservation/servers/gateway/sessions"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-redis/redis"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

var server = "mssql441.c2mbdnajn2pb.us-east-1.rds.amazonaws.com"
var user = "admin"
var password = "info441ishard"
var database = "INFO441A5"
var port = "1433"

var signingKey = "JusticsFromAbove"

var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;port=%s", server, user, password, database, port)
var db, connectErr = sql.Open("mssql", connString)

var msstore = U.NewMsSqlStore(db)
var redisaddr = "127.0.0.1:6379"
var client = redis.NewClient(&redis.Options{
	Addr: redisaddr,
})

var redisStore = S.NewRedisStore(client, time.Hour)

var ctx = HandlerContext{
	SessionKey:   signingKey,
	UserStore:    msstore,
	SessionStore: redisStore,
}

func TestUsersHandlers(t *testing.T) {
	cases := []struct {
		name             string
		method           string
		requestHeader    string
		expectStatusCode int
		user             *U.NewUser
	}{
		{
			"Not POST Method",
			"GET",
			"",
			405,
			&U.NewUser{"aaa@a.com", "info441", "info441", "Hedonist", "John", "Smith"},
		}, {
			"Not JSON HEADER",
			"POST",
			"text",
			http.StatusUnsupportedMediaType,
			&U.NewUser{"aaa@a.com", "info441", "info441", "Hedonist", "John", "Smith"},
		}, {
			"INVALID USER CASE",
			"POST",
			"application/json",
			http.StatusBadRequest,
			&U.NewUser{"fff@f.com", "info", "info44", "Evil", "John", "Smith"},
		}, {
			"Regular CASE",
			"POST",
			"application/json",
			http.StatusCreated,
			&U.NewUser{"fff@f.com", "info441", "info441", "Evil", "John", "Smith"},
		},
	}

	for _, c := range cases {
		jsonDat, _ := json.Marshal(c.user)
		req, _ := http.NewRequest(c.method, "/v1/users", bytes.NewBuffer(jsonDat))
		req.Header.Set("Content-Type", c.requestHeader)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctx.UsersHandler)
		handler.ServeHTTP(rr, req)

		code := rr.Code

		if code != c.expectStatusCode {
			t.Errorf("Expect Error Code is %d but we get %d", c.expectStatusCode, code)
		} else {
			if c.expectStatusCode == http.StatusCreated {
				newUser := UserLite{}

				reqResult := rr.Result()
				defer reqResult.Body.Close()

				body, _ := ioutil.ReadAll(reqResult.Body)

				json.Unmarshal(body, &newUser)
				userName := newUser.UserName

				userInSQL, getErr := ctx.UserStore.GetByUserName(userName)
				if getErr != nil {
					t.Errorf("Wrong Update: %v", getErr)
				}
				if newUser.ID != userInSQL.ID {
					t.Errorf("Wrong User In request, expect id as %d but get %d", userInSQL.ID, newUser.ID)
				}
				ctx.UserStore.Delete(userInSQL.ID)
			}
		}

	}
}

func TestSpecificUserHandler(t *testing.T) {
	user1 := &U.User{-1, "aaa@a.com", []byte("xbujjjhb"), "Justice", "John", "Smith", "dnwin"}
	user2 := &U.User{-1, "bbb@b.com", []byte("xbujjjhb"), "Evil", "John", "Smith", "dnwin"}

	newUser1, _ := ctx.UserStore.Insert(user1)
	newUser2, _ := ctx.UserStore.Insert(user2)

	rr := httptest.NewRecorder()

	authForUser1, _ := S.BeginSession(ctx.SessionKey, ctx.SessionStore, newUser1, rr)

	cases := []struct {
		name             string
		method           string
		requestHeader    string
		expectStatusCode int
		id               string
		auth             string
		update           *U.Updates
	}{
		{
			"Not GET/PATCH",
			"POST",
			"",
			http.StatusMethodNotAllowed,
			"",
			"",
			nil,
		}, {
			"Not Auth",
			"PATCH",
			"application/json",
			http.StatusUnauthorized,
			strconv.FormatInt(newUser2.ID, 10),
			"",
			&U.Updates{"Alpha", "Beta"},
		}, {
			"USER NOT FOUND",
			"GET",
			"",
			http.StatusForbidden,
			"-1",
			authForUser1.String(),
			nil,
		}, {
			"REGULAR GET",
			"GET",
			"",
			http.StatusOK,
			strconv.FormatInt(newUser1.ID, 10),
			authForUser1.String(),
			nil,
		}, {
			"PATCH NOT MATCH",
			"PATCH",
			"application/json",
			http.StatusForbidden,
			strconv.FormatInt(newUser2.ID, 10),
			authForUser1.String(),
			&U.Updates{"Alpha", "Beta"},
		}, {
			"PATCH WRONG HEADER",
			"PATCH",
			"text",
			http.StatusUnsupportedMediaType,
			strconv.FormatInt(newUser1.ID, 10),
			authForUser1.String(),
			&U.Updates{"Gamma", "Delta"},
		}, {
			"REG CASE: ID",
			"PATCH",
			"application/json",
			http.StatusOK,
			strconv.FormatInt(newUser1.ID, 10),
			authForUser1.String(),
			&U.Updates{"Zeta", "Epsilon"},
		}, {
			"REG CASE: ME",
			"PATCH",
			"application/json",
			http.StatusOK,
			"me",
			authForUser1.String(),
			&U.Updates{"Zeta", "Epsilon"},
		},
	}

	for _, c := range cases {
		//t.Logf(schemeBearer + c.auth)
		var reqContent io.Reader = nil
		if c.update != nil {
			jsonDat, _ := json.Marshal(c.update)
			reqContent = bytes.NewBuffer(jsonDat)
		}

		req, _ := http.NewRequest(c.method, "/v1/users/"+c.id, reqContent)
		req.Header.Add("Content-Type", c.requestHeader)
		req.Header.Add(headerAuthorization, schemeBearer+c.auth)
		rr = httptest.NewRecorder()

		handler := http.HandlerFunc(ctx.SpecificUserHandler)
		handler.ServeHTTP(rr, req)

		code := rr.Code

		if code != c.expectStatusCode {
			t.Errorf("%v - Expect Error Code is %d but we get %d", c.name, c.expectStatusCode, code)
		}

		if code == c.expectStatusCode && c.expectStatusCode == http.StatusOK {
			newUser := UserLite{}
			reqResult := rr.Result().Body
			decoder := json.NewDecoder(reqResult)
			decoder.Decode(&newUser)

			expectID := newUser1.ID

			if c.id != "me" {
				expectID, _ = strconv.ParseInt(c.id, 10, 64)
			}

			if newUser.ID != expectID {
				t.Errorf("%v - Wrong User In request, expect id as %d but get %d", c.name, newUser1.ID, newUser.ID)
			}
			if c.method == "PATCH" {
				userInSQL, _ := ctx.UserStore.GetById(newUser1.ID)
				if userInSQL.LastName != c.update.LastName || userInSQL.FirstName != c.update.FirstName {
					t.Errorf("%v - Does not update user to First Name %s, Last Name %s", c.name, c.update.FirstName, c.update.LastName)
				}
			}
		}

	}

	ctx.UserStore.Delete(newUser1.ID)
	ctx.UserStore.Delete(newUser2.ID)
	ctx.SessionStore.Delete(authForUser1)
}

func TestSessionsHandler(t *testing.T) {
	newUser := &U.NewUser{
		"aaa@a.com", "INFO441A4", "INFO441A4",
		"uwlaziestperson1", "Justice", "Evil",
	}

	user, _ := newUser.ToUser()
	userInDB, _ := ctx.UserStore.Insert(user)

	cases := []struct {
		name             string
		method           string
		requestHeader    string
		expectStatusCode int
		credential       *U.Credentials
		ip               string
		XForward         string
	}{
		{
			"WRONG METHOD",
			"GET",
			"",
			http.StatusMethodNotAllowed,
			nil,
			"",
			"",
		}, {
			"WRONG CONTENT TYPE",
			"POST",
			"text",
			http.StatusUnsupportedMediaType,
			&U.Credentials{"aaa@a.com", "INFO441A4"},
			"",
			"",
		}, {
			"WRONG EMAIL",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"bbb@b.com", "INFO441A4"},
			"",
			"",
		}, {
			"WRONG PWD",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "info441a4"},
			"",
			"",
		}, {
			"CORRECT CASE Normal",
			"POST",
			"application/json",
			http.StatusCreated,
			&U.Credentials{"aaa@a.com", "INFO441A4"},
			"192.168.1.1",
			"192.168.1.2, 192.367.2.2",
		}, {
			"Reset: 1 - Fail",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "info441a4"},
			"",
			"",
		}, {
			"Reset: 2 - Fail",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "info441a4"},
			"",
			"",
		}, {
			"Reset: 3 - Fail",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "info441a4"},
			"",
			"",
		}, {
			"Reset: 4 - Fail",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "info441a4"},
			"",
			"",
		}, {
			"Reset: 5 - Success",
			"POST",
			"application/json",
			http.StatusCreated,
			&U.Credentials{"aaa@a.com", "INFO441A4"},
			"192.168.1.1",
			"192.168.1.2, 192.367.2.2",
		}, {
			"Reset: 6 - Fail",
			"POST",
			"application/json",
			http.StatusCreated,
			&U.Credentials{"aaa@a.com", "INFO441A4"},
			"192.168.1.1",
			"192.168.1.2, 192.367.2.2",
		}, {
			"Reset: 7 - Success",
			"POST",
			"application/json",
			http.StatusCreated,
			&U.Credentials{"aaa@a.com", "INFO441A4"},
			"192.168.1.1",
			"192.168.1.2, 192.367.2.2",
		}, {
			"Block: 1 - Fail",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "info441a4"},
			"",
			"",
		}, {
			"Block: 2 - Fail",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "info441a4"},
			"",
			"",
		}, {
			"Block: 3 - Fail",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "info441a4"},
			"",
			"",
		}, {
			"Block: 4 - Fail",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "info441a4"},
			"",
			"",
		}, {
			"Block: 5 - Fail",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "info441a4"},
			"",
			"",
		}, {
			"Block: 6 - Block",
			"POST",
			"application/json",
			http.StatusUnauthorized,
			&U.Credentials{"aaa@a.com", "INFO441A4"},
			"192.168.1.1",
			"192.168.1.2, 192.367.2.2",
		},
	}

	for _, c := range cases {

		jsonDat, _ := json.Marshal(c.credential)
		req, _ := http.NewRequest(c.method, "/v1/sessions", bytes.NewBuffer(jsonDat))
		req.RemoteAddr = c.ip
		req.Header.Set("Content-Type", c.requestHeader)
		req.Header.Add("X-Forwarded-For", c.XForward)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctx.SessionsHandler)
		handler.ServeHTTP(rr, req)

		code := rr.Code

		if code != c.expectStatusCode {
			t.Errorf("%v: Expect Error Code is %d but we get %d",
				c.name, c.expectStatusCode, code)
		}

		newUser := UserLite{}
		reqResult := rr.Result().Body
		decoder := json.NewDecoder(reqResult)
		decoder.Decode(&newUser)

		if code == c.expectStatusCode && c.expectStatusCode == http.StatusCreated {
			if userInDB.ID != newUser.ID {
				t.Errorf("%v: Wrong User In request, expect id as %d but get %d",
					c.name, userInDB.ID, newUser.ID)
			}
		}

	}
	ctx.SessionStore.RemoveFailRecord("aaa@a.com")
	ctx.UserStore.Delete(userInDB.ID)
}

func TestSpecificSessionHandler(t *testing.T) {
	user := &U.User{-1, "zzz@z.com", []byte("xbujjjhb"), "Test", "John", "Smith", "dnwin"}
	newUser, _ := ctx.UserStore.Insert(user)
	rr := httptest.NewRecorder()
	authForNewUser, _ := S.BeginSession(ctx.SessionKey, ctx.SessionStore, newUser, rr)

	cases := []struct {
		name             string
		method           string
		path             string
		expectStatusCode int
	}{
		{
			"Non DELETE Method",
			"POST",
			"/v1/sessions/mine",
			http.StatusMethodNotAllowed,
		}, {
			"Incorrect Path",
			"DELETE",
			"/v1/sessions/wrong",
			http.StatusForbidden,
		}, {
			"Regular Cases",
			"DELETE",
			"/v1/sessions/mine",
			-1,
		},
	}

	for _, c := range cases {
		req, _ := http.NewRequest(c.method, c.path, nil)
		req.Header.Add(headerAuthorization, schemeBearer+authForNewUser.String())

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(ctx.SpecificSessionHandler)
		handler.ServeHTTP(rr, req)

		code := rr.Code

		if c.expectStatusCode != -1 && code != c.expectStatusCode {
			t.Errorf("%v: Expect Error Code is %d but we get %d",
				c.name, c.expectStatusCode, code)
		}

		if code == -1 {
			reqResult, _ := ioutil.ReadAll(req.Body)
			if string(reqResult) != "signed out" {
				t.Errorf("%v: unexpected message delivered, expected signed out but got %d", c.name, reqResult)
			}

			sessionUser := U.User{}
			_, getErr := S.GetState(req, ctx.SessionKey, ctx.SessionStore, &sessionUser)
			if getErr == nil {
				t.Errorf("%v: expected session to be deleted but didn't", c.name)
			}
		}
	}
	ctx.UserStore.Delete(newUser.ID)
}
