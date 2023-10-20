from flask import Flask, request, Response
import requests
import os

app = Flask(__name__)

@app.route('/auth', methods=['GET'])
def auth():
    access_token = request.headers.get('X-Forwarded-Access-Token')
    user = request.headers.get('X-Forwarded-User')

    auth_url = os.getenv("ACCESS_TOKEN_ENDPOINT")
    res = requests.get(auth_url, headers={"Authorization": f"{access_token}"})
    if res.status_code != 200:
        return Response(status=401)
    env_name = os.getenv("USER")

    if user != env_name:
        return Response(status=401)

    return Response(status=200)

if __name__ == '__main__':
    app.run(port=5000)