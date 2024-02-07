package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func (s *LogServiceImpl) Write(ctx context.Context, msg string) error {

	f, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("ошибка открытия файла: %s\n", err)
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%v - %s \n", time.Now(), msg)); err != nil {
		log.Printf("ошибка записи лога: %s\n", err)
		return err
	}

	return nil
}
