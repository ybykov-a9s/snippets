package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

const ValidToken = "MD7)(HGTT#@$LK!PPO0"
const InvalidToken = "__Illegal_access_token"

type SearchResult struct {
	Status int
	Err    string
}

type TestCaseConnectToServer struct {
	ID         string
	Connection *SearchClient
	IsError    bool
	ErrorMsg   string
}

type TestCaseLimitOffset struct {
	ID       string
	Request  *SearchRequest
	IsError  bool
	ErrorMsg string
}

type Person struct {
	ID            int    `xml:"id"`
	GUID          string `xml:"guid"`
	IsActive      bool   `xml:"isActive"`
	Balance       string `xml:"balance"`
	Picture       string `xml:"picture"`
	Age           int    `xml:"age"`
	EyeColor      string `xml:"eyeColor"`
	FirstName     string `xml:"first_name"`
	LastName      string `xml:"last_name"`
	Gender        string `xml:"gender"`
	Company       string `xml:"company"`
	Email         string `xml:"email"`
	Phone         string `xml:"phone"`
	Address       string `xml:"address"`
	About         string `xml:"about"`
	Registered    string `xml:"registered"`
	FavoriteFruit string `xml:"favoriteFruit"`
}

type Persons struct {
	Version string   `xml:"version,attr"`
	List    []Person `xml:"row"`
}

