package ai

import (
	"errors"
	"net/http"

	"omniforge-api/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GenerateText(c *gin.Context) {
	var input GenerateTextRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			MessageInvalidRequest,
		)
		return
	}

	if err := ValidateGenerateTextRequest(&input); err != nil {
	switch {
	case errors.Is(err, ErrPromptRequired):
		response.Error(
			c,
			http.StatusBadRequest,
			MessagePromptRequired,
		)

	case errors.Is(err, ErrPromptTooLong):
		response.Error(
			c,
			http.StatusBadRequest,
			MessagePromptTooLong,
		)

	case errors.Is(err, ErrInvalidOutputFormat):
		response.Error(
			c,
			http.StatusBadRequest,
			MessageInvalidOutputFormat,
		)

	case errors.Is(err, ErrInvalidPromptType):
		response.Error(
			c,
			http.StatusBadRequest,
			MessageInvalidPromptType,
		)

	case errors.Is(err, ErrStructuredStreamingNotValid):
		response.Error(
			c,
			http.StatusBadRequest,
			MessageStructuredStreamInvalid,
		)
	}

	return
}

	if input.Stream {
		h.streamText(c, input)
		return
	}

	if input.OutputFormat == OutputFormatStructured {
		result, err := h.service.GenerateStructured(
			c.Request.Context(),
			input,
		)

		if err != nil {
			response.Error(
				c,
				http.StatusInternalServerError,
				MessageGenerateFailed,
			)
			return
		}

		response.Success(
			c,
			http.StatusOK,
			MessageGenerateSuccess,
			result,
		)
		return
	}

	result, err := h.service.GenerateText(
		c.Request.Context(),
		input,
	)

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			MessageGenerateFailed,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		MessageGenerateSuccess,
		result,
	)
}

func (h *Handler) streamText(
	c *gin.Context,
	input GenerateTextRequest,
) {
	c.Header(
		"Content-Type",
		"text/event-stream",
	)

	c.Header(
		"Cache-Control",
		"no-cache",
	)

	c.Header(
		"Connection",
		"keep-alive",
	)

	err := h.service.StreamText(
		c.Request.Context(),
		input,
		func(delta string) error {
			c.SSEvent(
				"message",
				delta,
			)

			c.Writer.Flush()

			return nil
		},
	)

	if err != nil {
		c.SSEvent(
			"error",
			MessageGenerateFailed,
		)

		c.Writer.Flush()
	}
}