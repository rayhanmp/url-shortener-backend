package db

import (
	"os"
	"github.com/joho/godotenv"
	supa "github.com/nedpals/supabase-go"
)

var Supabase *supa.Client

func CreateClient() {
	err:= godotenv.Load()

	url:= os.Getenv("SUPABASE_URL")
	key:= os.Getenv("SUPABASE_KEY")
	Supabase = supa.CreateClient(url, key)
}
