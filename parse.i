//Module name used in python
%module parse
//File contents copied to wrapper file
%{
	//Required for numpy integration
	#define SWIG_FILE_WITH_INIT
	#include "parse.h"
%}

//Use numpy arrays
%include "numpy.i"
%init %{
	import_array();
%}

%apply (char* STRING) {(char* string)}
%apply (int* INPLACE_ARRAY1, int DIM1) {(int* output, int n)}

%include "parse.h"