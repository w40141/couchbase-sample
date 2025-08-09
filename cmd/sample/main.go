// Package main
package main

import (
	"fmt"
	"log"
	"net/http"
)

// hellHandlerは /hell エンドポイントへのリクエストを処理します。
func hellHandler(w http.ResponseWriter, r *http.Request) {
	// GETメソッド以外は許可しない
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}
	// レスポンスとして "Hello, World!" を書き込みます。
	fmt.Fprintf(w, "Hello, World!")
}

// healthcheckHandlerは /healthcheck エンドポイントへのリクエストを処理します。
func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// GETメソッド以外は許可しない
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}
	// HTTPステータスコードを200 OKに設定します。
	w.WriteHeader(http.StatusOK)
	// レスポンスボディとして "OK" を書き込みます。
	w.Write([]byte("OK"))
}

// pingHandlerは /ping エンドポイントへのリクエストを処理します。
func pingHandler(w http.ResponseWriter, r *http.Request) {
	// GETメソッド以外は許可しない
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}
	// HTTPステータスコードを200 OKに設定します。
	w.WriteHeader(http.StatusOK)
	// レスポンスボディとして "pong" を書き込みます。
	w.Write([]byte("pong"))
}

// main関数はプログラムのエントリーポイントです。
func main() {
	// /hell のリクエストを hellHandler で処理するように設定します。
	http.HandleFunc("/hell", hellHandler)
	// /healthcheck のリクエストを healthcheckHandler で処理するように設定します。
	http.HandleFunc("/healthcheck", healthcheckHandler)
	// /ping のリクエストを pingHandler で処理するように設定します。
	http.HandleFunc("/ping", pingHandler)

	// サーバーが起動することを示すメッセージを出力します。
	fmt.Println("Server is running on http://localhost:8080")

	// 8080ポートでサーバーを起動します。起動に失敗した場合はエラーを出力し、プログラムを終了します。
	log.Fatal(http.ListenAndServe(":8080", nil))
}
