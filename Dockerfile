FROM python:3.9.5-slim-buster

RUN pip install pipenv

ADD ./app /application
ADD requirements.txt /application

WORKDIR /application

RUN pip install -r requirements.txt

RUN pip install gunicorn[gevent]

EXPOSE 5000

CMD gunicorn --worker-class gevent --workers 8 --bind 0.0.0.0:5000 wsgi:app --max-requests 10000 --timeout 5 --keep-alive 5 --log-level info
