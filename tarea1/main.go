package main


import (
  "database/sql"
  "fmt"
  "html/template"
  "net/http"
  "time"



  _ "github.com/lib/pq"
)

// Comment blablalba
type Comment struct {
  Time     sql.NullString
  IP       sql.NullString
  Message  sql.NullString
}

// Params blablabla
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
  password = "hola"
  dbname   = "test"
)


func index(w http.ResponseWriter, r *http.Request) {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

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
    // handle this error better than this
    // log.Errorf(ctx, "Getting comments: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't get latest comments. Refresh?"
		indexTemplate.Execute(w, params)
  }
  
  // insert
  // info := ""
  // sqlStatement := `
  //   INSERT INTO info(time, IP, message)
  //   VALUES ($1, $2, $3)
  //   RETURNING message`
  // err = db.QueryRow(sqlStatement, "7:41", "88889", "jhonnattan").Scan(&info)
  // if err != nil {
    //   panic(err)
    // }
    // fmt.Println("New record ID is", info)
    
  
  // var comments []Comment
  for rows.Next() {
    var comment Comment
    err = rows.Scan(&comment.Time, &comment.IP, &comment.Message)
    // err = rows.Scan(&id, &firstName)
    if err != nil {
      // handle this error
      panic(err)
    }
    params.Comments = append(params.Comments, comment)
    // params.Comments = append([]Comment{comment}, params.Comments...)
    // comments = append(comments, comment)
    // params.Comments = comments //append(params.Comments, comment)
    // comments = append(comments, comment)
    // fmt.Println(comment)
  }
  // reverse
  for i, j := 0, len(params.Comments)-1; i < j; i, j = i+1, j-1 {
      params.Comments[i], params.Comments[j] = params.Comments[j], params.Comments[i]
  }
  // fmt.Println(params.Comments)
  if r.Method == "GET" {
    indexTemplate.Execute(w, params)
    return
  }
  // params.Comments = append([]Comment{comment}, params.Comments...)
  comment := Comment{
		IP:  sql.NullString{ String: r.RemoteAddr, Valid: true },
    Message:  sql.NullString{ String: r.FormValue("comment"), Valid: true },
    Time:  sql.NullString{ String: time.Now().Format("15:04"), Valid: true },
    
		// Message: r.FormValue("comment"),
		// Time:  time.Now().Format("15:04"),
  }
  
  // insert
  info := ""
  sqlStatement = `
  INSERT INTO info(time, IP, message)
  VALUES ($1, $2, $3)
  RETURNING message`

  err = db.QueryRow(sqlStatement, comment.Time, comment.IP, comment.Message).Scan(&info)
  // fmt.Println("New record ID is", info)
  if err != nil {
    // panic(err)
    w.WriteHeader(http.StatusInternalServerError)
    params.Notice = "Couldn't add new post. Try again?"
    params.Message = comment.Message.String // Preserve their message so they can try again.
    indexTemplate.Execute(w, params)
    return
  }
  // params.Comments = append(params.Comments, comment)
  params.Comments = append([]Comment{comment}, params.Comments...)
	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", comment.IP.String)
  // fmt.Println(params.Notice, "lalalala")
  indexTemplate.Execute(w, params)
  return
}

func main() {
    http.HandleFunc("/", index)
    if err := http.ListenAndServe(":8080", nil); err != nil {
      panic(err)
    }
}
  // get any error encountered during iteration
  // err = rows.Err()
  // if err != nil {
  //   panic(err)
  // }
  // switch err {
  //   case sql.ErrNoRows:
  //     fmt.Println("No rows were returned!")
  //     return
  //   case nil:
  //     fmt.Println(comment)
  //   default:
  //     panic(err)
  // }



