package api

import (
	"encoding/json"
	"io/ioutil"
	"log" // エラーロギング用のパッケージ
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

// initはパッケージの初期化時に一度だけ実行される
func init() {
	// 乱数生成器のシードを設定
	rand.Seed(time.Now().UnixNano())

	// JSONファイルから店舗データを読み込む
	jsonFile, err := os.Open("public/json/food_stores.json")
	if err != nil {
		// ファイルを開けなかった場合、エラーをログに出力してアプリケーションを停止する
		log.Fatalf("FATAL: public/json/food_stores.json を開けませんでした: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		// ファイルを読み込めなかった場合、エラーをログに出力してアプリケーションを停止する
		log.Fatalf("FATAL: public/json/food_stores.json を読み込めませんでした: %v", err)
	}

	err = json.Unmarshal(byteValue, &stores)
	if err != nil {
		// JSONのパースに失敗した場合、エラーをログに出力してアプリケーションを停止する
		log.Fatalf("FATAL: public/json/food_stores.json のパースに失敗しました: %v", err)
	}

	// ファイルが正しく読み込まれた場合に成功メッセージをログに出力する
	log.Println("正常に", len(stores), "件の店舗情報を public/json/food_stores.json から読み込みました")
}

// ProcessMenuStep はフォームの入力値に基づいて、ユーザーフローの次のステップを決定する
func ProcessMenuStep(params map[string]string) map[string]interface{} {
	step, ok := params["step"]
	// 'step' が存在しない場合は、最初のアクセスとみなす
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

		// 特別な「かき氷」ルートを処理する
		if q3Str == "かき氷" {
			return map[string]interface{}{
				"CurrentStep": "ice_flavor",
				"Index":       idx,
				"Flavors":     []string{"いちご", "メロン", "ブルーハワイ", "オレンジ"},
			}
		}

		// 通常の結果表示を処理する
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

	// 不正なステップの場合は、最初の質問にデフォルトで戻す
	return map[string]interface{}{"CurrentStep": "q1"}
}

// getFoods は質問3の食べ物の選択肢リストを決定する
func getFoods(q1, q2 string) []string {
	// ユーザーが「サプライズ」を望み、「誰かと」一緒にいる場合は、特別な食べ物リストを表示する
	if q1 == "2" && q2 == "1" {
		return correctFoods
	}

	// それ以外の場合は、ダミーの選択肢をランダムに生成する
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

