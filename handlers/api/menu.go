package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

// Store は店舗情報を保持する構造体
type Store struct {
	Name  string `json:"name"`
	Genre string `json:"genre"`
	URL   string `json:"url"`
}

// グローバル変数：アプリケーション全体で利用するデータを保持
var (
	// genreToStoreMapは、ジャンル名から店舗情報への高速なアクセスのために使用
	genreToStoreMap map[string]Store
	// allGenresは、JSONから読み込んだ全てのユニークなジャンルを保持
	allGenres []string
)

// correctFoodsは、特定の条件で表示される正解の食べ物リスト
var correctFoods = []string{"鰻", "かき氷", "タイ料理", "和食"}

// foodToIndexMapは、食べ物の名前をインデックス計算用の数値に変換するために使用
var foodToIndexMap = map[string]int{
	"鰻":    0,
	"かき氷":  1,
	"タイ料理": 2,
	"和食":   3,
}

// initはパッケージの初期化時に一度だけ実行される
func init() {
	// JSONファイルから店舗データを読み込む
	jsonFile, err := os.Open("public/json/food_stores.json")
	if err != nil {
		log.Fatalf("FATAL: public/json/food_stores.json を開けませんでした: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("FATAL: public/json/food_stores.json を読み込めませんでした: %v", err)
	}

	var stores []Store
	if err := json.Unmarshal(byteValue, &stores); err != nil {
		log.Fatalf("FATAL: public/json/food_stores.json のパースに失敗しました: %v", err)
	}

	// 読み込んだデータを効率的に利用できるよう、マップを初期化
	genreToStoreMap = make(map[string]Store)
	uniqueGenres := make(map[string]bool)
	for _, s := range stores {
		if _, ok := genreToStoreMap[s.Genre]; !ok {
			genreToStoreMap[s.Genre] = s
		}
		if !uniqueGenres[s.Genre] {
			uniqueGenres[s.Genre] = true
			allGenres = append(allGenres, s.Genre)
		}
	}

	log.Println("正常に", len(stores), "件の店舗情報を public/json/food_stores.json から読み込みました")
}

// ProcessMenuStep はフォームの入力値に基づいて、ユーザーフローの次のステップを決定する
func ProcessMenuStep(params map[string]string) map[string]interface{} {
	step, ok := params["step"]
	if !ok {
		return map[string]interface{}{"CurrentStep": "q1"}
	}

	q1Val := params["q1"]
	q2Val := params["q2"]

	switch step {
	case "q1":
		if q1Val == "" {
			return map[string]interface{}{
				"CurrentStep": "q1",
				"Error":       "選択してください。",
			}
		}
		return map[string]interface{}{
			"CurrentStep": "q2",
			"Q1":          q1Val,
		}

	case "q2":
		if q2Val == "" {
			return map[string]interface{}{
				"CurrentStep": "q2",
				"Q1":          q1Val,
				"Error":       "選択してください。",
			}
		}
		return map[string]interface{}{
			"CurrentStep": "q3",
			"Q1":          q1Val,
			"Q2":          q2Val,
			"Foods":       getFoods(q1Val, q2Val),
		}

	case "q3":
		q3Val := params["q3"]
		if q3Val == "" {
			return map[string]interface{}{
				"CurrentStep": "q3",
				"Q1":          q1Val,
				"Q2":          q2Val,
				"Foods":       getFoods(q1Val, q2Val),
				"Error":       "選択してください。",
			}
		}

		// 特別な「かき氷」ルートを処理する
		if q3Val == "かき氷" {
			q1, _ := strconv.Atoi(q1Val)
			q2, _ := strconv.Atoi(q2Val)
			q3 := foodToIndexMap["かき氷"]
			idx := q1*16 + q2*4 + q3

			return map[string]interface{}{
				"CurrentStep": "ice_flavor",
				"Index":       idx,
				"Flavors":     []string{"いちご", "メロン", "ブルーハワイ", "オレンジ"},
			}
		}

		// 通常の結果表示を処理する
		recommendedStore, storeExists := genreToStoreMap[q3Val]
		if !storeExists {
			log.Printf("ERROR: おすすめの店舗が見つかりません genre=%s", q3Val)
			return map[string]interface{}{"CurrentStep": "q1", "Error": "おすすめ店舗が見つかりませんでした。"}
		}

		resultData := map[string]interface{}{
			"CurrentStep": "result",
			"Result":      recommendedStore.Name + "のお店をオススメします！",
			"StoreURL":    recommendedStore.URL,
		}

		return resultData

	case "ice_flavor":
		// このステップは、ブラウザからの通常のフォーム送信では到達しません。
		// JavaScriptからのAPI呼び出しが /api/orders を直接叩くため、このロジックは主に画面表示の再生成用です。
		q1, _ := strconv.Atoi(q1Val)
		q2, _ := strconv.Atoi(q2Val)
		q3 := foodToIndexMap["かき氷"]
		idx := q1*16 + q2*4 + q3

		return map[string]interface{}{
			"CurrentStep": "ice_flavor",
			"Index":       idx,
			"Flavors":     []string{"いちご", "メロン", "ブルーハワイ", "オレンジ"},
		}
	}

	// 不正なステップの場合は、最初の質問にデフォルトで戻す
	return map[string]interface{}{"CurrentStep": "q1"}
}

// getFoods は質問3の食べ物の選択肢リストを決定する
func getFoods(q1, q2 string) []string {
	if q1 == "2" && q2 == "1" {
		return correctFoods
	}

	seed, err := strconv.ParseInt(q1+q2, 10, 64)
	if err != nil {
		seed = 1
	}
	r := rand.New(rand.NewSource(seed))

	shuffledGenres := make([]string, len(allGenres))
	copy(shuffledGenres, allGenres)

	r.Shuffle(len(shuffledGenres), func(i, j int) {
		shuffledGenres[i], shuffledGenres[j] = shuffledGenres[j], shuffledGenres[i]
	})

	correctFoodSet := make(map[string]bool)
	for _, food := range correctFoods {
		correctFoodSet[food] = true
	}

	var result []string
	for _, genre := range shuffledGenres {
		if !correctFoodSet[genre] {
			result = append(result, genre)
			if len(result) == 4 {
				break
			}
		}
	}
	return result
}

// --- ここから下は注文APIのロジック ---

// RegisterIceOrderAPIRoutes は注文APIのエンドポイントを登録します。
func RegisterIceOrderAPIRoutes(e *echo.Echo) {
	e.POST("/api/orders", createOrder)
}

// createOrder はクライアントからの注文リクエストを受け取り、外部APIへ中継します。
func createOrder(c echo.Context) error {
	// ▼▼▼ ハードコードされたIDを環境変数から読み込むように修正 ▼▼▼
	storeID := os.Getenv("STORE_ID")
	if storeID == "" {
		log.Println("エラー: 環境変数 STORE_ID が設定されていません。")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "サーバーの設定エラーです。"})
	}
	// ▲▲▲ ここまで修正 ▲▲▲

	// クライアントからのリクエストボディを読み込む
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストボディが不正です"})
	}

	// 外部APIのエンドポイントURLを構築
	url := fmt.Sprintf("https://kakigori-api.fly.dev/v1/stores/%s/orders", storeID)

	// 外部APIへの新しいリクエストを作成
	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		c.Logger().Errorf("注文APIプロキシ: リクエスト構築エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "リクエストの構築に失敗しました"})
	}
	// ヘッダーを設定
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SaikyoUI/1.0 (+echo)")

	// HTTPクライアントを作成し、リクエストを送信
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.Logger().Errorf("注文APIプロキシ: 外部APIへのリクエストエラー: %v", err)
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "外部APIへのリクエストに失敗しました"})
	}
	defer resp.Body.Close()

	// 外部APIからのレスポンスボディを読み込む
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Errorf("注文APIプロキシ: レスポンスボディの読み込みエラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "外部APIのレスポンス読み込みに失敗しました"})
	}

	// 外部APIからのステータスコードとレスポンスボディを、そのままクライアントに返す
	return c.Blob(resp.StatusCode, "application/json", respBody)
}

