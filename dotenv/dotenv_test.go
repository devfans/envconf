package dotenv_test

import (
	"os"
	"testing"

	"github.com/devfans/envconf/dotenv"
)

func TestEnv(t *testing.T) {
	t.Log(os.Getenv("a"))
	t.Log(os.Getenv("b"))
	t.Log(dotenv.Int("a"))
	t.Log(dotenv.Uint("b"))
	t.Log(dotenv.Bool("c"))
	t.Log(dotenv.Bool("d"))
	t.Log(os.Getenv("test"))
	t.Log(dotenv.String("test"))
	t.Log(dotenv.EnvConf().Get("a"))
	t.Log(dotenv.EnvConf().Get("b"))
}