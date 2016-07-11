
Wikifeat   
========

[![Build Status](https://travis-ci.org/rhinoman/wikifeat.svg?branch=master)](https://travis-ci.org/rhinoman/wikifeat)

### Introduction
Wikifeat is an open source collaboration platform built around the ever-popular [Wiki](http://wikipedia.org/wiki/Wiki) concept.  It is meant to be extensible and highly usable.  Many enterprise collaboration platforms may be powerful and extensible, but they are often difficult to develop for and have an unfriendly, overly complicated user experience.  Wikifeat uses a microservices architecture which not only facilitates scalability, but also allows for the addition of custom services to extend the core system.

![wikifeat_screenshot_sm](https://cloud.githubusercontent.com/assets/1859198/11432240/9207d3b8-9477-11e5-909f-fbf62e627e62.png)

### Goals

The overarching goals for the Wikifeat project are:

- **Usability** - The core Wikifeat system should provide a simple, intuitive wiki experience with a
clean UI an a simple, easy to learn markup language (i.e., markdown).
- **Extensibility** - It should be possible to extend Wikifeat with additional services and front-end
plugins to provide new functionality and integrate Wikifeat with other systems.
- **Scalability** - It should be possible to horizontally scale a Wikifeat installation without great
difficulty.
 
### Key Features

#### Wikis

- Allows for the creation of multiple 'wikis', each with its own access controls
- Markdown is used for the markup language.
- Users are able to comment on individual pages.

Wiki pages are edited using simple Markdown -- specifically, [Commonmark](http://commonmark.org/).  A screenshot of the editor interface is shown here:

![wikifeat_edit_sm](https://cloud.githubusercontent.com/assets/1859198/11432232/5fb70082-9477-11e5-904c-c3b5a83d0a82.png)

#### REST API
The wikifeat core services each expose a REST API to facilitate extensibility and integration with other services.

#### Plugin Support
Wikifeat has the ability to load javascript plugins to extend the functionality of the web application.

#### Scalability
Wikifeat's architecture is designed with scalability in mind.  The core system consists of a set of microservices registered to an Etcd service registry.  Multiple instances of each services may be running across multiple machines, with service discovery handled by your Etcd cluster.

### Status

I'm no longer actively developing Wikifeat and have moved on to other projects.  I will accept pull requests, so please feel free to contribute and/or fork.  

For an example of a running Wikifeat installation, see https://www.wikifeat.org

See the [Technical Overview][1] for a more detailed introduction.

Documentation
-------------

Please see the [Wikifeat](https://www.wikifeat.org) website for more information and documentation.  Documentation is also a work in progress.

Contributing
------------

Contributions are most welcome.  See https://www.wikifeat.org/app/wikis/wikifeat/pages/contributing for details.

  [1]: https://www.wikifeat.org/app/wikis/wikifeat/pages/technical-overview
