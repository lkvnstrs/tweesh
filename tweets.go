package main

import (
    "encoding/base64"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type Timestamp time.Time

type Tweet struct {
    Id      string  `json:"id_str"`
    Text    string  `json:"text"`
}

// Used to avoid recusion in UnmarshalJSON for Tweet.
type tweet Tweet

// String returns a string version of a Tweet.
func (t *Tweet) AsWords() []Word {

    text := strings.TrimSpace(TrimMentions(t.Text))
    split := strings.Split(text, " ")
    words := make([]Word, 0)

    for _, s := range split {

        if s == "" || s == " " {
            continue
        }

        words = append(words, Word{s: s, id: t.Id})
    }

    return words
}

// UnmarshalJSON implements the json.Unmarshaller interface
// for a Tweet.
// func (t *Tweet) UnmarshalJSON(b []byte) error {
//     var tmp tweet

//     if err := json.Unmarshal(b, &tmp); err == nil {
//         return err
//     }

//     t.Id = tmp.Id
//     t.Text = TrimMentions(tmp.Text)

//     return nil
// }

// UnmarshalJSON implements the json.Unmarshaller interface
// for a Timestamp.
func (t *Timestamp) UnmarshalJSON(b []byte) error {

    s := string(b)
    v, err := time.Parse(time.RubyDate, s[1:len(s)-1])
    if err != nil {
        return err
    }

    *t = Timestamp(v)
    return nil
}

// TweetGetter is an interface for objects that can GetTweets.
type TweetGetter interface {

    // GetTweets retrieves the last n tweets for a given user. 
    GetTweets(username string, n int) []Tweet
}

// tweetGetter is an implementation of the TweetGetter interface.
type tweetGetter struct {
    Token string
}

func NewTweetGetter(key, secret string) (*tweetGetter, error) {
    token, err := GetBearerToken(key, secret)
    if err != nil {
        return nil, err
    }

    return &tweetGetter{Token: token}, nil
}

// GetTweets retrieves the last n tweets for a given user.
func (tf *tweetGetter) GetTweets(username string, n int) ([]Tweet, error) {

    var tweets []Tweet
    client := &http.Client{}
    url := "https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=" + username

    req, err := http.NewRequest("GET", url + "&count=" + strconv.Itoa(n), nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("Authorization", "Bearer " + tf.Token)

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if err = json.Unmarshal(body, &tweets); err != nil {
        return nil, err
    }

    return tweets, nil
}

// GetBearerToken returns a bearer token given a 
// consumer key and secret.
func GetBearerToken(key, secret string) (string, error) {

    // Format the credentials
    bearer_cred := []byte(key + `:` + secret)
    encoded_cred := base64.StdEncoding.EncodeToString(bearer_cred)

    client := &http.Client{}

    req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token?grant_type=client_credentials", nil)
    if err != nil {
        return "", err
    }

    req.Header.Add("Authorization", "Basic " + encoded_cred)

    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // Unwrap JSON
    var b struct { AccessToken string `json:"access_token"` }
    
    if err := json.Unmarshal(body, &b); err != nil {
        return "", err
    }

    return b.AccessToken, nil
}

// TrimMentions trims mentions from the beginning of a tweet.
func TrimMentions(s string) string {

    var skipToNextSpace bool = false

    at := []rune("@")[0]
    space := []rune(" ")[0]
    rune_s := []rune(s)

    for i,r := range rune_s {

        if skipToNextSpace {
            if r == space {
                skipToNextSpace = false
            }
            continue
        }

        if r == at {
            skipToNextSpace = true
            continue
        }

        return string(rune_s[i:])
    }

    return s
}
