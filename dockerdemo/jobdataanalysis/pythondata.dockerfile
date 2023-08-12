FROM jupyter/datascience-notebook
WORKDIR /app/
COPY demo.py .
ENTRYPOINT ["python"]
CMD ["demo.py"]

