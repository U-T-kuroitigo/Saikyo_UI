package web

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

// 店舗情報を保持する構造体を定義
type Store struct {
	Name  string
	Genre string
	URL   string
}

// メニューページのルートを登録
func RegisterMenuPageRoutes(e *echo.Echo) {
	e.GET("/menu", HandleMenu)
	e.POST("/menu", HandleMenu)
}

// データベースの役割を果たすデータ
var stores = []Store{
	{Name: "メキシカンフード", Genre: "メキシカンフード", URL: "https://www.toranomonhills.com/gourmet_shops/7537.html"},
	{Name: "鮎ラーメン +", Genre: "ラーメン", URL: "https://www.toranomonhills.com/toranomonyokocho/1022.html"},
	{Name: "AM STRAM GRAM", Genre: "タルト", URL: "https://www.toranomonhills.com/gourmet_shops/3617.html"},
	{Name: "意気な寿司処阿部", Genre: "寿司", URL: "https://www.toranomonhills.com/gourmet_shops/0008.html"},
	{Name: "餃子マニア", Genre: "中華", URL: "https://www.toranomonhills.com/gourmet_shops/3642.html"},
	{Name: "CARAVAN", Genre: "ピザ", URL: "https://www.toranomonhills.com/gourmet_shops/7077.html"},
	{Name: "Cassolo", Genre: "イタリアン", URL: "https://www.toranomonhills.com/gourmet_shops/0096.html"},
	{Name: "焼千房 虎ノ門", Genre: "鉄板焼き", URL: "https://www.toranomonhills.com/gourmet_shops/0025.html"},
	{Name: "スタバ", Genre: "カフェ", URL: "https://www.toranomonhills.com/gourmet_shops/0022.html"},
	{Name: "スコインター", Genre: "タイ料理", URL: "https://www.toranomonhills.com/gourmet_shops/5575.html"},
	{Name: "鰻まえはら", Genre: "鰻", URL: "https://www.toranomonhills.com/gourmet_shops/0027.html"},
	{Name: "京菓匠 鶴屋吉信", Genre: "和菓子", URL: "https://www.toranomonhills.com/gourmet_shops/0041.html"},
	{Name: "BeBu(ビブ)", Genre: "ハンバーガー", URL: "https://www.toranomonhills.com/gourmet_shops/0011.html"},
	{Name: "なるたけ虎ノ門", Genre: "和食", URL: "https://www.toranomonhills.com/gourmet_shops/0093.html"},
}

// 正解ルートの食べ物リスト
var correctFoods = []string{"鰻", "かき氷", "タイ料理", "和食"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// getFoodsは質問3の選択肢を決定するヘルパー関数
func getFoods(q1, q2 string) []string {
	// 質問1が「サプライズが欲しい」かつ質問2が「誰かと楽しく」の場合のみ、かき氷ルートへ
	if q1 == "2" && q2 == "1" {
		return correctFoods
	} else {
		allGenres := make([]string, 0, len(stores))
		for _, s := range stores {
			allGenres = append(allGenres, s.Genre)
		}

		shuffled := make([]string, len(allGenres))
		copy(shuffled, allGenres)
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

		uniqueGenres := make(map[string]bool)
		for _, food := range correctFoods {
			uniqueGenres[food] = true
		}

		var result []string
		for _, genre := range shuffled {
			if !uniqueGenres[genre] && len(result) < 4 {
				result = append(result, genre)
			}
		}
		return result
	}
}

// HandleMenuは全ての質問と結果を処理するハンドラ
func HandleMenu(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		data := map[string]interface{}{
			"CurrentStep": "q1",
		}
		return c.Render(http.StatusOK, "menu.html", data)
	}

	step := c.FormValue("step")

	switch step {
	case "q1":
		data := map[string]interface{}{
			"CurrentStep": "q2",
			"Q1":          c.FormValue("q1"),
		}
		return c.Render(http.StatusOK, "menu.html", data)

	case "q2":
		q1 := c.FormValue("q1")
		q2 := c.FormValue("q2")

		data := map[string]interface{}{
			"CurrentStep": "q3",
			"Q1":          q1,
			"Q2":          q2,
			"Foods":       getFoods(q1, q2),
		}
		return c.Render(http.StatusOK, "menu.html", data)

	case "q3":
		q1Str := c.FormValue("q1")
		q2Str := c.FormValue("q2")
		q3Str := c.FormValue("q3")

		q1, _ := strconv.Atoi(q1Str)
		q2, _ := strconv.Atoi(q2Str)
		q3 := -1
		if q3Str == "鰻" {
			q3 = 0
		}
		if q3Str == "かき氷" {
			q3 = 1
		}
		if q3Str == "タイ料理" {
			q3 = 2
		}
		if q3Str == "和食" {
			q3 = 3
		}

		idx := q1*16 + q2*4 + q3

		if q3Str == "かき氷" {
			data := map[string]interface{}{
				"CurrentStep": "ice_flavor",
				"Index":       idx,
				"Flavors":     []string{"いちご", "メロン", "ブルーハワイ", "オレンジ"},
			}
			return c.Render(http.StatusOK, "menu.html", data)
		} else {
			var recommendedStore Store
			for _, s := range stores {
				if s.Genre == q3Str {
					recommendedStore = s
					break
				}
			}

			data := map[string]interface{}{
				"CurrentStep": "result",
				"Index":       idx,
				"Result":      recommendedStore.Name + "のお店をオススメします！",
				"StoreURL":    recommendedStore.URL,
			}
			return c.Render(http.StatusOK, "menu.html", data)
		}

	case "ice_flavor":
		flavor := c.FormValue("flavor")
		idxStr := c.FormValue("idx")

		data := map[string]interface{}{
			"CurrentStep": "result",
			"Index":       idxStr,
			"Result":      "おめでとうございます！" + flavor + "かき氷が注文されました！",
		}
		return c.Render(http.StatusOK, "menu.html", data)
	}

	return c.String(http.StatusBadRequest, "Invalid step")
}
