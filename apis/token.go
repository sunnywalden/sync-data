package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// Helper func:  Read input from specified file or stdin
func loadData(p string) ([]byte, error) {
	if p == "" {
		return nil, fmt.Errorf("No path specified")
	}

	var rdr io.Reader
	if p == "-" {
		rdr = os.Stdin
	} else if p == "+" {
		return []byte("{}"), nil
	} else {
		if f, err := os.Open(p); err == nil {
			rdr = f
			defer f.Close()
		} else {
			return nil, err
		}
	}
	return ioutil.ReadAll(rdr)
}

// Print a json object in accordance with the prophecy (or the command line options)
func printJSON(j interface{}) error {
	var out []byte
	var err error

	if *flagCompact == false {
		out, err = json.MarshalIndent(j, "", "    ")
	} else {
		out, err = json.Marshal(j)
	}

	if err == nil {
		fmt.Println(string(out))
	}

	return err
}


// Verify a token and output the claims.  This is a great example
// of how to verify and view a token.
func verifyToken() error {
	// get the token
	tokData, err := loadData(*flagVerify)
	if err != nil {
		return fmt.Errorf("Couldn't read token: %v", err)
	}

	// trim possible whitespace from token
	tokData = regexp.MustCompile(`\s*$`).ReplaceAll(tokData, []byte{})
	if *flagDebug {
		fmt.Fprintf(os.Stderr, "Token len: %v bytes\n", len(tokData))
	}

	// Parse the token.  Load the key from command line option
	token, err := jwt.Parse(string(tokData), func(t *jwt.Token) (interface{}, error) {
		data, err := loadData(*flagKey)
		if err != nil {
			return nil, err
		}
		if isEs() {
			return jwt.ParseECPublicKeyFromPEM(data)
		} else if isRs() {
			return jwt.ParseRSAPublicKeyFromPEM(data)
		}
		return data, nil
	})

	// Print some debug data
	if *flagDebug && token != nil {
		fmt.Fprintf(os.Stderr, "Header:\n%v\n", token.Header)
		fmt.Fprintf(os.Stderr, "Claims:\n%v\n", token.Claims)
	}

	// Print an error if we can't parse for some reason
	if err != nil {
		return fmt.Errorf("Couldn't parse token: %v", err)
	}

	// Is token invalid?
	if !token.Valid {
		return fmt.Errorf("Token is invalid")
	}

	// Print the token details
	if err := printJSON(token.Claims); err != nil {
		return fmt.Errorf("Failed to output claims: %v", err)
	}

	return nil
}

// Create, sign, and output a token.  This is a great, simple example of
// how to use this library to create and sign a token.
func signToken() error {
	// get the token data from command line arguments
	tokData, err := loadData(*flagSign)
	if err != nil {
		return fmt.Errorf("Couldn't read token: %v", err)
	} else if *flagDebug {
		fmt.Fprintf(os.Stderr, "Token: %v bytes", len(tokData))
	}

	// parse the JSON of the claims
	var claims jwt.MapClaims
	if err := json.Unmarshal(tokData, &claims); err != nil {
		return fmt.Errorf("Couldn't parse claims JSON: %v", err)
	}

	// add command line claims
	if len(flagClaims) > 0 {
		for k, v := range flagClaims {
			claims[k] = v
		}
	}

	// get the key
	var key interface{}
	key, err = loadData(*flagKey)
	if err != nil {
		return fmt.Errorf("Couldn't read key: %v", err)
	}

	// get the signing alg
	alg := jwt.GetSigningMethod(*flagAlg)
	if alg == nil {
		return fmt.Errorf("Couldn't find signing method: %v", *flagAlg)
	}

	// create a new token
	token := jwt.NewWithClaims(alg, claims)

	// add command line headers
	if len(flagHead) > 0 {
		for k, v := range flagHead {
			token.Header[k] = v
		}
	}

	if isEs() {
		if k, ok := key.([]byte); !ok {
			return fmt.Errorf("Couldn't convert key data to key")
		} else {
			key, err = jwt.ParseECPrivateKeyFromPEM(k)
			if err != nil {
				return err
			}
		}
	} else if isRs() {
		if k, ok := key.([]byte); !ok {
			return fmt.Errorf("Couldn't convert key data to key")
		} else {
			key, err = jwt.ParseRSAPrivateKeyFromPEM(k)
			if err != nil {
				return err
			}
		}
	}

	if out, err := token.SignedString(key); err == nil {
		fmt.Println(out)
	} else {
		return fmt.Errorf("Error signing token: %v", err)
	}

	return nil
}

// showToken pretty-prints the token on the command line.
func showToken() error {
	// get the token
	tokData, err := loadData(*flagShow)
	if err != nil {
		return fmt.Errorf("Couldn't read token: %v", err)
	}

	// trim possible whitespace from token
	tokData = regexp.MustCompile(`\s*$`).ReplaceAll(tokData, []byte{})
	if *flagDebug {
		fmt.Fprintf(os.Stderr, "Token len: %v bytes\n", len(tokData))
	}

	token, err := jwt.Parse(string(tokData), nil)
	if token == nil {
		return fmt.Errorf("malformed token: %v", err)
	}

	// Print the token details
	fmt.Println("Header:")
	if err := printJSON(token.Header); err != nil {
		return fmt.Errorf("Failed to output header: %v", err)
	}

	fmt.Println("Claims:")
	if err := printJSON(token.Claims); err != nil {
		return fmt.Errorf("Failed to output claims: %v", err)
	}

	return nil
}

func isEs() bool {
	return strings.HasPrefix(*flagAlg, "ES")
}

func isRs() bool {
	return strings.HasPrefix(*flagAlg, "RS") || strings.HasPrefix(*flagAlg, "PS")
}