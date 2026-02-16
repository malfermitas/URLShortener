package out

type KeyGenerator interface {
	// Generate возвращает новый ключ (без гарантии уникальности — проверка ложится на сервис).
	Generate() string
}
