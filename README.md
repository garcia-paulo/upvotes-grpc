# upvotes-grpc

Hello, thank you for looking at this humble repository. The project is intended as a coding test for a job application at [Klever](https://klever.io/).
The project is a gRPC upvotes service, where a authenticated user can access posts created by other users and upvote them. Pretty simple right?

# Starting up

If you just want to access the API. You can do so via the ip address: **_164.92.121.144:8080_**.

Alternatively, the project is read to be built and the commands needed to run the project are in the _Makefile_ in the root directory. If you don't have Make installed just go to `/server` and run `go run .`. Of course, you'll need **go** installed.

To test the endpoint I used [BloomRPC](https://github.com/bloomrpc/bloomrpc) which functions as Postman or Imsomnia, but for gRPCs.
All you'll need to do is:

- Import the protos
- Create a user
- Put the returned token in the metadata (the tab at the bottom) as **"authorization"**

![image](https://user-images.githubusercontent.com/79415003/165133620-e374b652-edcd-4873-b640-ff4fa53b2178.png)

And that's it, you'll have access to all of post methods. It's simple but I hope you'll like it.
