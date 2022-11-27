package grpcserver

import (
	"AlexSarva/go-shortener/constant"
	"AlexSarva/go-shortener/internal/app"
	"AlexSarva/go-shortener/models"
	pb "AlexSarva/go-shortener/proto"
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/caarlos0/env/v6"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestShortenerServer_Batch(t *testing.T) {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Println(GlobalContainerErr)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)
	database := app.NewStorage()
	delCh := make(chan models.DeleteURL)
	md := metadata.New(map[string]string{"user_id": "6967689f-f072-4607-b720-114cea99293f"})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	type fields struct {
		UnimplementedShortenerServer pb.UnimplementedShortenerServer
		Database                     *app.Database
		DelChan                      chan models.DeleteURL
	}
	type args struct {
		ctx context.Context
		in  *pb.CorrelationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.CorrelationResponse
		wantErr bool
	}{
		{
			name: "BATCH positive test #1",
			fields: fields{
				Database: database,
				DelChan:  delCh,
			},
			args: args{
				ctx: ctx,
				in: &pb.CorrelationRequest{
					OriginalUrls: []*pb.OriginalUrlElement{
						{CorrelationId: "1", OriginalUrl: "github2.com"},
						{CorrelationId: "2", OriginalUrl: "googlasdecom"},
						{CorrelationId: "3", OriginalUrl: "yahoo2.com"},
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortenerServer{
				UnimplementedShortenerServer: tt.fields.UnimplementedShortenerServer,
				Database:                     tt.fields.Database,
				DelChan:                      tt.fields.DelChan,
			}

			_, err := s.Batch(tt.args.ctx, tt.args.in)

			if (err != nil) != tt.wantErr {
				t.Errorf("Batch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestShortenerServer_Delete(t *testing.T) {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Println(GlobalContainerErr)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)
	database := app.NewStorage()
	delCh := make(chan models.DeleteURL)
	md := metadata.New(map[string]string{"user_id": "6967689f-f072-4607-b720-114cea99293f"})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	type fields struct {
		UnimplementedShortenerServer pb.UnimplementedShortenerServer
		Database                     *app.Database
		DelChan                      chan models.DeleteURL
	}
	type args struct {
		ctx context.Context
		in  *pb.DeleteRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.DeleteResponse
		wantErr bool
	}{
		{
			name: "DELETE positive test #1",
			fields: fields{
				Database: database,
				DelChan:  delCh,
			},
			args: args{
				ctx: ctx,
				in: &pb.DeleteRequest{
					ShortUrls: []*pb.ShortURL{
						{Url: "http://localhost:8090/EgSfzXZERY"},
						{Url: "http://localhost:8090/ObAnaGyLJE"},
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortenerServer{
				UnimplementedShortenerServer: tt.fields.UnimplementedShortenerServer,
				Database:                     tt.fields.Database,
				DelChan:                      tt.fields.DelChan,
			}
			_, err := s.Delete(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestShortenerServer_GetAllURLs(t *testing.T) {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Println(GlobalContainerErr)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)
	database := app.NewStorage()
	delCh := make(chan models.DeleteURL)
	md := metadata.New(map[string]string{"user_id": "6967689f-f072-4607-b720-114cea99293f"})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	type fields struct {
		UnimplementedShortenerServer pb.UnimplementedShortenerServer
		Database                     *app.Database
		DelChan                      chan models.DeleteURL
	}
	type args struct {
		ctx context.Context
		in  *emptypb.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.AllURLsResponse
		wantErr bool
	}{
		{
			name: "GetAllURLs positive test #1",
			fields: fields{
				Database: database,
				DelChan:  delCh,
			},
			args: args{
				ctx: ctx,
				in:  &emptypb.Empty{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortenerServer{
				UnimplementedShortenerServer: tt.fields.UnimplementedShortenerServer,
				Database:                     tt.fields.Database,
				DelChan:                      tt.fields.DelChan,
			}
			_, err := s.GetAllURLs(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestShortenerServer_GetOriginalURL(t *testing.T) {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Println(GlobalContainerErr)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)
	database := app.NewStorage()
	delCh := make(chan models.DeleteURL)
	md := metadata.New(map[string]string{"user_id": "6967689f-f072-4607-b720-114cea99293f"})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	type fields struct {
		UnimplementedShortenerServer pb.UnimplementedShortenerServer
		Database                     *app.Database
		DelChan                      chan models.DeleteURL
	}
	type args struct {
		ctx context.Context
		in  *pb.OriginalURLRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.OriginalURLResponse
		wantErr bool
	}{
		{
			name: "GetShortURL positive test #1",
			fields: fields{
				Database: database,
				DelChan:  delCh,
			},
			args: args{
				ctx: ctx,
				in: &pb.OriginalURLRequest{
					ShortUrl: &pb.ShortURL{
						Url: "http://localhost:8080/wREIXHPdEv",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortenerServer{
				UnimplementedShortenerServer: tt.fields.UnimplementedShortenerServer,
				Database:                     tt.fields.Database,
				DelChan:                      tt.fields.DelChan,
			}
			_, err := s.GetOriginalURL(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOriginalURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestShortenerServer_GetShortURL(t *testing.T) {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Println(GlobalContainerErr)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)
	database := app.NewStorage()
	delCh := make(chan models.DeleteURL)
	md := metadata.New(map[string]string{"user_id": "6967689f-f072-4607-b720-114cea99293f"})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	type fields struct {
		UnimplementedShortenerServer pb.UnimplementedShortenerServer
		Database                     *app.Database
		DelChan                      chan models.DeleteURL
	}
	type args struct {
		ctx context.Context
		in  *pb.ShortURLRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.ShortURLResponse
		wantErr bool
	}{
		{
			name: "GetShortURL positive test #1",
			fields: fields{
				Database: database,
				DelChan:  delCh,
			},
			args: args{
				ctx: ctx,
				in: &pb.ShortURLRequest{
					OriginalUrl: &pb.OriginalURL{
						Url: "github.com",
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortenerServer{
				UnimplementedShortenerServer: tt.fields.UnimplementedShortenerServer,
				Database:                     tt.fields.Database,
				DelChan:                      tt.fields.DelChan,
			}
			_, err := s.GetShortURL(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShortURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestShortenerServer_Ping(t *testing.T) {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Println(GlobalContainerErr)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)
	database := app.NewStorage()
	delCh := make(chan models.DeleteURL)
	md := metadata.New(map[string]string{"user_id": "6967689f-f072-4607-b720-114cea99293f"})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	type fields struct {
		UnimplementedShortenerServer pb.UnimplementedShortenerServer
		Database                     *app.Database
		DelChan                      chan models.DeleteURL
	}
	type args struct {
		ctx context.Context
		in  *emptypb.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.PingResponse
		wantErr bool
	}{
		{
			name: "Ping positive test #1",
			fields: fields{
				Database: database,
				DelChan:  delCh,
			},
			args: args{
				ctx: ctx,
				in:  &emptypb.Empty{},
			},
			want:    &pb.PingResponse{Status: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortenerServer{
				UnimplementedShortenerServer: tt.fields.UnimplementedShortenerServer,
				Database:                     tt.fields.Database,
				DelChan:                      tt.fields.DelChan,
			}
			got, err := s.Ping(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ping() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortenerServer_Stata(t *testing.T) {
	var cfg models.Config
	// Приоритет будет у ФЛАГОВ
	// Загружаем конфиг из переменных окружения
	err := env.Parse(&cfg)
	GlobalContainerErr := constant.BuildContainer(cfg)
	if GlobalContainerErr != nil {
		log.Println(GlobalContainerErr)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", cfg)
	database := app.NewStorage()
	delCh := make(chan models.DeleteURL)
	md := metadata.New(map[string]string{"user_id": "6967689f-f072-4607-b720-114cea99293f"})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	type fields struct {
		UnimplementedShortenerServer pb.UnimplementedShortenerServer
		Database                     *app.Database
		DelChan                      chan models.DeleteURL
	}
	type args struct {
		ctx context.Context
		in  *emptypb.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.StataResponse
		wantErr bool
	}{
		{
			name: "Stata positive test #1",
			fields: fields{
				Database: database,
				DelChan:  delCh,
			},
			args: args{
				ctx: ctx,
				in:  &emptypb.Empty{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortenerServer{
				UnimplementedShortenerServer: tt.fields.UnimplementedShortenerServer,
				Database:                     tt.fields.Database,
				DelChan:                      tt.fields.DelChan,
			}
			_, err := s.Stata(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
