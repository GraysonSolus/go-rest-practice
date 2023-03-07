# go-rest-practice

Things to do next:

1. Consider refactoring the controller/service breakdown. Perhaps controller should contain logic that parses request data, rather than separate service.
2. Abstract away some of the service layer concrete usages to make it more testable.
3. Initialize a proper logger in main.
4. Configure endpoint properties through environment variables. Currently things like min/max height and width, the server address and port, etc., are hardcoded - want to avoid at all costs.
5. Additional validation and checks - check file size, maybe require checksums passed in the request. 
6. Auth - implement JWT unary interceptor to validate user data
