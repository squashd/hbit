package http

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/SQUASHD/hbit/auth"
	"github.com/SQUASHD/hbit/config"
)

func SetUpAPIGateway(
	authSvc auth.Service,
	jwtConf config.JwtOptions,
) (http.Handler, error) {
	var errs []error

	userSvcProxy, userErr := getUserSvcProxyFromEnv()
	featsSvcProxy, featErr := getFeatSvcProxyFromEnv()
	rpgSvcProxy, rpgErr := getRPGSvcProxyFromEnv()
	taskSvcProxy, taskErr := getTaskSvcProxyFromEnv()

	for _, err := range []error{userErr, featErr, rpgErr, taskErr} {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	// Auth handler is set here since we do authentication at the gateway level
	authHandler := newAuthHandler(authSvc, jwtConf)

	// Entry point for the gateway
	// Centralizing authentication and registration
	// These routes are not authenticated
	entry := http.NewServeMux()
	gateway := http.NewServeMux()

	// Catch all route
	entry.HandleFunc("/", notFound)
	gateway.Handle("/auth/", http.StripPrefix("/auth", entry))

	// Auth / unprotected routes
	entry.HandleFunc("POST /login", authHandler.Login)

	// Orchestrator for registration
	registrationOrch := NewRegistrationOrchestrator(authSvc, os.Getenv("RPG_SVC_URL"), &http.Client{})
	entry.HandleFunc("POST /register", registrationOrch.OrchestrateRegistration)
	entry.HandleFunc("GET /verify", Verify(authSvc, jwtConf))

	authMiddleware := JwtAuthRouterMiddleware(authSvc, jwtConf)

	// Settings and soonTM profiles
	gateway.Handle("/user/", authMiddleware(http.StripPrefix("/user", userSvcProxy)))

	// Achievements/feats driven by events
	gateway.Handle("/feats/", authMiddleware(http.StripPrefix("/feats", featsSvcProxy)))

	// RPG service - driven by task completion but managed independently
	gateway.Handle("/rpg/", authMiddleware(http.StripPrefix("/rpg", rpgSvcProxy)))

	// Orchestrator for tasks
	taskOrchestrator := NewTaskOrchestrationRouter(&http.Client{})
	// Since the tasks done and undo are the driver of events, they are handled by the orchestrator
	// to provide a unified response to the client.
	gateway.Handle("POST /tasks/{id}/{action}", authMiddleware(http.StripPrefix("/tasks", taskOrchestrator)))

	// Generic task management
	gateway.Handle("/tasks/", authMiddleware(http.StripPrefix("/tasks", taskSvcProxy)))

	return gateway, nil
}

func getUserSvcProxyFromEnv() (*httputil.ReverseProxy, error) {
	userSvcURL, err := url.Parse(os.Getenv("USER_SVC_URL"))
	if err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = userSvcURL.Scheme
			r.URL.Host = userSvcURL.Host
		},
	}, nil
}

func getFeatSvcProxyFromEnv() (*httputil.ReverseProxy, error) {
	featSvcURL, err := url.Parse(os.Getenv("FEAT_SVC_URL"))
	if err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = featSvcURL.Scheme
			r.URL.Host = featSvcURL.Host
		}}, nil
}

func getRPGSvcProxyFromEnv() (*httputil.ReverseProxy, error) {
	rpgSvcURL, err := url.Parse(os.Getenv("RPG_SVC_URL"))
	if err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = rpgSvcURL.Scheme
			r.URL.Host = rpgSvcURL.Host
		}}, nil
}

func getTaskSvcProxyFromEnv() (*httputil.ReverseProxy, error) {
	taskSvcURL, err := url.Parse(os.Getenv("TASK_SVC_URL"))
	if err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = taskSvcURL.Scheme
			r.URL.Host = taskSvcURL.Host
		}}, nil
}

// getUpdateSvcProxyFromEnv is no longer needed as we use an orchestrator to handle the
// updates based on the tasks done and undone.
// func getUpdateSvcProxyFromEnv() (*httputil.ReverseProxy, error) {
// 	updateSvcURL, err := url.Parse(os.Getenv("UPDATES_SVC_URL"))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &httputil.ReverseProxy{
// 		Director: func(r *http.Request) {
// 			r.URL.Scheme = updateSvcURL.Scheme
// 			r.URL.Host = updateSvcURL.Host
// 		},
// 	}, nil
//
// }
