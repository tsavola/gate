// Copyright (c) 2017 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package webapi contains definitions useful for accessing the HTTP and
// websocket APIs.  See https://github.com/tsavola/gate/blob/master/Web.md for
// general documentation.
package webapi

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/tsavola/gate/internal/serverapi"
)

// Version of the Gate webserver API.
const Version = 0

// Name of the module reference source and associated content hash algorithm.
const ModuleRefSource = "sha384"

// Algorithm for converting module content to reference key.  A reference key
// string can be formed by encoding a hash digest with base64.RawURLEncoding.
const ModuleRefHash = crypto.SHA384

// Request URL paths.
const (
	PathVersions   = "/gate/"                            // Available API versions.
	Path           = "/gate/v0"                          // The API.
	PathModule     = Path + "/module"                    // Base of relative module URIs.
	PathModules    = PathModule + "/"                    // Module sources.
	PathModuleRefs = PathModules + ModuleRefSource + "/" // Module reference keys.
	PathInstances  = Path + "/instance/"                 // Instance ids.
)

// Query parameters.
const (
	ParamAction   = "action"
	ParamFunction = "function" // For call or launch action.
	ParamInstance = "instance" // For call or launch action.
	ParamDebug    = "debug"    // For call, launch or resume action.
)

// Actions on modules.  Ref action can be combined with call or launch in a
// single request (action parameter appears twice).
const (
	ActionRef    = "ref"    // Put (reference), post (source) or websocket (call/launch).
	ActionUnref  = "unref"  // Post (reference).
	ActionCall   = "call"   // Put (reference), post (any) or websocket (any).
	ActionLaunch = "launch" // Put (reference), post (any).
)

// Actions on instances.
const (
	ActionIO       = "io"       // Post or websocket.
	ActionStatus   = "status"   // Post.
	ActionWait     = "wait"     // Post.
	ActionSuspend  = "suspend"  // Post.
	ActionResume   = "resume"   // Post.
	ActionSnapshot = "snapshot" // Post.
)

// HTTP request headers.
const (
	HeaderAuthorization = "Authorization" // "Bearer" JSON Web Token.
)

// HTTP request or response headers.
const (
	HeaderContentLength = "Content-Length"
	HeaderContentType   = "Content-Type"
)

// HTTP response headers.
const (
	HeaderLocation = "Location"        // Absolute module ref path.
	HeaderInstance = "X-Gate-Instance" // UUID.
	HeaderStatus   = "X-Gate-Status"   // Status of instance as JSON.
	HeaderDebug    = "X-Gate-Debug"
)

// The supported module content type.
const ContentTypeWebAssembly = "application/wasm"

// The supported key type.
const KeyTypeOctetKeyPair = "OKP"

// The supported elliptic curve.
const KeyCurveEd25519 = "Ed25519"

// The supported signature algorithm.
const SignAlgEdDSA = "EdDSA"

// The supported authorization type.
const AuthorizationTypeBearer = "Bearer"

// JSON Web Key.
type PublicKey struct {
	Kty string `json:"kty"`           // Key type.
	Crv string `json:"crv,omitempty"` // Elliptic curve.
	X   string `json:"x,omitempty"`   // Base64url-encoded unpadded public key.
}

// PublicKeyEd25519 creates a JWK for a JWT header.
func PublicKeyEd25519(publicKey []byte) *PublicKey {
	return &PublicKey{
		Kty: KeyTypeOctetKeyPair,
		Crv: KeyCurveEd25519,
		X:   base64.RawURLEncoding.EncodeToString(publicKey),
	}
}

// JSON Web Token header.
type TokenHeader struct {
	Alg string     `json:"alg"`           // Signature algorithm.
	JWK *PublicKey `json:"jwk,omitempty"` // Public side of signing key.
}

// TokenHeaderEdDSA creates a JWT header.
func TokenHeaderEdDSA(publicKey *PublicKey) *TokenHeader {
	return &TokenHeader{
		Alg: SignAlgEdDSA,
		JWK: publicKey,
	}
}

// MustEncode to a JWT component.
func (header *TokenHeader) MustEncode() []byte {
	serialized, err := json.Marshal(header)
	if err != nil {
		panic(err)
	}

	encoded := make([]byte, base64.RawURLEncoding.EncodedLen(len(serialized)))
	base64.RawURLEncoding.Encode(encoded, serialized)
	return encoded
}

// JSON Web Token payload.
type Claims struct {
	Exp   int64    `json:"exp"`             // Expiration time.
	Aud   []string `json:"aud,omitempty"`   // https://authority/api
	Nonce string   `json:"nonce,omitempty"` // Unique during expiration period.
}

// Status response header.
type Status struct {
	State  string `json:"state,omitempty"`
	Cause  string `json:"cause,omitempty"`
	Result int    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
	Debug  string `json:"debug,omitempty"`
}

func (status Status) String() (s string) {
	if status.State == "" {
		if status.Error == "" {
			return "error"
		} else {
			return fmt.Sprintf("error: %s", status.Error)
		}
	}

	s = status.State

	if status.Cause != "" {
		s = fmt.Sprintf("%s abnormally: %s", s, status.Cause)
	} else if status.State == "terminated" {
		s = fmt.Sprintf("%s with result %d", s, status.Result)
	}

	if status.Error != "" {
		s = fmt.Sprintf("%s; error: %s", s, status.Error)
	}
	return
}

// Response to PathModuleRefs request.
type ModuleRefs struct {
	Modules []ModuleRef `json:"modules"`
}

// An item in a ModuleRefs response.
type ModuleRef = serverapi.ModuleRef

// Response to a PathInstances request.
type Instances struct {
	Instances []InstanceStatus `json:"instances"`
}

// An item in an Instances response.
type InstanceStatus struct {
	Instance string `json:"instance"`
	Status   Status `json:"status"`
}

// ActionCall websocket request message.
type Call struct {
	Authorization string `json:"authorization,omitempty"`
	ContentType   string `json:"content_type,omitempty"`
	ContentLength int64  `json:"content_length,omitempty"`
}

// Reply to Call message.
type CallConnection struct {
	Location string `json:"location,omitempty"` // Absolute module ref path.
	Instance string `json:"instance,omitempty"` // UUID.
	Debug    string `json:"debug,omitempty"`
}

// ActionIO websocket request message.
type IO struct {
	Authorization string `json:"authorization"`
}

// Reply to IO message.
type IOConnection struct {
	Connected bool   `json:"connected"`
	Debug     string `json:"debug,omitempty"`
}

// Second and final text message on successful ActionCall or ActionIO websocket
// connection.
type ConnectionStatus struct {
	Status Status `json:"status"` // Instance status after disconnection.
}

// FunctionRegexp matches valid a function name.
var FunctionRegexp = regexp.MustCompile("^[A-Za-z0-9-._]{1,31}$")
