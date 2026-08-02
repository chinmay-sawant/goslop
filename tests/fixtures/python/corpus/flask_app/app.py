from flask import Flask, request, send_file, send_from_directory

app = Flask(__name__)
ROOT = "/var/data"


@app.route("/dl")
def dl_file():
    path = request.args["path"]
    return send_file(path)


@app.route("/safe")
def dl_dir():
    name = request.args["name"]
    return send_from_directory(ROOT, name)
