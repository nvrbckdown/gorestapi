package app

import (
	"../models"
	u "../utils"
	"context"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/new", "/api/user/login"} // Список эндпоинтов
		requestPath := r.URL.Path                               // Текущий путь запроса

		//Проверяем, не требует ли запрос аутентификации, обслуживаем запрос, если он не нужен
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Получение токена

		if tokenHeader == "" { //Токен отсутсвует, возвращаем 403 http-код Unauthorized
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //Токен обычно поставляется в формате `Bearer {token-body}`, мы проверяем, соответствует ли полученный токен этому требованию
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Получаем вторую часть токена
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Неправильный токен, как правильно, возвращается 403-код
			response := u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid { //токен недействителен, возможно, не подписан на этом сервере
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		//Всё прошло хорошо, продолжаем выполнение запроса
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //передать управление следующиму обработчику!
	})
}
