package main

import (
	pb "AlexSarva/go-shortener/proto"
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	id := uuid.NewString()
	urls := []*pb.BaseURL{
		{Url: "google.com", UserId: id},
		{Url: "googlecom", UserId: id},
		{Url: "yandex.com", UserId: id},
	}
	for _, url := range urls {
		// добавляем пользователей
		resp, err := c.GetShortURL(context.Background(), &pb.BaseURLRequest{
			BaseUrl: url,
		})
		if err != nil {
			log.Println("Ошибка: ", err)
		} else {
			log.Printf("%+v\n", resp)
		}
		//if (&pb.UserURLResponse{}) != resp {
		//	log.Printf("%+v\n", resp)
		//}
	}

}
