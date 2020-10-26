#!flask/bin/python
from flask import Flask

app = Flask(__name__)

@app.route('/')
def index():
    return "Welcome to my bockchain Node!"

if __name__ == '__main__':
    app.run(debug=True)