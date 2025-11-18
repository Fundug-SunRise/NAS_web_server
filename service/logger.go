package service

import (
	"fmt"
	"os"
	"path"
	"time"
)

type Logger struct {
	logfile string
}

func (l *Logger) PrintINFO(content string) {
	l.writeToLog(content, "INFO")
}
func (l *Logger) PrintFATAL(content string) {
	l.writeToLog(content, "FATAL")
}

func (l *Logger) writeToLog(content, typeinfo string) {

	if l.logfile == "" {
		now := time.Now()
		l.createLogFile(now.Format("02-01-2006_15-04-05"))
	}

	file, err := os.OpenFile(l.logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %v\n", err)
		return
	}
	defer file.Close()

	now := time.Now()

	data := fmt.Sprintf("%s %s -> %s \n", now.Format("02.01.2006 15:04:05"), typeinfo, content)

	_, err = file.WriteString(data)
	if err != nil {
		fmt.Printf("Ошибка при записи: %v\n", err)
		return
	}

}

func (l *Logger) createLogFile(name string) {
	dir := path.Join("log")
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Ошибка при создании директории: %v\n", err)
		return
	}

	path := path.Join("log", name+".log")

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Ошибка при создания файла: %v\n", err)
		return
	}
	defer file.Close()

	l.logfile = path

	fmt.Printf("LogFile Назначен: %s", l.logfile)
}
