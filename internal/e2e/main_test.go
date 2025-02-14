package e2e

import (
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/wDRxxx/avito-shop/internal/config"
	"github.com/wDRxxx/avito-shop/internal/e2e/models"
)

var users []models.User
var apiURL string

func init() {
	e2eConfigPath := "../../e2e.env"

	err := config.Load(e2eConfigPath)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	apiURL = "http://" + net.JoinHostPort(os.Getenv("HTTP_HOST"), os.Getenv("HTTP_PORT"))
	fmt.Println(apiURL)
}

// Runs like this, to have data for tests
func TestE2E(t *testing.T) {
	t.Run("TestAuth-1", TAuth)
	t.Run("TestAuth-2", TAuth)

	t.Run("TestBuy", TBuy)
	t.Run("TestSendCoin", TSendCoin)
	t.Run("TestInfo", TInfo)
}
