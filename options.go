package tg

// SendConfig holds configuration for sending a message.
type SendConfig struct {
	ThreadID           int
	Markup             Markup
	Silent             bool
	ProtectContent     bool
	ReplyToMessageID   int
	DisableLinkPreview bool
}

// SendOption configures SendConfig.
type SendOption func(*SendConfig)

// WithThreadID directs the message to a specific forum topic or thread.
func WithThreadID(id int) SendOption {
	return func(c *SendConfig) {
		c.ThreadID = id
	}
}

// WithMarkup attaches an inline keyboard, reply keyboard, or reply action.
func WithMarkup(m Markup) SendOption {
	return func(c *SendConfig) {
		c.Markup = m
	}
}

// WithSilent sends the message without sound/notification.
func WithSilent(silent bool) SendOption {
	return func(c *SendConfig) {
		c.Silent = silent
	}
}

// WithProtectContent disables forwarding and saving of the message.
func WithProtectContent(protect bool) SendOption {
	return func(c *SendConfig) {
		c.ProtectContent = protect
	}
}

// WithReplyTo marks the message as a reply to a specific message ID.
func WithReplyTo(id int) SendOption {
	return func(c *SendConfig) {
		c.ReplyToMessageID = id
	}
}

// WithDisableLinkPreview disables web page previews for links in this message.
func WithDisableLinkPreview(disable bool) SendOption {
	return func(c *SendConfig) {
		c.DisableLinkPreview = disable
	}
}

// ApplySendOptions compiles a SendConfig from functional options.
func ApplySendOptions(opts ...SendOption) SendConfig {
	var cfg SendConfig
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

// EditConfig holds configuration for editing an existing message.
type EditConfig struct {
	Markup             Markup
	DisableLinkPreview bool
}

// EditOption configures EditConfig.
type EditOption func(*EditConfig)

// WithEditMarkup updates the message's attached keyboard markup.
func WithEditMarkup(m Markup) EditOption {
	return func(c *EditConfig) {
		c.Markup = m
	}
}

// WithEditDisableLinkPreview toggles link previews for the edited message.
func WithEditDisableLinkPreview(disable bool) EditOption {
	return func(c *EditConfig) {
		c.DisableLinkPreview = disable
	}
}

// ApplyEditOptions compiles an EditConfig from functional options.
func ApplyEditOptions(opts ...EditOption) EditConfig {
	var cfg EditConfig
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

// CallbackConfig holds configuration for answering a callback query.
type CallbackConfig struct {
	Text      string
	ShowAlert bool
	URL       string
	CacheTime int
}

// CallbackOption configures CallbackConfig.
type CallbackOption func(*CallbackConfig)

// WithAlert shows an alert modal dialog instead of a transient toast banner.
func WithAlert(alert bool) CallbackOption {
	return func(c *CallbackConfig) {
		c.ShowAlert = alert
	}
}

// WithCallbackText sets the notification text displayed to the user.
func WithCallbackText(text string) CallbackOption {
	return func(c *CallbackConfig) {
		c.Text = text
	}
}

// WithCallbackURL opens a URL client-side (for games or external flows).
func WithCallbackURL(url string) CallbackOption {
	return func(c *CallbackConfig) {
		c.URL = url
	}
}

// WithCacheTime sets client caching time for callback query results in seconds.
func WithCacheTime(seconds int) CallbackOption {
	return func(c *CallbackConfig) {
		c.CacheTime = seconds
	}
}

// ApplyCallbackOptions compiles a CallbackConfig from functional options.
func ApplyCallbackOptions(opts ...CallbackOption) CallbackConfig {
	var cfg CallbackConfig
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}
