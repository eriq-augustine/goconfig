package goconfig;

import (
   "encoding/json"
   "io/ioutil"
   "os"
   "testing"

   "github.com/stretchr/testify/assert"

   log "github.com/eriq-augustine/golog"
)

func writeOptionsFile(file *os.File) (bool, error) {
   encoder := json.NewEncoder(file);

   err := encoder.Encode(options);
   if (err != nil) {
      log.ErrorE("Unable to encode options", err);
      return false, err;
   }

   return true, nil;
}

func copyMap(oldMap map[string]interface{}) map[string]interface{} {
   var newMap map[string]interface{} = make(map[string]interface{});
   for key, value := range oldMap {
      newMap[key] = value;
   }
   return newMap;
}

func assertMapEquals(assert *assert.Assertions, expected map[string]interface{}, actual map[string]interface{}) {
   assert.Equal(len(expected), len(actual));

   for key, value := range(actual) {
      // Need to check the maps (objects) specially.
      mapVal, ok := value.(map[string]interface{});
      if (ok) {
         expectedMapVal, ok := expected[key].(map[string]interface{});

         if (!assert.True(ok, "Expected value is a map where actual value is not")) {
            continue;
         }

         assertMapEquals(assert, expectedMapVal, mapVal);
      } else {
         // Other type mismatches will fall through to here (including expected map where actual is not).
         assert.EqualValues(expected[key], value);
      }
   }
}

func TestLoadFile(t *testing.T) {
   assert := assert.New(t);

   options["int"] = 1;
   options["bool"] = true;
   options["string"] = "string";
   options["type"] = map[string]interface{}{"A":1};

   oldOptions := copyMap(options);

   file, err := ioutil.TempFile(os.TempDir(), "config_test");
   if (!assert.Nil(err)) {
      assert.FailNow("Could not get temp file: %v", err);
   }

   defer file.Close();
   defer os.Remove(file.Name());

   ok, err := writeOptionsFile(file);
   assert.True(ok);
   assert.Nil(err);

   Reset();

   LoadFile(file.Name());
   assertMapEquals(assert, oldOptions, options);
}

func TestGet(t *testing.T) {
   assert := assert.New(t);

   options["int"] = 1;
   options["bool"] = true;
   options["string"] = "string";
   options["type"] = struct{A int}{1};

   assert.Equal(1, Get("int"));
   assert.Equal(true, Get("bool"));
   assert.Equal("string", Get("string"));
   assert.Equal(struct{A int}{1}, Get("type"));
}

func TestBasicGetTyped(t *testing.T) {
   assert := assert.New(t);

   options["int"] = 1;
   options["bool"] = true;
   options["string"] = "string";

   assert.Equal(1, GetInt("int"));
   assert.Equal(true, GetBool("bool"));
   assert.Equal("string", GetString("string"));
}

func TestGetDefault(t *testing.T) {
   assert := assert.New(t);

   options["int"] = 1;
   options["bool"] = true;
   options["string"] = "string";
   options["type"] = struct{A int}{1};

   // Non-default value.
   assert.Equal(1, GetIntDefault("int", 99));
   assert.Equal(true, GetBoolDefault("bool", false));
   assert.Equal("string", GetStringDefault("string", "otherString"));
   assert.Equal(struct{A int}{1}, GetDefault("type", nil));

   // Default value.
   assert.Equal(99, GetIntDefault("foo", 99));
   assert.Equal(false, GetBoolDefault("bar", false));
   assert.Equal("otherString", GetStringDefault("baz", "otherString"));
   assert.Nil(GetDefault("clown", nil));
}


func TestHas(t *testing.T) {
   assert := assert.New(t);

   options["int"] = 1;
   options["bool"] = true;
   options["string"] = "string";
   options["type"] = struct{A int}{1};

   assert.Equal(true, Has("int"));
   assert.Equal(true, Has("bool"));
   assert.Equal(true, Has("string"));
   assert.Equal(true, Has("type"));
}

func TestReset(t *testing.T) {
   assert := assert.New(t);

   options["int"] = 1;
   options["bool"] = true;
   options["string"] = "string";
   options["type"] = struct{A int}{1};

   Reset();

   assert.Equal(false, Has("int"));
   assert.Equal(false, Has("bool"));
   assert.Equal(false, Has("string"));
   assert.Equal(false, Has("type"));
}
