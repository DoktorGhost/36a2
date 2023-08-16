package typeStruct

import "time"

type Post struct {
	ID      int       // номер записи
	Title   string    // заголовок публикации
	Content string    // содержание публикации
	PubTime time.Time // время публикации
	Link    string    // ссылка на источник
}

/*
func (p *Post) FormatPubTime() string {
	t := time.Unix(p.PubTime, 0)
	return t.Format("2006-01-02 15:04:05")
}
*/

func (p *Post) GetSummary() string {
	// Вернуть краткое описание новости, например, первые 100 символов содержания
	maxLength := 100
	if len(p.Content) > maxLength {
		return p.Content[:maxLength] + "..."
	}
	return p.Content
}

func NewPost(title, content, link string, pubTime time.Time) Post {
	return Post{
		Title:   title,
		Content: content,
		PubTime: pubTime,
		Link:    link,
	}
}
