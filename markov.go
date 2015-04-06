// Adapted from "Codewalk: Generating arbitrary text: a Markov chain algorithm"
// https://golang.org/doc/codewalk/markov/

package main

import (
    "math/rand"
    "strings"
)

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
    chain       map[string][]string
    prefixLen   int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
    return &Chain{make(map[string][]string), prefixLen}
}

// AddTweets parses the given Tweets into the Markov chain.
func (c *Chain) AddTweets(tweets *[]Tweet) {

    for _, t := range *tweets {

        line := TrimMentions(t.Text)
        if len(line) < 20 {
            continue
        }
        
        p := make(Prefix, c.prefixLen)

        for _, s := range strings.Split(line, " ") {
            key := p.String()
            c.chain[key] = append(c.chain[key], s)
            p.Shift(s)
        }
    }
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
    p := make(Prefix, c.prefixLen)
    var words []string
    for i := 0; i < n; i++ {
        choices := c.chain[p.String()]
        if len(choices) == 0 {
            break
        }

        next := choices[rand.Intn(len(choices))]
        words = append(words, next)
        p.Shift(next)
    }

    return strings.Join(words, " ")
}