package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inquizarus/nagg/pkg/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestThatJWTCopyClaimToHeaderWorksWhenClaimIsPresent(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTAiLCJ0eXAiOiJKV1QifQ.eyJ1c2VyIjoiZDFiNDdjZWYtMzc0ZC00MDY2LThiOTAtZTViODZhMmY1MWMxIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4Nzg3L2FwaS9hdXRoL2xvZ2luIiwiYXVkIjpbImh0dHA6Ly9sb2NhbGhvc3Q6ODc4Ny9hcGkvbm90ZWJvb2tzIl0sImV4cCI6MTY4MDk3NDc4NiwiaWF0IjoxNjgwMzY5OTg2fQ.ip2iYhFGXIEPnURXApdhY-OsCaV7SdxL55c4iswhfcPXD7uWjBuaGPWOT-7ZaZPpEqPoFOQAW4JDIgeNEd3Unq8RqJWCuaLgVJwT9n0CQqkh5-m_dH4RJVQ6iSanJvghSfEllwYucdpze1uiJA2oOuPXmZxgYTTuqqV6a5541d2RSJ8QfnwzMoDnCSkOTlCEmqauoSHVMgTTtl_BkmU4Kf4QU2ouB2jtLOf6ZElXAMBDaMug4_xsFVuQVic3D3oSzm_VFcR9t1OB96x0ywrnR-J5oDuyDxchwgpKKPLIVjv7uoVN6qNRsCwxTh-7--CGkJUuXyaDzxkqzLK8xTtEA9fzppDbeKDC9OjmMz7XVroOFlxuuuHJDzSZLi0LmDTdMoXJiYbDzsKRLtaeXITFHS30xot0clov67t1w4zqTZkPsTCizEwskdwcEn9o0SffjQV9mU0uBSMitzkHpI0gxEJ1MphjyCWrWfRAkURXsOhg7xkC81Vdxbe31GOkeSSCGz4AqoxRb-fMARbS0vZ4ME9_fPUDfPReNDHrMf4Xff1GiHCeHJRLoQy--g9Gnik_qbudwJRBrXZSKV7P4j6kEgL1d_36yZyr0kRFuNKxVtIOfA6ZoppAO92rCsXg7t-AhvdE6oXVrf-h6jZtN-z2MUCfdw_gy9NEX21hngG2538"
	mw := middlewares.MakeJWTCopyClaimToHeaderMiddleware("user", "x-user")
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	mw(nil).ServeHTTP(w, r)

	assert.Equal(t, "d1b47cef-374d-4066-8b90-e5b86a2f51c1", r.Header.Get("x-user"))
}

func TestThatJWTCopyClaimToHeaderWorksWhenTokenIsNotPresent(t *testing.T) {
	mw := middlewares.MakeJWTCopyClaimToHeaderMiddleware("user", "x-user")
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mw(nil).ServeHTTP(w, r)

	assert.Equal(t, "", r.Header.Get("x-user"))
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}
