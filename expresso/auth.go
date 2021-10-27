package expresso

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/justtaldevelops/expresso/expresso/packet"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// sessionServer is the main URL for accessing parts of the Mojang session server.
const sessionServer = "https://sessionserver.mojang.com/session/minecraft/"

// authResponse is a response for authentication from Mojang containing important information such as UUIDs.
type  authResponse struct {
	UUID string `json:"id"`
	Name string `json:"name"`
	Properties interface{} `json:"properties,omitempty"`
}

// authenticatedWithMojang returns true if the player is authenticated with Mojang.
func authenticatedWithMojang(username string, sharedSecret []byte, encryptionRequest *packet.EncryptionRequest) (bool, authResponse) {
	var data authResponse

	params := &url.Values{}
	params.Set("username", username)
	params.Set("serverId", authDigest(encryptionRequest.ServerID, sharedSecret, encryptionRequest.PublicKey))

	req, err := http.NewRequest("GET", sessionServer + "hasJoined", nil)
	req.URL.RawQuery = params.Encode()
	if err != nil {
		return false, authResponse{}
	}
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, authResponse{}
	}
	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return false, authResponse{}
	}
	_ = httpResp.Body.Close()
	if len(respBody) < 1 {
		return false, authResponse{}
	}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return false, authResponse{}
	}
	if len(data.UUID) == 0 {
		return false, authResponse{}
	}
	return true, data
}

// authDigest computes a special SHA-1 digest required for Minecraft web
// authentication on Premium servers (online-mode=true).
// Source: http://wiki.vg/Protocol_Encryption#Server
// Reference: https://gist.github.com/toqueteos/5372776.
func authDigest(serverID string, sharedSecret []byte, publicKey rsa.PublicKey) string {
	publicKeyEncoded, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}

	h := sha1.New()
	h.Write([]byte(serverID))
	h.Write(sharedSecret)
	h.Write(publicKeyEncoded)
	hash := h.Sum(nil)

	// Check for negative hashes
	negative := (hash[0] & 0x80) == 0x80
	if negative {
		hash = twosComplement(hash)
	}

	// Trim away zeroes
	res := strings.TrimLeft(fmt.Sprintf("%x", hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}

// twosComplement ...
func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = ^p[i]
		if carry {
			carry = p[i] == 0xff
			p[i]++
		}
	}
	return p
}
