package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/netutil"
	"golang.org/x/sync/errgroup"

	. "github-network/database/pkg"
	"github-network/database/pkg/db"
	"github-network/database/pkg/fetcher"
)

type StatusCode = int

var logger = log.New(os.Stdout, "[server]", log.Lshortfile)

type Server struct {
	port            uint16
	lock            bool
	limitConnection int
	respTime        time.Duration
}

func InitServer(port uint16) *Server {
	if err := db.ConnectDb(); err != nil {
		log.Fatal(err)
	}
	return &Server{
		port:            port,
		limitConnection: 1,     // 同時接続を1までにする
		lock:            false, // 処理をブロック
		respTime:        3000 * time.Millisecond,
	}
}

func (server *Server) Run() {
	logger.Println("server is running...")
	http.HandleFunc("/pulls", server.pullRequestHandler)
	http.HandleFunc("/issues", server.issuesHandler)
	http.HandleFunc("/users", server.usersHandler)
	http.HandleFunc("/all", server.allDataHandler)
	//log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", server.port), nil))
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	l = netutil.LimitListener(l, server.limitConnection)
	http.Serve(l, nil)
}

func (server *Server) batchProcessing(r *http.Request, f func(*http.Request) StatusCode) StatusCode {
	server.lock = true
	fmt.Printf("\n")
	logger.Println("server locked! >> batch process started!")
	statusCode := make(chan StatusCode)
	go func() {
		statusCode <- f(r)
		server.lock = false
		fmt.Printf("\n")
		logger.Println("\nserver unlocked!")
	}()

	select {
	case code := <-statusCode:
		return code
	case <-time.After(server.respTime):
		return http.StatusAccepted
	}
}

func (server *Server) pullRequestHandler(w http.ResponseWriter, r *http.Request) {
	if server.lock == true {
		logger.Println("blocked request (server is locked)!")
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, "try later!")
		return
	}
	statusCode := server.batchProcessing(r, processPullRequestsData)
	w.WriteHeader(statusCode)

	switch statusCode {
	case http.StatusOK:
		io.WriteString(w, "done!")
	case http.StatusAccepted:
		io.WriteString(w, "processing!")
	default:
		io.WriteString(w, "error occurred!")
	}
}

func processPullRequestsData(r *http.Request) StatusCode {
	// parse request
	urlValues := r.URL.Query()
	owner := urlValues.Get("owner")
	repo := urlValues.Get("repo")
	if owner == "" || repo == "" {
		return http.StatusBadRequest
	}

	return _processPullRequestsData(owner, repo)
}

func _processPullRequestsData(owner string, repo string) StatusCode {
	eg := &errgroup.Group{}

	for page := 1; ; page++ {
		fmt.Printf(
			"processing: pull_requests, owner = %s, repo = %s, page = %d\r",
			owner, repo, page,
		)
		prs := PullRequests(owner, repo, page)
		// data fetch
		if err := prs.Fetch(); err != nil {
			return http.StatusBadRequest
		}

		eg.Go(func() error {
			// insert into database
			if err := prs.InsertIntoDb(); err != nil {
				return err
			}
			return nil
		})

		if len(*prs.PullRequests) < fetcher.ResultsPerPage {
			break
		}
		time.Sleep(fetcher.SleepTime)
	}
	fmt.Println("")

	if err := eg.Wait(); err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (server *Server) issuesHandler(w http.ResponseWriter, r *http.Request) {
	if server.lock == true {
		logger.Println("blocked request (server is locked)!")
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, "try later!")
		return
	}
	statusCode := server.batchProcessing(r, processIssuesData)
	w.WriteHeader(statusCode)

	switch statusCode {
	case http.StatusOK:
		io.WriteString(w, "done!")
	case http.StatusAccepted:
		io.WriteString(w, "processing!")
	default:
		io.WriteString(w, "error occurred!")
	}
}

func processIssuesData(r *http.Request) StatusCode {
	urlValues := r.URL.Query()
	owner := urlValues.Get("owner")
	repo := urlValues.Get("repo")
	if owner == "" || repo == "" {
		return http.StatusBadRequest
	}

	return _processIssuesData(owner, repo)
}

func _processIssuesData(owner string, repo string) StatusCode {
	eg := &errgroup.Group{}

	for page := 1; ; page++ {
		fmt.Printf(
			"processing: issues, owner = %s, repo = %s, page = %d\r",
			owner, repo, page,
		)
		iss := Issues(owner, repo, page)
		// data fetch
		if err := iss.Fetch(); err != nil {
			return http.StatusBadRequest
		}

		eg.Go(func() error {
			// insert into database
			if err := iss.InsertIntoDb(); err != nil {
				return err
			}
			return nil
		})

		if len(*iss.Issues) < fetcher.ResultsPerPage {
			break
		}
		time.Sleep(fetcher.SleepTime)
	}

	if err := eg.Wait(); err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (server *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
	if server.lock == true {
		logger.Println("blocked request (server is locked)!")
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, "try later!")
		return
	}
	statusCode := server.batchProcessing(r, processUsersData)
	w.WriteHeader(statusCode)

	switch statusCode {
	case http.StatusOK:
		io.WriteString(w, "done!")
	case http.StatusAccepted:
		io.WriteString(w, "processing!")
	default:
		io.WriteString(w, "error occurred!")
	}
}

func processUsersData(_ *http.Request) StatusCode {
	return _processUsersData()
}

func _processUsersData() StatusCode {
	user_logins, err := db.QueryUserLogins()
	if err != nil {
		return http.StatusInternalServerError
	}
	var wg sync.WaitGroup
	errChan := make(chan error)

	for _, user_login := range user_logins {
		fmt.Printf(
			"processing: user, user_login = %s",
			user_login,
		)
		user := User(user_login)
		// data fetch
		if err := user.Fetch(); err != nil {
			return http.StatusBadRequest
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := user.InsertIntoDb(); err != nil {
				errChan <- err
			}
		}()

		time.Sleep(fetcher.SleepTime)
	}

	for err := range errChan {
		if err != nil {
			return http.StatusInternalServerError
		}
	}
	return http.StatusOK
}

func (server *Server) allDataHandler(w http.ResponseWriter, r *http.Request) {
	statusCode := server.batchProcessing(r, processAllData)
	w.WriteHeader(statusCode)

	switch statusCode {
	case http.StatusOK:
		io.WriteString(w, "done!")
	case http.StatusAccepted:
		io.WriteString(w, "processing!")
	default:
		io.WriteString(w, "error occurred!")
	}
}

func processAllData(r *http.Request) StatusCode {
	// parse request
	urlValues := r.URL.Query()
	owner := urlValues.Get("owner")
	repo := urlValues.Get("repo")
	if owner == "" || repo == "" {
		return http.StatusBadRequest
	}

	var statusCode = http.StatusOK
	if code := _processPullRequestsData(owner, repo); code != 200 {
		statusCode = code
	}
	time.Sleep(fetcher.SleepTime)
	if code := _processIssuesData(owner, repo); code != 200 {
		statusCode = code
	}
	time.Sleep(fetcher.SleepTime)
	if code := _processUsersData(); code != 200 {
		statusCode = code
	}

	return statusCode
}
