package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"encoding/json"
	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"

	DB_HOST_DEFAULT = "35.237.202.244" // Google cloud VM
	//DB_HOST_DEFAULT = "localhost"
	DB_PORT_DEFAULT = "5432"
	DB_USER_DEFAULT = "stellar"
	DB_PASS_DEFAULT = "123"
	DB_NAME_DEFAULT = "test"
)

func main() {
	initDb()
	defer db.Close()

	http.HandleFunc("/api/createTestData", createTestData)
	http.HandleFunc("/api/index", indexHandler)
	//http.HandleFunc("/api/repo/", repoHandler)
	log.Fatal(http.ListenAndServe("localhost:8100", nil))
}

////////////////////////////////////////////////////

// repository contains the details of a repository
type repositorySummary struct {
	ID         int
	Name       string
	Owner      string
	TotalStars int
}

type repositories struct {
	Repositories []repositorySummary
}

// indexHandler calls `queryRepos()` and marshals the result as JSON
func indexHandler(w http.ResponseWriter, r *http.Request) {
	repos := repositories{}

	err := queryRepos(&repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	fmt.Fprintf(w, string(out))
}
func queryRepos(repos *repositories) error {
	rows, err := db.Query(`
		SELECT
			id,
			repository_owner,
			repository_name,
			total_stars
		FROM repositories
		ORDER BY total_stars DESC`)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		repo := repositorySummary{}
		err = rows.Scan(
			&repo.ID,
			&repo.Owner,
			&repo.Name,
			&repo.TotalStars,
		)
		if err != nil {
			return err
		}
		repos.Repositories = append(repos.Repositories, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////////

// indexHandler calls `queryRepos()` and marshals the result as JSON
func createTestData(w http.ResponseWriter, r *http.Request) {
	err := createtestdataRepositories()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Done")
}
func createtestdataRepositories() error {
	return makeQuery(`
		INSERT INTO repositories (id, repository_owner, repository_name, total_stars)
		(
			SELECT new_data.id, new_data.repository_owner, new_data.repository_name, new_data.total_stars
			FROM (
				SELECT 1 as id, 'Owner 1' as repository_owner, 'Name 1' as repository_name, 11 as total_stars
				UNION ALL 
				SELECT 2 as id, 'Owner 2' as repository_owner, 'Name 2' as repository_name, 12 as total_stars
				UNION ALL 
				SELECT 3 as id, 'Owner 3' as repository_owner, 'Name 3' as repository_name, 13 as total_stars
				UNION ALL 
				SELECT 4 as id, 'Owner 4' as repository_owner, 'Name 4' as repository_name, 14 as total_stars
				UNION ALL 
				SELECT 5 as id, 'Owner 5' as repository_owner, 'Name 5' as repository_name, 15 as total_stars
				UNION ALL 
				SELECT 6 as id, 'Owner 6' as repository_owner, 'Name 6' as repository_name, 16 as total_stars
			) as new_data
			LEFT JOIN repositories as old_data 
			ON new_data.id = old_data.id
			WHERE old_data.id is NULL 
		);`)
}

////////////////////////////////////////////////////

//// repository contains the details of a repository
//type repository struct {
//	ID              int
//	Name            string
//	Owner           string
//	RepoAge         int
//	Initialized     bool
//	CommitsPerMonth string
//	StarsPerMonth   string
//	TotalStars      int
//}
//
//// Error handling types
//
//type errRepoNotInitialized string
//func (e errRepoNotInitialized) Error() string {
//	return string(e)
//}
//type errRepoNotFound string
//func (e errRepoNotFound) Error() string {
//	return string(e)
//}
//
//// parseParams accepts a req and returns the `num` path tokens found after the `prefix`.
//// returns an error if the number of tokens are less or more than expected
//func parseParams(req *http.Request, prefix string, num int) ([]string, error) {
//	url := strings.TrimPrefix(req.URL.Path, prefix)
//	params := strings.Split(url, "/")
//	if len(params) != num || len(params[0]) == 0 || len(params[1]) == 0 {
//		return nil, fmt.Errorf("Bad format. Expecting exactly %d params", num)
//	}
//	return params, nil
//}
//
//// repoHandler processes the response by parsing the params, then calling
//// `query()`, and marshaling the result in JSON format, sending it to
//// `http.ResponseWriter`.
//func repoHandler(w http.ResponseWriter, r *http.Request) {
//	repo := repository{}
//	params, err := parseParams(r, "/api/repo/", 2)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusUnauthorized)
//		return
//	}
//	repo.Owner = params[0]
//	repo.Name = params[1]
//
//	data, err := queryRepo(&repo)
//	if err != nil {
//		switch err.(type) {
//		case errRepoNotFound:
//			http.Error(w, err.Error(), 404)
//		case errRepoNotInitialized:
//			http.Error(w, err.Error(), 401)
//		default:
//			http.Error(w, err.Error(), 500)
//		}
//		return
//	}
//
//	out, err := json.Marshal(data)
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//		return
//	}
//
//	fmt.Fprintf(w, string(out))
//}
//
//// queryRepo first fetches the repository, and if nothing is wrong
//// it returns the result of fetchData()
//func queryRepo(repo *repository) (*repoData, error) {
//	err := fetchRepo(repo)
//	if err != nil {
//		return nil, err
//	}
//	return fetchData(repo)
//}
//
//// fetchData calls utility functions to collect data from
//// the database, builds and returns the `RepoData` value
//func fetchData(repo *repository) (*repoData, error) {
//	data := repoData{}
//	err := fetchMonthlyData(repo, &data)
//	if err != nil {
//		return nil, err
//	}
//	err := fetchWeeklyData(repo, &data)
//	if err != nil {
//		return nil, err
//	}
//	err := fetchYearlyData(repo, &data)
//	if err != nil {
//		return nil, err
//	}
//	err := fetchTimelineData(repo, &data)
//	if err != nil {
//		return nil, err
//	}
//	err := fetchOwnerData(repo, &data)
//	if err != nil {
//		return nil, err
//	}
//	return &data, nil
//}
//
//// fetchRepo given a Repository value with name and owner of the repo
//// fetches more details from the database and fills the value with more
//// data
//func fetchRepo(repo *repository) error {
//	if len(repo.Name) == 0 {
//		return fmt.Errorf("Repository name not correctly set")
//	}
//	if len(repo.Owner) == 0 {
//		return fmt.Errorf("Repository owner not correctly set")
//	}
//	sqlStatement := `
//		SELECT
//			id,
//			initialized,
//            repository_created_months_ago
//		FROM repositories
//		WHERE repository_owner=$1 and repository_name=$2
//        LIMIT 1;`
//	row := db.QueryRow(sqlStatement, repo.Owner, repo.Name)
//	err := row.Scan(&repo.ID, &repo.Initialized, &repo.RepoAge)
//	if err != nil {
//		switch err {
//		case sql.ErrNoRows:
//			//locally handle SQL error, abstract for caller
//			return errRepoNotFound("Repository not found")
//		default:
//			return err
//		}
//	}
//
//	if !repo.Initialized {
//		return errRepoNotInitialized("Repository not initialized")
//	}
//	if repo.RepoAge < 3 {
//		return errRepoNotInitialized("Repository not initialized")
//	}
//	return nil
//}
//
//// fetchOwnerData given a Repository object with the `Owner` value
//// it fetches information about it from the database
//func fetchOwnerData(repo *repository, data *repoData) error {
//	if len(repo.Owner) == 0 {
//		return fmt.Errorf("Repository owner not correctly set")
//	}
//	sqlStatement := `
//		SELECT
//			id,
//			name,
//            COALESCE(description, ''),
//            COALESCE(avatar_url, ''),
//            COALESCE(github_id, ''),
//            added_by,
//            enabled,
//            COALESCE(installation_id, ''),
//            repository_selection
//        FROM organizations
//        WHERE name=$1
//        ORDER BY id DESC LIMIT 1`
//	row := db.QueryRow(sqlStatement, repo.Owner)
//	err := row.Scan(
//		&data.Owner.ID,
//		&data.Owner.ID,
//		&data.Owner.ID,
//		&data.Owner.ID,
//		&data.Owner.ID,
//		&data.Owner.ID,
//		&data.Owner.ID,
//		&data.Owner.ID,
//		)
//}
////////////////////////////////////////////////////

func initDb() {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport], config[dbuser], config[dbpass], config[dbname])

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected! Checking tables ...")

	err = checkTables(config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables checked/prepared")
}
func checkTables(config map[string]string) error {
	return makeQuery(`
		CREATE TABLE if not exists ` + config[dbname] + `.public.repositories (
			id int,
    		repository_owner varchar(255),
    		repository_name varchar(255),
    		total_stars int
		);`)
}

func dbConfig() map[string]string {
	conf := make(map[string]string)

	host, ok := os.LookupEnv(dbhost)
	if !ok {
		//panic("DBHOST environment variable required but not set")
		host = DB_HOST_DEFAULT
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		//panic("DBPORT environment variable required but not set")
		port = DB_PORT_DEFAULT
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		//panic("DBUSER environment variable required but not set")
		user = DB_USER_DEFAULT
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		//panic("DBPASS environment variable required but not set")
		password = DB_PASS_DEFAULT
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		//panic("DBNAME environment variable required but not set")
		name = DB_NAME_DEFAULT
	}

	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name

	return conf
}

////////////////////////////////////////////////////
func makeQuery(query string) error {
	rows, err := db.Query(query)
	if rows != nil {
		defer rows.Close()
		if rows.Next() {
			fmt.Println("rows.Next()")
		}
	}
	if err != nil {
		return err
	}
	return nil
}
