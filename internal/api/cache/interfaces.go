package cache

type Storage interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Del(key string)
}