func MyResponder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	key := r.FormValue("query")
	if key == "__broken_json" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provocated bad JSON without finishing curly brace here"`) //broken json
		return
	}

	key = r.FormValue("query")
	if key == "__bad_request" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provocated bad request"}`) //bad request
		return
	}

	key = r.FormValue("query")
	if key == "__cycled_redirect" {
		http.Redirect(w, r, r.URL.String(), http.StatusSeeOther) // Cycled redirect
		return
	}

	key = r.FormValue("query")
	if key == "__short_server_timeout" {
		time.Sleep(2 * time.Second)
		http.Redirect(w, r, r.URL.String(), http.StatusSeeOther) // Pause is longer than http timeout
		return
	}

	key = r.FormValue("order_by")
	Order, _ := strconv.Atoi(key)
	if Order < -1 || Order > 1 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "ErrorBadOrderField"}`)
		return
	}

	key = r.FormValue("order_field")

	if key != "" && key != "ID" && key != "Age" && key != "Name" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "StatusInternalServerError"}`)
		return
	}

	at := r.Header.Get("AccessToken")
	if at == "*" {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "StatusInternalServerError"}`)
		return
	} else if at != ValidToken {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error" : "StatusUnauthorized"}`)
		return
	}

	// Final part
	req := SearchRequest{}
	req.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	req.Offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	req.Query = r.URL.Query().Get("query")
	req.OrderField = r.URL.Query().Get("order_field")
	req.OrderBy, _ = strconv.Atoi(r.URL.Query().Get("order_by"))

	PersonsList, _ := LoadXML()

	data, err := SearchXML(PersonsList, req)
	if err != nil {
		fmt.Printf("MyResponder() - Error calling SearchXML: %v", err)
		return
	}

	// Sending bad JSON after AUTH StatusOK with
	if req.Query == "__error_marshal" {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{`)
		return
	}

	w.Write(data)
}

func LoadXML() (v *Persons, e error) {
	xmlData, err := ioutil.ReadFile("./dataset.xml")
	if err != nil {
		fmt.Printf("LoadXML() - Error loading XML file: %v", err)
		return nil, err
	}

	err = xml.Unmarshal(xmlData, &v)
	if err != nil {
		fmt.Printf("LoadXML() - Error Unmarshal XML: %v", err)
		return nil, err
	}
	return v, nil
}

func SearchXML(v *Persons, r SearchRequest) (resp []byte, e error) {
	sResult := []Person{}

	for _, i := range v.List {
		if r.Query != "" {
			if strings.Contains(i.FirstName+" "+i.LastName, r.Query) || strings.Contains(i.About, r.Query) {
				sResult = append(sResult, i)
			}
		} else {
			sResult = append(sResult, i)
		}
	}

	switch r.OrderField {
	case "Age":
		{
			if r.OrderBy == OrderByAsc {
				sort.Slice(sResult[:], func(i, j int) bool {
					return sResult[i].Age < sResult[j].Age
				})
			} else if r.OrderBy == OrderByDesc {
				sort.Slice(sResult[:], func(i, j int) bool {
					return sResult[i].Age > sResult[j].Age
				})
			}
		}
	case "ID":
		{
			if r.OrderBy == OrderByAsc {
				sort.Slice(sResult[:], func(i, j int) bool {
					return sResult[i].ID < sResult[j].ID
				})
			} else if r.OrderBy == OrderByDesc {
				sort.Slice(sResult[:], func(i, j int) bool {
					return sResult[i].ID > sResult[j].ID
				})
			}
		}
	default:
		{
			if r.OrderBy == orderAsc {
				sort.Slice(sResult, func(i, j int) bool {
					if sResult[i].LastName < sResult[j].LastName {
						return true
					}
					if sResult[i].LastName > sResult[j].LastName {
						return false
					}
					return sResult[i].FirstName < sResult[j].FirstName
				})
			} else if r.OrderBy == OrderByDesc {
				sort.Slice(sResult, func(i, j int) bool {
					if sResult[i].LastName < sResult[j].LastName {
						return true
					}
					if sResult[i].LastName > sResult[j].LastName {
						return false
					}
					return sResult[i].FirstName > sResult[j].FirstName
				})
			}
		}
	}

	// Offset
	if len(sResult) >= r.Offset {
		sResult = sResult[r.Offset:]
	} else {
		sResult = []Person{}
	}

	//  Limit
	if len(sResult) > r.Limit {
		sResult = sResult[:r.Limit]
	} else {
	}

	result, err := json.Marshal(sResult)
	if err != nil {
		fmt.Printf("SearchXML() - Error Marshal XML: %v", err)
		return nil, err
	}
	return result, nil
}

func SearchServer() (ts *httptest.Server) {
	ts = httptest.NewUnstartedServer(http.HandlerFunc(MyResponder))
	ts.Config = &http.Server{
		ReadTimeout:  1000 * time.Millisecond,
		WriteTimeout: 1000 * time.Millisecond,
		IdleTimeout:  1000 * time.Millisecond,
		Handler:      http.HandlerFunc(MyResponder),
	}
	ts.Start()
	return ts
}

func TestAccessToken(t *testing.T) {
	cases := []TestCaseConnectToServer{
		{
			ID: "Invalid token",
			Connection: &SearchClient{
				URL:         "",
				AccessToken: InvalidToken,
			},
			IsError:  true,
			ErrorMsg: "Bad AccessToken",
		},
		{
			ID: "Absent token",
			Connection: &SearchClient{
				URL:         "",
				AccessToken: "",
			},
			IsError:  true,
			ErrorMsg: "Bad AccessToken",
		},
	}

	req := SearchRequest{
		Limit:      1,
		Offset:     0,
		Query:      "",
		OrderField: "Age",
		OrderBy:    OrderByAsc,
	}

	ts := SearchServer()
	defer ts.Close()

	for caseNum, item := range cases {

		srv := &SearchClient{
			URL:         ts.URL,
			AccessToken: item.Connection.AccessToken,
		}

		_, err := srv.FindUsers(req)

		if item.IsError && fmt.Sprintf("%s", err) != item.ErrorMsg {
			t.Errorf("[%d] Error: %#v", caseNum, err)
		}
	}

	// trigger http.StatusInternalServerError by sending "*" as AccessToken
	srv := &SearchClient{
		URL:         ts.URL,
		AccessToken: "*",
	}

	_, err := srv.FindUsers(req)
	if fmt.Sprintf("%s", err) != "SearchServer fatal error" {
		t.Errorf("Error: %#v", err)
	}
}

func TestLimitOffset(t *testing.T) {
	cases := []TestCaseLimitOffset{
		{
			ID: "Limit less than 0",
			Request: &SearchRequest{
				Limit:      -1,
				Offset:     0,
				Query:      "",
				OrderField: "0",
				OrderBy:    1,
			},
			IsError:  true,
			ErrorMsg: "limit must be > 0",
		},
		{
			ID: "Limit greater than 25",
			Request: &SearchRequest{
				Limit:      30,
				Offset:     0,
				Query:      "",
				OrderField: "0",
				OrderBy:    1,
			},
			IsError:  false,
			ErrorMsg: "",
		}, {
			ID: "Offset less than 0",
			Request: &SearchRequest{
				Limit:      1,
				Offset:     -1,
				Query:      "",
				OrderField: "0",
				OrderBy:    1,
			},
			IsError:  true,
			ErrorMsg: "offset must be > 0",
		},
	}

	ts := SearchServer()
	defer ts.Close()

	srv := &SearchClient{
		URL:         ts.URL,
		AccessToken: ValidToken,
	}

	for caseNum, item := range cases {
		req := item.Request
		_, err := srv.FindUsers(*req)

		if item.IsError && fmt.Sprintf("%s", err) != item.ErrorMsg {
			t.Errorf("[%d] Error: %#v", caseNum, err)
		}
	}

}

func TestClientRequest(t *testing.T) {
	cases := []TestCaseLimitOffset{
		{
			ID: "Bad server JSON",
			Request: &SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "__broken_json",
				OrderField: "0",
				OrderBy:    1,
			},
			IsError:  true,
			ErrorMsg: "cant unpack error json:",
		},
		{
			ID: "Too negative OrderBy",
			Request: &SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Joe",
				OrderField: "",
				OrderBy:    -2,
			},
			IsError:  true,
			ErrorMsg: "OrderFeld",
		}, {
			ID: "Too positive OrderBy",
			Request: &SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Joe",
				OrderField: "",
				OrderBy:    2,
			},
			IsError:  true,
			ErrorMsg: "OrderFeld",
		},
		{
			ID: "Bad request",
			Request: &SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "__bad_request",
				OrderField: "0",
				OrderBy:    1,
			},
			IsError:  true,
			ErrorMsg: "",
		},
		{
			ID: "Cycled Redirect",
			Request: &SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "__cycled_redirect",
				OrderField: "0",
				OrderBy:    1,
			},
			IsError:  true,
			ErrorMsg: "",
		},
		{
			ID: "Server timneout",
			Request: &SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "__short_server_timeout",
				OrderField: "0",
				OrderBy:    1,
			},
			IsError:  true,
			ErrorMsg: "",
		},
	}

	ts := SearchServer()
	defer ts.Close()

	srv := &SearchClient{
		URL:         ts.URL,
		AccessToken: ValidToken,
	}

	for caseNum, item := range cases {
		req := item.Request
		_, err := srv.FindUsers(*req)

		if item.IsError && !strings.Contains(fmt.Sprintf("%s", err), item.ErrorMsg) {
			t.Errorf("[%d] Error: %#v", caseNum, err)
		}
	}
}

func TestLoadAndUnmarshalExternalXML(t *testing.T) {
	_, err := LoadXML()
	if err != nil {
		t.Errorf("Error: %#v", err)
	}

	cases := []TestCaseLimitOffset{
		{
			ID: "Sort by descending age",
			Request: &SearchRequest{
				Limit:      24,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Age",
				OrderBy:    OrderByDesc,
			},
			IsError:  false,
			ErrorMsg: "",
		},
		{
			ID: "Sort by ascending Name",
			Request: &SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "",
				OrderField: "Age",
				OrderBy:    OrderByAsc,
			},
			IsError:  false,
			ErrorMsg: "",
		},
		{
			ID: "Sort by ascending ID",
			Request: &SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "cillum",
				OrderField: "Name",
				OrderBy:    OrderByDesc,
			},
			IsError:  false,
			ErrorMsg: "",
		},
		{
			ID: "Error Marshaling",
			Request: &SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "__error_marshal",
				OrderField: "Name",
				OrderBy:    OrderByDesc,
			},
			IsError:  true,
			ErrorMsg: "unexpected end of JSON input",
		},
	}
	ts := SearchServer()
	defer ts.Close()

	srv := &SearchClient{
		URL:         ts.URL,
		AccessToken: ValidToken,
	}

	for caseNum, item := range cases {
		req := item.Request
		_, err := srv.FindUsers(*req)

		if err != nil {
			if item.IsError && !strings.Contains(fmt.Sprintf("%s", err), item.ErrorMsg) {
				t.Errorf("[%d] Error: %#v", caseNum, err)
			}
		}
		// Last part with checking client results is not implemented yet
	}
}
