package bots

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"

	"todo-backend/models"
)

type Bots struct {
	DB            *gorm.DB
	TgBot         *gotgbot.Bot
	Discord       *discordgo.Session
	TgToken       string
	DiscordToken  string
	OwnerTelegram string
	OwnerDiscord  string
}

func InitBots(db *gorm.DB) (*Bots, error) {
	b := &Bots{
		DB:            db,
		TgToken:       os.Getenv("TG_BOT_TOKEN"),
		DiscordToken:  os.Getenv("DISCORD_TOKEN"),
		OwnerTelegram: strings.TrimSpace(os.Getenv("OWNER_TG_ID")),
		OwnerDiscord:  strings.TrimSpace(os.Getenv("OWNER_DC_ID")),
	}
	if b.TgToken != "" {
		if err := b.initTelegram(); err != nil {
			return nil, err
		}
	}
	if b.DiscordToken != "" {
		if err := b.initDiscord(); err != nil {
			return nil, err
		}
	}
	go b.schedulerLoop()
	return b, nil
}

func (b *Bots) initTelegram() error {
	bot, err := gotgbot.NewBot(b.TgToken, nil)
	if err != nil {
		return err
	}
	dp := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dp, nil)
	log.Printf("Starting Telegram bot with username %s", bot.User.Username)
	go func() {
		err = updater.StartPolling(
			bot,
			&ext.PollingOpts{
				DropPendingUpdates: true,
				GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
					Timeout: 9,
					RequestOpts: &gotgbot.RequestOpts{
						Timeout: time.Second * 10,
					},
				},
			},
		)
		if err != nil {
			log.Println("failed to start polling:", err.Error())
		}
	}()
	b.TgBot = bot
	return nil
}

func (b *Bots) initDiscord() error {
	dg, err := discordgo.New("Bot " + b.DiscordToken)
	if err != nil {
		return err
	}
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Discord bot is up as", r.User.Username)
	})
	if err := dg.Open(); err != nil {
		return err
	}
	b.Discord = dg
	return nil
}

func (b *Bots) schedulerLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		<-ticker.C
		b.processReminders()
	}
}

func (b *Bots) processReminders() {
	now := time.Now()
	var todos []models.Todo
	if err := b.DB.Where("notify_enabled = ? AND completed = ?", true, false).Find(&todos).Error; err != nil {
		return
	}
	for _, t := range todos {
		if t.NextNotifyAt == nil {
			continue
		}
		if t.DueDate != nil && now.After(*t.DueDate) {
			continue
		}
		if t.NextNotifyAt.After(now) {
			continue
		}
		msg := buildMessage(&t)
		if t.TelegramEnabled && b.TgBot != nil && b.OwnerTelegram != "" {
			chatID, err := strconv.ParseInt(b.OwnerTelegram, 10, 64)
			if err == nil {
				_, _ = b.TgBot.SendMessage(chatID, msg, nil)
			} else {
				log.Println("invalid OWNER_TG_ID:", err)
			}
		}
		if t.DiscordEnabled && b.Discord != nil && b.OwnerDiscord != "" {
			channel, err := b.Discord.UserChannelCreate(b.OwnerDiscord)
			if err == nil {
				_, _ = b.Discord.ChannelMessageSend(channel.ID, msg)
			} else {
				log.Println("failed to create discord dm channel:", err)
			}
		}
		if t.NotifyFrequency > 0 && t.NextNotifyAt != nil {
			next := t.NextNotifyAt.Add(time.Duration(t.NotifyFrequency) * time.Minute)
			if t.DueDate != nil && next.After(*t.DueDate) {
				t.NextNotifyAt = nil
			} else {
				t.NextNotifyAt = &next
			}
		} else {
			t.NextNotifyAt = nil
		}
		_ = b.DB.Save(&t)
	}
}

func buildMessage(t *models.Todo) string {
	strDue := "no due date"
	if t.DueDate != nil {
		strDue = t.DueDate.Format(time.RFC1123)
	}
	return fmt.Sprintf("Reminder: %s\nDue: %s\nPriority: %s\nTask ID: %s", t.Title, strDue, t.Priority, t.ID)
}
