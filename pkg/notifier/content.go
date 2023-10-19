package notifier

type Content struct {
	Title string
	Data  []byte
	// push使用
	ClickType string
	URL       string
}
