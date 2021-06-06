#Setup SWIG interface file
from distutils.core import *
from distutils      import sysconfig

#Importing numpy arrays
import numpy

#Get the numpy include directory
try:
    numpy_include = numpy.get_include()
except AttributeError:
    numpy_include = numpy.get_numpy_include()
  
#Parse extension module
_parse = Extension(name='_parse',sources=["parse.i","parse.c"], include_dirs=[numpy_include],) 
  
#Setup stage
setup(name="parse", 
      version="1.0", 
      ext_modules=[_parse]) 