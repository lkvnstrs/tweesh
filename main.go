package main

import (
    "bufio"
    "flag"  
    "fmt"
    "math/rand"
    "os"
    "time"
    "sync"
)

func main() {
   
    numWords := flag.Int("words", 40, "maximum number of words to print")
    prefixLen := flag.Int("prefix", 2, "prefix lenght in words")

    flag.Parse()

    rand.Seed(time.Now().UnixNano())

    key := os.Getenv("TWITTER_KEY")
    secret := os.Getenv("TWITTER_SECRET")

    if key == "" || secret == "" {
        panic("Credentials not set.")
    }

    tg, err := NewTweetGetter(key, secret)
    if err != nil {
        fmt.Println(err)
    }

    c := NewChain(*prefixLen)
    // accounts := []string{"arjunblj", "lkvnstrs", "ctbeiser", "v_mohankumar", "phrmsilva", "ZZZZUnit"}
    accounts := []string{"2chainz", "TheRock", "LilTunechi", "Drake", "FrencHMonTanA", "kanyewest", "SnoopDogg", "yungleann"}

    var wg sync.WaitGroup

    wg.Add(len(accounts))

    for _, a := range accounts {
        go func(a string) {
            tweets, _ := tg.GetTweets(a, 200)
            c.AddTweets(&tweets)
            wg.Done()
        }(a)
    }
    
    wg.Wait()

    in := ""
    reader := bufio.NewReader(os.Stdin)

    for in != "exit" {
        text := c.Generate(*numWords)
        fmt.Println(text)
        fmt.Println("Press enter to coninute (exit to end)...")
        in, _ = reader.ReadString('\n')
    }
}