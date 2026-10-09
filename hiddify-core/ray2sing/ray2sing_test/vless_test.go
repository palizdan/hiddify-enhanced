package ray2sing_test

import (
	"context"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/hiddify/ray2sing/ray2sing"
	"github.com/stretchr/testify/require"
)

func TestVless(t *testing.T) {

	url := "vless://25da296e-1d96-48ae-9867-4342796cd742@172.67.149.95:443?encryption=none&fp=chrome&host=vless.229feb8b52a0e7e117ea76f8b591bcb3.workers.dev&path=%2F%3Fed%3D2048&security=tls&sni=vless.229feb8b52a0e7e117ea76f8b591bcb3.workers.dev&type=ws#رایگان | VLESS | @Helix_Servers | US🇺🇸 | 0️⃣1️⃣"

	// Define the expected JSON structure
	expectedJSON := `
	{
		"outbounds": [
		  {
			"type": "vless",
			"tag": "رایگان | VLESS | @Helix_Servers | US🇺🇸 | 0️⃣1️⃣ § 0",
			"server": "172.67.149.95",
			"server_port": 443,
			"uuid": "25da296e-1d96-48ae-9867-4342796cd742",
			"tls": {
			  "enabled": true,
			  "server_name": "vless.229feb8b52a0e7e117ea76f8b591bcb3.workers.dev",
			  "alpn": "http/1.1",
			  "utls": {
				"enabled": true,
				"fingerprint": "chrome"
			  }
			},
			"transport": {
			  "type": "ws",
			  "path": "/",
			  "headers": {
				"Host": "vless.229feb8b52a0e7e117ea76f8b591bcb3.workers.dev"
			  },
			  "max_early_data": 2048,
			  "early_data_header_name": "Sec-WebSocket-Protocol"
			},
			"packet_encoding": "xudp"
		  }
		]
	  }
	`
	ray2sing.CheckUrlAndJson(url, expectedJSON, t)
}

func TestVlessEncryptionWithXHTTP(t *testing.T) {
	privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	require.NoError(t, err)
	publicKey := base64.RawURLEncoding.EncodeToString(privateKey.PublicKey().Bytes())
	encryption := "mlkem768x25519plus.native.0rtt.100-111-1111." + publicKey
	link := fmt.Sprintf(
		"vless://409f106a-b2f2-4416-b186-5429c9979cd9@pasarguard.example:443?encryption=%s&security=tls&sni=pasarguard.example&type=xhttp&host=cdn.example.com&path=%%2Fpasarguard&mode=auto#PasarGuard",
		url.QueryEscape(encryption),
	)

	configJSON, err := ray2sing.Ray2Singbox(context.Background(), link, false)
	require.NoError(t, err)
	var config struct {
		Outbounds []struct {
			Type       string `json:"type"`
			Encryption string `json:"encryption"`
			Transport  struct {
				Type string `json:"type"`
				Mode string `json:"mode"`
				Host string `json:"host"`
				Path string `json:"path"`
			} `json:"transport"`
		} `json:"outbounds"`
	}
	require.NoError(t, json.Unmarshal(configJSON, &config))
	require.Len(t, config.Outbounds, 1)
	require.Equal(t, "vless", config.Outbounds[0].Type)
	require.Equal(t, encryption, config.Outbounds[0].Encryption)
	require.Equal(t, "xhttp", config.Outbounds[0].Transport.Type)
	require.Equal(t, "auto", config.Outbounds[0].Transport.Mode)
	require.Equal(t, "cdn.example.com", config.Outbounds[0].Transport.Host)
	require.Equal(t, "/pasarguard", config.Outbounds[0].Transport.Path)
}
