# StockProfiling
Application that creates stock profile for the user and displays the current stock value when requested. 

##Objective 
This application allows the user to perform the following,
1. Buying Stocks: Creates the user profile given the stock type that user wants to procure 
2. Displaying Stocks: User can view his profile and the current value of the stock

##How to Execute?

1. Change the directory to the destination folder
2. Run server and client seperately.

<pre> 
To run the server :
go run server.go

To run the client:
go run client.go 
</pre>

##About the application 
1. Buying Stocks:
<pre>
Run the client program using the command,
go run client.go GOOG:50%,YHOO:50% 2000

In the above command, GOOG:50% => refers to the stock type and percentage of the amount to be invested 
                      2000     => refers to the total budget for purchasing the stocks
                      
The client program can be executed without the command line arguments. 
Just follow the instructions after running the client program
</pre>

2. Displaying Stocks:
<pre>
go run client.go 12345 

In the above command, 12345 => refers to the user id/client id/trade id 

When the above command is executed, user portfolio will be displayed on the prompt. 

The response from running the above command :
Stocks : GOOG:100+$520.25,YHOO:200:-$30.40 

The +/- sign will be displayed in the response if there is a change in current stock value

</pre>

##Notes 
1. The folder stockprofiling_http uses HTTP Json rpc interface in the server for buying stocks and checking client portfolio and 
stockprofiling_tcp uses TCP interface for the application. 
2. Gorilla mux library is used for dispatching the incoming requests from the user


