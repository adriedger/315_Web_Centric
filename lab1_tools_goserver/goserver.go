package main

import(
    "fmt"
    "net/http"
    "os"
    "time"
    "io"
    "bytes"
)

var variable *bytes.Buffer = &bytes.Buffer{}

func handleVariable(w http.ResponseWriter, r *http.Request){
    switch r.Method{
    case "GET":
        _, err := io.Copy(w, variable)
        if err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "GET: %s\n", variable)
    case "POST":
        _, err := io.Copy(w, variable)
        if err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "POST: %s\n", variable)
    default:
            http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
            return
    }
}

func handleDate(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "%d-%d-%d\n", time.Now().Year(), time.Now().Month(), time.Now().Day())
}

func handleError(w http.ResponseWriter, r *http.Request){
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func handler(w http.ResponseWriter, r *http.Request){
    fmt.Println("Got one!")
}

func main(){
    http.HandleFunc("/", handler)
    http.HandleFunc("/api/v1/error", handleError)
    http.HandleFunc("/api/v1/date", handleDate)
    http.HandleFunc("/api/v1/variable", handleVariable)

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}
