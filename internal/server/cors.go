package server

import "net/http"

// 모델 C: 프론트가 외부(GitHub Pages 등)에 있고 데이터만 러너로 교차출처 호출한다.
// /api/* 만 CORS 를 허용한다. /_admin 은 절대 CORS 를 주지 않는다(로컬 전용).
//
// 출처가 교사마다 달라 열거 불가 → Allow-Origin: *. 쿠키/인증을 쓰지 않는 공개 데이터
// API 라 * 가 안전하다(Allow-Credentials 를 켜지 않으므로 자격증명 전송도 불가).
// CORS 는 브라우저가 요청을 "보내게만" 할 뿐, 동사 허용은 여전히 프리셋이 강제한다.

func setCORS(w http.ResponseWriter) {
	h := w.Header()
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
	h.Set("Access-Control-Allow-Headers", "Content-Type")
	h.Set("Access-Control-Max-Age", "600")
}

// withCORS 는 /api 핸들러를 감싸 모든 응답(성공·에러 모두)에 CORS 헤더를 얹는다.
// 헤더는 본 핸들러가 WriteHeader 를 호출하기 전에 설정되므로 그대로 전송된다.
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		h(w, r)
	}
}

// handlePreflight 는 /api/* 의 OPTIONS(프리플라이트)에 CORS 헤더 + 204 로 즉시
// 응답한다. 프리셋/본문 처리는 하지 않는다(실제 요청에서 프리셋이 강제됨).
func (s *Server) handlePreflight(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	w.WriteHeader(http.StatusNoContent)
}
