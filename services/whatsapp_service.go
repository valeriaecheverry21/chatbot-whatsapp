package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/user/chatbot-whatsapp/config"
	"github.com/user/chatbot-whatsapp/models"
)

type WhatsAppService struct {
	httpClient *http.Client
	cfg        *config.WhatsAppConfig
	log        *zap.Logger
}

func NewWhatsAppService(cfg *config.WhatsAppConfig, log *zap.Logger) *WhatsAppService {
	return &WhatsAppService{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		cfg:        cfg,
		log:        log,
	}
}

type SendTextInput struct {
	Phone   string
	Message string
	Preview bool
}

func (s *WhatsAppService) SendText(ctx context.Context, input SendTextInput) (*models.Message, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                input.Phone,
		"type":              "text",
		"text": map[string]interface{}{
			"preview_url": input.Preview,
			"body":        input.Message,
		},
	}

	resp, err := s.callAPI(ctx, body)
	if err != nil {
		return nil, err
	}

	s.log.Info("text message sent",
		zap.String("phone", input.Phone),
		zap.String("whatsapp_id", resp.Messages[0].ID),
	)

	return &models.Message{
		WhatsAppID:  resp.Messages[0].ID,
		Direction:   models.DirectionOutbound,
		MessageType: models.MessageTypeText,
		Status:      models.StatusSent,
		Body:        input.Message,
	}, nil
}

type SendImageInput struct {
	Phone   string
	Caption string
	URL     string
}

func (s *WhatsAppService) SendImage(ctx context.Context, input SendImageInput) (*models.Message, error) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                input.Phone,
		"type":              "image",
		"image": map[string]string{
			"link":    input.URL,
			"caption": input.Caption,
		},
	}

	resp, err := s.callAPI(ctx, body)
	if err != nil {
		return nil, err
	}

	return &models.Message{
		WhatsAppID:  resp.Messages[0].ID,
		Direction:   models.DirectionOutbound,
		MessageType: models.MessageTypeImage,
		Status:      models.StatusSent,
		Body:        input.Caption,
		MediaURL:    input.URL,
	}, nil
}

type SendTemplateInput struct {
	Phone    string
	Name     string
	Language string
	Params   []string
}

func (s *WhatsAppService) SendTemplate(ctx context.Context, input SendTemplateInput) (*models.Message, error) {
	components := []map[string]interface{}{}

	if len(input.Params) > 0 {
		var params []map[string]interface{}
		for _, p := range input.Params {
			params = append(params, map[string]interface{}{
				"type": "text",
				"text": p,
			})
		}
		components = append(components, map[string]interface{}{
			"type":       "body",
			"parameters": params,
		})
	}

	payload := map[string]interface{}{
		"name":       input.Name,
		"language":   map[string]string{"code": input.Language},
	}

	if len(components) > 0 {
		payload["components"] = components
	}

	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                input.Phone,
		"type":              "template",
		"template":          payload,
	}

	resp, err := s.callAPI(ctx, body)
	if err != nil {
		return nil, err
	}

	s.log.Info("template sent",
		zap.String("phone", input.Phone),
		zap.String("template", input.Name),
		zap.String("whatsapp_id", resp.Messages[0].ID),
	)

	return &models.Message{
		WhatsAppID:     resp.Messages[0].ID,
		Direction:      models.DirectionOutbound,
		MessageType:    models.MessageTypeTemplate,
		Status:         models.StatusSent,
		TemplateName:   input.Name,
		TemplateLang:   input.Language,
		TemplateParams: strings.Join(input.Params, "|"),
	}, nil
}

func (s *WhatsAppService) MarkAsRead(ctx context.Context, messageID string) error {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageID,
	}

	_, err := s.callAPI(ctx, body)
	if err != nil {
		s.log.Warn("failed to mark as read",
			zap.String("message_id", messageID),
			zap.Error(err),
		)
	}

	return nil
}

type apiResponse struct {
	Messages []apiMessage `json:"messages"`
	Error    *apiError    `json:"error,omitempty"`
}

type apiMessage struct {
	ID string `json:"id"`
}

type apiError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (s *WhatsAppService) callAPI(ctx context.Context, payload map[string]interface{}) (*apiResponse, error) {
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages",
		s.cfg.APIVersion, s.cfg.PhoneNumberID)

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.cfg.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result apiResponse
	if err := json.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("whatsapp api error: %s (code %d)",
			result.Error.Message, result.Error.Code)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(respData))
	}

	return &result, nil
}
