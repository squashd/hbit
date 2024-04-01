package http

import "net/http"

func (s *Server) RegisterRoutes() http.Handler {
	mainRouter := http.NewServeMux()
	mainRouter.HandleFunc("GET /achievements", s.AuthMiddleware(s.handleAchievementsGet))

	mainRouter.HandleFunc("POST /login", s.handleLogin)
	mainRouter.HandleFunc("POST /register", s.handleRegister)
	mainRouter.HandleFunc("POST /revoke", s.AdminMiddleware(s.handleRevoke))

	mainRouter.HandleFunc("GET /tasks", s.AuthMiddleware(s.handleTaskFindAll))
	mainRouter.HandleFunc("POST /tasks", s.AuthMiddleware(s.handleTaskCreate))
	mainRouter.HandleFunc("PUT /tasks", s.AuthMiddleware(s.handleTaskUpdate))
	mainRouter.HandleFunc("DELETE /tasks", s.AuthMiddleware(s.handleTaskDelete))

	mainRouter.HandleFunc("GET /characters/{id}", s.AuthMiddleware(s.handleCharacterGet))
	mainRouter.HandleFunc("POST /characters", s.AuthMiddleware(s.handleCharacterCreate))
	mainRouter.HandleFunc("PUT /characters/{id}", s.AuthMiddleware(s.handleCharacterUpdate))
	mainRouter.HandleFunc("DELETE /characters/{id}", s.AuthMiddleware(s.handleCharacterDelete))

	mainRouter.HandleFunc("GET /quests", s.AuthMiddleware(s.handleQuestGetAll))

	mainRouter.HandleFunc("GET /achievements", s.AuthMiddleware(s.handleAchievementsGet))

	mainRouter.HandleFunc("GET /settings", s.AuthMiddleware(s.handleSettingsGet))
	mainRouter.HandleFunc("PUT /settings", s.AuthMiddleware(s.handleSettingsUpdate))

	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("GET /characters", s.AuthMiddleware(s.handleCharacterGetAll))
	mainRouter.Handle("/admin/", s.AdminRouterMiddleware(adminRouter))

	return mainRouter
}
