package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	address := flag.String("address", ":3333", "Сетевой адресс HTTP")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	server := &http.Server{
		Addr:     *address,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Println("Server is listening on", *address)
	err := server.ListenAndServe()
	app.errorLog.Fatal(err)
}
