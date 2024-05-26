package cookie

import "net/http"

func SetSecureCookie(w http.ResponseWriter, key, value string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode, // Lax is required because of the OAuth redirect. Otherwise cookies are blocked on return.
	}
	http.SetCookie(w, cookie)
}

func ExpireCookie(w http.ResponseWriter, key string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}
