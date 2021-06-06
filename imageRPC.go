package main

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"path/filepath"
	"gocv.io/x/gocv"
	"github.com/streadway/amqp"
)

//Function declared to parse errors
func failError(err error, errmsg string) {
	if err != nil {
		log.Fatalf("%s: %s", errmsg, err)
	}
}

//Helper functions to get and set pixel values for channel depth of 3
type Vecb []uint8

//Adapted from https://github.com/hybridgroup/gocv/issues/339
//Other Matrix methods to extract pixel color result in overflow
//Copies the CV_8U channels into val
func GetVecbAt(mat gocv.Mat, row int, col int) Vecb {
	ch := mat.Channels()
	val := make(Vecb, ch)

	for c := 0; c < ch; c++ {
		val[c] = mat.GetUCharAt(row, col * ch + c)
	}

	return val
}

//Sets 3 channels at index in matrix to val's 3 uint8 values
func (val Vecb) SetVecbAt(mat gocv.Mat, row int, col int) {
	ch := mat.Channels()

	for c := 0; c < ch; c++ {
		mat.SetUCharAt(row, col * ch + c, val[c])
	}
}

//Function to upscale an image
func scaleImg(filename string, factor int) string {
	//Filename path
	path := "/gorpc/Images/" + filename

	//Extract matrix of image file data
	pixels := gocv.IMRead(path, 1)

	//Get dimensions from picture
	dimensions := pixels.Size()
	height := dimensions[0]
	width := dimensions[1]

	//Get type of mat
	matType := pixels.Type()

	//Create a new mat to fill
	bigMat := gocv.NewMatWithSize(height * factor, width * factor, matType)


	//Created a wait group to sync row concurrency
	wg := sync.WaitGroup{}
	wg.Add(height)

	//Iterate through array in rows
	for i := 0; i < height; i++ {
		//Go Routine call to fill in rows of big matrix concurrently
		go func (i int){
			//Sync waitgroup
			defer wg.Done()
			for j := 0; j < width; j++ {
				val := GetVecbAt(pixels, i, j)

				//Iterate through larger matrix
				for a := i * factor; a < i * factor + factor; a++ {
					for b := j * factor; b < j * factor + factor; b++ {
							val.SetVecbAt(bigMat, a, b)
					}
				}
			}
		}(i)
	}

	//Big matrix is finished when all rows processed
	wg.Wait()

	//Remove extension from filename
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	//Rename the new scaled image
	newName := name + "S" + ext

	//New path
	newPath := "/gorpc/Images/" + newName

	//Save the new image from matrix in local directory
	gocv.IMWrite(newPath, bigMat)

	return newName
}

//Function to multiply copies of an image
func copyImg(filename string, rowX int, colX int) string {
	//Filename path
	path := "/gorpc/Images/" + filename

	//Extract matrix of image file data
	pixels := gocv.IMRead(path, 1)

	//Get dimensions from picture
	dimensions := pixels.Size()
	height := dimensions[0]
	width := dimensions[1]

	//Get type of mat
	matType := pixels.Type()

	//Create a new mat to fill
	bigMat := gocv.NewMatWithSize(height * rowX, width * colX, matType)

	//Created a wait group to sync filling images
	wg := sync.WaitGroup{}
	wg.Add(rowX * colX)

	//Fill in image copies by relative index on matrix
	for i := 0; i < rowX; i++{
		for j := 0; j < colX; j++{
			go func (i int, j int){
				//Decrement counter if an image copy is made
				defer wg.Done()

				//Iterate over original image and store in new index copy
				for x := 0; x < height; x++{
					for y := 0; y < width; y++{
						val := GetVecbAt(pixels, x, y)

						val.SetVecbAt(bigMat, (i*height) + x, (j*width) + y)
					}
				}
			}(i, j)
		}
	}

	//Wait till all copies are filled in
	wg.Wait()

	//Remove extension from filename
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	//Rename the new scaled image
	newName := name + "C" + ext

	//New path
	newPath := "/gorpc/Images/" + newName

	//Save the new image from matrix in local directory
	gocv.IMWrite(newPath, bigMat)

	return newName
}

