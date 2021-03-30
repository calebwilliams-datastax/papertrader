package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

/*
{
    "data": {
        "users": {
            "values": [
                {
                    "id": "16170764946e37",
                    "name": "system",
                    "email": "system@papertrader.com",
                    "created": "2021-03-29T22:54:54.100674-05:00"
                }
            ]
        }
    }
}
*/

func GenerateID() string {
	n := strconv.Itoa(int(time.Now().Unix()))
	r := strconv.Itoa(rand.Int())
	hash := sha256.Sum256([]byte(r))
	readable := hex.EncodeToString(hash[:2])
	id := fmt.Sprintf("%s%s", n, readable)
	return id
}

func GenerateKey() string {
	id := GenerateID()
	hash := sha256.Sum256([]byte(id))
	key := hex.EncodeToString(hash[:16])
	return key
}

func ToJson(v interface{}) string {
	r, err := json.Marshal(&v)
	if err != nil {
		return ""
	}
	return string(r)
}

func InsertMutation(obj, table string, v interface{}) string {
	values, ret := GraphQLValues(v)
	return fmt.Sprintf(`{"query": "mutation insert {%s: insert%s(value: %s) {value %s}}"}`, obj, table, values, ret)
}

func UpdateByID(obj, table string, v interface{}) string {
	values, ret := GraphQLValues(v)
	return fmt.Sprintf(`{"query" : "mutation updateOne%s {%s: update%s(value: "%s", ifExists: true ) {value %s}}"}`, obj, obj, table, values, ret)

}

func QueryAll(table string, v interface{}) string {
	_, values := GraphQLValues(v)
	return fmt.Sprintf(`{"query": "query all {%s (value:{}){values %s}}"}`, table, values)
}

func QueryByID(obj, id string, v interface{}) string {
	_, ret := GraphQLValues(v)
	return fmt.Sprintf(`{"query" : "query byid {%s (value: {id:\"%s\"}) {values %s}"}`, obj, id, ret)
}

func QueryByValues(table string, query map[string]string, v interface{}) string {
	q, _ := GraphQLValues(query)
	_, values := GraphQLValues(v)
	return fmt.Sprintf(`{"query" : "query byvalues {%s (value:%s) {values %s }}"}`, table, q, values)
}

func DeleteByID(obj, id string) string {
	return fmt.Sprintf(`{"query" : "mutation delete%s {PaP: delete%s(value: {id:\"%s\"}, ifExists: true ) {value {id}}}"}`, obj, obj, id)
}

func GraphQLValues(v interface{}) (string, string) {
	fmt.Printf("parings : %+v\n", v)
	//graphql formatting hack
	//everythings strings... so thats fun.
	i := strings.Builder{}
	b := strings.Builder{}
	m := map[string]string{}
	marshal, _ := json.Marshal(v)
	_ = json.Unmarshal(marshal, &m)
	b.WriteString("{")
	i.WriteString("{")
	for key, value := range m {
		i.WriteString(fmt.Sprintf(`%s `, key))
		fmt.Printf(`%s:"%s"%s`, key, value, "\n")
		b.WriteString(fmt.Sprintf(`%s:\"%s\"`, key, value))
	}
	b.WriteString("}")
	i.WriteString("}")
	return b.String(), i.String()
}

/*
mutation deleteOneBook {
  PaP: deletebook(value: {title:"Pride and Prejudice", author: "Jane Austen"}, ifExists: true ) {
    value {
      title
    }
  }
}
query oneBoko {
    book (value: {title:"Moby Dick"}) {
      values {
      	title
      	author
      }
    }
}
*/
