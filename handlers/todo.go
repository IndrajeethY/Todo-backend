package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"todo-backend/models"
)

type CreateTodoInput struct {
	Title           string     `json:"title" binding:"required"`
	Description     string     `json:"description"`
	DueDate         *time.Time `json:"due_date"`
	Priority        string     `json:"priority"`
	NotifyEnabled   *bool      `json:"notify_enabled"`
	NotifyFrequency *int       `json:"notify_frequency"`
	OrderIndex      *int       `json:"order_index"`
	TelegramEnabled bool       `json:"telegram_enabled"`
	DiscordEnabled  bool       `json:"discord_enabled"`
}

type UpdateTodoInput struct {
	Title           *string    `json:"title"`
	Description     *string    `json:"description"`
	DueDate         *time.Time `json:"due_date"`
	Priority        *string    `json:"priority"`
	Completed       *bool      `json:"completed"`
	NotifyEnabled   *bool      `json:"notify_enabled"`
	NotifyFrequency *int       `json:"notify_frequency"`
	OrderIndex      *int       `json:"order_index"`
	TelegramEnabled *bool      `json:"telegram_enabled"`
	DiscordEnabled  *bool      `json:"discord_enabled"`
}

func RegisterTodoRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	listHandler := func(c *gin.Context) {
		uid := c.GetString("user_id")
		var todos []models.Todo
		query := db.Where("user_id = ?", uid)
		if p := c.Query("priority"); p != "" {
			query = query.Where("priority = ?", p)
		}
		if s := c.Query("completed"); s != "" {
			switch s {
			case "true":
				query = query.Where("completed = ?", true)
			case "false":
				query = query.Where("completed = ?", false)
			}
		}
		if c.Query("due_before") != "" {
			if ttime, err := time.Parse(time.RFC3339, c.Query("due_before")); err == nil {
				query = query.Where("due_date <= ?", ttime)
			}
		}
		if err := query.Order("order_index asc, due_date asc").Find(&todos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch todos"})
			return
		}
		c.JSON(http.StatusOK, todos)
	}

	createHandler := func(c *gin.Context) {
		var in CreateTodoInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		uid := c.GetString("user_id")
		notifyEnabled := true
		if in.NotifyEnabled != nil {
			notifyEnabled = *in.NotifyEnabled
		}
		freq := 60
		if in.NotifyFrequency != nil {
			freq = *in.NotifyFrequency
		}
		t := models.Todo{
			UserID:          uid,
			Title:           in.Title,
			Description:     in.Description,
			DueDate:         in.DueDate,
			Priority:        in.Priority,
			Completed:       false,
			NotifyEnabled:   notifyEnabled,
			NotifyFrequency: freq,
			TelegramEnabled: in.TelegramEnabled,
			DiscordEnabled:  in.DiscordEnabled,
		}
		if in.OrderIndex != nil {
			t.OrderIndex = *in.OrderIndex
		}
		if t.Priority == "" {
			t.Priority = "medium"
		}
		if err := db.Create(&t).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, t)
	}

	getHandler := func(c *gin.Context) {
		uid := c.GetString("user_id")
		id := c.Param("id")
		var t models.Todo
		if err := db.Where("id = ? AND user_id = ?", id, uid).First(&t).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch"})
			return
		}
		c.JSON(http.StatusOK, t)
	}

	updateHandler := func(c *gin.Context) {
		uid := c.GetString("user_id")
		id := c.Param("id")
		var t models.Todo
		if err := db.Where("id = ? AND user_id = ?", id, uid).First(&t).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch"})
			return
		}
		var in UpdateTodoInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if in.Title != nil {
			t.Title = *in.Title
		}
		if in.Description != nil {
			t.Description = *in.Description
		}
		if in.DueDate != nil {
			t.DueDate = in.DueDate
		}
		if in.Priority != nil {
			t.Priority = *in.Priority
		}
		if in.Completed != nil {
			t.Completed = *in.Completed
			if t.Completed {
				t.NextNotifyAt = nil
			}
		}
		if in.NotifyEnabled != nil {
			t.NotifyEnabled = *in.NotifyEnabled
		}
		if in.NotifyFrequency != nil {
			t.NotifyFrequency = *in.NotifyFrequency
		}
		if in.OrderIndex != nil {
			t.OrderIndex = *in.OrderIndex
		}
		if in.TelegramEnabled != nil {
			t.TelegramEnabled = *in.TelegramEnabled
		}
		if in.DiscordEnabled != nil {
			t.DiscordEnabled = *in.DiscordEnabled
		}
		if err := db.Save(&t).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
			return
		}
		c.JSON(http.StatusOK, t)
	}

	deleteHandler := func(c *gin.Context) {
		uid := c.GetString("user_id")
		id := c.Param("id")
		if err := db.Where("id = ? AND user_id = ?", id, uid).Delete(&models.Todo{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
			return
		}
		c.Status(http.StatusNoContent)
	}

	completeHandler := func(c *gin.Context) {
		uid := c.GetString("user_id")
		id := c.Param("id")
		now := time.Now()
		if err := db.Model(&models.Todo{}).
			Where("id = ? AND user_id = ?", id, uid).
			Updates(map[string]interface{}{"completed": true, "updated_at": now, "next_notify_at": nil}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark complete"})
			return
		}
		var t models.Todo
		if err := db.Where("id = ? AND user_id = ?", id, uid).First(&t).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"message": "marked complete"})
			return
		}
		c.JSON(http.StatusOK, t)
	}

	reorderHandler := func(c *gin.Context) {
		uid := c.GetString("user_id")
		var payload []struct {
			ID    string `json:"id" binding:"required"`
			Index int    `json:"index" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tx := db.Begin()
		for _, p := range payload {
			if err := tx.Model(&models.Todo{}).Where("id = ? AND user_id = ?", p.ID, uid).Update("order_index", p.Index).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reorder"})
				return
			}
		}
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reorder"})
			return
		}
		c.Status(http.StatusOK)
	}

	rg.GET("/", listHandler)
	rg.GET("", listHandler)
	rg.POST("/", createHandler)
	rg.POST("", createHandler)
	rg.GET("/:id", getHandler)
	rg.PUT("/:id", updateHandler)
	rg.DELETE("/:id", deleteHandler)
	rg.POST("/:id/complete", completeHandler)
	rg.PATCH("/reorder", reorderHandler)
}
