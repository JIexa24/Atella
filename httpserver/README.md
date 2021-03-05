/api/v1/test/code/200  
  Methods: Get  
  Return: Code 200  
  Params: None  

/api/v1/test/code/401  
  Methods: Get  
  Return: Code 401  
  Params: None  

/api/v1/test/code/404  
  Methods: Get  
  Return: Code 404  
  Params: None  

/api/v1/test/code/405  
  Methods: Get  
  Return: Code 405  
  Params: None  

/api/v1/test/code/500  
  Methods: Get  
  Return: Code 500  
  Params: None  

/api/v1/test/code/501  
  Methods: Get  
  Return: Code 501  
  Params: None  

/api/v1/ping
  Methods: Get  
  Return: JSON, {"response":"pong" }
  Params: None  

/api/v1/echo/method  
  Methods: Get  
  Return: JSON, {"response":"<Request method>"}   
  Params: None  
  
/api/v1/get/host  
  Methods: Get  
  Return: JSON, {"response":{"hostname":"<hostname>","version":"<version>"}}
  Params: None  
  
/api/v1/get/vector  
  Methods: Get  
  Return: JSON, {"response":{"list":{<vector>}}}  
  Params: None  

/api/v1/get/version  
  Methods: Get  
  Return: JSON, {"response":"<version>"}  
  Params: None  

/api/v1/get/hostname  
  Methods: Get  
  Return: JSON, {"response":"<hostname>"}  
  Params: None  

/api/v1/set/vector  
  Methods: POST
  Return: JSON, {"response":"<Error>"} if error was cause  
  Params:  
    - Security Code in header X-Atella-Auth  
    - Header Content-Type with value "application/x-www-form-urlencoded"  
    - Form parameter: hostname, that contains client hostname  
    - Form parameter: vector, that contains client vector in json format  