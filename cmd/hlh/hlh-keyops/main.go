package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/helpers/keyops/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
	"github.com/number571/hidden-lake/internal/utils/keys"
)

var (
	gFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of application"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about application"),
		flag.NewFlagBuilder("-g", "--generate").
			WithDescription("generate (private / public) key pair"),
		flag.NewFlagBuilder("-p", "--path").
			WithDescription("set path to (private / public) key pair").
			WithDefinedValue("."),
		flag.NewFlagBuilder("-o", "--output").
			WithDescription("output encapsulated secret file").
			WithDefinedValue(settings.CPathEncKey),
		flag.NewFlagBuilder("-e", "--encapsulate").
			WithDescription("encapsulate secret by public key").
			WithDefinedValue(""),
		flag.NewFlagBuilder("-d", "--decapsulate").
			WithDescription("decapsulate secret by private key").
			WithDefinedValue(""),
	).Build()
)

func main() {
	args := os.Args[1:]
	if ok := gFlags.Validate(args); !ok {
		fmt.Fprintln(os.Stderr, "args invalid")
		os.Exit(1)
	}

	if countArgs(args) > 1 {
		fmt.Fprintln(os.Stderr, "more then one arg")
		os.Exit(2)
	}

	if gFlags.Get("-v").GetBoolValue(args) {
		fmt.Println(build.GetVersion())
		return
	}

	if gFlags.Get("-h").GetBoolValue(args) {
		help.Println(settings.CAppFullName, settings.CAppDescription, gFlags)
		return
	}

	inputPath := strings.TrimSuffix(gFlags.Get("-p").GetStringValue(args), "/")
	if err := os.MkdirAll(inputPath, 0700); err != nil {
		fmt.Fprintln(os.Stderr, "mkdir path")
		os.Exit(3)
	}

	if gFlags.Get("-g").GetBoolValue(args) {
		if err := generateKeyPair(inputPath); err != nil {
			fmt.Fprintln(os.Stderr, "generate key pair")
			os.Exit(4)
		}
		return
	}

	if gFlags.Get("-e").GetBoolValue(args) {
		encPubKeyPath := gFlags.Get("-e").GetStringValue(args)
		pubKey, err := readPubKey(encPubKeyPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "load public key")
			os.Exit(5)
		}
		encKey, secret, err := pubKey.GetKEMPubKey().Encapsulate()
		if err != nil {
			fmt.Fprintln(os.Stderr, "encapsulate op")
			os.Exit(6)
		}
		outEncKeyPath := gFlags.Get("-o").GetStringValue(args)
		if err := os.WriteFile(outEncKeyPath, encKey, 0600); err != nil {
			fmt.Fprintln(os.Stderr, "write encapsulated key")
			os.Exit(7)
		}
		fmt.Printf("%X", secret)
		return
	}

	if gFlags.Get("-d").GetBoolValue(args) {
		decPrvKeyPath := gFlags.Get("-d").GetStringValue(args)
		privKey, _, err := loadKeyPair(inputPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "load key pair (not generated?)")
			os.Exit(8)
		}
		encKey, err := os.ReadFile(decPrvKeyPath) //nolint:gosec
		if err != nil {
			fmt.Fprintln(os.Stderr, "read encapsulated key")
			os.Exit(9)
		}
		secret, err := privKey.GetKEMPrivKey().Decapsulate(encKey)
		if err != nil {
			fmt.Fprintln(os.Stderr, "decapsulate op")
			os.Exit(10)
		}
		fmt.Printf("%X", secret)
		return
	}

	help.Println(settings.CAppFullName, settings.CAppDescription, gFlags)
}

func countArgs(args []string) int {
	count := 0
	actions := []bool{
		gFlags.Get("-v").GetBoolValue(args),
		gFlags.Get("-h").GetBoolValue(args),
		gFlags.Get("-g").GetBoolValue(args),
		gFlags.Get("-e").GetBoolValue(args),
		gFlags.Get("-d").GetBoolValue(args),
	}
	for _, exist := range actions {
		if exist {
			count++
		}
	}
	return count
}

func readPubKey(pathToPubKey string) (asymmetric.IPubKey, error) {
	pubKeyStr, err := os.ReadFile(pathToPubKey) //nolint:gosec
	if err != nil {
		return nil, err
	}
	pubKey := asymmetric.LoadPubKey(string(pubKeyStr))
	if pubKey == nil {
		return nil, errors.New("load public key") //nolint:err113
	}
	return pubKey, nil
}

func generateKeyPair(inputPath string) error {
	keyPath := filepath.Join(inputPath, settings.CPathKey)
	if _, err := os.Stat(keyPath); !os.IsNotExist(err) {
		return errors.New("private key already exist") //nolint:err113
	}
	privKey := asymmetric.NewPrivKey()
	if err := os.WriteFile(keyPath, []byte(privKey.ToString()), 0600); err != nil {
		return err
	}
	pubPath := filepath.Join(inputPath, settings.CPathPubKey)
	if _, err := keys.GetPubKey(privKey, pubPath); err != nil {
		return err
	}
	return nil
}

func loadKeyPair(inputPath string) (asymmetric.IPrivKey, asymmetric.IPubKey, error) {
	keyPath := filepath.Join(inputPath, settings.CPathKey)
	privKeyStr, err := os.ReadFile(keyPath) //nolint:gosec
	if err != nil {
		return nil, nil, err
	}
	privKey := asymmetric.LoadPrivKey(string(privKeyStr))
	if privKey == nil {
		return nil, nil, errors.New("load private key") //nolint:err113
	}
	pubKeyPath := filepath.Join(inputPath, settings.CPathPubKey)
	pubKey, err := keys.GetPubKey(privKey, pubKeyPath)
	if err != nil {
		return nil, nil, err
	}
	return privKey, pubKey, nil
}
