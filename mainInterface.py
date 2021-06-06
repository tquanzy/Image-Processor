import os
import numpy as np
from PIL import Image
import parse
import pika
import uuid

#RPC Client component adapted from https://www.rabbitmq.com/tutorials/tutorial-six-python.html
class ImageRpcClient(object):
    def __init__(self):
        self.connection = pika.BlockingConnection(
            pika.ConnectionParameters(host='rabbitmq'))

        self.channel = self.connection.channel()

        result = self.channel.queue_declare(queue='', exclusive=True)
        self.callback_queue = result.method.queue

        self.channel.basic_consume(
            queue=self.callback_queue,
            on_message_callback=self.on_response,
            auto_ack=True)

    def on_response(self, ch, method, props, body):
        if self.corr_id == props.correlation_id:
            self.response = body

    def call(self, ImageStr):
        self.response = None
        self.corr_id = str(uuid.uuid4())
        self.channel.basic_publish(
            exchange='',
            routing_key='rpc_queue',
            properties=pika.BasicProperties(
                reply_to=self.callback_queue,
                correlation_id=self.corr_id,
            ),
            body=ImageStr)
        while self.response is None:
            self.connection.process_data_events()

        return self.response.decode('utf8')

#Main function call
def main():
    #Set up RPC client
    imageProcess_rpc = ImageRpcClient()

    #Array to pass to C to filter instructions
    x = np.zeros((3), dtype=np.intc)

    #Ask for user input image name
    print("Please specify the filename (including extension) of a local image file you would like to process.")
    name = input("> ")

    #Error check if file exists
    path = './Images/' + name

    if not os.path.isfile(path):
        print("This image file does not exist in local directory.")
        return

    #Ask for user instruction
    print("Please input an image processing instruction.")
    print("Your options are: scale x, copy x x, recurse")
    print("x is a placeholder for an integer value (ie. copy 3 2)")

    instruct = input("> ")

    #Call C function to parse instruction
    parse.readInstruction(instruct, x)

    #Error check
    if (-1 in x):
        print("Your instruction is not valid.")
        return

    #Bundle message as a string
    message = name + " " + str(x[0]) + " " + str(x[1]) + " " + str(x[2])

    #Call RPC function to process message
    response = imageProcess_rpc.call(message)
    print(response, "created.")

    #Deprecated command because of Docker container permissions
    #Show the image created
    #with Image.open("./Images/" + response) as img:
    #    img.show()

if __name__ == "__main__":
    main()
