// Adapted from "Codewalk: Generating arbitrary text: a Markov chain algorithm"
// https://golang.org/doc/codewalk/markov/

package main

import (
    "errors"
    "math/rand"
    "strings"
)

// Word is a single word in a Prefix and Chain.
// s is the string value of the Word.
// id is a unique identifier for the source of the Word.
type Word struct {
    s   string
    id  string
}

// Prefix is a Markov chain prefix of one of more words.
type Prefix []string

// String returns the Prefix as a string (for us as a mpa key).
func (p Prefix) String() string {
    return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
    copy(p, p[1:])
    p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixs to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can havemutliple suffixes.
type Chain struct {
    chain       map[string][]Word
    prefixLen   int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
    return &Chain{make(map[string][]Word), prefixLen}
}

// AddTweets parses the given WordSources into the Markov chain.
func (c *Chain) AddWords(words []Word) {
        
    if words == nil {
       return
    }

    p := make(Prefix, c.prefixLen)

    for _, w := range words {

        key := p.String()
        c.chain[key] = append(c.chain[key], w)
        p.Shift(w.s)
    }
}

// Generate returns a string of at most n words generated from Chain.
// Ensures a generated phrase has some source variation.
func (c *Chain) Generate(n int) (string, error) {
    
    var current_id string
    var words []string

    for {

        p := make(Prefix, c.prefixLen)
        found_id_variation := false

        for i := 0; i < n; i++ {

            choices := c.chain[p.String()]
            if len(choices) == 0 {
                break
            }

            next := choices[rand.Intn(len(choices))]

            words = append(words, next.s)

            if !found_id_variation {
                found_id_variation = current_id != "" && current_id != next.id
                current_id = next.id
            }

            p.Shift(next.s)
        }

        if found_id_variation {
            return strings.Join(words, " "), nil
        }

        // clear vars
        words = words[:0]
        current_id = ""
    }

    return "", errors.New("No unique strings found")    
    
}