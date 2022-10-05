package utils

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkShortURLGenerator(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	length := 10
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ShortURLGenerator(length)
	}
}

func BenchmarkCreateShortUrl(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	length := 10
	path := "http://localhost:8080"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer() // останавливаем таймер
		link := ShortURLGenerator(length)
		b.StartTimer() // возобновляем таймер
		CreateShortURL(path, link)
	}
}

func BenchmarkValidateShortURL(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	length := 10
	path := "http://localhost:8080"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer() // останавливаем таймер
		id := ShortURLGenerator(length)
		link := CreateShortURL(path, id)
		b.StartTimer() // возобновляем таймер
		ValidateShortURL(link, path, length)
	}

}

func TestValidateURL(t *testing.T) {
	type args struct {
		rawText string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test #1",
			args: args{
				rawText: "yandex.ru",
			},
			want: true,
		},
		{
			name: "positive test #2",
			args: args{
				rawText: "https://www.google.com",
			},
			want: true,
		},
		{
			name: "positive test #3",
			args: args{
				rawText: "https://google",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateURL(tt.args.rawText); got != tt.want {
				t.Errorf("ValidateURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateShortURL(t *testing.T) {
	type args struct {
		rawText string
		path    string
		n       int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive test #1",
			args: args{
				rawText: "http://localhost:8080/sdfhjsdwer",
				path:    "http://localhost:8080",
				n:       10,
			},
			want: true,
		},
		{
			name: "positive test #2",
			args: args{
				rawText: "http://localhost:8080/sdfhjsdw",
				path:    "http://localhost:8080",
				n:       10,
			},
			want: false,
		},
		{
			name: "positive test #3",
			args: args{
				rawText: "http://localhost:9090/sdfhjsdwas",
				path:    "http://localhost:8080",
				n:       10,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateShortURL(tt.args.rawText, tt.args.path, tt.args.n); got != tt.want {
				t.Errorf("ValidateShortURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
