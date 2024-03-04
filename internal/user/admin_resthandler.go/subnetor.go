package adminresthandler

import (
	"net"
	"net/http"
)

// IPHeader заголовок с реальным IP адреса клиента.
const IPHeader = "X-Real-IP"

// Subnet проверка IP адреса клиента на вхождение в доверенную подсеть.
func (h *AdminRESTHandler) checkSubnet(next http.Handler) http.Handler {
	subnetFn := func(res http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get(IPHeader)
		if ip == "" || !h.subnet.Contains(net.ParseIP(ip)) {
			res.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(res, r)
	}
	return http.HandlerFunc(subnetFn)
}