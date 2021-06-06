# Image Processor
 Multi-language implementation of an image processor

Description:

The overall goal of the project is to emulate a shell for image processing with different languages. I created a project in Python, C, and Go that can take in an image filename (placed in Images folder) through a text UI and can process the image to return a result image (scaling the image to be larger, copying the image in many iterations in a matrix, or recursing the image like a golden ratio diagram). The idea is a proof in concept for the importance of multi-language communication in applications such as web graphics or mobile displays, since the transformation of large image data can be very taxing and should be passed off to a systems language for faster processing rather than a scripting language that codes a website such as Javascript or Python.

Languages:

I used Python, C, and Go in my project. Python was used as the main user interface of the system. The Python script asks for user input filename for an image file as well as a specific command (ie. scale 5 or copy 7 9 or recurse). This is to represent the front end of an application that needs to manipulate an image (ie. a website builder that needs to resize images for different screen formats). The Python code then calls on C functions to parse the string instruction into a numeric array. C is used for processing and error handling of text instructions in this system. Go is the backend of my system. It is an RPC server that receives messages from Python containing a filename and an int array representing a processing request. The Go code then calls upon OpenCV libraries to parse the image file into a matrix format. I manipulate the matrix integers using goroutines to speed up the processing of the image into larger sizes. The resultant image is then saved in Images folder and the name is returned to the Python client.

Methods:

I called upon C functions using SWIG to interface C code into a Python library. This code is contained in: parse.c, parse.h, parse.i, setup.py, and makefile. Python communicates with the Go code through the use of a message queue using RabbitMQ. The Go server receives RPC requests from the Python client and replies with the finished filename.

Instructions to Run:

\*\*\*If running docker commands results in a &quot;Couldn&#39;t connect to Docker daemon&quot; error, all docker commands must be prepended with sudo (ie. &quot;sudo docker-compose up -d&quot;) to execute properly\*\*\*

1. Copy my Project Directory &quot;Final&quot;

2. &quot;cd Final&quot;

3. &quot;docker-compose up -d&quot;

This creates the three Docker containers, Rabbitmq server, the Go RPC server, and the Python interface client, in the background. The first time running this line will install all required dependencies and take some time.

Future runs of this command will be much faster!

4. &quot;docker exec -it pyinterface /bin/bash&quot;

5. &quot;cd pyin&quot;

6. &quot;python3 mainInterface.py&quot;

You may have to wait ~10 seconds if the first time calling this code causes a RabbitMQ error. That is because this code cannot execute until the Rabbitmq server is fully set up in the background.

7. Follow Python prompts

To Exit Testing:

8. &quot;exit&quot;

9. &quot;docker-compose down&quot;