package wget

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// Проверка доступности файла
func IsFileExists(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

// Создание файла
func CreateFile(filePath string) (*os.File, error) {
	// Если файл по данному пути существует - возвращается ошибка
	if IsFileExists(filePath) {
		return nil, fmt.Errorf("file %s already exist", filePath)
	}

	// Создание всех необходимых директорий для файла
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return nil, err
	}

	// Создание файла
	return os.Create(filePath)
}

// Форматирование URL
func FormatURL(url, contentURL string) string {
	if strings.HasPrefix(url, "/") {
		url = strings.TrimRight(contentURL, "/") + "/" + strings.TrimLeft(url, "/")
	}
	return strings.ToLower(url)
}

// Парсинг HTML и извлечение URL
func GetURLsFromHTML(htmlContent io.Reader) []string {
	tokenizer := html.NewTokenizer(htmlContent)
	URLs := make([]string, 0, 10)

	// Проверка корректности URL
	isCorrectURL := func(url string) bool {
		prefixes := []string{"https://", "http://", "/"}
		for _, prefix := range prefixes {
			if strings.HasPrefix(url, prefix) {
				return true
			}
		}
		return false
	}

	// Итерационный проход всех токенов
	for tType := tokenizer.Next(); tType != html.ErrorToken; tType = tokenizer.Next() {
		if tType != html.StartTagToken {
			continue
		}
		token := tokenizer.Token()

		// Если тег токена не <a> - пропуск
		if token.Data != "a" {
			continue
		}

		for _, attr := range token.Attr {
			// Если аттрибут не href - пропуск
			if attr.Key != "href" {
				continue
			}

			// Если корректный URL - добавление в список
			if isCorrectURL(attr.Val) {
				URLs = append(URLs, attr.Val)
				break
			}
		}
	}
	return URLs
}

type Downloader struct {
	visitedLinks map[string]struct{}
	host         string
	otherHosts   bool
}

// Создание/форматирование пути файла, где будет храниться файл обрабатываемого URL
func makeDirectoryPath(host, filepath, extension string) string {
	filepath = strings.TrimSuffix(filepath, "/")
	if filepath == "" {
		filepath = "index"
	}
	if !strings.HasSuffix(filepath, extension) {
		filepath += extension
	}
	return path.Join(host, filepath)
}

// Рекурсивная загрузка файлов по URL
func (downloader *Downloader) Download(url *url.URL, depth int) (int, error) {
	// Если глубина < 0 - выход из рекурсии
	if depth < 0 {
		return len(downloader.visitedLinks), nil
	}
	// Добавление текущего URL в множество посещенных ссылок
	downloader.visitedLinks[strings.TrimSuffix(url.Host+url.Path, "/")] = struct{}{}

	urlString := fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, url.Path)
	fmt.Println("downloading:", urlString)

	// Получение ответа на GET запрос обрабатываемого URL
	response, err := http.Get(urlString)
	if err != nil {
		return len(downloader.visitedLinks), err
	}
	defer response.Body.Close()

	// Получение расширения ответа по Content-Type
	extensions, err := mime.ExtensionsByType(response.Header.Get("Content-Type"))
	// Если произошла ошибка или расширение не распознано - возврат ошибки
	if err != nil || len(extensions) == 0 {
		return len(downloader.visitedLinks), errors.New("can't detect file extension: " + url.Host + url.Path)
	}

	// Создание файла для сохранения body ответа
	file, err := CreateFile(makeDirectoryPath(url.Host, url.Path, extensions[len(extensions)-1]))
	if err != nil {
		return len(downloader.visitedLinks), err
	}
	defer file.Close()

	// Чтение данных из body
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return len(downloader.visitedLinks), err
	}

	reader := bytes.NewReader(data)
	// Сохранение ответа в файл
	if _, err := io.Copy(file, reader); err != nil {
		return len(downloader.visitedLinks), err
	}
	// Перемещение указателя для повторого чтения из bytes.Reader
	if _, err := reader.Seek(0, 0); err != nil {
		return len(downloader.visitedLinks), err
	}

	// Для каждой ссылки из полученного ответа
	for _, htmlURLString := range GetURLsFromHTML(reader) {
		htmlURL, err := url.Parse(FormatURL(htmlURLString, urlString))
		if err != nil {
			return len(downloader.visitedLinks), err
		}
		// Если не производится скачивание с других хостов, а URL имеет другой адрес - пропуск
		if !downloader.otherHosts && htmlURL.Host != downloader.host {
			continue
		}
		// Если URL уже был посещён - пропуск
		if _, ok := downloader.visitedLinks[strings.TrimRight(htmlURL.Host+htmlURL.Path, "/")]; ok {
			continue
		}
		// Запуск рекурсионной загрузки рассматриваемого URL
		if _, err := downloader.Download(htmlURL, depth-1); err != nil {
			return len(downloader.visitedLinks), err
		}
	}
	return len(downloader.visitedLinks), nil
}

func NewDownloader(host string, otherHosts bool) *Downloader {
	return &Downloader{
		visitedLinks: make(map[string]struct{}),
		host:         host,
		otherHosts:   otherHosts,
	}
}

// Реализация утилиты WGET
func WGET(options Options) (int, error) {
	downloader := NewDownloader(options.BaseURL.Host, options.OtherHosts)
	if options.Recursive {
		return downloader.Download(options.BaseURL, options.RecursionDepth)
	}
	return downloader.Download(options.BaseURL, 0)
}
