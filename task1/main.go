package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

const (
	urlAPI = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1"
	update = 10 * time.Minute
)

type Coin struct { //можно еще добавить fields, но посчитала что эти самые основные.
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Price  float64 `json:"current_price"`
}

func main() {
	symbol := flag.String("symbol", "", "Symbol of the specified cryptocurrency") //Добавить возможность получать курс для определенной криптовалюты. Я решила показать через flags где рядом мы будем указывать нужный нам символ. Если не указать, получим просто все крипто данные.
	flag.Parse()                                                                  //разбор аргументов командной строки.
	updateCrypto(*symbol)
}

func getCryptoData() ([]Coin, error) { //функция чтобы получить данные о криптовалютах. Сначала отправляет запрос к API, получает нужные данные и декодирует из JSON.
	response, err := http.Get(urlAPI) //Сначала отправляет HTTP запрос к API, тут так же resp.Body открывается для чтения. Получает ответ и идет обработка ошибки если возникнет.
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() //Когда отправили запрос, открылся resp.Body, тут важно его закрывать и освободить ресурсы. Сразу после завершения getCryptoData вызывается эта функция, так как defer отложенная функция.

	var coins []Coin                                    //coins будет хранить данные о cryptocurrency
	err = json.NewDecoder(response.Body).Decode(&coins) //используя пакет мы декодируем данные в формате JSON в структуру coins. Это восстановливает информацию к её изначальному виду, только после этого программа может работать напрямую
	if err != nil {
		return nil, err
	}
	return coins, nil
}

func updateCrypto(symbol string) { //чтобы обновлять данные о cryptocurrency. А параметр symbol чтобы определить конкретную криптовалюту если понадобится.
	for {
		coins, err := getCryptoData() //вызываем функцию чтобы получить данные
		if err != nil {
			fmt.Println("Error when retrieving cryptocurrency data ", err)
			continue //если даже возникнет ошибка, в след итерации мы заново попытаемся получить данные
		}

		if symbol != "" { //если указать и symbol будет не пустой возвращает нужный символ
			found := false //перед каждым новым поиском криптовалюта не найдена
			for _, coin := range coins {
				if coin.Symbol == symbol {
					found = true //когда мы уже нашли криптовалюту с нужным символом
					fmt.Printf("Cryptocurrency: %s (%s), Current price: %.0f\n", coin.Name, coin.Symbol, coin.Price)
					break
				}
			}
			if !found { //выполняется если указанный symbol не найден
				fmt.Println("Cryptocurrency with symbol", symbol, "not found")
			}
		} else { //если ничего не указать и symbol будет пустой возвращает все данные
			for _, coin := range coins {
				fmt.Printf("Cryptocurrency: %s (%s), Current price: %.0f\n", coin.Name, coin.Symbol, coin.Price)
			}
		}
		time.Sleep(update) //чтобы обновлять не чаще чем раз в 10 минут
	}
}

//example: go run main.go -symbol=btc
