FROM python:3
RUN pip install requests
COPY python-requests.py /
ENTRYPOINT [ "python" ]
