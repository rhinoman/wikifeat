Technical Overview
==================

This guide is intended to give a high-level overview of the key components and technologies used by Wikifeat.

Architecture
------------

Wikifeat uses a [microservices](https://wikipedia.org/wiki/Microservices) architecture, with each 
service focusing on a few specific tasks.  This allows for considerable flexibilty when it comes to 
flexibility and scalability.  For instance, multiple instances of each service may be run on one or 
more machines with minimal configuration needed, and only the services experiencing high load need 
be duplicated.

The services are able to 'discover' one another by querying a service broker.  As of this writing, 
Wikifeat uses (and requires) [Etcd](https://github.com/coreos/etcd) as its service broker.  Each 
service maintains a local cache (which is refreshed periodically) with the locations of other 
services in the system.  When a service needs to make a request to a different service, it selects 
one of the locations in its local cache in a load-balanced fashion

The services communicate with each other via HTTP using REST APIs.  There are pros and cons to 
this method, with one of the chief benefits being flexibility; additional 'plugin' services can be 
implemented in virtually any programming language (The core services are written in 
[Go](http://golang.org).  Another benefit of handling communication this way is faster development, 
due to the fact that most developers are already familiar with communication via HTTP, which would 
not be the case with something like [0mq](http://zeromq.org).  The primary con of using HTTP for 
inter-service communication is there some amount of overhead associated with making HTTP requests 
that would not be present using a lower-level socket-based communication method.  We believe the 
gains in flexibility and accessibility to be gained by using HTTP outweigh the relatively minor 
cost in performance.

Of course, it is not difficult to imagine custom Wikifeat installations running numerous custom 
plugins that communicate amongst themselves using other protocols, condescending to use HTTP only 
when making requests to the core services :)

Data Persistence
----------------

The Wikifeat core services use [CouchDB](http://couchdb.apache.org) to provide data persistence 
(primarily, wiki documents) and user management.  CouchDB is a document-oriented NoSQL database that 
fits well with the 'wiki' data model.  It's great scalability features make a good fit for creating 
highly-available systems with Wikifeat.  It's user and session management features are also 
leveraged to provide user authentication and authorization.  

Custom services are of course free to use other database solutions for their own data.  CouchDB might 
not be the best choice for, say, Geospatial Data, not to mention creating plugins to 'integrate' 
other systems with Wikifeat.

Web Application
---------------

Wikifeat includes a front-end web application written using [Backbone.js](http://backbonejs.org) 
and [Marionette](http://marionettejs.org).  The web application may be extended via plugins.  
Plugins may add additional front-end functionality by, for example, serving as the user interface 
for a custom back-end service.  Also, content plugins may be used to allow the embedding of additional 
media/data-types within wiki pages.


