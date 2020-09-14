FROM ubuntu:18.04

LABEL com.example.version="testy"
LABEL vendor1="testi"
LABEL vendor2=ZENITH\ Incorporated
LABEL com.example.release-date="2015-02-12"
LABEL com.example.version.is-production=""

COPY . /app
RUN make /app
CMD python /app/app.py
