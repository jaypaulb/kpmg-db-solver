package canvus

import (
	"context"
	"fmt"
)

// ServerConfig represents the server configuration.
type ServerConfig struct {
	Access         string                `json:"access,omitempty"`
	Authentication *AuthenticationConfig `json:"authentication,omitempty"`
	Email          *EmailConfig          `json:"email,omitempty"`
	ExternalURL    string                `json:"external_url,omitempty"`
	ServerName     string                `json:"server_name,omitempty"`
}

// AuthenticationConfig represents authentication settings for the server.
type AuthenticationConfig struct {
	DomainAllowList      []string        `json:"domain_allow_list,omitempty"`
	Password             *PasswordConfig `json:"password,omitempty"`
	QRCode               *QRCodeConfig   `json:"qr_code,omitempty"`
	RequireAdminApproval bool            `json:"require_admin_approval,omitempty"`
	SAML                 *SAMLConfig     `json:"saml,omitempty"`
}

// PasswordConfig represents password authentication settings.
type PasswordConfig struct {
	Enabled       bool `json:"enabled,omitempty"`
	SignUpEnabled bool `json:"sign_up_enabled,omitempty"`
}

// QRCodeConfig represents QR code authentication settings.
type QRCodeConfig struct {
	Enabled bool `json:"enabled,omitempty"`
}

// SAMLConfig represents SAML authentication settings.
type SAMLConfig struct {
	ACSURL             string `json:"acs_url,omitempty"`
	Enabled            bool   `json:"enabled,omitempty"`
	IDPCertFingerPrint string `json:"idp_cert_finger_print,omitempty"`
	IDPEntityID        string `json:"idp_entity_id,omitempty"`
	IDPTargetURL       string `json:"idp_target_url,omitempty"`
	NameIDFormat       string `json:"name_id_format,omitempty"`
	SignUpEnabled      bool   `json:"sign_up_enabled,omitempty"`
	SPEntityID         string `json:"sp_entity_id,omitempty"`
}

// EmailConfig represents email server settings.
type EmailConfig struct {
	MailReplyToAddress              string `json:"mail_reply_to_address,omitempty"`
	MailReplyToName                 string `json:"mail_reply_to_name,omitempty"`
	MailSenderAddress               string `json:"mail_sender_address,omitempty"`
	MailSenderName                  string `json:"mail_sender_name,omitempty"`
	SMTPAllowSelfSignedCertificates bool   `json:"smtp_allow_self_signed_certificates,omitempty"`
	SMTPHost                        string `json:"smtp_host,omitempty"`
	SMTPPassword                    string `json:"smtp_password,omitempty"`
	SMTPPort                        int    `json:"smtp_port,omitempty"`
	SMTPSecurity                    string `json:"smtp_security,omitempty"`
	SMTPUsername                    string `json:"smtp_username,omitempty"`
}

// GetServerConfig retrieves the server configuration from the Canvus API.
func (s *Session) GetServerConfig(ctx context.Context) (*ServerConfig, error) {
	var config ServerConfig
	err := s.doRequest(ctx, "GET", "server-config", nil, &config, nil, false)
	if err != nil {
		return nil, fmt.Errorf("GetServerConfig: %w", err)
	}
	return &config, nil
}

// UpdateServerConfig updates the server configuration in the Canvus API.
func (s *Session) UpdateServerConfig(ctx context.Context, req ServerConfig) (*ServerConfig, error) {
	var config ServerConfig
	err := s.doRequest(ctx, "PATCH", "server-config", req, &config, nil, false)
	if err != nil {
		return nil, fmt.Errorf("UpdateServerConfig: %w", err)
	}
	return &config, nil
}

// SendTestEmail sends a test email to the current user via the Canvus API.
func (s *Session) SendTestEmail(ctx context.Context) error {
	return s.doRequest(ctx, "POST", "server-config/send-test-email", nil, nil, nil, false)
}
