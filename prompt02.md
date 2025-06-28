You are a Golang developer. 

Add a function to main.go that take a string as an input and return an array of string as an ouput. 
The input is a string that describes a web path.
The output array will be Datadog metrics name.
The string contains dynamic range that are identified by < and >. 
For instance a valid value for input is /bm/branding/buying-guide/landing-pages/<uuid:landing_page_uuid>
When the input param is /bm/branding/buying-guide/landing-pages/<uuid:landing_page_uuid> 
then the function returns an array with :
- get_bm/branding/buying-guide/landing-pages/_uuid:landing_page_uuid
- post_bm/branding/buying-guide/landing-pages/_uuid:landing_page_uuid
- options_bm/branding/buying-guide/landing-pages/_uuid:landing_page_uuid

For each input param, the returned array size is 3 elements.
The first element is for a GET HTTP request.
The second element is for a POST HTTP request.
The third element is for a OPTIONS HTTP request.

The input dynamic param <uuid:landing_page_uuid> is transformed to _uuid:landing_page_uuid

The function should be unit tested with different values.


