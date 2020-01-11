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

// Name of the module reference source and associated content hash algorithm.
const ModuleRefSource = "sha384"

// Algorithm for converting module content to reference key.  A reference key
// string can be formed by encoding a hash digest with base64.RawURLEncoding.
const ModuleRefHash = crypto.SHA384

// Request URL paths.
const (
	Path           = "/gate/"               // The API.
	PathModule     = "/gate/module"         // Base of relative module URIs.
	PathModules    = "/gate/module/"        // Module sources.
	PathModuleRefs = "/gate/module/sha384/" // Module reference keys.
	PathInstances  = "/gate/instance/"      // Instance ids.
)

// Query parameters.
const (
	ParamAction   = "action"
	ParamFunction = "function" // For call, launch or resume action.
	ParamInstance = "instance" // For call or launch action.
	ParamDebug    = "debug"    // For call, launch or resume action.
)

// Actions on modules.  ActionRef can be combined with ActionCall or
// ActionLaunch in a single request (ParamAction appears twice in the URL).
const (
	ActionRef    = "ref"    // Put (reference), post (source) or websocket (call/launch).
	ActionUnref  = "unref"  // Post (reference).
	ActionCall   = "call"   // Put (reference), post (any) or websocket (any).
	ActionLaunch = "launch" // Put (reference), post (any).
)

// Actions on instances.  ActionWait can be combined with ActionKill or
// ActionSuspend in a single request (ParamAction appears twice in the URL).
const (
	ActionIO       = "io"       // Post or websocket.
	ActionStatus   = "status"   // Post.
	ActionWait     = "wait"     // Post.
	ActionKill     = "kill"     // Post.
	ActionSuspend  = "suspend"  // Post.
	ActionResume   = "resume"   // Post.
	ActionSnapshot = "snapshot" // Post.
	ActionDelete   = "delete"   // Post.
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

// Instance state enumeration.
const (
	StateRunning    = "RUNNING"
	StateSuspended  = "SUSPENDED"
	StateHalted     = "HALTED"
	StateTerminated = "TERMINATED"
	StateKilled     = "KILLED"
)

// Instance state cause enumeration.  Empty value means that the cause is a
// normal one (e.g. client action, successful completion).
//
// The cause enumeration is open-ended: new values may appear in the future.
const (
	CauseNormal = ""

	// Abnormal causes for StateSuspended:
	CauseCallStackExhausted = "CALL_STACK_EXHAUSTED"
	CauseABIDeficiency      = "ABI_DEFICIENCY"

	// Abnormal causes for StateKilled:
	CauseUnreachable                   = "UNREACHABLE"
	CauseMemoryAccessOutOfBounds       = "MEMORY_ACCESS_OUT_OF_BOUNDS"
	CauseIndirectCallIndexOutOfBounds  = "INDIRECT_CALL_INDEX_OUT_OF_BOUNDS"
	CauseIndirectCallSignatureMismatch = "INDIRECT_CALL_SIGNATURE_MISMATCH"
	CauseIntegerDivideByZero           = "INTEGER_DIVIDE_BY_ZERO"
	CauseIntegerOverflow               = "INTEGER_OVERFLOW"
	CauseABIViolation                  = "ABI_VIOLATION"
	CauseInternal                      = "INTERNAL"
)

// Status response header.
type Status struct {
	State  string `json:"state,omitempty"`
	Cause  string `json:"cause,omitempty"`
	Result int    `json:"result,omitempty"` // Meaningful if StateHalted or StateTerminated.
	Error  string `json:"error,omitempty"`  // Optional details for abnormal causes.
	Debug  string `json:"debug,omitempty"`
}

func (status Status) String() (s string) {
	switch {
	case status.State == "":
		if status.Error == "" {
			return "error"
		} else {
			return fmt.Sprintf("error: %s", status.Error)
		}

	case status.Cause != "":
		s = fmt.Sprintf("%s abnormally: %s", status.State, status.Cause)

	case status.State == StateHalted || status.State == StateTerminated:
		s = fmt.Sprintf("%s with result %d", status.State, status.Result)

	default:
		s = status.State
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
	Instance  string `json:"instance"`
	Status    Status `json:"status"`
	Transient bool   `json:"transient,omitempty"`
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

// FunctionRegexp matches a valid function name.
var FunctionRegexp = regexp.MustCompile("^[A-Za-z0-9-._]{1,31}$")
