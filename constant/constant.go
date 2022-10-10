package constant

import (
	"AlexSarva/go-shortener/models"
	"errors"
	"fmt"

	"github.com/sarulabs/di"
)

var GlobalContainer di.Container

func BuildContainer(cfg models.Config) error {
	builder, _ := di.NewBuilder()
	if err := builder.Add(di.Def{
		Name:  "server-config",
		Build: func(ctn di.Container) (interface{}, error) { return cfg, nil }}); err != nil {
		return errors.New(fmt.Sprint("Ошибка инициализации контейнера", err))
	}
	GlobalContainer = builder.Build()
	return nil
}
