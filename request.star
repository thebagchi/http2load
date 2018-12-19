load('json.star', 'json')

def MakeRequest() :
  return {
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