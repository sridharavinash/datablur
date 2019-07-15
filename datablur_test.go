package datablur

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const filewithpath = "/tmp/datablur.csv"

func setupFile() {
	f, err := os.OpenFile(filewithpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("foo1,bar1\nabc1,def1\nonlyinfile,someval\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func TestSubstituteBlur(t *testing.T) {
	lookup := map[string]string{
		"foo":       "bar",
		"abc":       "def",
		"onlyinmem": "some_mem_val",
	}

	t.Run("valid substitution", func(t *testing.T) {
		s := &Substitute{lookupTable: lookup}
		got, ok := s.blur("foo")
		want := "bar"

		assert.Equal(t, want, got)
		assert.True(t, ok)
	})
	t.Run("invalid substitution", func(t *testing.T) {
		s := &Substitute{lookupTable: lookup}
		got, ok := s.blur("invalid")
		want := "invalid"

		assert.Equal(t, want, got)
		assert.False(t, ok)

	})

	t.Run("valid substitution in file", func(t *testing.T) {
		setupFile()
		s := &Substitute{lookupFile: filewithpath}
		got, ok := s.blur("foo1")
		want := "bar1"
		assert.Equal(t, want, got)
		assert.True(t, ok)
	})
	t.Run("invalid substitution in file", func(t *testing.T) {
		setupFile()
		s := &Substitute{lookupFile: filewithpath}
		got, ok := s.blur("fooinvalid")
		want := "fooinvalid"
		assert.Equal(t, want, got)
		assert.False(t, ok)
	})
	t.Run("set lookuptable and lookupFile, finds the value in lookuptable", func(t *testing.T) {
		setupFile()
		s := &Substitute{lookupFile: filewithpath, lookupTable: lookup}
		got, ok := s.blur("onlyinmem")
		want := "some_mem_val"
		assert.Equal(t, want, got)
		assert.True(t, ok)
	})
}

func TestRot13Blur(t *testing.T) {
	t.Run("valid substitution", func(t *testing.T) {
		r := &Rot13{}
		got, ok := r.blur("Vg jbexf 13!")
		want := "It works 13!"

		assert.Equal(t, want, got)
		assert.True(t, ok)
	})
}
