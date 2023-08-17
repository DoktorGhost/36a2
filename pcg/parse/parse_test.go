package parse_test

import (
	"testing"

	"GoNews/pcg/parse"
	"GoNews/pcg/typeStruct"

	"github.com/stretchr/testify/assert"
)

func TestParseRSS(t *testing.T) {
	// Создание фикстуры с RSS-потоком
	fixture := []byte(`
		<rss version="2.0">
			<channel>
				<item>
					<title>Test Title 1</title>
					<description>Test Description 1</description>
					<pubDate>2023-08-16T10:00:00Z</pubDate>
					<link>http://example.com/test1</link>
				</item>
				<item>
					<title>Test Title 2</title>
					<description>Test Description 2</description>
					<pubDate>2023-08-16T11:00:00Z</pubDate>
					<link>http://example.com/test2</link>
				</item>
			</channel>
		</rss>
	`)

	// Парсинг фикстуры с RSS-потоком
	posts, err := parse.ParseRSSFixture(string(fixture))
	assert.NoError(t, err, "Unexpected error")

	// Проверка количества полученных постов
	assert.Len(t, posts, 2, "Unexpected number of posts")

	// Проверка значений постов
	expectedPosts := []typeStruct.Post{
		{
			Title:   "Test Title 1",
			Content: "Test Description 1",
			PubTime: "2023-08-16T10:00:00Z",
			Link:    "http://example.com/test1",
		},
		{
			Title:   "Test Title 2",
			Content: "Test Description 2",
			PubTime: "2023-08-16T11:00:00Z",
			Link:    "http://example.com/test2",
		},
	}

	for i, post := range posts {
		assert.Equal(t, expectedPosts[i], post, "Post values do not match")
	}
}
