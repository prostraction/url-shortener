package memory

import (
	"log"
	"strings"
	"testing"
)

func OneLetterDubTest(hashMap map[string]string, t *testing.T) {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-/_"
	for i := 0; i < 1000*len(alphabet); i++ {
		want := strings.Repeat(string(rune(alphabet[i%len(alphabet)])), i/len(alphabet)+1)
		hash, err := ToHash(hashMap, want)
		if err != nil {
			t.Fatal(err)
		}
		url, err := FromHash(hashMap, hash)
		if err != nil {
			t.Fatal(err)
		}
		if url != want {
			t.Fatal("hash mismatch!", url, " != ", want)
		}
	}
}
func MoveAlphabetTest(hashMap map[string]string, t *testing.T) {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-/_"
	for i := 0; i < 300000; i++ {
		want := (alphabet + alphabet)[:100]
		want = want[i%100:] + want[:i%100]
		hash, err := ToHash(hashMap, want)

		if err != nil {
			t.Fatal(err, want)
		}
		url, err := FromHash(hashMap, hash)
		if err != nil {
			t.Fatal(err)
		}
		if url != want {
			t.Fatal("hash mismatch!", url, " != ", want)
		}
		if i%len(alphabet) == 0 {
			runeAlphabet := []rune(alphabet)
			char := runeAlphabet[i/len(alphabet)%len(alphabet)]
			if char < 126 {
				runeAlphabet[i/len(alphabet)%len(alphabet)]++
			} else {
				runeAlphabet[i/len(alphabet)%len(alphabet)] = 'a'
			}
			alphabet = string(runeAlphabet)
		}
	}
}
func TwoLetterTest(hashMap map[string]string, t *testing.T) {
	wantArr := []string{"qo", "po", "ro", "fo", "fa", "pa", "qa", "pu", "qu", "qupu", "fafo", "qaqu", "qo-", "po-", "ro-", "fo-", "fa-", "pa-", "qa-", "pu-", "qu-", "qupu-", "fafo-", "qaqu-"}
	for _, word := range wantArr {
		for i := 0; i < 500; i++ {
			want := word
			want = strings.Repeat(want, i+1)
			hash, err := ToHash(hashMap, want)
			if err != nil {
				t.Fatal(want, hash, err)
			}
			url, err := FromHash(hashMap, hash)
			if err != nil {
				t.Fatal(url, hash, err)
			}
			if url != want {
				t.Fatal("hash mismatch!", url, " != ", want)
			}
		}
	}
}

func SentenceTest(hashMap map[string]string, t *testing.T) {
	wantArr := []string{
		"testing word", "word testing",
		"testing-word", "word-testing",
		"one word", "word one",
		"one-word", "word-one",
		"para rapa", "rapa para",
		"para-rapa", "rapa-para"}
	for _, want := range wantArr {
		hash, err := ToHash(hashMap, want)
		if err != nil {
			t.Fatal(want, hash, err)
		}
		url, err := FromHash(hashMap, hash)
		if err != nil {
			t.Fatal(url, hash, err)
		}
		if url != want {
			t.Fatal("hash mismatch!", url, " != ", want)
		}
	}
}

func TestHashTableText(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	hashMap := make(map[string]string)
	log.Println("OneLetterDubTest")
	OneLetterDubTest(hashMap, t)
	log.Println("MoveAlphabetTest")
	MoveAlphabetTest(hashMap, t)
	log.Println("TwoLetterTest")
	TwoLetterTest(hashMap, t)
	log.Println("SentenceTest")
	SentenceTest(hashMap, t)
}
