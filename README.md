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

    
Oauth2 with 3 leg auth
----------------------

1. AuthClient1 (Service Provider) Requests credentials. Commonly it is clientID and ClientSecret stored securily.
2. AuthClient2 (End User) Creates credentials. Client id, password, 2 factor token. 
3. AuthClient1 requests authorized redirect url to pass to end user to make authorization. 
    Both clients dealing with same Identity Provider 
    /oauth2/authorize endpoint 

4. Once authorized we get RefreshToken to get access to AccessToken (ServiceProvider stores refresh token securily)
5. AuthClient1 Requests access token from IDP on behalf of AuthClient2
    /oauth2/token endpoint 

6. IDP checks if refresh token still valid and issues access token


RefreshToken doesn't expire because invalidated only by time of unpredictable event. 
It should be secure random. But anyway we have to have generator for this. 

AccessToken is expiring one and used to authetnicate against some services. 
    Service may everytime ask IDP to validate the access token 
    Service may trust some public key than token should be signed by this public key (JWT)


Once we want to replace AuthClient1 with something bigger and abstract we need to store RefreshToken somewhere and share. 
For example, our AuthClient1 got 3rd party RefreshToken and shared with us. 

    Side notice 
        1. Our user should have been authenticated with our IDP and third party IDP  
        2. Our user requests refresh token for us and let us know it. 
            It happens actually by requestion our api do call to 3rd party IDP and when that IDP gets confirmation from our user 
            we are obtaining the RefreshToken 
            