//Function to recursively draw an image
func recurseImg(filename string) string {
	//Filename path
	path := "/gorpc/Images/" + filename

	//Extract matrix of image file data
	pixels := gocv.IMRead(path, 1)

	//Get dimensions from picture
	dimensions := pixels.Size()
	height := dimensions[0]
	width := dimensions[1]

	//Get type of mat
	matType := pixels.Type()

	//Create a new mat to fill
	bigMat := gocv.NewMatWithSize(height * 3, width * 5, matType)

	//Created a wait group to sync filling images
	wg := sync.WaitGroup{}
	wg.Add(3)

	//Concurrently call recursive drawing
	go func(){
		defer wg.Done()
		//Fill in two size 1 copies
		for x := 0; x < height; x++{
			for y := 0; y < width; y++{
				val := GetVecbAt(pixels, x, y)

				val.SetVecbAt(bigMat, x, y)
				val.SetVecbAt(bigMat, x, y + width)
			}
		}
	}()
	

	go func(){
		defer wg.Done()
		//Fill in a size 2 copy
		for x := 0; x < height; x++{
			for y := 0; y < width; y++{
				val := GetVecbAt(pixels, x, y)

				//Iterate through larger matrix
				for a := x * 2; a < x * 2 + 2; a++ {
					for b := y * 2; b < y * 2 + 2; b++ {
						val.SetVecbAt(bigMat, a + height, b)
					}
				}
			}
		}
	}()
	
	go func(){
		defer wg.Done()
		//Fill in a size 3 copy
		for x := 0; x < height; x++{
			for y := 0; y < width; y++{
				val := GetVecbAt(pixels, x, y)

				//Iterate through larger matrix
				for a := x * 3; a < x * 3 + 3; a++ {
					for b := y * 3; b < y * 3 + 3; b++ {
						val.SetVecbAt(bigMat, a, b + 2 * width)
					}
				}
			}
		}
	}()
	
	//Sync the drawing of image arrays
	wg.Wait()

	//Remove extension from filename
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	//Rename the new scaled image
	newName := name + "R" + ext

	//New path
	newPath := "/gorpc/Images/" + newName

	//Save the new image from matrix in local directory
	gocv.IMWrite(newPath, bigMat)

	return newName
}

//RPC Server component adapted from https://www.rabbitmq.com/tutorials/tutorial-six-go.html
//Main RPC procedure
func main(){
	//Connect to local host server
	port, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failError(err, "Failed to connect to RabbitMQ.")
	defer port.Close()

	//Create a channel on port
	ch, err := port.Channel()
	failError(err, "Failed to open a channel on port.")
	defer ch.Close()

	//Declare a callback queue
	queue, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable flag
		false,       // auto delete flag
		false,       // exclusive flat
		false,       // no-wait flag
		nil,         // extra arguments
	)
	failError(err, "Failed to create a queue for channel.")

	//Set QoS
	err = ch.Qos(
		1,
		0,
		false,
		)
	failError(err, "Failed to set channel QoS.")

	//Register a client
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failError(err, "Failed to register a consumer")

	//Channel repeats receiving to keep server running
	loop := make(chan bool)

	//Process message and reply
	go func() {
		for d := range msgs {
			//Receive a message from client
			phrase := string(d.Body)

			//Convert phrase into fields
			phraseArr := strings.Fields(phrase)

			//Extract array elements
			a, _ := strconv.Atoi(phraseArr[1])
			b, _ := strconv.Atoi(phraseArr[2])
			c, _ := strconv.Atoi(phraseArr[3])

			var response string

			//Perform an image processing function based on input
			switch{
			case a == 1:
				log.Printf("Scaling %s\n", phraseArr[0])
				response = scaleImg(phraseArr[0], b)
			case a == 2:
				log.Printf("Copying %s\n", phraseArr[0])
				response = copyImg(phraseArr[0], b, c)

			case a == 3:
				log.Printf("Recursing %s\n", phraseArr[0])
				response = recurseImg(phraseArr[0])
			}


			err = ch.Publish(
				"",          // exchange
				d.ReplyTo, // routing key
				false,       // mandatory
				false,       // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(response),
				})
			failError(err, "Failed to publish a reply.")

			d.Ack(false)
		}
	}()

	//Waiting on messages
	log.Printf("Waiting on RPC Message")
	<-loop
}