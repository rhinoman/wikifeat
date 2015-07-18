User Service
============

The User Service API provides an interface for managing 
users.

GET /users
  - Returns a list of users
	- Query Parameters
	  - pageNum - Page Number
		- numPerPage - Number of records to return

POST /users
  - Create a new user

GET /users/{user-id}
  - Returns a user record

PUT /users/{user-id}
  - Update a user record
	- Header Parameters
	  - If-Match - The last revision of the user record

PUT /users/{user-id}/grant_role
  - Grants a role to a user

PUT /users/{user-id}/revoke_role
  - Revokes a user role

DELETE /users/{user-id}
  - Deletes a user 

POST /users/login
  - Creates a new user session

DELETE /users/login
  - Destroys a user session (ie., logout)

GET /users/current_user
  - Returns the user record associated with the provided
	  authentication header(s)

