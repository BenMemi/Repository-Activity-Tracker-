package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	database "tracker/database"

	//Import googles golang github API
	"github.com/google/go-github/github"

	//Import GORM (go ORM) to interact with the database
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	//
)

const (
//TODO: ENV Variables
//Need username and PAT (personal access token) to access github API for traffic API

//Repositories

// Database connection String
)

var (
	githubUsername = ""
	githubPassword = ""
	dsn            = ""
)

func main() {
	githubUsername = os.Getenv("GITHUB_USERNAME")
	githubPassword = os.Getenv("GITHUB_PASSWORD")
	dsn = os.Getenv("DATABASE_URL")

	for {
		//Need a background context
		ctx := context.Background()

		//Create a new github client with authentication
		tp := github.BasicAuthTransport{
			Username: strings.TrimSpace(githubUsername),
			Password: strings.TrimSpace(githubPassword),
		}
		client := github.NewClient(tp.Client())

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		} else {
			fmt.Println("Connected to database")
			fmt.Println(db)
		}

		db.AutoMigrate(&database.Clone{})
		db.AutoMigrate(&database.View{})
		db.AutoMigrate(&database.Path{})
		db.AutoMigrate(&database.Referral{})

		//Get the traffic data for the repository
		trafficClones, _, err := client.Repositories.ListTrafficClones(ctx, "balancer-labs", "balancer-v2-monorepo", &github.TrafficBreakdownOptions{})
		fmt.Println("TRAFFIC CLONES")
		fmt.Println("///////////////////////////////////////")
		for _, clone := range trafficClones.Clones {
			fmt.Println("timestamp: ", clone.Timestamp)
			fmt.Println("count: ", *clone.Count)
			fmt.Println("uniques: ", *clone.Uniques)
			clone := database.Clone{
				Day:     clone.Timestamp.Time,
				Count:   *clone.Count,
				Uniques: *clone.Uniques,
			}
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&clone)
		}
		fmt.Println("///////////////////////////////////////")

		trafficViewers, _, err := client.Repositories.ListTrafficViews(ctx, "balancer-labs", "balancer-v2-monorepo", &github.TrafficBreakdownOptions{})
		fmt.Println("TRAFFIC VIEWS")
		fmt.Println("///////////////////////////////////////")
		for _, viewer := range trafficViewers.Views {
			fmt.Println("timestamp: ", viewer.Timestamp)
			fmt.Println("count: ", *viewer.Count)
			fmt.Println("uniques: ", *viewer.Uniques)
			viewer := database.View{
				Day:     viewer.Timestamp.Time,
				Count:   *viewer.Count,
				Uniques: *viewer.Uniques,
			}
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&viewer)
		}
		fmt.Println("///////////////////////////////////////")

		//fmt.Println(response)

		trafficPaths, _, err := client.Repositories.ListTrafficPaths(ctx, "balancer-labs", "balancer-v2-monorepo")
		fmt.Println("TRAFFIC PATHS")
		fmt.Println("///////////////////////////////////////")
		for _, path := range trafficPaths {
			fmt.Println("path: ", *path.Path)
			fmt.Println("title: ", *path.Title)
			fmt.Println("count: ", *path.Count)
			fmt.Println("uniques: ", *path.Uniques)
			path := database.Path{
				Path:    *path.Path,
				Title:   *path.Title,
				Count:   *path.Count,
				Uniques: *path.Uniques,
				Day:     time.Now().Truncate(24 * time.Hour),
			}
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&path)
		}
		fmt.Println("///////////////////////////////////////")

		trafficReferrals, _, err := client.Repositories.ListTrafficReferrers(ctx, "balancer-labs", "balancer-v2-monorepo")
		fmt.Println("TRAFFIC REFERRALS")
		fmt.Println("///////////////////////////////////////")
		for _, referral := range trafficReferrals {
			fmt.Println("Referrer: ", *referral.Referrer)
			fmt.Println("count: ", *referral.Count)
			fmt.Println("Unique: ", *referral.Uniques)
			referral := database.Referral{
				Referrer: *referral.Referrer,
				Count:    *referral.Count,
				Uniques:  *referral.Uniques,
				Day:      time.Now().Truncate(24 * time.Hour),
			}
			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&referral)
		}

		fmt.Println("///////////////////////////////////////")

		fmt.Println(err)
		time.Sleep(24 * time.Hour)
	}
}
