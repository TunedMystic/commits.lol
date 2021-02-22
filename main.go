package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/integrii/flaggy"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/tunedmystic/commits.lol/app/clients/github"
	"github.com/tunedmystic/commits.lol/app/config"
	"github.com/tunedmystic/commits.lol/app/db"
	"github.com/tunedmystic/commits.lol/app/pipeline"
	"github.com/tunedmystic/commits.lol/app/server"
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
	cmdRunServer := flaggy.NewSubcommand("server")
	cmdRunServer.Description = "Run the application server"
	flaggy.AttachSubcommand(cmdRunServer, 1)

	// The 'fetch-commits' subcommand.
	cmdFetchCommits := flaggy.NewSubcommand("fetch-commits")
	cmdFetchCommits.Description = "Fetch commits by date range"
	cmdFetchCommits.String(&fetchCommitsFromDate, "f", "from", "AuthorDate from")
	cmdFetchCommits.String(&fetchCommitsToDate, "t", "to", "AuthorDate to")
	flaggy.AttachSubcommand(cmdFetchCommits, 1)

	// The 'limits' subcommand.
	cmdLimits := flaggy.NewSubcommand("limits")
	cmdLimits.Description = "Check API rate limits"
	flaggy.AttachSubcommand(cmdLimits, 1)

	flaggy.Parse()

	if len(os.Args) < 2 {
		flaggy.ShowHelp("")
		return
	}

	utils.SetupLogging()
	flushSentry := utils.SetupSentry()
	defer flushSentry()

	if cmdRunServer.Used {
		RunTasks()
		RunServer()
	}

	if cmdFetchCommits.Used {
		utils.MustParseDate(fetchCommitsFromDate)
		utils.MustParseDate(fetchCommitsToDate)
		FetchCommits(fetchCommitsFromDate, fetchCommitsToDate)
	}

	if cmdLimits.Used {
		CheckRateLimits()
	}
}

// RunServer ...
func RunServer() {
	zap.S().Info("[run] server")
	db := db.NewSqliteDB(config.App.DatabaseName)
	defer db.Close()

	s := server.NewServer(&db)

	addr := fmt.Sprintf("0.0.0.0:%v", config.App.Port)
	zap.S().Info("Server is running on ", addr)
	log.Fatal(http.ListenAndServe(addr, s.Routes()))
}

// RunTasks ...
func RunTasks() {
	zap.S().Info("[run] periodic tasks")
	c := cron.New()
	c.AddFunc("@every 60m", func() {
		to := time.Now().UTC()
		from := to.AddDate(0, 0, -3) // 3 days back.
		FetchCommits(from.Format("2006-01-02"), to.Format("2006-01-02"))
	})
	c.Start()
}

// FetchCommits ...
func FetchCommits(fromDate, toDate string) {
	zap.S().Infof("[run] fetch-commits from %s to %s", fromDate, toDate)
	db := db.NewSqliteDB(config.App.DatabaseName)
	defer db.Close()

	options := github.CommitSearchOptions{
		FromDate: fromDate,
		ToDate:   toDate,
		Sort:     github.SortDesc,
	}

	// Run the commit pipeline with randomly fetched searchTerms.
	p := pipeline.Commits(&db)
	p.WithOptions(options)
	p.WithRandomSearchTerms()
	p.Run()
	zap.S().Info("[done] fetch-commits")
}

// CheckRateLimits ...
func CheckRateLimits() {
	zap.S().Infof("[run] limits")
	c := github.NewClient()

	response, err := c.RateLimits()
	if err != nil {
		log.Fatal(err)
	}

	limitsDisplay, _ := json.MarshalIndent(response.Resources, "", "   ")

	zap.S().Infof("Github Rate Limits\n%s", string(limitsDisplay))
}
