import time
from typing import Any
from typing import Callable as callable
from functools import partial

import httpx

METHODS = ["get", "options", "head", "post", "put", "patch", "delete"]
SESSION = httpx.Client()

def _req(
        url: str,
        ra: callable[[httpx.Response], int]=None,
        method: str = "get",
        session: httpx.Client = SESSION,
        *args: list[Any],
        **kwargs: dict[str, Any]
    ) -> httpx.Response:
    """Custom request function with retry after capabilities for 429s.

    Args:
        url (str): URL to send the request to.
        ra (Callable[[httpx.Response], int], optional): Retry after function,
            receives the Response object and returns the seconds before
            retrying. Defaults to None.
        method (str, optional): The request method. Defaults to "get".
        session (httpx.Client, optional): Session client. Defaults to SESSION.

    Returns:
        httpx.Response: Response object.
    """
    resp = getattr(session, method)(url, follow_redirects=True, *args, **kwargs)
    if resp.status_code == 429:
        if ra is not None:
            time.sleep(ra(resp))
            return _req(url, ra, method, session, *args, **kwargs)
    return resp

class req:
    pass

req_dict = {
    req: {},
}

for r, kw in req_dict.items():
    for i in METHODS:
        setattr(r, i, partial(_req, method=i, **kw))