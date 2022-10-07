package main

import (
	"context"
	"fmt" //for printing
	"log"
	"os"
	"strings" //String conversion library
	"time"    //for time

	database "tracker/database" //our local db module for dealing with our db, mainly schema in there

	"github.com/google/go-github/github" //Import googles golang github API

	"github.com/profclems/go-dotenv" //Import dotenv library to deal with env variables before CICD is needed

	//Import GORM (go ORM) to interact with the database
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	//Repository struct to hold the repository name and the owner of the repo
	//TODO: Expand this to include additional PATs maybe?
	Owner string
	Repo  string
}

// Globaly used variables
var (
	//github username
	githubUsername = ""
	//github PAT (personal access token)
	githubPassword = ""
	//database connection string
	dsn = ""

	//Github Repository List
	//ADD YOUR REPOS HERE <--------------
	repositories = []Repository{
		{
			Owner: "balancer-labs",
			Repo:  "balancer-v2-monorepo",
		},
	}
)

func main() {
	//load repos

	//Load the .evn
	err := dotenv.LoadConfig()
	if err != nil {
		//panic if we cannot load the .env
		fmt.Println("error loading .env file")
	}

	//grab the .env variables, careful this will silently fail!
	githubPassword = dotenv.GetString("GITHUB_PASSWORD")
	githubUsername = dotenv.GetString("GITHUB_USERNAME")
	dsn = dotenv.GetString("DATABASE_URL")

	//check if all the variables exist in the .env, if not panic and send an error message
	if githubPassword == "" || githubUsername == "" || dsn == "" {
		//try getting from os as well (if in production)
		githubPassword = os.Getenv("GITHUB_PASSWORD")
		githubUsername = os.Getenv("GITHUB_USERNAME")
		dsn = os.Getenv("DATABASE_URL")
		if githubPassword == "" || githubUsername == "" || dsn == "" {
			panic("Missing .env variables!")
		}
	}

	//infinite loop of calls and authentication to github and db
	for {
		//Need a background context, just standard stuff
		ctx := context.Background()

		//Create a new github http client that is authenticated with our PAT
		tp := github.BasicAuthTransport{
			Username: strings.TrimSpace(githubUsername),
			Password: strings.TrimSpace(githubPassword),
		}
		//this is the client object
		client := github.NewClient(tp.Client())

		//start a connection to the database with out connection string
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		//panic if we cannot connect to the database
		if err != nil {
			panic("failed to connect database")
		} else {
			//or else we are good to go
			fmt.Println("Connected to database")
			fmt.Println(db)
		}

		//Migrate schemas to the database (make/edit tables)
		//TODO: this should go through an array, not hard coded
		db.AutoMigrate(&database.Clone{})
		db.AutoMigrate(&database.View{})
		db.AutoMigrate(&database.Path{})
		db.AutoMigrate(&database.Referral{})

		for _, repo := range repositories {
			//Get the traffic data for the repository with an API request
			trafficClones, _, err := client.Repositories.ListTrafficClones(ctx, repo.Owner, repo.Repo, &github.TrafficBreakdownOptions{})
			fmt.Println("TRAFFIC CLONES")
			fmt.Println("///////////////////////////////////////")

			if err != nil {
				log.Fatalln("could not make the API request", err)
			}

			//TODO should be a function not copy pasting code
			//Loop through all the days where clones happened and add them to the database and print them out
			for _, clone := range trafficClones.Clones {
				fmt.Println("timestamp: ", clone.Timestamp)
				fmt.Println("count: ", *clone.Count)
				fmt.Println("uniques: ", *clone.Uniques)
				clone := database.Clone{
					Day:        clone.Timestamp.Time,
					Count:      *clone.Count,
					Uniques:    *clone.Uniques,
					Repository: repo.Repo,
				}
				db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&clone)
			}
			fmt.Println("///////////////////////////////////////")

			//Get the traffic data for the repository with an API request
			trafficViewers, _, err := client.Repositories.ListTrafficViews(ctx, repo.Owner, repo.Repo, &github.TrafficBreakdownOptions{})
			fmt.Println("TRAFFIC VIEWS")
			fmt.Println("///////////////////////////////////////")
			if err != nil {
				log.Fatalln("could not make the API request", err)
			}
			//Loop through all the days where views happened and add them to the database and print them out
			for _, viewer := range trafficViewers.Views {
				fmt.Println("timestamp: ", viewer.Timestamp)
				fmt.Println("count: ", *viewer.Count)
				fmt.Println("uniques: ", *viewer.Uniques)
				viewer := database.View{
					Day:        viewer.Timestamp.Time,
					Count:      *viewer.Count,
					Uniques:    *viewer.Uniques,
					Repository: repo.Repo,
				}
				db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&viewer)
			}
			fmt.Println("///////////////////////////////////////")

			//Get the traffic data for the repository with an API request
			trafficPaths, _, err := client.Repositories.ListTrafficPaths(ctx, repo.Owner, repo.Repo)
			if err != nil {
				log.Fatalln("could not make the API request", err)
			}
			fmt.Println("TRAFFIC PATHS")
			fmt.Println("///////////////////////////////////////")
			for _, path := range trafficPaths {
				fmt.Println("path: ", *path.Path)
				fmt.Println("title: ", *path.Title)
				fmt.Println("count: ", *path.Count)
				fmt.Println("uniques: ", *path.Uniques)

				//Loop through all the days where paths happened and add them to the database and print them out
				path := database.Path{
					Path:       *path.Path,
					Title:      *path.Title,
					Count:      *path.Count,
					Uniques:    *path.Uniques,
					Day:        time.Now().Truncate(24 * time.Hour), //Gets the day of the path was retrieved
					Repository: repo.Repo,
				}
				db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&path)
			}
			fmt.Println("///////////////////////////////////////")

			//same as above but with referrals
			trafficReferrals, _, err := client.Repositories.ListTrafficReferrers(ctx, repo.Owner, repo.Repo)
			if err != nil {
				log.Fatalln("could not make the API request", err)
			}
			fmt.Println("TRAFFIC REFERRALS")
			fmt.Println("///////////////////////////////////////")
			for _, referral := range trafficReferrals {
				fmt.Println("Referrer: ", *referral.Referrer)
				fmt.Println("count: ", *referral.Count)
				fmt.Println("Unique: ", *referral.Uniques)
				referral := database.Referral{
					Referrer:   *referral.Referrer,
					Count:      *referral.Count,
					Uniques:    *referral.Uniques,
					Day:        time.Now().Truncate(24 * time.Hour),
					Repository: repo.Repo,
				}
				db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&referral)
			}

			fmt.Println("///////////////////////////////////////")
		}
		//sleep for 24 hours and then do it again
		time.Sleep(24 * time.Hour)
	}
}
