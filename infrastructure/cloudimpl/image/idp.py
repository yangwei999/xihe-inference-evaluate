import os

import requests
from typing import Dict, List
from unittest import mock

from jupyter_server.auth import Authorizer, IdentityProvider, User
from tornado.httputil import url_concat
from tornado.log import app_log
from traitlets import Instance, default

class XiheUser(User):
    """Subclass User to store JupyterHub data"""

    # not a dataclass field

    def __init__(self, name):
        super().__init__(username=name)

class XiheIdentityProvider(IdentityProvider):
    """Identity Provider for JupyterHub OAuth

    Replacement for JupyterHub's HubAuthenticated mixin
    """

    def get_user(self, handler):
        name = handler._headers.get("X-Forwarded-User")
        ak = handler._headers.get("X-Forwarded-Access-Token")
        if name is None:
            print("missing user")
            return None
        if ak is None:
            print("missing access token")
            return None
        auth_url = os.getenv("ACCESS_TOKEN_ENDPOINT")
        if auth_url:
            print(f"missing auth url")
            return None
        res = requests.get(auth_url, headers={"Authorization": f"{ak}"})
        if res.status_code != 200:
            print(f"get user failed: {res.status_code}")
            return None
        env_name = os.getenv("USER")

        return XiheUser(env_name) if env_name == name else None

c = get_config()  # noqa

c.ServerApp.identity_provider_class = XiheIdentityProvider
c.ServerApp.IdentityProvider.token = ""
c.ServerApp.password = ""