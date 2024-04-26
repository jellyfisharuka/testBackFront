package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type InstaInfluencer struct { //структура которая представляет данные об Influencer
	Rank                string
	Name                string
	Username            string
	Category            string
	Subscribers         string
	Audience            string
	AuthenticEngagement string
	TotalEngagement     string
}

// эта функция чтобы создать новый экземпляр структуры  InstaInfluencer
func NewInstaInfluencer(rank, name, username, category, subscribers, audience, authentic, engagement string) *InstaInfluencer {
	return &InstaInfluencer{
		Rank:                rank,
		Name:                name,
		Username:            username,
		Category:            category,
		Subscribers:         subscribers,
		Audience:            audience,
		AuthenticEngagement: authentic,
		TotalEngagement:     engagement,
	}
}

// Эта функция парсит данные из элемента goquery.Selection в InstagramInfluencer. Возвращает уже NewInstagramInfluencer
func ParseInfluencer(item *goquery.Selection) *InstaInfluencer {
	rank := item.Find(".rank").Text()
	name := item.Find(".contributor__title").Text()
	username := item.Find(".contributor__name-content").Text()
	category := item.Find(".category").Text()
	subscribers := item.Find(".subscribers").Text()
	audience := item.Find(".audience").Text()
	authentic := item.Find(".authentic").Text()
	engagement := item.Find(".engagement").Text()

	return NewInstaInfluencer(rank, name, username, category, subscribers, audience, authentic, engagement)
}

// Эта функция чтобы записывать данные о InstaInfluencer
func (i *InstaInfluencer) WriteToCSV(writer *csv.Writer) {
	writer.Write([]string{
		i.Rank,
		i.Name,
		i.Username,
		i.Category,
		i.Subscribers,
		i.Audience,
		i.AuthenticEngagement,
		i.TotalEngagement,
	})
}

func main() {
	url := "https://hypeauditor.com/top-instagram-all-russia/"

	response, err := http.Get(url) //Сначала отправляет HTTP запрос, тут так же resp.Body открывается для чтения.
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close() //Когда отправили запрос, открылся resp.Body, тут важно его закрывать и освободить ресурсы.

	if response.StatusCode != 200 { //если запрос не успешный то status code error
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body) //новый документ из HTML страницы
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("users.csv") //новый CSV файл и там будем сохранять наши данные
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file) //для записи данных создается обьект csv.Writer
	defer writer.Flush()          //метод Flush нужна чтобы убедиться, что все данные точно были записаны.

	doc.Find(".row").Each(func(index int, item *goquery.Selection) { //итерация по классу "row"
		//influencer := ParseInfluencer(item)
		if index == 0 { //первый элемент идет как заголовок в CSV
			writer.Write([]string{"Rank", "Name", "Username", "Category", "Subscribers", "Audience", "Authentic Engagement", "Total Engagement"})
		}
		influencer := ParseInfluencer(item) //вызываем функцию для парсинга
		influencer.WriteToCSV(writer)       //вызываем функцию чтобы записывать данные в CSV
	})

	fmt.Println("Data has been written to CSV.")
}
