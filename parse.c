#include <stdio.h>
#include <string.h>
#include <stdlib.h>

//Function for parsing instructions
void readInstruction(char* string, int* output, int n){
	//Incorrect array length passed in
	if (n != 3){
		output[0] = -1; //Change first number to negative to prompt error
		return;
	}

	//Counter number of phrases
	int phraseCount = 0;
	for (int i = 0; i < (int) strlen(string); i++){
		//Count phrases by incrementing when previous character is a space but current is not
		if ((i == 0 || string[i-1] == ' ') && string[i] != ' '){
			phraseCount++;
		}
	}

	//Allocate space for storing array of strings
	char** instruction;
	instruction = malloc((phraseCount) * sizeof(char*));
	for (int i = 0; i < phraseCount; i++){
		instruction[i] = (char*) malloc(100 * sizeof(char));
	}

	//Iterate through the string, skipping spaces and copying phrases
	const char skip[2] = " ";
	char* phrase;
	phrase = strtok(string, skip);

	int index = 0;
	while (phrase != NULL)
	{
		strcpy(instruction[index], phrase);
		//printf("%s\n", instruction[index]);
		phrase = strtok(NULL, skip);
		index++;
	}
    
	//Check first instruction and assign array values accordingly
	if (strcmp(instruction[0], "scale") == 0){
		output[0] = 1;
		if (atoi(instruction[1]) != 0){ //Check for numeric
			output[1] = atoi(instruction[1]);
		} else {
			output[1] = -1;
		}
	} else if (strcmp(instruction[0], "copy") == 0){
		output[0] = 2;
		if (atoi(instruction[1]) != 0){ 
			output[1] = atoi(instruction[1]);
		} else {
			output[1] = -1;
		}
		if (atoi(instruction[2]) != 0){
			output[2] = atoi(instruction[2]);
		} else {
			output[2] = -1;
		}
	} else if (strcmp(instruction[0], "recurse") == 0){
		output[0] = 3;
	} else { //Invalid instruction input
        output[0] = -1;
    }

	//Free pointers
	for (int i = 0; i < phraseCount; i++){
      free(instruction[i]);
    }
    free(instruction);

	return;
}

/*//Testing functionality of function call to parse instructions
int main(){
	int nums[3] = {0, 0, 0};
	char temp[] = "copy 3 3";
	readInstruction(temp, nums, 3);
	for (int i = 0; i < 3; i++){
		printf("%d\n", nums[i]);
	}
	return 0;
}*/
