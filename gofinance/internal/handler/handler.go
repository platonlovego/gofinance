package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gofinance/internal/service"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.register)
			auth.POST("/login", h.login)
		}

		protected := api.Group("/", h.authMiddleware)
		{
			protected.GET("/categories", h.getCategories)

			tx := protected.Group("/transactions")
			{
				tx.POST("", h.createTransaction)
				tx.GET("", h.getTransactions)
				tx.GET("/:id", h.getTransaction)
				tx.PUT("/:id", h.updateTransaction)
				tx.DELETE("/:id", h.deleteTransaction)
			}

			protected.GET("/summary", h.getSummary)
		}
	}

	return r
}

// authMiddleware проверяет JWT из заголовка Authorization: Bearer <token>
func (h *Handler) authMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
		return
	}

	userID, err := h.svc.ParseToken(parts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Set("userID", userID)
	c.Next()
}

func userIDFromCtx(c *gin.Context) int {
	id, _ := c.Get("userID")
	return id.(int)
}
