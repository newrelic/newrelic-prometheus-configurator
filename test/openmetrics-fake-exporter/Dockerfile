FROM python:slim-buster

WORKDIR /app

COPY . .
RUN pip3 install -r requirements.txt

CMD [ "python3", "fake-exporter.py"]
