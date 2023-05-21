package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/aiviaio/go-binance/v2"

	"test_aivia/internal/config"
	async_map "test_aivia/pkg"
)

type App struct {
	httpServer *http.Server
}

func NewApp(cfg *config.Config) (*App, error) {
	return &App{
		httpServer: &http.Server{
			Addr:           ":" + cfg.AppPort,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}, nil
}

func (a *App) Run() error {

	var wg sync.WaitGroup
	ctx := context.Background()
	client := binance.NewClient("", "")
	info, err := client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return err
	}

	var symbolsSlice []string
	for _, symbol := range info.Symbols {
		if len(symbolsSlice) >= 5 {
			break
		}
		symbolsSlice = append(symbolsSlice, symbol.Symbol)
	}

	counters := async_map.NewCounters()
	ch := make(chan async_map.Counters)
	for _, symbol := range symbolsSlice {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			listPrices, err := client.NewListPricesService().Symbol(symbol).Do(ctx)
			if err != nil {
				return
			}
			for _, listPrice := range listPrices {
				counters.Store(listPrice.Symbol, listPrice.Price)
				ch <- *counters
			}
		}(symbol)
	}
	for value := range ch {
		for k, v := range value.LoadAll() {
			fmt.Println(k, v)
		}
	}

	go func() {
		if err = a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to listen and serve", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting service...")
	close(ch)
	wg.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("Shutting down server...")
	return a.httpServer.Shutdown(ctx)
}
