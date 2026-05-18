// Copyright © yanghy. All Rights Reserved.
//
// # Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0

package notification

import (
	"fmt"
	"github.com/energye/energy/v3/platform/notification"
	"github.com/energye/energy/v3/platform/notification/types"
	"github.com/energye/lcl/api"
	"sync"
	"time"
)

var (
	gNotification types.INotification
	once          = sync.Once{}
)

func Init() {
	once.Do(func() {
		gNotification = notification.New()
		gNotification.SetOnNotificationResponse(func(result types.Result) {

		})
		api.SetOnReleaseCallback(func() {
			// close gNotification
		})
	})
}

func Info(title, content string) {
	if gNotification != nil {
		opts := types.Options{
			ID:       fmt.Sprintf("success-%d", time.Now().Unix()),
			Title:    "✅ " + title,
			Subtitle: title,
			Body:     content,
		}
		err := gNotification.SendNotification(opts)
		if err != nil {

		}
	}
}

func Error(title, content string) {
	if gNotification != nil {
		opts := types.Options{
			ID:       fmt.Sprintf("error-%d", time.Now().Unix()),
			Title:    "❌ " + title,
			Subtitle: title,
			Body:     content,
		}
		err := gNotification.SendNotification(opts)
		if err != nil {

		}
	}
}
