package http

import "net/http"

func (s *serverMonolith) registerAuthRoutes(router *http.ServeMux, authHandler *AuthHandler) {
	router.HandleFunc("POST /login", authHandler.Login)
	router.HandleFunc("POST /register", authHandler.Register)
	router.HandleFunc("POST /revoke", s.AdminMiddleware(authHandler.Revoke))
}

func (s *serverMonolith) registerTaskRoutes(router *http.ServeMux, taskHandler *TaskHandler) {
	router.HandleFunc("GET /tasks", s.AuthMiddleware(taskHandler.FindAll))
	router.HandleFunc("POST /tasks", s.AuthMiddleware(taskHandler.Create))
	router.HandleFunc("PUT /tasks", s.AuthMiddleware(taskHandler.Update))
	router.HandleFunc("DELETE /tasks", s.AuthMiddleware(taskHandler.Delete))
}

func (s *serverMonolith) registerCharacterRoutes(router *http.ServeMux, charHandler *CharacterHandler) {
	router.HandleFunc("GET /characters/{id}", s.AuthMiddleware(charHandler.CharacterGet))
	router.HandleFunc("POST /characters", s.AuthMiddleware(charHandler.CharacterCreate))
	router.HandleFunc("PUT /characters/{id}", s.AuthMiddleware(charHandler.CharacterUpdate))
	router.HandleFunc("DELETE /characters/{id}", s.AuthMiddleware(charHandler.CharacterDelete))
}

func (s *serverMonolith) registerQuestRoutes(router *http.ServeMux, questHandler *QuestHandler) {
	router.HandleFunc("GET /quests", s.AuthMiddleware(questHandler.GetAll))
}

func (s *serverMonolith) registerAchievementRoutes(router *http.ServeMux, achHandler *AchievementHandler) {
	router.HandleFunc("GET /achievements", s.AuthMiddleware(achHandler.AchievementsGet))
}

func (s *serverMonolith) registerUserRoutes(router *http.ServeMux, userHandler *UserHandler) {
	router.HandleFunc("GET /settings", s.AuthMiddleware(userHandler.SettingsGet))
	router.HandleFunc("PUT /settings", s.AuthMiddleware(userHandler.SettingsUpdate))
}

func (s *serverMonolith) RegisterRoutes() http.Handler {
	router := http.NewServeMux()

	authHandler := NewAuthHandler(s.authSvc)
	s.registerAuthRoutes(router, authHandler)

	taskHandler := NewTaskHandler(s.taskSvc)
	s.registerTaskRoutes(router, taskHandler)

	charHandler := NewCharacterHandler(s.charSvc)
	s.registerCharacterRoutes(router, charHandler)

	achHandler := NewAchievementHandler(s.achSvc)
	s.registerAchievementRoutes(router, achHandler)

	questHandler := NewQuestHandler(s.questSvc)
	s.registerQuestRoutes(router, questHandler)

	userHandler := NewUserHandler(s.userSvc)
	s.registerUserRoutes(router, userHandler)

	adminRouter := http.NewServeMux()

	characterHandler := NewCharacterHandler(s.charSvc)
	adminRouter.HandleFunc("GET /characters", s.AuthMiddleware(characterHandler.CharacterGetAll))
	router.Handle("/admin/", http.StripPrefix("/admin", s.AdminRouterMiddleware(adminRouter)))

	return router
}
