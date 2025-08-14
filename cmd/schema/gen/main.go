package main

import (
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
)

type Name struct {
	Name string `json:"name"`
	Sex  string `json:"sex,omitempty"`
}

func main() {
	jsch := makeJSONSchema[Name]()

	explore(jsch)

	person := Name{Name: "Siu Yin"}

	jschResolved, err := jsch.Resolve(nil)
	if err != nil {
		log.Fatal("resolve: ", err)
	}

	if err := jschResolved.Validate(&person); err != nil {
		log.Fatal("resolve error: ", err)
	}
	fmt.Println("validated")
}

func makeJSONSchema[T any]() *jsonschema.Schema {
	jsch, err := jsonschema.For[T]()
	if err != nil {
		log.Fatal("for: ", err)
	}
	return jsch
}

func explore(jsch *jsonschema.Schema) {
	fmt.Printf("%v %v:\n", jsch.Type, jsch.Properties)
	for _, r := range jsch.Required {
		fmt.Printf("\t%s: %s\n", r, jsch.Properties[r].Type)
	}
}
