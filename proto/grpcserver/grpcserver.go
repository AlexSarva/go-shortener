package grpcserver

import (
	"AlexSarva/go-shortener/constant"
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/storage"
	"AlexSarva/go-shortener/utils"
	"context"
	// импортируем пакет со сгенерированными protobuf-файлами
	pb "AlexSarva/go-shortener/proto"

	"google.golang.org/grpc/status"
)

const ShortLen int = 10

// ShortenerServer поддерживает все необходимые методы сервера.
type ShortenerServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedShortenerServer

	// используем database
	Database *app.Database
}

// GetShortURL реализует интерфейс добавления ссылки для сокращения.
func (s *ShortenerServer) GetShortURL(ctx context.Context, in *pb.BaseURLRequest) (*pb.BaseURLResponse, error) {
	var response pb.BaseURLResponse

	cfg := constant.GlobalContainer.Get("server-config").(models.Config)

	if utils.ValidateURL(in.BaseUrl.Url) {
		id := utils.ShortURLGenerator(ShortLen)
		shortURL := utils.CreateShortURL(cfg.BaseURL, id)
		dbErr := s.Database.Repo.InsertURL(id, in.BaseUrl.Url, shortURL, in.BaseUrl.UserId)
		if dbErr != nil {
			if dbErr == storage.ErrDuplicatePK {
				existShortURL, _ := s.Database.Repo.GetURLByRaw(in.BaseUrl.Url)

				response.Url = existShortURL.ShortURL
				return &response, nil
			} else {
				return nil, status.Errorf(500,
					dbErr.Error())
			}
		}

		newShortURL, _ := s.Database.Repo.GetURL(id)
		response.Url = newShortURL.ShortURL
		return &response, nil

	} else {
		return nil, status.Errorf(400,
			"check valid url please")
	}
}
