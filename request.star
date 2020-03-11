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
    },
    "expect" : 200
  }

# Global Script Variables
repeat = 1
requests = json.ToJSON([
  {
    "method" : "POST",
    "path" : "http://localhost:8080/api/v1/object",
    "headers" : {
      "Content-Type" : [
        "application/json"
      ]
    },
    "body" : json.ToJSON({
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
    "body" : json.ToJSON([
      "item1",
      "item2"
    ])
  },
  MakeRequest()
])