"Auth to Go"
============

Composition of security libraries joined into single modular server. 

There are a lot solutions for authentication and security, but it is 
problem to find free maintained authentication server which may work 
with custom DB schema and easily can be adopted. 
"Auth to Go" is convenient wrap intend to solve all this problems. 

philosophy behind the project: everything is packed and ready to 
be "eaten". 

Authentication flows
====================

Login Page (or End User Authentication)
---------------------------------------

[Password hashing and Salting](https://www.maketecheasier.com/password-hashing-encryption/)
[Oauth2 vision on end user authentication](https://oauth.net/articles/authentication/)

Password (secret) OAuth2
------------------------

https://www.oauth.com/oauth2-servers/access-tokens/password-grant/

In terms of OAuth2 - it means get token by providing username and password. 
Main use case is login end user. 

If you are building your own Identity Provider or similar thing, it may require to 
authenticat some mobild application before granting access to some other application. 

Another use case may be just login page: 
    1. User accesses login page and sends credentials 
    2. Server running the login page adds token to its request to the Auth Server 
    3. Auth server verifies that login request is done from trusted application and generates token
        for user web application. 

    

