load('json.star', 'json')

#print available functions in the module json
print(dir(json))

def MakeRequest() :
  return {
    "method" : "GET",
    "path" : "http://localhost:8080/api/v1/from/function",
    "headers" : {
      "Content-Type" : [
        "application/json"
      ]
    }
  }

repeat = 1
requests = json.ToJson([
  {
    "method" : "POST",
    "path" : "http://localhost:8080/api/v1/object",
    "headers" : {
      "Content-Type" : [
        "application/json"
      ]
    },
    "body" : json.ToJson({
      "key" : "value"
    })
  },
  {
    "method" : "POST",
    "path" : "http://localhost:8080/api/v1/array",
    "headers" : {
      "Content-Type" : [
        "application/json"
      ]
    },
    "body" : json.ToJson([
      "item1",
      "item2"
    ])
  },
  MakeRequest()
])