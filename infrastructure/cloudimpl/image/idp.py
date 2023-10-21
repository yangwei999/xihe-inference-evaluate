import os

import requests

from jupyter_server.auth import IdentityProvider, User
from tornado.log import app_log


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
        ak = handler.request.headers.get("X-Auth-Request-Access-Token")
        if ak is None:
            app_log.error("missing access token")
            return None

        auth_url = os.getenv("ACCESS_TOKEN_ENDPOINT")
        if auth_url is None:
            app_log.error(f"missing auth url")
            return None

        res = requests.get(auth_url, headers={"Authorization": f"{ak}"})
        if res.status_code != 200:
            app_log.error(f"get user failed: {res.status_code}")
            return None

        user = res.json().get("username")
        if user is None:
            app_log.error(f"missing user data")
            return None

        env_name = os.getenv("USER")
        app_log.info(f"{user} vs {env_name}")
        return XiheUser(env_name) if env_name == user else None

c = get_config()  # noqa

c.ServerApp.identity_provider_class = XiheIdentityProvider
c.ServerApp.IdentityProvider.token = ""
c.ServerApp.password = ""