Wiki API
========

The Wiki API provides an interface to manage wiki records 
and pages.

GET /wikis
 - Returns a list of wikis 
 - Query Parameters
  - pageNum - Page Number
	- numPerPage - Number of records to return
	- memberOnly - Only show wikis for which user is a member

POST /wikis
 - Creates a new wiki
 
GET /wikis/{wiki-id}
 - Fetches the wiki record for the sepcified wiki-id

GET /wikis/slug/{wiki-slug}
 - Fetches a wiki record by its slug

PUT /wikis/{wiki-id}
 - Updates a wiki record
 - Header Parameters
  - If-Match - The last revision of the wiki record 

DELETE /wikis/{wiki-id}
 - Deletes a wiki record.  Also deletes all of the wiki's
   pages and data and removes the wiki's database.
