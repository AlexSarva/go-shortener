package constant

import (
	"AlexSarva/go-shortener/models"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/sarulabs/di"
)

var GlobalContainer di.Container

func BuildContainer(cfg models.Config) error {
	builder, _ := di.NewBuilder()
	var workIPNet *net.IPNet
	if err := builder.Add(di.Def{
		Name:  "server-config",
		Build: func(ctn di.Container) (interface{}, error) { return cfg, nil }}); err != nil {
		return errors.New(fmt.Sprint("Ошибка инициализации контейнера", err))
	}

	if cfg.TrustedSubnet != "" {
		_, ipNet, cidrErr := net.ParseCIDR(cfg.TrustedSubnet)
		if cidrErr != nil {
			log.Fatal(cidrErr)
		}
		workIPNet = ipNet
	}

	if err := builder.Add(di.Def{
		Name: "ip-net",
		Build: func(ctn di.Container) (interface{}, error) {
			return workIPNet, nil
		}}); err != nil {
		return errors.New(fmt.Sprint("Ошибка инициализации контейнера", err))
	}

	GlobalContainer = builder.Build()
	return nil
}
