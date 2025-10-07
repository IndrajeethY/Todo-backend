package bots

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
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
	log.Printf("Starting Telegram bot with username %s", bot.User.Username)
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
    now := time.Now().UTC()
    var todos []models.Todo
    if err := b.DB.Where("notify_enabled = ? AND completed = ?", true, false).Find(&todos).Error; err != nil {
        return
    }
    for _, t := range todos {
        if t.NextNotifyAt == nil {
            now := time.Now().UTC()
            t.NextNotifyAt = &now
            _ = b.DB.Save(&t)
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
                _, _ = b.TgBot.SendMessage(chatID, msg, &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
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
	return fmt.Sprintf(
		"**Reminder:** `%s`\n**Due:** `%s`\n**Priority:** `%s`\n**Task ID:** `%s`",
		t.Title,
		strDue,
		t.Priority,
		t.ID,
	)
}
