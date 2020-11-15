package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/integrii/flaggy"
	_ "github.com/mattn/go-sqlite3" // sqlite

	"github.com/tunedmystic/commits.lol/app/clients/github"
	"github.com/tunedmystic/commits.lol/app/config"
	"github.com/tunedmystic/commits.lol/app/db"
	"github.com/tunedmystic/commits.lol/app/server"
	"github.com/tunedmystic/commits.lol/app/services"
	"github.com/tunedmystic/commits.lol/app/utils"
)

func main() {
	flaggy.SetName("commits.lol")
	flaggy.SetDescription("A web app that collects funny commit messages")

	// Options for subcommands.
	todayDate := time.Now().UTC().Format("2006-01-02")
	fetchCommitsFromDate := todayDate
	fetchCommitsToDate := todayDate

	// The 'run-server' subcommand.
	cmdRunServer := flaggy.NewSubcommand("run-server")
	flaggy.AttachSubcommand(cmdRunServer, 1)

	// The 'fetch-commits' subcommand.
	cmdFetchCommits := flaggy.NewSubcommand("fetch-commits")
	cmdFetchCommits.Description = "Fetch commits by date range"
	cmdFetchCommits.String(&fetchCommitsFromDate, "f", "from", "AuthorDate from")
	cmdFetchCommits.String(&fetchCommitsToDate, "t", "to", "AuthorDate to")
	flaggy.AttachSubcommand(cmdFetchCommits, 1)

	flaggy.Parse()

	if len(os.Args) < 2 {
		flaggy.ShowHelp("")
	}

	if cmdRunServer.Used {
		RunServer()
	}

	if cmdFetchCommits.Used {
		utils.MustParseDate(fetchCommitsFromDate)
		utils.MustParseDate(fetchCommitsToDate)
		FetchCommits(fetchCommitsFromDate, fetchCommitsToDate)
	}
}

// RunServer ...
func RunServer() {
	db := db.NewSqliteDB(config.App.DatabaseName)
	s := server.NewServer(db)

	addr := fmt.Sprintf("0.0.0.0:%v", config.App.Port)
	fmt.Printf("Running server on %v ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.Routes()))
}

// FetchCommits ...
func FetchCommits(fromDate, toDate string) {
	db := db.NewSqliteDB(config.App.DatabaseName)

	options := github.CommitSearchOptions{
		FromDate: fromDate,
		ToDate:   toDate,
		Sort:     github.SortDesc,
	}

	// Run the commit pipeline with randomly fetched searchTerms.
	services.Commits(db).WithOptions(options).WithRandomSearchTerms().Run()
}
