package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Store は店舗情報を保持する構造体
type Store struct {
	Name  string `json:"name"`
	Genre string `json:"genre"`
	URL   string `json:"url"`
}

var stores []Store
var correctFoods = []string{"鰻", "かき氷", "タイ料理", "和食"}

// init関数はパッケージの初期化時に一度だけ実行される
func init() {
	// 乱数のシードを設定
	rand.Seed(time.Now().UnixNano())

	// JSONファイルから店舗データを読み込む
	jsonFile, err := os.Open("json/food_stores.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &stores)
}

// ProcessMenuStep はフォームの入力値を受け取り、次に表示すべき内容を決定して返す
func ProcessMenuStep(params map[string]string) map[string]interface{} {
	step, ok := params["step"]
	// stepがなければ最初の質問画面を表示
	if !ok {
		return map[string]interface{}{"CurrentStep": "q1"}
	}

	switch step {
	case "q1":
		return map[string]interface{}{
			"CurrentStep": "q2",
			"Q1":          params["q1"],
		}

	case "q2":
		q1 := params["q1"]
		q2 := params["q2"]
		return map[string]interface{}{
			"CurrentStep": "q3",
			"Q1":          q1,
			"Q2":          q2,
			"Foods":       getFoods(q1, q2),
		}

	case "q3":
		q1Str := params["q1"]
		q2Str := params["q2"]
		q3Str := params["q3"]

		q1, _ := strconv.Atoi(q1Str)
		q2, _ := strconv.Atoi(q2Str)
		q3 := -1
		if q3Str == "鰻" {
			q3 = 0
		} else if q3Str == "かき氷" {
			q3 = 1
		} else if q3Str == "タイ料理" {
			q3 = 2
		} else if q3Str == "和食" {
			q3 = 3
		}

		idx := q1*16 + q2*4 + q3

		// かき氷ルートの場合
		if q3Str == "かき氷" {
			return map[string]interface{}{
				"CurrentStep": "ice_flavor",
				"Index":       idx,
				"Flavors":     []string{"いちご", "メロン", "ブルーハワイ", "オレンジ"},
			}
		}

		// 通常の結果表示
		var recommendedStore Store
		for _, s := range stores {
			if s.Genre == q3Str {
				recommendedStore = s
				break
			}
		}
		return map[string]interface{}{
			"CurrentStep": "result",
			"Index":       idx,
			"Result":      recommendedStore.Name + "のお店をオススメします！",
			"StoreURL":    recommendedStore.URL,
		}

	case "ice_flavor":
		flavor := params["flavor"]
		idxStr := params["idx"]
		return map[string]interface{}{
			"CurrentStep": "result",
			"Index":       idxStr,
			"Result":      "おめでとうございます！" + flavor + "かき氷が注文されました！",
		}
	}

	// 不正なステップの場合は最初の質問に戻す
	return map[string]interface{}{"CurrentStep": "q1"}
}

// getFoods は質問3の選択肢を決定するヘルパー関数
func getFoods(q1, q2 string) []string {
	// 質問1が「サプライズが欲しい」かつ質問2が「誰かと楽しく」の場合のみ、かき氷ルートへ
	if q1 == "2" && q2 == "1" {
		return correctFoods
	}

	// それ以外はダミーの選択肢を生成
	allGenres := make([]string, 0, len(stores))
	for _, s := range stores {
		allGenres = append(allGenres, s.Genre)
	}

	rand.Shuffle(len(allGenres), func(i, j int) { allGenres[i], allGenres[j] = allGenres[j], allGenres[i] })

	uniqueGenres := make(map[string]bool)
	for _, food := range correctFoods {
		uniqueGenres[food] = true
	}

	var result []string
	for _, genre := range allGenres {
		if !uniqueGenres[genre] && len(result) < 4 {
			result = append(result, genre)
		}
	}
	return result
}
