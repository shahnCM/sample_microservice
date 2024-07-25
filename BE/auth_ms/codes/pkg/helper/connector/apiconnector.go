package connector

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ConnectorDto struct {
	Https   bool   `json:"https"`
	Host    string `json:"host"`
	Proxy   string `json:"proxy"`
	Port    string `json:"port"`
	Path    string `json:"path"`
	ApiUrl  string `json:"api_url"`
	AuthKey string `json:"auth_key"`
}

func Connector(msName string, path string) (*ConnectorDto, error) {
	msName = strings.ToLower(msName)

	switch msName {

	case "card":
		return cardMsConnector(path)

	case "catalog":
		return catalogMsConnector(path)

	case "commonui":
		return commonuiMsConnector(path)

	default:
		return nil, fiber.ErrUnprocessableEntity
	}

}

func cardMsConnector(path string) (*ConnectorDto, error) {
	host := os.Getenv("HOST_CARD_MS")
	port := os.Getenv("PORT_CARD_MS")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8001"
	}

	return &ConnectorDto{
		Host:    host,
		Port:    port,
		Path:    path,
		ApiUrl:  host + ":" + port + "/" + path,
		AuthKey: "",
	}, nil
}

func catalogMsConnector(path string) (*ConnectorDto, error) {
	host := os.Getenv("HOST_CATALOG_MS")
	port := os.Getenv("PORT_CATALOG_MS")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8001"
	}

	return &ConnectorDto{
		Host:    host,
		Port:    port,
		Path:    path,
		ApiUrl:  host + ":" + port + "/" + path,
		AuthKey: "",
	}, nil
}

func commonuiMsConnector(path string) (*ConnectorDto, error) {
	host := os.Getenv("HOST_COMMONUI_MS")
	port := os.Getenv("PORT_COMMONUI_MS")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8001"
	}

	return &ConnectorDto{
		Host:    host,
		Port:    port,
		Path:    path,
		ApiUrl:  host + ":" + port + "/" + path,
		AuthKey: "",
	}, nil
}
