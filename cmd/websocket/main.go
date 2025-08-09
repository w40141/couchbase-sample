// Package main provides a simple WebSocket client to subscribe to changes from a Sync Gateway.
package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Sync GatewayのURLとデータベース名を環境変数から取得
	// 例: SYNC_GATEWAY_URL="ws://localhost:4984/db/"
	syncGatewayURL := os.Getenv("SYNC_GATEWAY_URL")
	if syncGatewayURL == "" {
		log.Fatal("SYNC_GATEWAY_URL環境変数が設定されていません。")
	}

	// WebSocket接続の開始
	u, err := url.Parse(syncGatewayURL)
	if err != nil {
		log.Fatal("URLの解析に失敗しました:", err)
	}

	log.Printf("Sync Gatewayの変更を購読します: %s", u.String())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(u.String()+"_changes?feed=websocket", nil)
	if err != nil {
		log.Fatal("WebSocket接続に失敗しました:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("読み込みエラー:", err)
				return
			}
			log.Printf("変更を受信しました: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			// 定期的にPingを送信して接続を維持する
			err := c.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("Ping送信エラー:", err)
				return
			}
			log.Printf("Pingを送信しました at %v", t)

		case <-interrupt:
			log.Println("割り込み信号を受信しました。接続をクローズします。")

			// 接続をきれいに閉じる
			err := c.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				log.Println("Closeメッセージ送信エラー:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// _changesリクエストのWebSocketプロトコルについて
// Sync Gatewayの_changes?feed=websocketエンドポイントは、
// 接続後に特別なリクエストボディを送信する必要はありません。
// 接続が確立されると、Sync Gatewayはドキュメントの変更が発生するたびに、
// JSON形式の変更通知を自動的に送信してきます。
