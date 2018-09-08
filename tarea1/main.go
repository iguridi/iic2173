package main


import (
  "database/sql"
  "fmt"
  "html/template"
  "net/http"
  "time"
  "os"

  _ "github.com/lib/pq"
)

// Comment the object to be showed inthe web app
type Comment struct {
  Time     sql.NullString
  IP       sql.NullString
  Message  sql.NullString
}

// Params the params to send to the view
type Params struct {
	Notice string
	Message string
	IP string

	Comments []Comment

}

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)


const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  dbname   = "test"
)


func index(w http.ResponseWriter, r *http.Request) {
  if os.Getenv("PGPASSWORD") == "" {
    panic("PGPASSWORD env variable not set")
  }
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, "hola", dbname)

  params := Params{}
  

  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  // get
  sqlStatement := `SELECT * FROM info`
  rows, err := db.Query(sqlStatement)
  defer rows.Close()
  if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't get latest comments. Refresh?"
		indexTemplate.Execute(w, params)
  }

  for rows.Next() {
    var comment Comment
    err = rows.Scan(&comment.Time, &comment.IP, &comment.Message)
    if err != nil {
      panic(err)
    }
    params.Comments = append(params.Comments, comment)
  }

  // reverse
  for i, j := 0, len(params.Comments)-1; i < j; i, j = i+1, j-1 {
      params.Comments[i], params.Comments[j] = params.Comments[j], params.Comments[i]
  }

  if r.Method == "GET" {
    indexTemplate.Execute(w, params)
    return
  }
  
  comment := Comment{
		IP:  sql.NullString{ String: r.RemoteAddr, Valid: true },
    Message:  sql.NullString{ String: r.FormValue("comment"), Valid: true },
    Time:  sql.NullString{ String: time.Now().Format("15:04"), Valid: true },
  }
  
  // insert
  info := ""
  sqlStatement = `
  INSERT INTO info(time, IP, message)
  VALUES ($1, $2, $3)
  RETURNING message`

  err = db.QueryRow(sqlStatement, comment.Time, comment.IP, comment.Message).Scan(&info)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    params.Notice = "Couldn't add new post. Try again?"
    params.Message = comment.Message.String // Preserve their message so they can try again.
    indexTemplate.Execute(w, params)
    return
  }

  params.Comments = append([]Comment{comment}, params.Comments...)
	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", comment.IP.String)
  indexTemplate.Execute(w, params)
  return
}

func main() {
    http.HandleFunc("/", index)
    if err := http.ListenAndServe(":8080", nil); err != nil {
      panic(err)
    }
}



