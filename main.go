package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Ginのデフォルト設定でルーター作成（ログ・リカバリ付き）
	r := gin.Default()

	// ヘルスチェック用エンドポイント
	r.GET("/health", func(c *gin.Context) {
		// JSONで {"status": "ok"} を返す
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// ポート8080でHTTPサーバー起動
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
