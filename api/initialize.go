package api

func InitCache() {
	go func() {
		new(APITheme).ThemeCacheInit()
	}()
}
