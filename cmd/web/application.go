package main

import (
	"log"
	"sync"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	mutex    sync.RWMutex
}
