# upvotes-grpc
Hello, thank you for looking at this humble repository. The project is intended as a coding test for a job application at [Klever](https://klever.io/).
The project is a gRPC upvotes service, where a authenticated user can access posts created by other users and upvote them. Pretty simple right?

# Starting up
The project is read to be built and the commands needed to run the project are in the *Makefiles* through the project. If you don't have Make installed just go to `/server` and run `go run .`. Of course, you'll need **go** installed.

To test the endpoint I used [BloomRPC](https://github.com/bloomrpc/bloomrpc) which functions as Postman or Imsomnia, but for gRPCs.
All you'll need to do is:
* Import the protos
* Create a user
* Put the returned token in the metadata (the tab on the bottom) as **"authorization"**

![BloomRPC](https://user-images.githubusercontent.com/79415003/165082011-af7bc8b7-21b5-4646-8e97-03464559c242.png)

And that's it, you'll have access to all post methods. It's simple but I hope you'll like it.
