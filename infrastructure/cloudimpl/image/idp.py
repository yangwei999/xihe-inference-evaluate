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

        pool_id = os.getenv("USER_POOL")
        if pool_id == None:
            app_log.error(f"missing pool id")
            return None

        res = requests.get(auth_url, headers={"authorization": f"{ak}", "x-authing-userpool-id": f"{pool_id}"})
        if res.status_code != 200:
            app_log.error(f"get user failed: {res.status_code}")
            return None

        user_data = res.json().get("data")
        if user_data is None:
            app_log.error(f"missing user data in response")
            return None

        user = user_data.get("username")
        if user is None:
            app_log.error(f"missing username in user_data")
            return None

        env_name = os.getenv("USER")
        if env_name == user:
            return XiheUser(env_name)

        app_log.warning(f"requester {user} missmatch")
        return None

c = get_config()  # noqa

c.ServerApp.identity_provider_class = XiheIdentityProvider
c.ServerApp.IdentityProvider.token = ""
c.ServerApp.password = ""