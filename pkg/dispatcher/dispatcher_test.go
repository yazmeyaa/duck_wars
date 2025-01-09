package dispatcher_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/yazmeyaa/duck_wars/internal/models"
	"github.com/yazmeyaa/duck_wars/pkg/config"
	"github.com/yazmeyaa/duck_wars/pkg/database"
	"github.com/yazmeyaa/duck_wars/pkg/dispatcher"
	"github.com/yazmeyaa/duck_wars/pkg/sheduler"
)

func TestDispatcher(t *testing.T) {
	stop, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := config.New("../../.env")
	db, err := database.New(cfg)
	if err != nil {
		logger.Error(err.Error())
		t.FailNow()
		return
	}
	defer db.Close()

	sheduler := sheduler.NewSheduler(db, logger)

	_, err = sheduler.SheduleTask(models.Task{
		Command:    "finish_adventure",
		Payload:    "forest 1233321",
		SheduledAt: time.Now().Add(time.Second),
	})

	if err != nil {
		logger.Error(err.Error())
		t.FailNow()
		return
	}
	d := dispatcher.NewDispatcher(stop, logger, db, 5)

	go d.Run()

	<-time.After(time.Second * 5)

}
