# Repository-Activity-Tracker-
Monitors and writes the activity (clones/forks/visits) of a repository. 

# ENV Setup 

1. Make a new env file in the base directory with these variables in lines 1-3
GITHUB_PASSWORD=your PAT
GITHUB_USERNAME=your username
DATABASE_URL=your database url (connection string)

then in lines 6-8 (starting in line 6)
export GITHUB_PASSWORD=your PAT
export GITHUB_USERNAME=your username
export DATABASE_URL=your database url (connection string)

The reason of this is is to copy the lines 6 to 8 from the file and make a start script for you (not tracked by github) in the terraform directory. Then you can just run up.sh and terraform to make a live - free instance that will write to your database. If you want to run local just go run tracker.go will do. 



# Steps 
0. Make GCP service account and get the key file and put the path to the key file in terraform/main.tf (line is marked)
1. Make .env as detailed above 
2. Run up.sh 
3. cd terraform
4. Run terraform init
5. Run terraform apply
