package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"
)

type Value struct {
	s    string
	kind string
}

type Storage struct {
	Inner  map[string]Value `json:"inner"`
	Sql    map[string][]int `json:"sqlmap"`
	logger *zap.Logger
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}

	defer logger.Sync()
	logger.Info("created new storage")

	return Storage{
		Inner:  make(map[string]Value),
		Sql:    make(map[string][]int),
		logger: logger,
	}, nil
}

func reverseInts(input []int) []int {
	if len(input) == 0 {
		return input
	}
	return append(reverseInts(input[1:]), input[0])
}

func (r Storage) LPUSH(key string, val []int) (err string) {
	fmt.Println("Start LLLLPush")
	for i := range val {
		_, ok := r.Sql[key]
		if !ok {
			r.Sql[key] = []int{val[i]}

		} else {
			temp := reverseInts(r.Sql[key])
			r.Sql[key] = reverseInts(append(temp, val[i]))

		}
		fmt.Println(r.Sql[key])

	}
	fmt.Println(r.Sql[key])
	return ""

}

func (r Storage) RPUSH(key string, val []int) (err string) {
	fmt.Println("Start RRRRPush")
	for i := range val {
		_, ok := r.Sql[key]
		if !ok {
			r.Sql[key] = []int{val[i]}

		} else {

			r.Sql[key] = append(r.Sql[key], val[i])

		}
		fmt.Println(r.Sql[key])

	}
	fmt.Println(r.Sql[key])
	return ""

}

func (r Storage) LPOP(key string, val []int) (toPopV []int, err string) {

	toPop := []int{}

	if len(val) == 0 {

		toPop = append(toPop, r.Sql[key][0])
		r.Sql[key] = r.Sql[key][1:]
	} else if len(val) == 1 {
		if val[0] >= len(r.Sql[key]) {

			toPop = append(toPop, r.Sql[key]...)
			r.Sql[key] = []int{}
			return toPop, ""
		}

		toPop = r.Sql[key][:val[0]]

		r.Sql[key] = r.Sql[key][val[0] : len(r.Sql[key])-1]

	} else {

		toPop = append(toPop, r.Sql[key][val[0]:val[1]+1]...)

		r.Sql[key] = append(r.Sql[key][0:val[0]], r.Sql[key][val[1]+1:len(r.Sql[key])]...)

	}

	return toPop, ""

}

func (r Storage) RADDTOSET(key string, val []int) (err string) {
	fmt.Println("Start RRRRPush")
	for i := range val {
		_, ok := r.Sql[key]
		if !ok {
			r.Sql[key] = []int{val[i]}

		} else {
			flag := false
			for _, v := range r.Sql[key] {
				if v == val[i] {
					flag = true
					break
				}
			}
			if !flag {
				r.Sql[key] = append(r.Sql[key], val[i])
			}

		}

	}
	fmt.Println(r.Sql[key])
	return ""
}
func (r Storage) LSET(key string, ind int, elem int) (err string) {
	_, ok := r.Sql[key]
	if ok {

		if len(r.Sql[key]) > ind {
			r.Sql[key][ind] = elem

			return fmt.Sprintf("(integer) %d", elem)
		} else {
			return "Index out of range"
		}
	}
	return "Index out of range1"
}

func (r Storage) LGET(key string, ind int) (res int, err string) {
	_, ok := r.Sql[key]
	if ok {
		fmt.Println("r.sql[key]")
		fmt.Println(len(r.Sql[key]))
		if len(r.Sql[key]) > ind {
			fmt.Println(r.Sql[key])
			return r.Sql[key][ind], ""
		} else {
			return -1, "Index out of range"
		}
	}
	return -1, "Index out of range1"
}

func (r Storage) SaveToJson() {
	jsonData, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Ошибка при маршализации:", err)
		return
	}

	file, err := os.Create("storage.json")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close() // Закрываем файл после завершения работы

	// Запись JSON в файл
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}

	fmt.Println(string(jsonData)) // Вывод JSON строки
}

func (r Storage) ReadJson() {
	// Открытие файла
	file, err := os.Open("storage.json")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close() // Закрываем файл после завершения работы

	n := 1
	fmt.Println(n)
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&r)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	// Вывод прочитанных данных
	fmt.Println(r)
}

func (r Storage) Set(key, value string) {
	switch GetType(value) {
	case "D":
		r.Inner[key] = Value{s: value, kind: "D"}
	case "Fl64":
		r.Inner[key] = Value{s: value, kind: "Fl64"}
	case "S":
		r.Inner[key] = Value{s: value, kind: "S"}
	}

	r.logger.Info("key set")
	r.logger.Sync()
}

func (r Storage) Get(key string) *string {
	res, ok := r.Inner[key]
	if !ok {
		return nil
	}

	return &res.s
}

func GetType(value string) string {
	if _, err := strconv.Atoi(value); err == nil {
		return "D"
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "Fl64"
	}
	return "S"
}
