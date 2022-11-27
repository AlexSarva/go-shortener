package grpcserver

import (
	"AlexSarva/go-shortener/constant"
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	pb "AlexSarva/go-shortener/proto"
	"AlexSarva/go-shortener/storage"
	"AlexSarva/go-shortener/utils"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const ShortLen int = 10

// ShortenerServer поддерживает все необходимые методы сервера.
type ShortenerServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedShortenerServer

	// используем database
	Database *app.Database
	// канал для удаления
	DelChan chan models.DeleteURL
}

// GetShortURL реализует интерфейс добавления ссылки для сокращения.
func (s *ShortenerServer) GetShortURL(ctx context.Context, in *pb.ShortURLRequest) (*pb.ShortURLResponse, error) {
	var response pb.ShortURLResponse
	userID, userIDErr := utils.GetUserID(ctx)
	if userIDErr != nil {
		userID = uuid.NewString()
	}

	log.Println("USER ID: ", userID)

	cfg := constant.GlobalContainer.Get("server-config").(models.Config)

	if !utils.ValidateURL(in.OriginalUrl.Url) {
		return nil, status.Errorf(codes.InvalidArgument,
			"check valid url please")
	}
	id := utils.ShortURLGenerator(ShortLen)
	shortURL := utils.CreateShortURL(cfg.BaseURL, id)
	dbErr := s.Database.Repo.InsertURL(id, in.OriginalUrl.Url, shortURL, userID)
	if dbErr != nil {
		if dbErr == storage.ErrDuplicatePK {
			existShortURL, _ := s.Database.Repo.GetURLByRaw(in.OriginalUrl.Url)

			response.ShortUrl = &pb.ShortURL{Url: existShortURL.ShortURL}
			return &response, nil
		} else {
			return nil, status.Errorf(500,
				dbErr.Error())
		}
	}

	newShortURL, _ := s.Database.Repo.GetURL(id)
	response.ShortUrl = &pb.ShortURL{Url: newShortURL.ShortURL}
	log.Printf("%+v\n", response.ShortUrl)
	return &response, nil

}

// GetOriginalURL реализует интерфейс получения исходной ссылки.
func (s *ShortenerServer) GetOriginalURL(ctx context.Context, in *pb.OriginalURLRequest) (*pb.OriginalURLResponse, error) {
	var response pb.OriginalURLResponse

	//cfg := constant.GlobalContainer.Get("server-config").(models.Config)

	splittedShortURL := strings.Split(in.ShortUrl.Url, "/")
	id := splittedShortURL[len(splittedShortURL)-1]

	res, er := s.Database.Repo.GetURL(id)
	if er != nil {
		return nil, status.Errorf(400,
			"No such short url in DB")
	}
	if res.Deleted == 1 {
		return nil, status.Errorf(codes.NotFound,
			"this short url have been deleted")
	}

	response.OriginalUrl = &pb.OriginalURL{Url: res.RawURL}

	return &response, nil

}

// Ping реализует интерфейс проверки доступа к БД
func (s *ShortenerServer) Ping(ctx context.Context, in *emptypb.Empty) (*pb.PingResponse, error) {
	ping := s.Database.Repo.Ping()
	if ping {
		return &pb.PingResponse{Status: true}, nil
	}
	return nil, status.Errorf(codes.Internal,
		"Internal server error")
}

// GetAllURLs реализует интерфейс получения всех сокращенных ссылок пользователя
func (s *ShortenerServer) GetAllURLs(ctx context.Context, in *emptypb.Empty) (*pb.AllURLsResponse, error) {
	var response pb.AllURLsResponse
	var responseURLs []*pb.UserURL

	userID, userIDErr := utils.GetUserID(ctx)
	if userIDErr != nil {
		return nil, status.Errorf(codes.Unauthenticated,
			"no user_id in request")
	}

	log.Printf("Ищем urls для пользователя %s", userID)
	res, er := s.Database.Repo.GetUserURLs(userID)
	if er != nil {
		if errors.Is(er, storage.ErrNoValues) {
			return nil, status.Errorf(codes.NotFound,
				"no data found in DB")
		}
		return nil, status.Errorf(codes.Internal,
			"Internal server error")
	}

	for _, userURL := range res {
		responseURLs = append(responseURLs, &pb.UserURL{
			ShortUrl:    userURL.ShortURL,
			OriginalUrl: userURL.RawURL,
		})
	}
	log.Printf("%+v\n", res)
	response.UserUrls = responseURLs
	return &response, nil
}

// Batch реализует интерфейс получения множества URL для сокращения
func (s *ShortenerServer) Batch(ctx context.Context, in *pb.CorrelationRequest) (*pb.CorrelationResponse, error) {
	var response pb.CorrelationResponse
	var responseURLs []*pb.ShortUrlElement
	var insertBatchURL []models.URL
	cfg := constant.GlobalContainer.Get("server-config").(models.Config)

	userID, userIDErr := utils.GetUserID(ctx)
	if userIDErr != nil {
		return nil, status.Errorf(codes.Unauthenticated,
			"no user_id in request")
	}

	for _, urlInfo := range in.OriginalUrls {
		id := utils.ShortURLGenerator(ShortLen)
		shortURL := utils.CreateShortURL(cfg.BaseURL, id)
		if utils.ValidateURL(urlInfo.OriginalUrl) {
			currentURLInsert := models.URL{
				ID:       id,
				RawURL:   urlInfo.OriginalUrl,
				ShortURL: shortURL,
				Created:  time.Now(),
				UserID:   userID,
			}
			currentURLResult := pb.ShortUrlElement{
				CorrelationId: urlInfo.CorrelationId,
				ShortUrl:      shortURL,
			}
			insertBatchURL = append(insertBatchURL, currentURLInsert)
			responseURLs = append(responseURLs, &currentURLResult)

		} else {
			currentURLResult := pb.ShortUrlElement{
				CorrelationId: urlInfo.CorrelationId,
				ShortUrl:      "not valid url",
			}
			responseURLs = append(responseURLs, &currentURLResult)
		}
	}

	dbErr := s.Database.Repo.InsertMany(insertBatchURL)
	if dbErr != nil {
		log.Println(dbErr)
	}

	response.ShortUrls = responseURLs
	return &response, nil
}

// Delete реализует интерфейс получения списка идентификаторов сокращённых URL для удаления
func (s *ShortenerServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	var response pb.DeleteResponse
	var deleteBatchURL []string

	userID, userIDErr := utils.GetUserID(ctx)
	if userIDErr != nil {
		return nil, status.Errorf(codes.Unauthenticated,
			"no user_id in request")
	}

	for _, url := range in.ShortUrls {
		splittedShortURL := strings.Split(url.Url, "/")
		id := splittedShortURL[len(splittedShortURL)-1]
		deleteBatchURL = append(deleteBatchURL, id)
	}

	go utils.AddDeleteURLs(models.DeleteURL{
		UserID: userID,
		URLs:   deleteBatchURL,
	}, s.DelChan)

	response.Status = true

	return &response, nil
}

// Stata реализует интерфейс получения статистики количетсва постов и пользователей
func (s *ShortenerServer) Stata(ctx context.Context, in *emptypb.Empty) (*pb.StataResponse, error) {
	var response pb.StataResponse

	res, er := s.Database.Repo.GetStat()
	if er != nil {
		return nil, status.Errorf(codes.Internal,
			"internal server error")
	}

	stata := pb.Stata{
		UrlsCnt:  int32(res.URLsCnt),
		UsersCnt: int32(res.UsersCnt),
	}

	response.Stata = &stata
	return &response, nil

}
