package dingding

type Message struct {
	MsgType  string       `json:"msgtype"`
	Markdown MarkdownBody `json:"markdown"`
}

type MarkdownBody struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
