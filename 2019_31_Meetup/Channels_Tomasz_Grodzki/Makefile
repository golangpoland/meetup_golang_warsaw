get_tex:
	# https://github.com/sbinet/present-tex
	present-tex *.slide > tmp.tex; \
	sed 's/\\usepackage\[latin1\]/\\usepackage\[utf8\]/g' tmp.tex > channel.tex ; \
	rm -f tmp.tex

get_pdf: get_tex
	pdflatex -shell-escape channel.tex

clear:
	