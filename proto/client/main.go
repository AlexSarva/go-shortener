package main

import (
	pb "AlexSarva/go-shortener/proto"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewShortenerClient(conn)

	// функция, в которой будем отправлять сообщения
	TestUsers(c)
}

func TestUsers(c pb.ShortenerClient) {
	// набор тестовых данных
	//id := uuid.NewString()
	//md := metadata.New(map[string]string{"user_id": uuid.NewString()})
	md := metadata.New(map[string]string{"user_id": "6967689f-f072-4607-b720-114cea99293f"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	//ctx := context.Background()
	urls := []*pb.OriginalURL{
		{Url: "github.com"},
		{Url: "googlecom"},
		{Url: "yahoo.com"},
	}
	var shortURLs []*pb.ShortURL
	log.Println("\nTEST: GetShortURL")
	for _, url := range urls {
		// добавляем пользователей
		resp, err := c.GetShortURL(ctx, &pb.ShortURLRequest{
			OriginalUrl: url,
		})
		if err != nil {
			log.Println("Ошибка: ", err)
		} else {
			log.Printf("%+v\n", resp)
			shortURLs = append(shortURLs, resp.ShortUrl)
		}
		//if (&pb.UserURLResponse{}) != resp {
		//	log.Printf("%+v\n", resp)
		//}
	}

	log.Println("\nTEST: GetOriginalURL")
	for _, shortURL := range shortURLs {
		//log.Println(shortURL)
		resp, err := c.GetOriginalURL(ctx, &pb.OriginalURLRequest{
			ShortUrl: shortURL,
		})
		if err != nil {
			log.Println("Ошибка: ", err)
		} else {
			log.Printf("%+v\n", resp)

		}
	}
	log.Println("\nTEST: Ping")
	resp, err := c.Ping(ctx, &emptypb.Empty{})
	if err != nil {
		log.Println("Ошибка: ", err)
	} else {
		log.Printf("%+v\n", resp)
	}

	log.Println("\nTEST: GetAllURLs")
	resp2, err2 := c.GetAllURLs(ctx, &emptypb.Empty{})
	if err2 != nil {
		log.Println("Ошибка: ", err2)
	} else {
		for _, elem := range resp2.UserUrls {
			log.Printf("%+v\n", elem)
		}
		//log.Printf("%+v\n", resp2)
	}

	log.Println("\nTEST: BATCH")
	correlationUrls := []*pb.OriginalUrlElement{
		{CorrelationId: "1", OriginalUrl: "github2.com"},
		{CorrelationId: "2", OriginalUrl: "googlasdecom"},
		{CorrelationId: "3", OriginalUrl: "yahoo2.com"},
	}

	resp3, err3 := c.Batch(ctx, &pb.CorrelationRequest{
		OriginalUrls: correlationUrls,
	})
	if err3 != nil {
		log.Println("Ошибка: ", err3)
	} else {
		for _, elem := range resp3.ShortUrls {
			log.Printf("%+v\n", elem)
		}
		//log.Printf("%+v\n", resp2)
	}

	log.Println("\nTEST: DELETE")
	urlsDel := []*pb.ShortURL{
		{Url: "http://localhost:8090/EgSfzXZERY"},
		{Url: "http://localhost:8090/ObAnaGyLJE"},
	}
	resp4, err4 := c.Delete(ctx, &pb.DeleteRequest{
		ShortUrls: urlsDel,
	})
	if err4 != nil {
		log.Println("Ошибка: ", err4)
	} else {
		log.Printf("%+v\n", resp4)
	}

	log.Println("\nTEST: STATA")
	resp5, err5 := c.Stata(ctx, &emptypb.Empty{})
	if err5 != nil {
		log.Println("Ошибка: ", err5)
	} else {
		log.Printf("%+v\n", resp5)
	}
}
