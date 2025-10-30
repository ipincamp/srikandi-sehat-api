package utils

import (
	"fmt"
	"io"
	"ipincamp/srikandi-sehat/config"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	InfoLogger  *log.Logger // Untuk App (prod) atau Debug (dev)
	ErrorLogger *log.Logger // Untuk Error dan Panic
	AuthLogger  *log.Logger // Khusus untuk event Autentikasi
)

// openLogFile adalah helper untuk membuka file log di direktori 'logs'
// dan membuatnya jika belum ada.
func openLogFile(filename string) *os.File {
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, 0755)
	}

	logFile, err := os.OpenFile(filepath.Join(logDir, filename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", filename, err)
	}
	return logFile
}

func InitLogger() {
	today := time.Now().Format("20060102")

	// 1. Setup Error Logger (YYYYMMDD_error.log)
	errorFilename := fmt.Sprintf("%s_error.log", today)
	errorFile := openLogFile(errorFilename)
	errorWriter := io.MultiWriter(os.Stderr, errorFile) // Tulis error ke Stderr dan file
	ErrorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)

	// 2. Setup Auth Logger (YYYYMMDD_auth.log)
	authFilename := fmt.Sprintf("%s_auth.log", today)
	authFile := openLogFile(authFilename)
	authWriter := io.MultiWriter(os.Stdout, authFile) // Tulis ke Stdout dan file
	AuthLogger = log.New(authWriter, "AUTH:  ", log.Ldate|log.Ltime|log.Lshortfile)

	// 3. Setup Info/Debug Logger berdasarkan Lingkungan
	var infoWriter io.Writer
	if config.Get("APP_ENV") == "production" {
		// Mode Produksi: YYYYMMDD_app.log
		appFilename := fmt.Sprintf("%s_app.log", today)
		appFile := openLogFile(appFilename)
		infoWriter = io.MultiWriter(os.Stdout, appFile)
		InfoLogger = log.New(infoWriter, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		// Mode Development: YYYYMMDD_debug.log
		debugFilename := fmt.Sprintf("%s_debug.log", today)
		debugFile := openLogFile(debugFilename)
		infoWriter = io.MultiWriter(os.Stdout, debugFile)
		InfoLogger = log.New(infoWriter, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func LogPanic(err interface{}) {
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, false)
	stackTrace := string(buf[:stackSize])

	ErrorLogger.Printf("Panic recovered: %v\nStack trace:\n%s", err, stackTrace)
}
