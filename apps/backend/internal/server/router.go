package server

// TODO:
//   // buildRouter membuat chi.Router, memasang middleware global, health check,
//   // dan mem-mount tiap module di bawah prefix /api/v1.
//   //
//   // Middleware global (pakai bawaan chi): RequestID, RealIP, Logger,
//   // Recoverer, Timeout, dan CORS bila perlu.
//   //
//   // Contoh mounting:
//   //   r.Route("/api/v1", func(r chi.Router) {
//   //       userMod.RegisterRoutes(r)      // rute publik (register/login)
//   //       r.Group(func(r chi.Router) {
//   //           r.Use(middleware.Authenticator(verifier))  // rute butuh login
//   //           questMod.RegisterRoutes(r)
//   //           scoringMod.RegisterRoutes(r)
//   //       })
//   //   })
//   func buildRouter(deps ...) http.Handler
