make:
	swig -python parse.i
	python3 setup.py build_ext --inplace

clean:
	rm -rf *.o *.so __pycache__ build parse.py parse_wrap.c